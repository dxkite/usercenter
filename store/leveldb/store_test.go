package leveldb

import (
	"dxkite.cn/usercenter/hash"
	"fmt"
	"testing"
)

func TestNewStore(t *testing.T) {
	s, err := NewDataStore("./tests/id")
	if err != nil {
		t.Error(err)
	}

	if err := s.Init(); err != nil {
		t.Error(err)
	}

	if uin, err := s.Create(hash.NewMd5Password("dxkite").String(), `{"name":"dxkite"}`); err != nil {
		t.Error(err)
	} else {
		fmt.Println("uin", uin)
	}
}

func TestDataStore_List(t *testing.T) {
	s, err := NewDataStore("./tests/id")
	if err != nil {
		t.Error(err)
	}
	if err := s.Init(0); err != nil {
		t.Error(err)
	}
	fmt.Println(s.Size())
	if l, err := s.List(1, 10); err != nil {
		t.Error(err)
	} else {
		fmt.Println(l)
	}
}
