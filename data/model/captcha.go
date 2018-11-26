package model

type Captcha struct {
	ID       string `json:"id" note:"验证码标识ID"`
	Value    string `json:"value" note:"验证码图片(base64)"`
	RsaKey   string `json:"rsaKey" note:"RSA公钥, 用于对登陆密码进行加密"`
	Required bool   `json:"required" note:"是否需要验证码"`
}

type CaptchaFilter struct {
	Mode   int `json:"mode" note:"验证码模式: 0-数字； 1-字母； 2-算术； 3-数字字母混合（默认）"`
	Length int `json:"length" note:"验证码长度, 默认4"`
	Width  int `json:"width" note:"图片宽度, 默认100像素"`
	Height int `json:"height" note:"图片高度, 默认30像素"`
}
