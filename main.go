package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"regexp"
	"strconv"
)

var dbname string = "journal"

func main() {
	checkDB()
	validateInput()

	if os.Args[1] == "log" {
		log(os.Args[2], os.Args[3], os.Args[4])
	} else if os.Args[1] == "list" {
		list()
	} else if os.Args[1] == "total" {
		total()
	} else if os.Args[1] == "hitlist" {
		hitlist()
	}
}

func printHelp() {
	fmt.Println("With this journal program you can keep track of who disturbs you and how much.\n")
	fmt.Println("journal list		Prints a list of all disturbances you have logged.")
	fmt.Println("journal hitlist		Prints a list of how many minutes each person have disturbed you.")
	fmt.Println("journal total		Prints the total number of minutes you have been disturbed.")
	fmt.Println("journal log		Logs a new disturbance, submit the arguments [name] [duration in minutes] [reason].")
	fmt.Println("			For example: journal log Sven 15 \"Wanted food again\"")
	os.Exit(0)
}

func printInvalidInput() {
	fmt.Println("You seem to have entered invalid input, please have a look at the manual.\n")
	printHelp()
}

func validateInput() {

	if len(os.Args) < 2 {
		printInvalidInput()
	} else if os.Args[1] == "log" {
		//check to make sure the user entered the correct number of arguments
		if len(os.Args) != 5 {
			printInvalidInput()
		}

		//check that the name only includes letters
		match, _ := regexp.MatchString("[[:alpha:]]", os.Args[2])
		if match == false {
			fmt.Println("The name you entered seems to not only contain letters")
			printInvalidInput()
		}

		//check to make sure the duration argument is a number
		if _, err := strconv.Atoi(os.Args[3]); err != nil {
			fmt.Println("The duration you have enters seems to no only have contained digits")
			printInvalidInput()
		}
	} else if (os.Args[1] == "help") || (os.Args[1] == "h") {
		printHelp()
	} else if (os.Args[1] == "list") || (os.Args[1] == "total") || (os.Args[1] == "hitlist") {
		return
	} else {
		printHelp()
	}

}

func log(name string, duration string, reason string) {

	db, err := sql.Open("mysql", "root:root@/"+dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//check if person exists in db
	var id int
	err = db.QueryRow("select colleagueid from person where name = ?", name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred, we will insert person
			_, err = db.Exec("INSERT person SET name = ?", name)
			if err != nil {
				panic(err)
			}
			//pick up the colleague id from the person just inserted
			err = db.QueryRow("select colleagueid from person where name = ?", name).Scan(&id)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	//insert disturbance
	_, err = db.Exec("INSERT entry SET colleagueid=?,duration=?,reason=?", id, duration, reason)
	if err != nil {
		panic(err)
	}
}

func total() {

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

	fmt.Println(totalduration)

}

func hitlist() {

	db, err := sql.Open("mysql", "root:root@/"+dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var colleagueid int
	var name string

	rows, err := db.Query("select * from person")

	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&colleagueid, &name)
		if err != nil {
			fmt.Println(err)
		}

		var total int
		err = db.QueryRow("select SUM(duration) from entry where colleagueid = ?", colleagueid).Scan(&total)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(name + " " + strconv.Itoa(total))

	}

	err = rows.Err()
	if err != nil {
		fmt.Println("Ojojoj3")
	}

}

func list() {
	db, err := sql.Open("mysql", "root:root@/"+dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var (
		id       int
		duration int
		reason   string
	)
	rows, err := db.Query("select * from entry")

	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &duration, &reason)
		if err != nil {
			fmt.Println(err)
		}

		var name string
		err = db.QueryRow("select name from person where colleagueid = ?", id).Scan(&name)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(name, duration, reason)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
}

func checkDB() {
	db, err := sql.Open("mysql", "root:root@/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
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
}
