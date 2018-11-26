package model

import (
	"errors"
	"strings"
)

type Login struct {
	Token   string `json:"token" note:"接口访问凭证" example:"7faf10b0bde847c9905c93966594c82b"`
	Account string `json:"account" required:"true" note:"账号名称"`
	Name    string `json:"name" note:"用户姓名"`
}

type LoginFilter struct {
	Account      string `json:"account" required:"true" note:"账号名称"`
	Password     string `json:"password" required:"true" note:"账号密码"`
	CaptchaId    string `json:"captchaId" required:"true" note:"验证码ID"`
	CaptchaValue string `json:"captchaValue" required:"true" note:"验证码"`
	Encryption   string `json:"encryption" note:"密码加密方法: 空-明文(默认); rsa-RSA密文(公钥通过调用获取验证码接口获取)"`
}

func (s *LoginFilter) Check(captchaRequired bool) error {
	if strings.TrimSpace(s.Account) == "" {
		return errors.New("账号为空")
	}
	if strings.TrimSpace(s.Password) == "" {
		return errors.New("密码为空")
	}
	if strings.TrimSpace(s.CaptchaId) == "" {
		return errors.New("验证码ID为空")
	}
	if captchaRequired {
		if strings.TrimSpace(s.CaptchaValue) == "" {
			return errors.New("验证码为空")
		}
	}

	return nil
}
