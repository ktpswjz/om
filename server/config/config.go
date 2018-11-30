package config

import (
	"encoding/json"
	"fmt"
	"github.com/ktpswjz/httpserver/http/server/configure"
	"github.com/ktpswjz/httpserver/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	path  string // file path to be loaded
	mutex sync.RWMutex
	args  *types.Args

	Name    string           `json:"name" note:"服务平台名称"`
	Log     Log              `json:"log" note:"日志"`
	Server  configure.Server `json:"server" note:"服务器"`
	Site    Site             `json:"site" note:"网站"`
	Proxy   Proxy            `json:"proxy" note:"代理"`
	Service Service          `json:"service" note:"服务"`
}

func NewConfig() *Config {
	return &Config{
		Name: "服务器管理平台",
		Log: Log{
			Folder: "",
			Level:  "error|warning|info",
		},
		Server: configure.Server{
			Http: configure.Http{
				Enabled: true,
				Address: "",
				Port:    "9618",
			},
			Https: configure.Https{
				Enabled: false,
				Address: "",
				Port:    "9613",
				Cert: configure.Certificate{
					File:     "",
					Password: "",
				},
			},
		},
		Site: Site{
			Root: "",
			Doc: SiteDoc{
				Enable: false,
				Root:   "",
			},
			Omw: SiteOmw{
				Root: "",
				Api: SiteOmwApi{
					Token: Token{
						Expiration: 30,
					},
				},
				Users: []SiteOmwUser{
					{
						Account:  "admin",
						Password: "1",
					},
				},
				Ldap: Ldap{
					Enable: false,
					Host:   "example.com",
					Port:   389,
					Base:   "dc=example,dc=com",
					Groups: make([]string, 0),
				},
				Setting: SiteOmwSetting{
					MenuVisible: true,
					Dashboard: SiteOmwSettingDashboard{
						TomcatVisible: true,
						ProxyVisible:  true,
					},
				},
			},
		},
		Proxy: Proxy{
			Enable: true,
			Http: ProxyServer{
				IP:      "",
				Port:    "80",
				Targets: make([]*ProxyTarget, 0),
			},
			Https: ProxyServer{
				IP:      "",
				Port:    "443",
				Targets: make([]*ProxyTarget, 0),
			},
		},
		Service: Service{
			Tomcats: make(ServiceTomcatCollection, 0),
			Others:  make(ServiceOtherCollection, 0),
			Jar: ServiceJar{
				Root:   "",
				Prefix: "jar-",
			},
		},
	}
}

func (s *Config) SetPath(path string) {
	s.path = path
}

func (s *Config) GetPath() string {
	return s.path
}

func (s *Config) SetArgs(args *types.Args) {
	s.args = args
}

func (s *Config) GetArgs() *types.Args {
	return s.args
}

func (s *Config) GetServer() *configure.Server {
	return &s.Server
}

func (s *Config) LoadFromFile(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

func (s *Config) SaveToFile(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	fileFolder := filepath.Dir(filePath)
	_, err = os.Stat(fileFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(fileFolder, 0777)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, string(bytes[:]))

	return err
}

func (s *Config) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

func (s *Config) FormatString() string {
	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes[:])
}
