package main

// package main is the first part of the executable run
// requires a main function

import (
	"fmt"      // useful for printing data to diff places
	"net/http" // used for web server or making web requests
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/contact" {
		fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>.")
	} else if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1> Welcome to my awesome website </h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 not found")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
