package controller

import (
	"fmt"
	"net/http"
)

// ErrorES - some ElasticSearch error
func ErrorES(msg string, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	if len(msg) <= 0 {
		msg = "ElasticSearch fatal error: unknown?"
	} else {
		msg = "ElasticSearch fatal error: " + msg
	}
	fmt.Fprint(w, msg)
}

// Error404 handles 404 - Page Not Found
func Error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "Not Found 404")
}

// Error500 handles 500 - Internal Server Error
func Error500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "Internal Server Error 500")
}

// InvalidToken handles CSRF attacks
func InvalidToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, `Your token <strong>expired</strong>, click <a href="javascript:void(0)" onclick="window.history.back()">here</a> to try again.`)
}
