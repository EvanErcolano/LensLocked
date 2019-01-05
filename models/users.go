package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource when a resource cannot be found in db
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided to a method like delete
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

const userPwPepper = "aaaafe93-7942-4e3d-a4fc-e295ba99d571"

// NewUserService takes care of setting up the db for the userService.
// If there is an error opening the db
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	// we can't defer the closing here because it closes when
	// the func goes out of scope. so it would close before we hand it off (bad)
	// instead we created  Close()
	return &UserService{
		db: db,
	}, nil
}

//UserService interacts with user objects
type UserService struct {
	db *gorm.DB
}

// DesctructiveReset drop a table and remigrates.
func (us *UserService) DesctructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

//AutoMigrate will attempt to automically migrate the user table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

//  The parenthesis before the function name is the Go way of defining the object
//  on which these functions will operate.  So this function is available on the
//  user service
// func (RECEIVER TYPE) funcName(params) (return types)

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
	return us.db.Create(user).Error
}

// Update the provided user with all of the data in the user object
func (us *UserService) Update(user *User) error {
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
}
