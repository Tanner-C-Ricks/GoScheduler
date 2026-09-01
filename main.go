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
	"BASE":        TEMPLATE_URL + "/base.html",
	"NAVIGATION":  TEMPLATE_URL + "/navigation.html",
	"HOME":        TEMPLATE_URL + "/home.html",
	"ACCOUNT":     TEMPLATE_URL + "/account.html",
	"ACCOUNTINFO": TEMPLATE_URL + "/login/account_info.html",
	"SCHEDULER":   TEMPLATE_URL + "/scheduler.html",
	"APPROVAL":    TEMPLATE_URL + "/approval.html",
	"SIGNUP":      TEMPLATE_URL + "/login/signup.html",
	"LOGIN":       TEMPLATE_URL + "/login/login.html",
	"LOGOUT":      TEMPLATE_URL + "/login/logout.html",
	"LOGGEDIN":    TEMPLATE_URL + "/login/logged_in.html",
	"LOGGEDOUT":   TEMPLATE_URL + "/login/logged_out.html",
	"ACCLIST":     TEMPLATE_URL + "/fragments/account_list.html",
}

func loadPage(files []string, w http.ResponseWriter, r *http.Request) {
	// fmt.Println(HTML_FILES["BASE"])
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		fmt.Println(fmt.Sprintf("There is no file found for path \"%v\"", files))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	return
}

func homePage(w http.ResponseWriter, r *http.Request) {
	loggedIn := checkLoggedIn(w, r)
	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES[loggedIn],
			HTML_FILES["HOME"],
		}

	loadPage(files, w, r)
	// w.Write([]byte("Home Page"))
}

func accountPage(w http.ResponseWriter, r *http.Request) {
	loggedIn := checkLoggedIn(w, r)

	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES[loggedIn],
			HTML_FILES["ACCOUNT"],
		}

	loadPage(files, w, r)
	// w.Write([]byte("Account Page"))
}

func schedulePage(w http.ResponseWriter, r *http.Request) {
	loggedIn := checkLoggedIn(w, r)

	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES[loggedIn],
			HTML_FILES["SCHEDULER"],
		}

	loadPage(files, w, r)
	// w.Write([]byte("Schedule Page"))
}

func approvalPage(w http.ResponseWriter, r *http.Request) {
	loggedIn := checkLoggedIn(w, r)

	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES[loggedIn],
			HTML_FILES["APPROVAL"],
		}

	loadPage(files, w, r)
	// w.Write([]byte("Approval Page"))
}

// func getAccountInfo(w http.ResponseWriter, r *http.Response) {

// 	acc := getAccount()

// 	files :=
// 		[]string{
// 			HTML_FILES["ACCOUNTINFO"],
// 		}

// 	ts, err := template.ParseFiles(files...)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	ts.Execute(w, acc)
// }

func getAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts := getAccounts()
	fmt.Println("Event handler called")

	files :=
		[]string{
			HTML_FILES["ACCLIST"],
		}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	ts.Execute(w, accounts)
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/{$}", homePage)
	mux.HandleFunc("/account/{$}", accountPage)
	mux.HandleFunc("GET /account/login/{$}", loginPage)
	mux.HandleFunc("POST /account/login/{$}", login)
	mux.HandleFunc("GET /account/signup/", signupPage)
	mux.HandleFunc("POST /account/signup/", signup)
	mux.HandleFunc("POST /account/logout/{$}", logout)
	mux.HandleFunc("POST /account/delete/", deleteAccountHandler)
	mux.HandleFunc("GET /account/accountInfo/", accountName)
	mux.HandleFunc("GET /accountlist/{$}", getAccountsHandler)
	mux.HandleFunc("/account/schedule/{$}", schedulePage)
	mux.HandleFunc("/account/schedule/approval/{$}", approvalPage)

	// err := http.ListenAndServe(":4000", mux)
	err := http.ListenAndServe(":4000", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("REQUEST:", r.Method, r.URL.Path)
			mux.ServeHTTP(w, r)
		},
	))
	fmt.Println(err)
}
