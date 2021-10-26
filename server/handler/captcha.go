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

var CaptchaHandler = func(a AntiBot, s CaptchaStore, cap *base64Captcha.Captcha) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ip := req.Header.Get("Client-Ip")
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
