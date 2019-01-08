package models

import (
	"errors"

	"lenslocked.com/hash"
	"lenslocked.com/rand"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource when a resource cannot be found in db
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided to a method like delete
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password is provided
	ErrInvalidPassword = errors.New("models: Password provided was invalid")
)

const userPwPepper = "aaaafe93-7942-4e3d-a4fc-e295ba99d571"
const hmacSecretKey = "secret-hmac-key"

// NewUserService takes care of setting up the db for the userService.
// If there is an error opening the db
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	// we can't defer the closing here because it closes when
	// the func goes out of scope. so it would close before we hand it off (bad)
	// instead we created  Close()
	return &UserService{
		db:   db,
		hmac: hmac,
	}, nil
}

//UserService interacts with user objects
type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the
// users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// ByID will look up the user with the id provided.
// 1 - user, nil 		(found user)
// 2 - nil, ErrNotFound (couldnt find user)
// 3 - nil, OtherError 	(db issues - send  a 500)
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail will look up the user with the email provided.
// 1 - user, nil 		(found user)
// 2 - nil, ErrNotFound (couldnt find user)
// 3 - nil, OtherError 	(db issues - send  a 500)
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with the given remember token
// this method will handle hashing the token for us
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	tokenHash := us.hmac.Hash(token)
	db := us.db.Where("remember_hash = ?", tokenHash)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// Authenticate can be used to authenticate a user with the
// provided email address and password.
// If the email address provided is invalid, this will return
//   nil, ErrNotFound
// If the password provided is invalid, this will return
//   nil, ErrInvalidPassword
// If the email and password are both valid, this will return
//   user, nil
// Otherwise if another error is encountered this will return
//   nil, error
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

// first will query using the provided gorm.Db and will get the
// first item returned and place it into dst. If nothing is returned
// then it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Create creates a user in the db via GORM
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)
	return us.db.Create(user).Error
}

// Update the provided user with all of the data in the user object
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
	return us.db.Save(user).Error
}

// Delete the provided user with all the data in the user object
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{
		Model: gorm.Model{ID: id},
	}
	return us.db.Delete(&user).Error
}

// Close closes the gorm db connection
func (us *UserService) Close() error {
	return us.db.Close()
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"` // - <- tells gorm to ignore this
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null; unique_index"`
}
