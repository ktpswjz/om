package config

type Token struct {
	Expiration int64 `json:"expiration" note:"凭证过期时间, 单位分钟, 默认30, 0表示永不过期"`
}
