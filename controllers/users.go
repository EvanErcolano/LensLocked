package controllers

import (
	"net/http"

	"lenslocked.com/views"
)

// NewUsers is used to create a new users controller.NewUsers
// This funtion will panic if the templates are not parsed correctly
// and shoudl be used only during initial setup
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

type Users struct {
	NewView *views.View
}

// this method is on the *User type
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}
