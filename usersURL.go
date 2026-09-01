package main

import (
	"fmt"
	"net/http"
)

func signupPage(w http.ResponseWriter, r *http.Request) {
	loggedIn := checkLoggedIn(w, r)

	if loggedIn == "LOGGEDIN" {
		w.Header().Set("HX-Redirect", "/account")
		w.WriteHeader(http.StatusOK)
		return
	}

	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES["LOGGEDOUT"],
			HTML_FILES["SIGNUP"],
		}

	loadPage(files, w, r)
}

func signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	acc := Account{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
	}

	success := createAccount(acc)
	if success != nil {
		w.Write([]byte("Failure"))
		fmt.Println(success)
	}
	// w.Write([]byte("Success"))
	// http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Println(acc)
	w.Header().Set("HX-Redirect", "/account/login/")
	w.WriteHeader(http.StatusOK)
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("No, here")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	acc := Account{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	acc = getAccount(acc)

	if acc.ID == 0 {
		w.Header().Set("HX-Redirect", "/account/signup")
		w.WriteHeader(http.StatusOK)
		return
	}

	addCookie(acc, w, r)

	// http.Redirect(w, r, "/account/", 308)
	// files :=
	// 	[]string{
	// 		HTML_FILES["LOGGEDIN"],
	// 	}

	// ts, err := template.ParseFiles(files...)
	// ts.Execute(w, acc)
	w.Header().Set("HX-Redirect", "/account/")
	w.WriteHeader(http.StatusOK)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Here")

	loggedIn := checkLoggedIn(w, r)

	if loggedIn == "LOGGEDIN" {
		w.Header().Set("HX-Redirect", "/account")
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Println("Logged Out")

	files :=
		[]string{
			HTML_FILES["BASE"],
			HTML_FILES["NAVIGATION"],
			HTML_FILES["LOGGEDOUT"],
			HTML_FILES["LOGIN"],
		}

	loadPage(files, w, r)
}

func logout(w http.ResponseWriter, r *http.Request) {
	removeCookie(w, r)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
	return
}

func deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := checkLoggedIn(w, r)
	if loggedIn != "LOGGEDIN" {
		return
	}
	acc := accountInfo(w, r)

	deleteAccount(acc)
	// acc := getAccountID(sessions[cookie.Value])
	removeCookie(w, r)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
	return
}

func checkLoggedIn(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// fmt.Println(err.Error())
		fmt.Println("Logged Out")
		return "LOGGEDOUT"
	}

	_, exists := sessions[cookie.Value]

	if exists {
		fmt.Println("Logged In")

		return "LOGGEDIN"
	}
	fmt.Println("Logged Out")

	return "LOGGEDOUT"
}
