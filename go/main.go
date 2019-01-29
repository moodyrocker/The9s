package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

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
	router.GET("/stats/:name", serveTemplate)
	log.Println("Listening on :80 ...")
	http.ListenAndServe(":80", router)
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	lis := filepath.Join("templates", "inc_list.go.html")
	id := filepath.Join("templates", "index.go.html")
	odd := filepath.Join("templates", "inc_odometer.go.html")
	tpli, err := template.ParseFiles(tpi, id, odd, lis)
	//log.Println(tpli) asdihasnpi
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

func serveTemplate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s := strings.Split(ps.ByName("name"), "_")
	counnt := s[0]
	vvvalue := []string{counnt, "_main.html"}
	// Join with no separator to create a compact string.
	joined := strings.Join(vvvalue, "")
	fmt.Println(joined)

	mf := filepath.Join("templates", "stats", joined)
	lp := filepath.Join("templates", "layout1.go.html")
	fp := filepath.Join("templates", "stats", ps.ByName("name"))
	tp := filepath.Join("templates", "inc.go.html")
	//    log.Println(fp)
	tpl, err := template.ParseFiles(lp, fp, tp, mf)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tpl.ExecuteTemplate(w, "layout3", nil)

}
