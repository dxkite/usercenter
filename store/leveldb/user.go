package leveldb

import (
	"dxkite.cn/usercenter/store"
	"encoding/json"
	"errors"
	"path"
	"strconv"
)

type UserStore struct {
	id   *DataStore
	name *NameStore
}

var (
	ErrUserExists    = errors.New("user exists")
	ErrUserNotExists = errors.New("user not exists")
)

func NewUserStore(p string) (s *UserStore, err error) {
	s = &UserStore{}
	if s.id, err = NewDataStore(path.Join(p, "data")); err != nil {
		return nil, err
	}
	if s.name, err = NewNameStore(path.Join(p, "name")); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *UserStore) Create(name, password string) (id uint64, err error) {
	eid, err := s.name.Get(name)
	if err != nil {
		return 0, err
	}
	if eid > 0 {
		return eid, ErrUserExists
	}
	id, err = s.id.Create(password, `{"name":`+strconv.QuoteToGraphic(name)+`}`)
	if err != nil {
		return
	}
	return id, s.name.Put(id, name)
}

func (s *UserStore) GetId(name string) (id uint64, err error) {
	return s.name.Get(name)
}

func (s *UserStore) Get(id uint64) (*store.UserInfo, error) {
	return s.id.Get(id)
}

func (s *UserStore) GetSize() (size uint64, err error) {
	return s.id.Size()
}

func (s *UserStore) List(start uint64, limit int) (list []*store.UserInfo, err error) {
	list, err = s.id.List(start, limit)
	if err != nil {
		return nil, err
	}
	item := &struct {
		Name string `json:"name"`
	}{}
	for i, v := range list {
		item.Name = ""
		_ = json.Unmarshal([]byte(v.Data), item)
		list[i].Name = item.Name
	}
	return list, nil
}

func (s *UserStore) Init(baseId uint64) error {
	return s.id.Init(baseId)
}

func (s *UserStore) SetBaseId(id uint64) error {
	return s.id.SetBaseId(id)
}

func (s *UserStore) SetPassword(name, password string) error {
	eid, err := s.name.Get(name)
	if err != nil {
		return err
	}
	if eid == 0 {
		return ErrUserNotExists
	}
	return s.id.Password(eid, password)
}
