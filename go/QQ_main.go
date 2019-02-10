package main

import (
	"database/sql"
	"fmt"

	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "/home/vagrant/The9s/the9s")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}

type Book struct {
	url_id		string
	url   		string
	short_url  string
	u_id		string
}

func main() {
	http.HandleFunc("/books", booksIndex)
	http.ListenAndServe(":80", nil)
}

func booksIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM url")
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer rows.Close()

	bks := make([]Book, 0)
	for rows.Next() {
		bk := Book{}
		err := rows.Scan(&bk.url_id, &bk.url, &bk.short_url, &bk.u_id) // order matters
		if err != nil {
			http.Error(w, http.StatusText(305), 305)
			return
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(200), 200)
		return
	}

	for _, bk := range bks {
		fmt.Fprintf(w, "%s, %s, %s, %s\n", bk.url_id, bk.url, bk.short_url, bk.u_id)
	}
}
