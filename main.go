package main

// package main is the first part of the executable run
// requires a main function

import (
	"fmt"
	"net/http" // used for web server or making web requests

	"lenslocked.com/controllers"
	"lenslocked.com/models"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	password = ""
	user     = "fenderjazzplayer"
	dbname   = "lenslocked_dev"
)

func main() {
	connString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	us, err := models.NewUserService(connString)
	must(err)
	defer us.Close()
	us.AutoMigrate()
	// us.DestructiveReset()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.HandleFunc("/signup", usersController.New).Methods("GET")
	r.HandleFunc("/signup", usersController.Create).Methods("POST")
	r.Handle("/login", usersController.LoginView).Methods("GET")
	r.HandleFunc("/login", usersController.Login).Methods("POST")
	fmt.Println("Starting the server on :3000.....")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
