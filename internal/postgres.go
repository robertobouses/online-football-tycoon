package internal

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	User     string
	Pass     string
	Host     string
	Port     string
	Database string
}

func NewPostgres(c DBConfig) (*sql.DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Pass, c.Host, c.Port, c.Database)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Println("error opening connection")
		return nil, err
	}
	log.Println("the connection has the correct credentials")
	return db, db.Ping()
}
