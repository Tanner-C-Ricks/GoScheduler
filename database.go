package main

import (
	// "math/rand"
	"net/http"
	"strconv"
)

var sessions = map[string]int{}

type Database struct {
	Name          string
	ColumnNames   []string
	FileURL       string
	CreationQuery string
}

func addCookie(user Account, w http.ResponseWriter, r *http.Request) {
	session_id := "super_secret" + strconv.Itoa(user.ID)
	sessions[session_id] = user.ID

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session_id,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func removeCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
		return
	}

	// acc := getAccountID(sessions[cookie.Value])

	delete(sessions, cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
	return
}
