package handler

import (
	"dxkite.cn/log"
	"dxkite.cn/usercenter/hash"
	"dxkite.cn/usercenter/store"
	"encoding/json"
	"net/http"
	"strconv"
)

type AntiBot interface {
	RequireVerifyCaptcha(ip string) bool
	VerifyCaptcha(ip string, answer string) bool
	CountSignFailed(ip string, val int) error
	ClearRequireCaptcha(ip string) error
}

// 登录配置
type SignResp struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

var SignInHandler = func(c AntiBot, us store.UserStore) http.Handler {
	return NewApiHandler(&SignResp{}, func(ctx *HttpContext, input interface{}) (interface{}, int, error) {
		r := input.(*SignResp)
		ip := ClientIp(ctx.request)

		// 检查验证码
		if c.RequireVerifyCaptcha(ip) {
			return nil, ErrCodeCaptcha, ErrCaptcha
		}

		if len(r.Name) == 0 || len(r.Password) == 0 {
			return nil, ErrCodeEmptyUserNameOrPassword, ErrEmptyUserNameOrPassword
		}

		id, err := us.GetId(r.Name)
		if err != nil {
			log.Error("user login", err)
			return nil, ErrCodeUserPasswordError, ErrUserPasswordError
		}

		user, err := us.Get(id)
		if err != nil {
			return nil, ErrCodeUserPasswordError, ErrUserPasswordError
		}

		if hash.VerifyPassword(user.PasswordHash, r.Password) {
			ctx.writer.Header().Set("uin", strconv.Itoa(int(id)))
			_ = c.ClearRequireCaptcha(ip)
		} else {
			_ = c.CountSignFailed(ip, 1)
			return nil, ErrCodeUserPasswordError, ErrUserPasswordError
		}
		return nil, 0, nil
	})
}

var UserInfoHandler = func(us store.UserStore) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, _ := strconv.Atoi(request.Header.Get("uin"))
		uid := uint64(id)
		user, err := us.Get(uid)
		if err != nil {
			log.Error(err)
			JsonError(writer, ErrUserNotExist.Error(), ErrCodeUserNotExist)
			return
		}
		d := &store.UserData{}
		if err := json.Unmarshal([]byte(user.Data), d); err != nil {
			log.Error("unmarshal user info", err)
		}
		resp := map[string]interface{}{}
		resp["id"] = id
		resp["name"] = d.Name
		WriteData(writer, 0, nil, resp)
	})
}
