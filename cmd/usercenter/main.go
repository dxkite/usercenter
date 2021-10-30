package main

import (
	"dxkite.cn/log"
	"dxkite.cn/usercenter/server"
	"dxkite.cn/usercenter/store/leveldb"
	"flag"
	"net/http"
	"os"
)

func main() {
	data := flag.String("data", "./data", "data save path")
	addr := flag.String("addr", ":8888", "listen addr")
	flag.Parse()

	us, err := leveldb.NewUserStore(*data)
	if err != nil {
		log.Error("open database error", err)
		os.Exit(1)
	}

	if err := us.Init(0); err != nil {
		log.Error("database init error", err)
		os.Exit(1)
	}

	log.Info("start server", *addr)
	_ = http.ListenAndServe(*addr, server.NewUserServer(us))
}
