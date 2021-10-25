package hash

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

func TestNewMd5PasswordHash(t *testing.T) {
	hash := "786c0de5aa291ea41fbb84053b84a540"
	salt := "8dFGZt7E"
	password := "dxkite"
	h := NewMd5PasswordHash(password, salt)
	if hex.EncodeToString(h.Hash) != hash {
		t.Error("password hash error", h.Hash, hash)
	}

	fmt.Println(h.Hash, hex.EncodeToString(h.Hash))
	m := string(h.Marshal())
	fmt.Println(m)

	hh := &PasswordHash{}
	if err := hh.Unmarshal([]byte(m)); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(hh, h) {
		t.Error("not equal", hh, h)
	}
	if !hh.Verify(password) {
		t.Error("verify password failed")
	}
}

func TestNewHexSaltPassword(t *testing.T) {
	hash := "786c0de5aa291ea41fbb84053b84a540"
	salt := "8dFGZt7E"
	password := "dxkite"
	h, err := NewHexSaltPassword(hash, salt)
	if err != nil {
		t.Error(err)
	}
	if !h.Verify(password) {
		t.Error("password verify error", hex.EncodeToString(h.Hash))
	}

	m := string(h.Marshal())
	fmt.Println(m)
	hh := &PasswordHash{}
	if err := hh.Unmarshal([]byte(m)); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(hh, h) {
		t.Error("not equal", hh, h)
	}
	if !hh.Verify(password) {
		t.Error("verify password failed")
	}
}

func TestNewMd5Password(t *testing.T) {
	password := "dxkite"
	h := NewMd5Password(password)
	if !h.Verify(password) {
		t.Error("verify password error")
	}
}
