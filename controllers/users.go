package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"

	"lenslocked.com/views"
)

// NewUsers is used to create a new users controller.NewUsers
// This funtion will panic if the templates are not parsed correctly
// and shoudl be used only during initial setup
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

// New is used to render the form where they can create a new user account
// GET / signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// SignupForm
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//Create is used to process the signup form to create a new account
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintln(w, user)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login is used to verify the provided email and password
// it will login the user in if they are correct
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address.")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid Password")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
	} // created our cookie but we must send it back via the response
	http.SetCookie(w, &cookie)
	fmt.Fprintln(w, user)

}
