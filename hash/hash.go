package hash

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strconv"
	"strings"
)

type Algorithm int

const Md5Repeat = 1024

const (
	Hash1024Md5 Algorithm = 0
)

type PasswordHash struct {
	Type Algorithm
	Hash []byte
	Salt []byte
}

func NewMd5Password(password string) *PasswordHash {
	salt := make([]byte, 8)
	_, _ = io.ReadFull(rand.Reader, salt)
	hashed := calcMd5(password, string(salt), Md5Repeat)
	return &PasswordHash{
		Type: Hash1024Md5,
		Hash: hashed,
		Salt: salt,
	}
}

func NewMd5PasswordHash(password, salt string) *PasswordHash {
	hashed := calcMd5(password, salt, Md5Repeat)
	return &PasswordHash{
		Type: Hash1024Md5,
		Hash: hashed,
		Salt: []byte(salt),
	}
}

func NewHexSaltPassword(password, salt string) (*PasswordHash, error) {
	b, err := hex.DecodeString(password)
	if err != nil {
		return nil, err
	}
	return &PasswordHash{
		Type: Hash1024Md5,
		Hash: b,
		Salt: []byte(salt),
	}, nil
}

func calcMd5(password, salt string, repeat int) []byte {
	hash := md5.New()
	hash.Reset()
	hash.Write([]byte(salt))
	hash.Write([]byte(password))
	hashed := hash.Sum(nil)
	for i := 1; i < repeat; i++ {
		hash.Reset()
		hash.Write(hashed)
		hashed = hash.Sum(nil)
	}
	return hashed
}

func (p *PasswordHash) Marshal() []byte {
	h := []string{
		strconv.Itoa(int(p.Type)),
		base64.RawStdEncoding.EncodeToString(p.Hash),
		base64.RawStdEncoding.EncodeToString(p.Salt),
	}
	return []byte(strings.Join(h, "$"))
}

func (p *PasswordHash) String() string {
	return string(p.Marshal())
}

func (p *PasswordHash) Unmarshal(b []byte) (err error) {
	hash := strings.SplitN(string(b), "$", 3)
	if n, err := strconv.Atoi(hash[0]); err != nil {
		return err
	} else {
		p.Type = Algorithm(n)
	}
	if p.Hash, err = base64.RawStdEncoding.DecodeString(hash[1]); err != nil {
		return err
	}
	if p.Salt, err = base64.RawStdEncoding.DecodeString(hash[2]); err != nil {
		return err
	}
	return nil
}

func (p *PasswordHash) Verify(password string) bool {
	if p.Type == Hash1024Md5 {
		hashed := calcMd5(password, string(p.Salt), Md5Repeat)
		return bytes.Equal(hashed, p.Hash)
	}
	return false
}

func VerifyPassword(hash, password string) bool {
	h := &PasswordHash{}
	if err := h.Unmarshal([]byte(hash)); err != nil {
		return false
	}
	return h.Verify(password)
}
