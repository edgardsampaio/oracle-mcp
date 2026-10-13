package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestHashAndVerifyPassword(t *testing.T) {
	encoded, err := HashPassword("senha-longa-ok")
	if err != nil {
		t.Fatal(err)
	}
	if !stringsHasArgon(encoded) {
		t.Fatalf("hash = %s", encoded)
	}
	ok, err := VerifyPassword("senha-longa-ok", encoded)
	if err != nil || !ok {
		t.Fatalf("verify ok=%v err=%v", ok, err)
	}
	ok, err = VerifyPassword("outra-senha-longa", encoded)
	if err != nil || ok {
		t.Fatalf("verify mismatch ok=%v err=%v", ok, err)
	}
}

func TestPasswordMinLength(t *testing.T) {
	if err := validatePassword("12345"); err == nil {
		t.Fatal("5 caracteres deveria falhar")
	}
	if err := validatePassword("123456"); err != nil {
		t.Fatal(err)
	}
}

func TestHashTokenStable(t *testing.T) {
	sum := sha256.Sum256([]byte("abc"))
	want := base64.RawURLEncoding.EncodeToString(sum[:])
	if HashToken("abc") != want {
		t.Fatal("hash instável")
	}
}

func stringsHasArgon(s string) bool {
	return len(s) > 10 && s[:10] == "$argon2id$"
}
