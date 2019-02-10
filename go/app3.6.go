package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

// this is the version with index page as template
var tpl *template.Template
var db *sql.DB
var err error

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.go.html"))

	//db, err = sql.Open("postgres", "postgres://bond:password@localhost/bookstore?sslmode=disable")
	db, err = sql.Open("sqlite3", "/home/vagrant/The9s/the9s")
	fmt.Println("You connected to your database.")

}

type uurl struct {
	Url       string
	Short_url string
}

func main() {
	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("static/"))
	router.GET("/", index)
	router.GET("/faq", faq)
	router.GET("/about", about)

	router.POST("/form", forminsert)
	router.GET("/form", form)

	router.GET("/tandc", tc)
	router.GET("/pp", pp)
	router.GET("/contact", contact)
	router.GET("/stats/:name", serveTemplate)
	log.Println("Listening on :80 ...")
	http.ListenAndServe(":80", router)
}

func forminsert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	lis := filepath.Join("templates", "inc_list.go.html")
	id := filepath.Join("templates", "form.go.html")

	url := r.FormValue("url")
	su := r.FormValue("short_url")

	tpli, err := template.ParseFiles(tpi, id, lis)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		// Log the detailed error
		//      log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tpli.ExecuteTemplate(w, "form", uurl{url, su})
}

func form(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	lis := filepath.Join("templates", "inc_list.go.html")
	id := filepath.Join("templates", "form.go.html")

	tpli, err := template.ParseFiles(tpi, id, lis)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		// Log the detailed error
		//      log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tpli.ExecuteTemplate(w, "form", nil)
}

// https://www.calhoun.io/intro-to-templates-p2-actions/
// https://www.thepolyglotdeveloper.com/2017/04/using-sqlite-database-golang-application/
func faq(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	lis := filepath.Join("templates", "inc_list.go.html")
	id := filepath.Join("templates", "faq.go.html")
	tpli, err := template.ParseFiles(tpi, id, lis)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		// Log the detailed error
		//      log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tpli.ExecuteTemplate(w, "faq1", nil)
}
func about(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	lis := filepath.Join("templates", "inc_list.go.html")
	id := filepath.Join("templates", "about.go.html")
	tpli, err := template.ParseFiles(tpi, id, lis)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		// Log the detailed error
		//      log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tpli.ExecuteTemplate(w, "about", nil)
}
func contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	lis := filepath.Join("templates", "inc_list.go.html")
	id := filepath.Join("templates", "contact.go.html")
	tpli, err := template.ParseFiles(tpi, id, lis)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		// Log the detailed error
		//      log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tpli.ExecuteTemplate(w, "contact", nil)
}
func tc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	lis := filepath.Join("templates", "inc_list.go.html")
	id := filepath.Join("templates", "tc.go.html")
	tpli, err := template.ParseFiles(tpi, id, lis)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		// Log the detailed error
		//      log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tpli.ExecuteTemplate(w, "tc", nil)
}
func pp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	lis := filepath.Join("templates", "inc_list.go.html")
	id := filepath.Join("templates", "pp.go.html")
	tpli, err := template.ParseFiles(tpi, id, lis)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		// Log the detailed error
		//      log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = tpli.ExecuteTemplate(w, "pp", nil)

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
	jschartid := []string{"chartjs", counnt, ".inc.go.html"}
	subsec := []string{"subsec", counnt, ".inc.go.html"}

	// Join with no separator to create a compact string.
	joined := strings.Join(vvvalue, "")
	jsjoined := strings.Join(jschartid, "")
	sssubsec := strings.Join(subsec, "")

	ss := filepath.Join("templates", "stats", sssubsec)
	mf := filepath.Join("templates", "stats", joined)
	js := filepath.Join("templates", "stats", jsjoined)
	lp := filepath.Join("templates", "layout1.go.html")
	fp := filepath.Join("templates", "stats", ps.ByName("name"))
	tp := filepath.Join("templates", "inc.go.html")
	//    log.Println(fp)
	tpl, err := template.ParseFiles(lp, fp, tp, mf, js, ss)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = tpl.ExecuteTemplate(w, "layout3", nil)
	if err != nil {
		log.Fatal(err)
	}
}
