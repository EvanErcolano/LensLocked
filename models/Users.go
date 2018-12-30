package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource when a resource cannot be found in db
	ErrNotFound = errors.New("models: resource not found")
)

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
func (us *UserService) DesctructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
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
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Close closes the gorm db connection
func (us *UserService) Close() error {
	return us.db.Close()
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null: unique_index"`
}
