package main

import (
	"dxkite.cn/usercenter/hash"
	"dxkite.cn/usercenter/store"
	"dxkite.cn/usercenter/store/leveldb"
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	OpAdd     = "add"
	OpAddHash = "add_hash"
	OpList    = "list"
)

func addUser(u store.UserStore, username, password string) error {
	if len(username) == 0 || len(password) == 0 {
		return errors.New("name or password empty")
	}
	id, err := u.Create(username, hash.NewMd5Password(password).String())
	if err == leveldb.ErrUserExists {
		fmt.Println("create user exists", username, "id", id)
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Println("create user", username, "id", id)
	return nil
}

func addUserHashSalt(u store.UserStore, username, md5Hash, salt string) error {
	if len(username) == 0 || len(md5Hash) == 0 || len(salt) == 0 {
		return errors.New("name or hash or salt empty")
	}
	h, err := hash.NewHexSaltPassword(md5Hash, salt)
	if err != nil {
		return errors.New(fmt.Sprint("parse hash error", err))
	}
	id, err := u.Create(username, h.String())
	if err == leveldb.ErrUserExists {
		fmt.Println("create user exists", username, "id", id)
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Println("create user", username, "id", id)
	return nil
}

func printList(u store.UserStore) error {
	size, err := u.GetSize()
	if err != nil {
		return err
	}
	fmt.Println("total user", size)
	list, err := u.List(1, int(size))
	if err != nil {
		return err
	}
	fmt.Println("Id\tName\tPasswordHash\tData")
	for _, v := range list {
		fmt.Println(v.Id, "\t", v.Name, "\t", v.PasswordHash, "\t", v.Data)
	}
	return nil
}

func main() {
	op := flag.String("op", "add", "operation")
	name := flag.String("name", "", "username")
	password := flag.String("password", "", "user password")
	md5Hash := flag.String("md5_hash", "", "user password 1024-md5-salt hash")
	salt := flag.String("salt", "", "user password hash salt")
	data := flag.String("data", "./data", "data save path")
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	us, err := leveldb.NewUserStore(*data)
	if err != nil {
		fmt.Println("open database error", err)
		os.Exit(1)
	}

	if err := us.Init(10000); err != nil {
		fmt.Println("database init error", err)
		os.Exit(1)
	}

	switch *op {
	case OpAdd:
		if err := addUser(us, *name, *password); err != nil {
			fmt.Println("create user error", *name, *password, err)
			os.Exit(1)
		}
	case OpAddHash:
		if err := addUserHashSalt(us, *name, *md5Hash, *salt); err != nil {
			fmt.Println("create user error", *name, *md5Hash, *salt, err)
			os.Exit(1)
		}
	case OpList:
		if err := printList(us); err != nil {
			fmt.Println("list user error", err)
			os.Exit(1)
		}
	}
}
