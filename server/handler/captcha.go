package handler

import (
	"github.com/mojocn/base64Captcha"
	"net/http"
)

type CaptchaStore interface {
	GetCaptchaId(ip string) (string, error)
	SetCaptchaId(ip, id string) error
}

type CaptchaResp struct {
	Required bool   `json:"required"`
	Type     string `json:"type"`
	Data     string `json:"data"`
}

type VerifyCaptcha struct {
	Answer string `json:"answer"`
}

var CaptchaHandler = func(a AntiBot, s CaptchaStore, cap *base64Captcha.Captcha) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ip := ClientIp(req)
		if id, err := s.GetCaptchaId(ip); err != nil {
			JsonError(w, JsonSystemError, http.StatusInternalServerError)
		} else {
			cap.Verify(id, "", true)
		}
		id, base64, err := cap.Generate()
		if err != nil {
			JsonError(w, JsonSystemError, http.StatusInternalServerError)
		}
		_ = s.SetCaptchaId(ip, id)
		d := &CaptchaResp{
			Required: a.RequireVerifyCaptcha(ip),
			Type:     "digit",
			Data:     base64,
		}
		WriteData(w, 0, nil, d)
	})
}

var VerifyCaptchaHandler = func(c AntiBot, s CaptchaStore) http.Handler {
	return NewApiHandler(&VerifyCaptcha{}, func(ctx *HttpContext, input interface{}) (interface{}, int, error) {
		r := input.(*VerifyCaptcha)
		ip := ClientIp(ctx.request)
		// 检查验证码
		if c.RequireVerifyCaptcha(ip) && c.VerifyCaptcha(ip, r.Answer) == false {
			return nil, ErrCodeCaptcha, ErrCaptcha
		}
		// 清除一次验证校验
		_ = c.CountSignFailed(ip, -1)
		// 清除当前验证码
		_ = s.SetCaptchaId(ip, "")
		return nil, 0, nil
	})
}
