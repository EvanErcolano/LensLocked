package main

// package main is the first part of the executable run
// requires a main function

import (
	"net/http" // used for web server or making web requests

	"lenslocked.com/controllers"

	"github.com/gorilla/mux"
)

func main() {
	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.HandleFunc("/signup", usersController.New).Methods("GET")
	r.HandleFunc("/signup", usersController.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
