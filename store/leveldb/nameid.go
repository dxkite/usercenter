package leveldb

import "github.com/syndtr/goleveldb/leveldb"

type NameStore struct {
	db *leveldb.DB
}

func NewNameStore(file string) (s *NameStore, err error) {
	s = &NameStore{}
	s.db, err = leveldb.OpenFile(file, nil)
	return s, err
}

func (s *NameStore) Put(id uint64, name string) error {
	i := Id(id)
	return s.db.Put([]byte(name), i.Marshal(), nil)
}

func (s *NameStore) Get(name string) (id uint64, err error) {
	i := Id(0)
	b, err := s.db.Get([]byte(name), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}
	if err := i.Unmarshal(b); err != nil {
		return 0, err
	}
	return uint64(i), nil
}
