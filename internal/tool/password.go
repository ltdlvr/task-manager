package tool

import "golang.org/x/crypto/bcrypt"

type Password struct{}

func NewPassword() *Password {
	return &Password{}
}

func (p *Password) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashed), err
}

func (p *Password) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
