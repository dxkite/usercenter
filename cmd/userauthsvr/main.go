package main

import (
	"context"
	gtcfg "dxkite.cn/gateway/config"
	"dxkite.cn/gateway/route"
	gtsvr "dxkite.cn/gateway/server"
	"dxkite.cn/log"
	"dxkite.cn/usercenter/config"
	"dxkite.cn/usercenter/server"
	"dxkite.cn/usercenter/store/leveldb"
	"flag"
	"net"
	"net/http"
	"os"
)

func init() {
	log.SetOutput(log.NewColorWriter())
	log.SetLogCaller(true)
	log.SetAsync(false)
	log.SetLevel(log.LMaxLevel)
}

func main() {

	ctx, _ := context.WithCancel(context.Background())
	conf := flag.String("conf", "./config.yml", "the config file")
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	cfg := gtcfg.NewConfig()

	if err := cfg.LoadFromFile(*conf); err != nil {
		log.Error(err)
	}

	userCfg := &config.Config{}
	if err := userCfg.LoadFromFile(*conf); err != nil {
		log.Error(err)
	}

	s := gtsvr.NewPortable(ctx, cfg)

	if userCfg.EnableUser {
		us, err := leveldb.NewUserStore(userCfg.DataPath)
		if err != nil {
			log.Error("open database error", err)
			os.Exit(1)
		}

		if err := us.Init(0); err != nil {
			log.Error("database init error", err)
			os.Exit(1)
		}

		userSvr := server.NewUserServer(us)
		userSvr = http.StripPrefix("/api/user", userSvr)
		s.Route.AddDynamicRoute("/api/user/signin", gtsvr.NewHandler(&route.RouteConfig{
			SignIn: true,
		}, userSvr))
		s.Route.AddDynamicRoute("/api/user/signout", gtsvr.NewHandler(&route.RouteConfig{
			SignOut: true,
		}, userSvr))
		s.Route.AddDynamicRoute("/api/user/captcha", gtsvr.NewHandler(&route.RouteConfig{}, userSvr))
		s.Route.AddDynamicRoute("/api/user/verify_captcha", gtsvr.NewHandler(&route.RouteConfig{}, userSvr))
		s.Route.AddDynamicRoute("/api/user", gtsvr.NewHandler(&route.RouteConfig{
			Sign: true,
		}, userSvr))
	}

	l, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Error(err)
	}

	log.Println("server start at", l.Addr(), "user-enable", userCfg.EnableUser)
	if err := s.Serve(l); err != nil {
		log.Error(err)
	}

	gtsvr.NewPortable(ctx, cfg)
}
