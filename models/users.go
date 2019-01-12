package models

import (
	"errors"
	"regexp"
	"strings"

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

	//ErrEmailRequired is returned when a user supplies a blank email address
	ErrEmailRequired = errors.New("models: email address is required")

	//ErrEmailInvalid is returned when a user supplies an invalid email address
	ErrEmailInvalid = errors.New("models: email address is not valid")

	// ErrEmailTaken is returned when the requested email address already
	// belongs to a different user during an update or create
	ErrEmailTaken = errors.New("models: email address is already taken")
)

const userPwPepper = "aaaafe93-7942-4e3d-a4fc-e295ba99d571"
const hmacSecretKey = "secret-hmac-key"

// User represents the user model stored in our db
// This is used for user accounts storing both email and passwords
// so users can login and gian access to their content
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"` // - <- tells gorm to ignore this
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null; unique_index"`
}

// UserDB is used to interact with the user database
//
// For pretty much all single user queries
// 1 - user, nil 		(found user)
// 2 - nil, ErrNotFound (couldnt find user)
// 3 - nil, OtherError 	(db issues - send  a 500)
//
// for single user queries, any error but ErrNotFund should
// result in a 500 error
type UserDB interface {
	// Methods for querying single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close db connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of methods used to manipulate and
// work with the user model
type UserService interface {
	// Authenticate will verify the provided email address and
	// password are correct, if true, the user corresponding to
	// that email will be returned. Otherwise, you will receive:
	// ErrNotFound, ErrInvalidPassword, or another error
	Authenticate(email, password string) (*User, error)
	UserDB
}

// NewUserService takes care of setting up the db for the userService.
// If there is an error opening the db
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
	}, nil
}

var _ UserService = &userService{}

//userService interacts with user objects
type userService struct {
	UserDB
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
func (us *userService) Authenticate(email, password string) (*User, error) {
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

// Ensure uservalidator implements this interface
var _ UserDB = &userValidator{}

type userValidatorFunc func(*User) error // is this an inteface?

func runUserValidationFuncs(user *User, fns ...userValidatorFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

// userValidator is a lyaer that validates things before htey go
// go to the db layer via nroalization and validation
// its methods will match whatever exists in userdb
type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

// byEmail will normalize the email address before calling
// ByEmail on the UserDb field.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValidationFuncs(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// ByRemember will hash the remember token and then call
// ByRemember on the subsequent UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	err := runUserValidationFuncs(&user, uv.hmacRemember)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create creates a user in the db via GORM
func (uv *userValidator) Create(user *User) error {
	// run all validation funcs and return any errors
	err := runUserValidationFuncs(user,
		uv.bcryptPassword,
		uv.setRememberIfUnset,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// Update will hash a remember token if it is provided.
func (uv *userValidator) Update(user *User) error {
	// run all validation funcs and return any errors
	err := runUserValidationFuncs(user,
		uv.bcryptPassword,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete validates the user id that is passed into Delete
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValidationFuncs(&user, uv.ensureIDGreaterThan(0))
	if err != nil {
		return nil
	}
	return uv.UserDB.Delete(id)
}

// bcryptPassword will hash a user's password with a predefined
// pepper (userPwPepper) and bcrypt if the password is not the empty
// string
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) ensureIDGreaterThan(n uint) userValidatorFunc {
	return userValidatorFunc(func(user *User) error {
		if user.ID <= n {
			return ErrInvalidID
		}
		return nil
	})
}

// normalizeEmail removes extra spaces and lowercases the provided email
func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.TrimSpace(user.Email)
	user.Email = strings.ToLower(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.UserDB.ByEmail(user.Email)
	if err == ErrNotFound {
		// Email address is not taken
		return nil
	}
	if err != nil {
		return err
	}

	// We found a user w/ this email address...
	// If the found user has the same ID as this user it is
	// an update and this is the same user.
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

// Ignored but allows us to check if userGorm ever stops
// satisfying the userdb interface
var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db: db,
	}, nil
}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail will look up the user with the email provided.
// 1 - user, nil 		(found user)
// 2 - nil, ErrNotFound (couldnt find user)
// 3 - nil, OtherError 	(db issues - send  a 500)
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with the given remember token
// this method expects the remember token to already by hashed
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	db := ug.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// Create creates a user in the db via GORM
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update the provided user with all of the data in the user object
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete the provided user with all the data in the user object
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close closes the gorm db connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset drops the user table and rebuilds it
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the
// users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
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
