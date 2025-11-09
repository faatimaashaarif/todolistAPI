package utils

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"github.com/gorilla/mux"

)



func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(HashPassword), nil
}

func ComparePassword(plainPass, HashedPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(HashedPass), []byte(plainPass))
	if err != nil {
		return err
	}
	return nil
}
func GetURLParam(r *http.Request, key string) string {
	vars := mux.Vars(r)
	return vars[key]
}
