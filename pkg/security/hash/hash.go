package hash

import (
	"github.com/marprin/postman-lib/pkg/security"
	"golang.org/x/crypto/bcrypt"
)

type (
	service struct {
	}
)

func NewHash() security.HashContract {
	return &service{}
}

func (s *service) CheckHashValidity(hash string, bareString string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(bareString))
}

func (s *service) GenerateHash(bareString string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(bareString), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	stringHash := string(hash)
	return &stringHash, nil
}
