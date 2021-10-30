package handler

import "errors"

const (
	ErrCodeEmptyUserNameOrPassword int = 50000 + iota
	ErrCodeUserPasswordError
	ErrCodeCaptcha
	ErrCodeUserNotExist
)

var (
	ErrEmptyUserNameOrPassword = errors.New("账号或者密码为空")
	ErrUserPasswordError       = errors.New("账号或者密码错误")
	ErrCaptcha                 = errors.New("验证码错误")
	ErrUserNotExist            = errors.New("用户不存在")
)
