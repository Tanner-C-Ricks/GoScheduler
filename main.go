package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var STATIC_URL = "./static"

var TEMPLATE_URL = STATIC_URL + "/HTML"

var DATABASE_URL = "./db/"

var HTML_FILES = map[string]string{
	"BASE":        TEMPLATE_URL + "/base.html",
	"NAVIGATION":  TEMPLATE_URL + "/navigation.html",
	"HOME":        TEMPLATE_URL + "/home.html",
	"ACCOUNT":     TEMPLATE_URL + "/account.html",
	"ACCOUNTINFO": TEMPLATE_URL + "/login/account_info.html",
	"SCHEDULE":    TEMPLATE_URL + "/scheduler.html",
	"APPROVAL":    TEMPLATE_URL + "/approval.html",
	"SIGNUP":      TEMPLATE_URL + "/login/signup.html",
	"SIGNUPADMIN": TEMPLATE_URL + "/login/signup_admin.html",
	"LOGIN":       TEMPLATE_URL + "/login/login.html",
	"LOGOUT":      TEMPLATE_URL + "/login/logout.html",
	"LOGGEDIN":    TEMPLATE_URL + "/login/logged_in.html",
	"LOGGEDOUT":   TEMPLATE_URL + "/login/logged_out.html",
	"ACCLIST":     TEMPLATE_URL + "/fragments/account_list.html",

	"TIMEBLOCK":       TEMPLATE_URL + "/fragments/time_block.html",
	"SAVEDTIMEBLOCK":  TEMPLATE_URL + "/fragments/time_block_saved.html",
	"SCHEDULEREDITOR": TEMPLATE_URL + "/scheduleEditor.html",

	"APPROVALLIST":   TEMPLATE_URL + "/fragments/approval_list.html",
	"APPROVALSTATUS": TEMPLATE_URL + "/fragments/approval_status.html",
	"NOTE":           TEMPLATE_URL + "/fragments/note.html",
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
	if loggedIn != "LOGGEDIN" {
		// w.Header().Set("HX-Redirect", "/account/login/")
		// w.WriteHeader(http.StatusOK)
		loginPage(w, r)
		fmt.Println("The user is not logged in.")

		return
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println("There was an error with the form", err)
		accountPage(w, r)
		return
	}
	id, err := strconv.Atoi(r.FormValue("ScheduleID"))

	if err != nil {
		fmt.Println("There is an error with the value sent", err)
		accountPage(w, r)
		return
	}

	sched := getScheduleFromSchedID(id)

	if !checkAdmin(w, r) && accountInfo(w, r).ID != sched.AccountID {
		// w.Header().Set("HX-Redirect", "/account/")
		// w.WriteHeader(http.StatusUnauthorized)
		accountPage(w, r)
		fmt.Println("The user is authorized to see schedule")
		return
	}

	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES[loggedIn],
			HTML_FILES["SCHEDULE"],
		}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		fmt.Println(fmt.Sprintf("There is no file found for path \"%v\"", files))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", sched)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	return
	// w.Write([]byte("Schedule Page"))
}

func scheduleEditorPage(w http.ResponseWriter, r *http.Request) {
	loggedIn := checkLoggedIn(w, r)
	if loggedIn != "LOGGEDIN" {
		// w.Header().Set("HX-Redirect", "/account/login/")
		// w.WriteHeader(http.StatusOK)
		loginPage(w, r)
		fmt.Println("The user is not logged in.")

		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println("There was an error with the form", err)
		accountPage(w, r)
		return
	}
	id, err := strconv.Atoi(r.FormValue("ScheduleID"))

	if err != nil {
		fmt.Println("There is an error with the value sent", err)
		accountPage(w, r)
		return
	}

	sched := getScheduleFromSchedID(id)
	fmt.Println(id)

	if !checkAdmin(w, r) && accountInfo(w, r).ID != sched.AccountID {
		// w.Header().Set("HX-Redirect", "/account/")
		// w.WriteHeader(http.StatusUnauthorized)
		accountPage(w, r)
		fmt.Println("The user is authorized to see schedule")
		return
	}

	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES[loggedIn],
			HTML_FILES["SCHEDULEREDITOR"],
		}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		fmt.Println(fmt.Sprintf("There is no file found for path \"%v\"", files))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", sched)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	return
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
	accounts := []Account{}
	if !checkAdmin(w, r) && checkLoggedIn(w, r) == "LOGGEDIN" {
		accounts = getAccounts()
	}

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

func getApprovalsHandler(w http.ResponseWriter, r *http.Request) {
	approvals := getApprovalsList(w, r)
	fmt.Println("List of Approvals Petitioned")

	files :=
		[]string{
			HTML_FILES["APPROVALLIST"],
		}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	ts.Execute(w, approvals)
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
	mux.HandleFunc("GET /account/signup/admin", signupAdminPage)
	mux.HandleFunc("POST /account/signup/", signup)
	mux.HandleFunc("POST /account/signup/admin", signupAdmin)

	mux.HandleFunc("POST /account/logout/{$}", logout)
	mux.HandleFunc("POST /account/delete/", deleteAccountHandler)
	mux.HandleFunc("GET /account/accountInfo/", accountInfoHandler)
	mux.HandleFunc("GET /accountlist/{$}", getAccountsHandler)
	mux.HandleFunc("POST /account/schedule/edit/{$}", schedulePage)
	mux.HandleFunc("POST /account/schedule/{$}", scheduleEditorPage)

	mux.HandleFunc("GET /account/schedule/newTimeBlock/{$}", newTimeBlock)
	mux.HandleFunc("POST /account/schedule/updateBlock/{$}", updateTimeBlock)
	mux.HandleFunc("GET /account/schedule/loadTimeBlocks/{$}", loadTimeBlocksHandler)
	mux.HandleFunc("GET /account/schedule/loadTimeBlocksEdit/{$}", loadTimeBlocksEditHandler)
	mux.HandleFunc("GET /account/schedule/send/{$}", sendScheduleForReview)
	mux.HandleFunc("GET /account/schedule/approval/unsubmitted/", unsubmittedSchedule)
	mux.HandleFunc("GET /account/schedule/scheduleID/", getScheduleID)

	mux.HandleFunc("GET /account/schedule/editor/{$}", scheduleEditorPage)
	mux.HandleFunc("GET /account/schedule/status/", getApprovalStatus)
	mux.HandleFunc("GET /account/schedule/note/", getNote)

	mux.HandleFunc("GET /admin/schedules/approvals/{$}", approvalPage)
	mux.HandleFunc("GET /admin/approvals/", getApprovalsHandler)
	mux.HandleFunc("POST /admin/schedules/approvals/approve/", approveSchedule)
	mux.HandleFunc("POST /admin/schedules/approvals/revision/", revisionSchedule)
	mux.HandleFunc("POST /admin/schedules/approvals/note/{$}", updateScheduleNote)

	// err := http.ListenAndServe(":4000", mux)
	err := http.ListenAndServe(":4000", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("REQUEST:", r.Method, r.URL.Path)
			mux.ServeHTTP(w, r)
		},
	))
	fmt.Println(err)
}
