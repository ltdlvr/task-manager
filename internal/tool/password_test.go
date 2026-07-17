package tool

import (
	"math/rand"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestVerifyPassword(t *testing.T) {
	t.Parallel()
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		t.Fatalf("generate password: %v", err)
	}
	pswdTool := NewPassword()
	password := string(buf)
	hashed, err := pswdTool.Hash(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if !pswdTool.Verify(password, hashed) {
		t.Error("VerifyPassword() = false; want true")
	}
	cost, err := bcrypt.Cost([]byte(hashed))
	if err != nil {
		t.Fatalf("get password hash cost: %v", err)
	}
	if cost != bcryptCost {
		t.Errorf("password hash cost = %d; want %d", cost, bcryptCost)
	}
}
