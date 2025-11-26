package database

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() error {
    dsn := "root:@tcp(127.0.0.1:3306)/login_google"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }

    DB = db
    return db.Ping()
}
