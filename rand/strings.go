package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

//Bytes will help us generate n random bytes and return error if there is one
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	// fill in b with random bytes
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String will generate a byte slce of size nBytes and then
// return a string that is the base64 url encoded version of
// that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", nil
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken is a helper function desined to generate
// remember tokens of a predetermined byte size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
