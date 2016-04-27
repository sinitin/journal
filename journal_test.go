package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestLog(t *testing.T) {
	dbname = "journal_test"
	ClearDB()
	log("sini", "7", "Unclear")
	log("yvo", "15", "Wanted help with parking")

	db, err := sql.Open("mysql", "root:root@/"+dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var totalduration int
	err = db.QueryRow("select SUM(duration) from entry").Scan(&totalduration)
	if err != nil {
		fmt.Println(err)
	}

	if totalduration != 22 {
		t.Error("Test failed")
	}

}

func ExampleList() {
	dbname = "journal_test"
	ClearDB()
	log("sini", "7", "Unclear")

	db, err := sql.Open("mysql", "root:root@/"+dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	list()

	// Output: sini 7 Unclear
}

func ClearDB() {

	db, err := sql.Open("mysql", "root:root@/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(" Pinging the db didnt work  ")
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + dbname)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS entry (colleagueid integer, duration integer, reason varchar(32) )")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS person (colleagueid integer NOT NULL AUTO_INCREMENT, name varchar(32), PRIMARY KEY (colleagueid))")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DELETE FROM entry")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DELETE FROM person")
	if err != nil {
		panic(err)
	}

}
