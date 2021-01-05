package utils

import (
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// contextKey defines a type for context keys shared in the app
type contextKey string
var CurrrentUserCtxKey contextKey = "currentUser"

func ToJSON(data interface{}) ([]byte, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return raw, err
}

// GenerateGUID Generates a GUID
func GenerateGUID() (string, error) {
	ns, err := uuid.NewV4()
	return strings.ReplaceAll(uuid.NewV5(ns, ns.String()).String(), "-", ""), err
}

// EncryptPassword Encrypts user password
func EncryptPassword(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(bytes), err
}

// CompareEncryptdPasswd Compares hashed password with possible plaintext password
func CompareHashAndPassword(hashedPasswd string, passwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(passwd))
}
