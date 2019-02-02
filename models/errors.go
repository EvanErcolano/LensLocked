package models

import "strings"

const (
	// ErrNotFound is returned when a resource when a resource cannot be found in db
	ErrNotFound modelError = "models: resource not found"
	// ErrPasswordIncorrect is returned when an invalid password is provided
	ErrPasswordIncorrect modelError = "models: Password provided was invalid"
	//ErrEmailRequired is returned when a user supplies a blank email address
	ErrEmailRequired modelError = "models: email address is required"
	//ErrEmailInvalid is returned when a user supplies an invalid email address
	ErrEmailInvalid modelError = "models: email address is not valid"
	// ErrEmailTaken is returned when the requested email address already
	// belongs to a different user during an update or create
	ErrEmailTaken modelError = "models: email address is already taken"
	// ErrPasswordTooShort is returned when an update or create is attemped with
	// a password that is less than the required number of characters
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"
	// ErrPasswordRequired is returned when a create is attempted without a password
	ErrPasswordRequired modelError = "models: password is required"
	// ErrTitleRequired is returned when a create or get on a gallery is attempted without a title
	ErrTitleRequired modelError = "models: the title of the gallery is required"

	//ErrUserIDRequired is returned when a create or get is attempted without a UserID
	ErrUserIDRequired privateError = "models: the userID is required"
	// ErrIDInvalid is returned when an invalid ID is provided to a method like delete
	ErrIDInvalid privateError = "models: ID provided was invalid"
	// ErrRememberTooShort when a rememebr token is not at least 32 bytes
	ErrRememberTooShort privateError = "models: Remember token must be at lest 32 bytes"
	// ErrRememberRequired is returned when a create or update is attempted
	// without a user remember token hash
	ErrRememberRequired privateError = "models: remember token is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
