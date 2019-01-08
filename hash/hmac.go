package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// NewHMAC creates and returns a new HMAC object
func NewHMAC(secretKey string) HMAC {
	h := hmac.New(sha256.New, []byte(secretKey))
	return HMAC{
		HMAC: h,
	}
}

// HMAC is a wrapper around the crypto/hmac package making it easier to use
// in our code
type HMAC struct {
	HMAC hash.Hash
}

// Hash will hash the provided input string using HMAC with
// the secret ey provided when the HMAC object was created
func (h HMAC) Hash(input string) string {
	h.HMAC.Reset()
	h.HMAC.Write([]byte(input))
	b := h.HMAC.Sum(nil)
	// base64 makes sure it is a valid utf8 string which is url safe
	return base64.URLEncoding.EncodeToString(b)
}
