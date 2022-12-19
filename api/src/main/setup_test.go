package main

import (
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rexlx/vapi/local/data"
)

var testApp settings
var mockDB sqlmock.Sqlmock

func TestMain(m *testing.M) {
	testDb, theMock, _ := sqlmock.New()
	mockDB = theMock
	defer testDb.Close()

	testApp = settings{
		config:      RuntimeParms{},
		infoLog:     log.New(os.Stdout, "testInf ", log.LUTC),
		errorLog:    log.New(os.Stdout, "testErr ", log.LUTC),
		models:      data.New(testDb), //the test db satisfies the sql.db
		environment: "dev",
	}

	os.Exit(m.Run())
}
