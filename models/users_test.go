package models

import (
	"fmt"
	"testing"
	"time"
)

// TODO: write tests for byID, by Email, Update, and delete

func testingUserService() (*UserService, error) {
	const (
		host     = "localhost"
		port     = 5432
		password = ""
		user     = "fenderjazzplayer"
		dbname   = "lenslocked_test"
	)

	connString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)

	us, err := NewUserService(connString)
	if err != nil {
		return nil, err
	}
	us.db.LogMode(false)
	// clear the user table between test
	us.DesctructiveReset()
	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		Name:  "Michael Scott",
		Email: "mscott@dundermifflin.com",
	}

	err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Errorf("Expected ID created than user, received %d", user.ID)
	}

	if time.Since(user.CreatedAt) > time.Duration(10*time.Second) {
		t.Errorf("Expected createdAt to be recent. Received %s", user.CreatedAt)
	}

	if time.Since(user.UpdatedAt) > time.Duration(10*time.Second) {
		t.Errorf("Expected createdAt to be recent. Received %s", user.UpdatedAt)
	}

}
