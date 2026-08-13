package auth

import (

	"github.com/alexedwards/argon2id"
)

func hashPassword(password string) (string, error) {
	hashedPasswords, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err 
	}

	return hashedPasswords, nil 
}


func validatePasswordHash(password string, hashedPassword string) (bool, error) {
	valid, _, err := argon2id.CheckHash(password, hashedPassword)
	if err != nil {
		return false, err 
	}

	return valid, nil 
}