package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rexlx/vapi/local/data"
	"github.com/rexlx/vapi/local/drivers"
)

type config struct {
	port    int
	logfile string
}

type settings struct {
	config       config
	infoLog      *log.Logger
	errorLog     *log.Logger
	models       data.Models
	environment  string
	runtimeParms string
}

func main() {
	var cfg config
	cfg.port = 8888
	cfg.logfile = "/Users/rexfitzhugh/vapi.log"

	file, err := os.OpenFile(cfg.logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	infoLog := log.New(file, "_info_ ", log.Ldate|log.Ltime)
	errorLog := log.New(file, "_error_ ", log.Ldate|log.Ltime)
	dsn := os.Getenv("DSN")
	environment := os.Getenv("ENV")
	runtimeParms := os.Getenv("CFG")

	db, err := drivers.GetPostgres(dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.SQL.Close()

	app := &settings{
		config:       cfg,
		infoLog:      infoLog,
		errorLog:     errorLog,
		models:       data.New(db.SQL),
		environment:  environment,
		runtimeParms: runtimeParms,
	}

	err = app.serve()
	if err != nil {
		log.Fatalln(err)
	}
}

func (app *settings) serve() error {
	app.infoLog.Printf("starting at %v on port..%v", time.Now(), app.config.port)
	app.infoLog.Println("logging to", app.config.logfile)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}
	return srv.ListenAndServe()
}
