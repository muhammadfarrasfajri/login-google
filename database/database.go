package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectMySQL() {

	var err error
	DB, err = sql.Open("mysql","root:@tcp(127.0.0.1:3306)/login_google")
	if err != nil {
		log.Fatal("DB error: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB unreachable: ", err)
	}

	log.Println("MySQL connected.")
}
