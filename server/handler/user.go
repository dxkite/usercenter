package handler

import (
	"dxkite.cn/log"
	"dxkite.cn/usercenter/hash"
	"dxkite.cn/usercenter/store"
	"errors"
	"net/http"
	"strconv"
)

const (
	ErrCodeEmptyUserNameOrPassword int = 50000 + iota
	ErrCodeUserNotExist
	ErrCodePasswordError
)

var (
	ErrEmptyUserNameOrPassword = errors.New("账号或者密码为空")
	ErrUserNotExist            = errors.New("用户不存在")
	ErrPasswordError           = errors.New("账号或者密码错误")
)

// 登录配置
type SignResp struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

var SignInHandler = func(us store.UserStore) http.Handler {
	return NewApiHandler(SignResp{}, func(w http.ResponseWriter, req *http.Request, input interface{}) (interface{}, int, error) {
		r := input.(*SignResp)
		if len(r.Name) == 0 || len(r.Password) == 0 {
			return nil, ErrCodeEmptyUserNameOrPassword, ErrEmptyUserNameOrPassword
		}
		id, err := us.GetId(r.Name)
		if err != nil {
			log.Error("user login", err)
			return nil, ErrCodeUserNotExist, ErrUserNotExist
		}
		user, err := us.Get(id)
		if err != nil {
			return nil, ErrCodeUserNotExist, ErrUserNotExist
		}
		if hash.VerifyPassword(user.PasswordHash, r.Password) {
			w.Header().Set("uin", strconv.Itoa(int(id)))
		} else {
			return nil, ErrCodePasswordError, ErrPasswordError
		}
		return nil, 0, nil
	})
}
