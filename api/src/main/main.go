package main

import (
	"encoding/json"
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
	logpath string
}

type settings struct {
	config       RuntimeParms
	infoLog      *log.Logger
	errorLog     *log.Logger
	models       data.Models
	environment  string
	runtimeParms string
}

type RuntimeParms struct {
	Logpath string `json:"logpath"`
	Port    int    `json:"port"`
	Env     string `json:"env"`
}

func main() {
	dsn := os.Getenv("DSN")
	environment := os.Getenv("ENV")
	runtimeConfig := os.Getenv("CFG")

	var cfg config
	var config RuntimeParms
	contents, err := os.ReadFile(runtimeConfig)
	if err != nil {
		log.Fatalln(err)
	}
	// this is where we unmarshal the contents into config
	err = json.Unmarshal(contents, &config)
	if err != nil {
		log.Fatalln(err)
	}
	cfg.port = 8888
	cfg.logpath = "/Users/rexfitzhugh/vapi.log"

	file, err := os.OpenFile(config.Logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	infoLog := log.New(file, "_info_ ", log.Ldate|log.Ltime)
	errorLog := log.New(file, "_error_ ", log.Ldate|log.Ltime)

	db, err := drivers.GetPostgres(dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.SQL.Close()

	app := &settings{
		config:       config,
		infoLog:      infoLog,
		errorLog:     errorLog,
		models:       data.New(db.SQL),
		environment:  environment,
		runtimeParms: runtimeConfig,
	}

	err = app.serve()
	if err != nil {
		log.Fatalln(err)
	}
}

func (app *settings) serve() error {
	app.infoLog.Printf("starting at %v on port..%v", time.Now(), app.config.Port)
	app.infoLog.Println("logging to", app.config.Logpath)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.Port),
		Handler: app.routes(),
	}
	return srv.ListenAndServe()
}
