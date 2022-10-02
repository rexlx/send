package drivers

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var conn = &DB{}

const maxConns = 11
const maxIdle = 5
const ttl = 5 * time.Minute

func GetPostgres(source string) (*DB, error) {
	d, err := sql.Open("pgx", source)
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(maxConns)
	d.SetMaxIdleConns(maxIdle)
	d.SetConnMaxLifetime(ttl)

	err = testDB(d)
	if err != nil {
		return nil, err
	}
	conn.SQL = d
	return conn, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	fmt.Println("db connection verfied")
	return err
}
