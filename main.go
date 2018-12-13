package main

// package main is the first part of the executable run
// requires a main function

import (
	"fmt"      // useful for printing data to diff places
	"net/http" // used for web server or making web requests
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1> Welcome to my awesome website </h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
