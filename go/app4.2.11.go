package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

var tpl *template.Template
var db *sql.DB
var err error

// var dbUsers = map[string]user{}      // user ID, user
//var dbSessions = map[string]string{} // session ID, user ID

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.go.html"))

	db, err = sql.Open("sqlite3", "/home/vagrant/The9s/the9s")
	fmt.Println("You connected to your database.")
}

type add_user struct {
	Email        string
	First_name   string
	Last_name    string
	Password     string
	Account_type string
}
type add_url struct {
	U_id      string
	Short_url string
	Long_url  string
}
type url struct {
	Url_id    string
	U_id      string
	Long_url  string
	Short_url string
}
type add_event struct {
	E_type   string
	Url_id   string
	Account  string
	Group_id string
	E_info   string
	E_status string
}
type event struct {
	E_id     string
	E_type   string
	Url_id   string
	Account  string
	Group_id string
	DATETIME string
	E_info   string
	E_status string
}

func main() {
	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("static/"))
	router.GET("/", index)
	router.GET("/signup", signup)
	router.POST("/signup", signup)

	router.POST("/url/insert", urlinsert)
	router.GET("/url", urlshow)
	router.GET("/url/delete", urldelete)
	router.GET("/url/edit", urledit)
	router.POST("/url/update", urlupdate)
	router.GET("/url/show", urlandevents)

	router.GET("/event", eventshow)
	router.POST("/event/insert", eventinsert)

	router.GET("/faq", faq)
	router.GET("/about", about)
	router.GET("/tandc", tc)
	router.GET("/pp", pp)
	router.GET("/contact", contact)

	router.GET("/stats/:name", serveTemplate)
	log.Println("Listening on :80 ...")
	http.ListenAndServe(":80", router)
}
func signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// get cookie if not set one
	if r.Method == http.MethodPost {

		// get form values
		au := add_user{}
		au.Email = r.FormValue("email")
		au.First_name = r.FormValue("first_name")
		au.Last_name = r.FormValue("last_name")
		au.Password = r.FormValue("password")
		au.Account_type = r.FormValue("account_type")
		_ = db
		_, err = db.Exec("INSERT INTO user(email, first_name, last_name, password, account_type) VALUES(?,?,?,?,?)", au.Email, au.First_name, au.Last_name, au.Password, au.Account_type)
		if err != nil {
			fmt.Printf("parsing failed: %s", err)
			return
		}
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
	}
	///user session and submission

	tpi := filepath.Join("templates", "inc.go.html")
	id := filepath.Join("templates", "signup.go.html")
	tpli, err := template.ParseFiles(tpi, id)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		// Log the detailed error
		//      log.Println(err.Error())
		fmt.Printf("parsing failed: %s", err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = tpli.ExecuteTemplate(w, "signup", nil)

}
func urldelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	url_id := r.FormValue("url")
	if url_id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	_ = db
	_, err = db.Exec("delete from url where url_id = ?", url_id)
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}

	http.Redirect(w, r, "/url", http.StatusSeeOther)
}
func urlupdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ul := url{}
	ul.Url_id = r.FormValue("url_id")
	ul.Short_url = r.FormValue("short_url")
	ul.Long_url = r.FormValue("long_url")

	_ = db
	_, err = db.Exec("update url set short_url = ?, url = ? where url_id  = ?", ul.Short_url, ul.Long_url, ul.Url_id)
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}
	http.Redirect(w, r, "/url", http.StatusSeeOther)
}
func urledit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	id := filepath.Join("templates", "urledit.go.html")
	tpli, err := template.ParseFiles(tpi, id)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	url_id := r.FormValue("url")
	if url_id == "" {
		http.Error(w, http.StatusText(900), 900)
		return
	}
	rows := db.QueryRow("select url_id, u_id, url, short_url from url where url_id = ?", url_id)
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}

	uurl := url{}
	err = rows.Scan(&uurl.Url_id, &uurl.U_id, &uurl.Long_url, &uurl.Short_url) // order matters
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}
	tpli.ExecuteTemplate(w, "form", uurl)
}
func urlinsert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ul := add_url{}
	ul.U_id = r.FormValue("u_id")
	ul.Short_url = r.FormValue("short_url")
	ul.Long_url = r.FormValue("long_url")

	_ = db
	_, err = db.Exec("INSERT INTO url(url, short_url, u_id) VALUES(?,?,?)", ul.Long_url, ul.Short_url, ul.U_id)
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}
	http.Redirect(w, r, "/url", http.StatusSeeOther)
}
func urlshow(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	id := filepath.Join("templates", "url.go.html")
	tpli, err := template.ParseFiles(tpi, id)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("select url_id, u_id, url, short_url from url order by url_id DESC")
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer rows.Close()

	urls := make([]url, 0)
	for rows.Next() {
		uurl := url{}
		err := rows.Scan(&uurl.Url_id, &uurl.U_id, &uurl.Long_url, &uurl.Short_url) // order matters
		if err != nil {
			log.Fatal(err)
			return
		}
		urls = append(urls, uurl)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return
	}
	tpli.ExecuteTemplate(w, "form", urls)
}
func urlandevents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	id := filepath.Join("templates", "urlshowall.go.html")
	tpli, err := template.ParseFiles(tpi, id)
	//log.Println(tpli) asdihasnpi
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	url_id := r.FormValue("url")
	if url_id == "" {
		http.Error(w, http.StatusText(900), 900)
		return
	}
	rows, err := db.Query("select * from event where url_id = ? order by group_id", url_id)
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}
	ues := make([]event, 0)
	for rows.Next() {
		uess := event{}
		err = rows.Scan(&uess.E_id, &uess.E_type, &uess.Url_id, &uess.Account, &uess.Group_id, &uess.DATETIME, &uess.E_info, &uess.E_status)
		if err != nil {
			fmt.Printf("parsing failed: %s", err)
			return
		}
		ues = append(ues, uess)
	}

	//https://dlintw.github.io/gobyexample/public/text-template.html
	//https://stackoverflow.com/questions/16985469/switch-or-if-elseif-else-inside-golang-html-templates
	//http://2016.8-p.info/post/06-18-go-html-template/	tpli.ExecuteTemplate(w, "form", ues)
	tpli.ExecuteTemplate(w, "form", ues)
}
func eventshow(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	id := filepath.Join("templates", "event.go.html")
	tpli, err := template.ParseFiles(tpi, id)

	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}

	url_id := r.FormValue("url")
	if url_id == "" {
		http.Error(w, http.StatusText(900), 900)
		return
	}
	row := db.QueryRow("select url_id, u_id, short_url from url where url_id = ?", url_id)
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}

	uurl := url{}
	err = row.Scan(&uurl.Url_id, &uurl.U_id, &uurl.Short_url)
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}

	tpli.ExecuteTemplate(w, "form", uurl)
}
func eventinsert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	now := time.Now()
	secs := now.Unix()
	fmt.Println(secs)

	adde := add_event{}
	adde.E_type = r.FormValue("e_type")
	adde.Url_id = r.FormValue("url_id")
	adde.Account = r.FormValue("account")
	adde.Group_id = r.FormValue("group_id")

	adde.E_info = r.FormValue("e_info")
	adde.E_status = r.FormValue("e_status")

	_ = db
	_, err = db.Exec("insert into event (e_type, url_id, account, group_id, e_status, e_info) values (?,?,?,?,?,?)",
		adde.E_type, adde.Url_id, adde.Account, adde.Group_id, adde.E_info, adde.E_status)
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}

	target := "/url/show?url=" + adde.Url_id
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func faq(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpi := filepath.Join("templates", "inc.go.html")
	id := filepath.Join("templates", "faq.go.html")
	tpli, err := template.ParseFiles(tpi, id)
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
	id := filepath.Join("templates", "about.go.html")
	tpli, err := template.ParseFiles(tpi, id)
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

	id := filepath.Join("templates", "contact.go.html")
	tpli, err := template.ParseFiles(tpi, id)
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
	id := filepath.Join("templates", "tc.go.html")
	tpli, err := template.ParseFiles(tpi, id)
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
	id := filepath.Join("templates", "pp.go.html")
	tpli, err := template.ParseFiles(tpi, id)
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
	id := filepath.Join("templates", "index.go.html")
	odd := filepath.Join("templates", "inc_odometer.go.html")
	tpli, err := template.ParseFiles(tpi, id, odd)
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

	url_id := r.FormValue("url")
	if url_id == "" {
		http.Error(w, http.StatusText(900), 900)
		return
	}
	rows, err := db.Query("select * from event where url_id = ? order by group_id DESC, datetime ", url_id)
	if err != nil {
		fmt.Printf("parsing failed: %s", err)
		return
	}
	ues := make([]event, 0)
	for rows.Next() {
		uess := event{}
		err = rows.Scan(&uess.E_id, &uess.E_type, &uess.Url_id, &uess.Account, &uess.Group_id, &uess.DATETIME, &uess.E_info, &uess.E_status)
		if err != nil {
			fmt.Printf("parsing failed: %s", err)
			return
		}
		ues = append(ues, uess)
	}

	err = tpl.ExecuteTemplate(w, "layout3", ues)
	if err != nil {
		log.Fatal(err)
	}
}
