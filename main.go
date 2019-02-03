package main

// package main is the first part of the executable run
// requires a main function

import (
	"fmt"
	"net/http" // used for web server or making web requests

	"lenslocked.com/controllers"
	"lenslocked.com/middleware"
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
	services, err := models.NewServices(connString)
	must(err)
	defer services.Close()

	services.AutoMigrate()
	// us.DestructiveReset()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(services.User)
	galleriesController := controllers.NewGalleries(services.Gallery)
	requireUserMw := middleware.RequireUser{services.User}

	r := mux.NewRouter()
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.HandleFunc("/signup", usersController.New).Methods("GET")
	r.HandleFunc("/signup", usersController.Create).Methods("POST")
	r.Handle("/login", usersController.LoginView).Methods("GET")
	r.HandleFunc("/login", usersController.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersController.CookieTest).Methods("GET")

	// Gallery routes
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesController.New)).Methods("GET")
	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesController.Create)).Methods("POST")
	fmt.Println("Starting the server on :3000.....")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
