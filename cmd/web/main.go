package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	retryAttempts := 0

	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready.")
			retryAttempts++
		} else {
			log.Println("connected to database!")
			return conn
		}

		if retryAttempts > 10 {
			return nil
		}

		log.Print("backing off for 1 second")
		time.Sleep(1 * time.Second)
		continue
	}

}

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to the database")
	}

	return conn
}

func main() {
	//connect to the database
	db := initDB()
	db.Ping()
}
