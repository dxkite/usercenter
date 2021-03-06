package server

import (
	"dxkite.cn/log"
	"dxkite.cn/usercenter/captcha"
	"dxkite.cn/usercenter/server/handler"
	"dxkite.cn/usercenter/store"
	"fmt"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

const MaxSignFailedTime = 3

type UserServer struct {
	*http.ServeMux
	// 验证码会话
	capId    map[string]string
	capStore base64Captcha.Store
	cap      *base64Captcha.Captcha
	// 统计登录失败次数
	signFailed map[string]int
}

func (s *UserServer) GetCaptchaId(ip string) (string, error) {
	return s.capId[ip], nil
}

func (s *UserServer) SetCaptchaId(ip, id string) error {
	s.capId[ip] = id
	return nil
}

func (s *UserServer) CountSignFailed(ip string, val int) error {
	s.signFailed[ip] += val
	fmt.Println(ip, s.signFailed[ip])
	return nil
}

func (s *UserServer) RequireVerifyCaptcha(ip string) bool {
	return s.signFailed[ip] > MaxSignFailedTime
}

func (s *UserServer) VerifyCaptcha(ip string, answer string) bool {
	id, ok := s.capId[ip]
	if !ok {
		return false
	}
	return s.cap.Verify(id, answer, true)
}

func (s *UserServer) ClearRequireCaptcha(ip string) error {
	delete(s.signFailed, ip)
	return nil
}

func NewUserServer(us store.UserStore) http.Handler {
	s := &UserServer{ServeMux: http.NewServeMux()}
	s.capStore = base64Captcha.DefaultMemStore
	s.cap = base64Captcha.NewCaptcha(&captcha.DigitConfig, s.capStore)
	s.capId = map[string]string{}
	s.signFailed = map[string]int{}

	// 登录
	s.Handle("/signin", handler.AllowMethod([]string{http.MethodPost},
		handler.SignInHandler(s, us)))

	// 登出
	s.Handle("/signout", handler.AllowMethod([]string{http.MethodPost},
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			uin := request.Header.Get("uin")
			log.Info("logout", uin)
			writer.Header().Set("uin", uin)
			writer.WriteHeader(http.StatusOK)
		})))

	// 获取验证码
	s.Handle("/captcha", handler.CaptchaHandler(s, s, s.cap))
	s.Handle("/verify_captcha", handler.VerifyCaptchaHandler(s, s))
	s.Handle("/profile", handler.UserInfoHandler(us))
	return s
}
