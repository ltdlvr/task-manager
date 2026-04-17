package tool

import (
	"math/rand"
	"testing"
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
}
