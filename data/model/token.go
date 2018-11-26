package model

import "time"

type Token struct {
	ID          string    `json:"id" note:"标识ID"`
	UserAccount string    `json:"userAccount" note:"用户账号"`
	UserName    string    `json:"userName" note:"用户姓名"`
	LoginIP     string    `json:"loginIp" note:"用户登陆IP"`
	LoginTime   time.Time `json:"loginTime" note:"登陆时间"`
	ActiveTime  time.Time `json:"activeTime" note:"最近激活时间"`
}

type TokenFilter struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	FunId    string `json:"funId"`
}
