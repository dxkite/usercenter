package server

import (
	"dxkite.cn/log"
	"dxkite.cn/usercenter/server/handler"
	"dxkite.cn/usercenter/store"
	"net/http"
)

type UserServer struct {
	*http.ServeMux
}

func NewUserServer(us store.UserStore) http.Handler {
	s := &UserServer{ServeMux: http.NewServeMux()}
	// 登录
	s.Handle("/signin", handler.SignInHandler(us))
	// 登出
	s.HandleFunc("/signout", func(writer http.ResponseWriter, request *http.Request) {
		uin := request.Header.Get("uin")
		log.Info("logout", uin)
		writer.Header().Set("uin", uin)
		writer.WriteHeader(http.StatusOK)
	})
	return s
}
