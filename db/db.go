package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
    db, err := sql.Open("postgres", "postgres://ec2admin:ec2mylove@localhost:5432/courses?sslmode=disable")
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}