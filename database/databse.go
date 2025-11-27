package database

import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
    dsn := "root:@tcp(127.0.0.1:3306)/login_google?parseTime=true"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Failed to connect DB: %v", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("DB ping failed: %v", err)
    }

    DB = db
    log.Println("Database connected")
}
