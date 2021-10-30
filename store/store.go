package store

import (
	"encoding/json"
)

type UserInfo struct {
	Id           uint64 `json:"id"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	Data         string `json:"data"`
}

type UserData struct {
	Name string `json:"name"`
}

func (i *UserInfo) String() string {
	d, _ := json.Marshal(i)
	return string(d)
}

type UserStore interface {
	Init(baseId uint64) error
	Create(name, password string) (id uint64, err error)
	GetId(name string) (id uint64, err error)
	Get(id uint64) (*UserInfo, error)
	GetSize() (size uint64, err error)
	List(id uint64, limit int) (list []*UserInfo, err error)
	SetPassword(name, password string) error
}
