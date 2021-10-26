package captcha

import "github.com/mojocn/base64Captcha"

const (
	NumberLength = 6
	ImageWidth   = 240
	ImageHeight  = 80
)

// 数字验证码
var DigitConfig = base64Captcha.DriverDigit{
	Height:   ImageHeight,
	Width:    ImageWidth,
	Length:   NumberLength,
	MaxSkew:  0,
	DotCount: 0,
}
