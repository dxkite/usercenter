package main

import (
	"dxkite.cn/usercenter/server"
	"dxkite.cn/usercenter/store/leveldb"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	data := flag.String("data", "./data", "data save path")
	flag.Parse()

	us, err := leveldb.NewUserStore(*data)
	if err != nil {
		fmt.Println("open database error", err)
		os.Exit(1)
	}

	if err := us.Init(0); err != nil {
		fmt.Println("database init error", err)
		os.Exit(1)
	}

	_ = http.ListenAndServe(":8888", server.NewUserServer(us))
}