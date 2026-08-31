package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var STATIC_URL = "./static"

var TEMPLATE_URL = STATIC_URL + "/HTML"

var HTML_FILES = map[string]string{
	"BASE":       TEMPLATE_URL + "/base.html",
	"NAVIGATION": TEMPLATE_URL + "/navigation.html",
	"HOME":       TEMPLATE_URL + "/home.html",
	"ACCOUNT":    TEMPLATE_URL + "/accounts.html",
	"SCHEDULER":  TEMPLATE_URL + "/scheduler.html",
	"APPROVAL":   TEMPLATE_URL + "/approval.html",
	"LOGIN":      TEMPLATE_URL + "/login/login.html",
	"LOGOUT":     TEMPLATE_URL + "/login/logout.html",
}

func loadPage(files []string, w http.ResponseWriter, r *http.Request) {
	fmt.Println(HTML_FILES["BASE"])
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		fmt.Println(fmt.Sprintf("There is no file found for path \"%d\"", files))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES["HOME"],
		}

	loadPage(files, w, r)
	// w.Write([]byte("Home Page"))
}

func accountPage(w http.ResponseWriter, r *http.Request) {
	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES["ACCOUNT"],
		}

	loadPage(files, w, r)
	// w.Write([]byte("Account Page"))
}

func schedulePage(w http.ResponseWriter, r *http.Request) {
	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES["SCHEDULER"],
		}

	loadPage(files, w, r)
	// w.Write([]byte("Schedule Page"))
}

func approvalPage(w http.ResponseWriter, r *http.Request) {
	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES["APPROVAL"],
		}

	loadPage(files, w, r)
	// w.Write([]byte("Approval Page"))
}

func loginPage(w http.ResponseWriter, r *http.Request) {

	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES["LOGIN"],
		}

	loadPage(files, w, r)
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/{$}", homePage)
	mux.HandleFunc("/account/{$}", accountPage)
	mux.HandleFunc("/account/schedule/{$}", schedulePage)
	mux.HandleFunc("/account/schedule/approval/{$}", approvalPage)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
