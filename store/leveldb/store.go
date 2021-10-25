package leveldb

import (
	"bytes"
	"dxkite.cn/usercenter/store"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type DataStore struct {
	db *leveldb.DB
}

type Id uint64

const basicId = 0

var (
	ErrSize    = errors.New("error size")
	ErrVersion = errors.New("error version")
)

func init() {
	gob.Register(&BasicData{})
}

type BasicData struct {
	PasswordHash string
	Data         string
	Version      int
}

func (b *BasicData) String() string {
	d, _ := json.Marshal(b)
	return string(d)
}

func (u *Id) Marshal() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(*u))
	return b
}

func (u *Id) Unmarshal(b []byte) error {
	if len(b) != 8 {
		return errors.New("error size")
	}
	*u = Id(binary.BigEndian.Uint64(b))
	return nil
}

func (s *DataStore) Create(password string, data string) (id uint64, err error) {
	size, preId, bd, err := s.getBase()
	if err != nil {
		return 0, err
	}
	dt := &BasicData{
		PasswordHash: password,
		Data:         data,
	}
	if err = s.saveBase(size+1, preId+1, bd); err != nil {
		return 0, err
	}
	if err := s.put(preId+1, dt); err != nil {
		return 0, err
	}
	return preId + 1, nil
}

func (s *DataStore) getBase() (size, id uint64, bd *BasicData, err error) {
	bd, err = s.get(basicId)
	if err != nil {
		return 0, 0, nil, err
	}
	b := []byte(bd.Data)
	if len(b) != 16 {
		return 0, 0, nil, ErrSize
	}
	size = binary.BigEndian.Uint64(b[:8])
	id = binary.BigEndian.Uint64(b[8:])
	return
}

func (s *DataStore) saveBase(size, id uint64, bd *BasicData) (err error) {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[:8], size)
	binary.BigEndian.PutUint64(b[8:], id)
	bd.Data = string(b)
	return s.put(basicId, bd)
}

func (s *DataStore) get(id uint64) (d *BasicData, err error) {
	i := Id(id)
	b, err := s.db.Get(i.Marshal(), nil)
	if err != nil {
		return nil, err
	}
	d = &BasicData{}
	if err = Decode(b, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DataStore) put(id uint64, dt *BasicData) (err error) {
	var data []byte
	i := Id(id)
	if bd, err := s.get(id); err != nil && err != leveldb.ErrNotFound {
		return err
	} else {
		if bd != nil && bd.Version != dt.Version {
			return ErrVersion
		}
	}
	dt.Version++
	if data, err = Encode(dt); err != nil {
		return err
	}
	if err := s.db.Put(i.Marshal(), data, nil); err != nil {
		return err
	}
	return nil
}

func (s *DataStore) Password(id uint64, password string) error {
	dt, err := s.get(id)
	if err != nil {
		return err
	}
	dt.PasswordHash = password
	return s.put(id, dt)
}

func (s *DataStore) List(start uint64, limit int) (list []*store.UserInfo, err error) {
	key := Id(start)
	iter := s.db.NewIterator(nil, nil)
	for ok := iter.Seek(key.Marshal()); ok; ok = iter.Next() {
		id := Id(0)
		if err = id.Unmarshal(iter.Key()); err != nil {
			return nil, err
		}
		dt := &BasicData{}
		if err = Decode(iter.Value(), dt); err != nil {
			return nil, err
		}
		list = append(list, &store.UserInfo{
			Id:           uint64(id),
			Data:         dt.Data,
			PasswordHash: dt.PasswordHash,
		})
		limit--
		if limit == 0 {
			break
		}
	}
	iter.Release()
	err = iter.Error()
	return
}

func (s *DataStore) Put(id uint64, data string) error {
	dt, err := s.get(id)
	if err != nil {
		return err
	}
	dt.Data = data
	return s.put(id, dt)
}

func (s *DataStore) Get(id uint64) (data *store.UserInfo, err error) {
	dt, err := s.get(id)
	if err != nil {
		return nil, err
	}
	return &store.UserInfo{
		Id:           uint64(id),
		Data:         dt.Data,
		PasswordHash: dt.PasswordHash,
	}, nil
}

func (s *DataStore) Size() (uint64, error) {
	size, _, _, err := s.getBase()
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (s *DataStore) Init(baseId uint64) error {
	_, _, _, err := s.getBase()
	if err == leveldb.ErrNotFound {
		return s.saveBase(0, baseId, &BasicData{})
	}
	return nil
}

func NewDataStore(file string) (s *DataStore, err error) {
	s = &DataStore{}
	s.db, err = leveldb.OpenFile(file, nil)
	return s, err
}

func Encode(dt *BasicData) ([]byte, error) {
	b := &bytes.Buffer{}
	e := gob.NewEncoder(b)
	if err := e.Encode(dt); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func Decode(b []byte, dt *BasicData) (err error) {
	e := gob.NewDecoder(bytes.NewReader(b))
	if err := e.Decode(dt); err != nil {
		return err
	}
	return nil
}
