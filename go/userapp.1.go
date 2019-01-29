package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

// this is the version with index page as template
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.go.html"))
}

func main() {
	router := httprouter.New()
	router.GET("/", index)

	//	http.HandleFunc("/stats/", serveTemplate)

	log.Println("Listening on :8888 ...")
	http.ListenAndServe(":8888", router)
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	id := filepath.Join("templates", "index.go.html")
	odd := filepath.Join("templates", "inc_odometer.go.html")
	tpli, err := template.ParseFiles(tpi, id, odd)
	//log.Println(tpli)
	if err != nil {
		// Log the detailed error
		//      log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tpli.ExecuteTemplate(w, "index1", nil)
}
