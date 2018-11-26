package config

type Log struct {
	Folder string `json:"folder" note:"文件夹路径，空则不输出到文件，输出至系统日至"`
	Level  string `json:"level" note:"输出等级，可选值：error | warning | info | trace | debug"`
}
