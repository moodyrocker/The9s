package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "/home/vagrant/The9s/the9s")
	checkErr(err)

	// query
	rows, err := db.Query("select url.url, user.u_id, url.url_id from url inner join user on url.u_id = user.u_id")
	checkErr(err)
	var urlurl string
	var useru_id int
	var urlurl_id int

	for rows.Next() {
		err = rows.Scan(&urlurl, &useru_id, &urlurl_id)
		checkErr(err)
		fmt.Printf("%v %v %v\n", urlurl, useru_id, urlurl_id)
	}
	rows.Close() //good habit to close
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
