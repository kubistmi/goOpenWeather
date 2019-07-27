package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	sqlFile, err := os.Open("connstr")
	if err != nil {
		log.Fatal(err)
	}

	sqlBytes, err := ioutil.ReadAll(sqlFile)
	if err != nil {
		log.Fatal(err)
	}

	SQLCONN := string(sqlBytes)

	db, err := sql.Open("postgres", SQLCONN)
	if err != nil {
		log.Fatal(err)
	}

	//age := 21
	//rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	rows, err := db.Query("SELECT * FROM pg_tables limit 10")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rows)

	db.Close()
}
