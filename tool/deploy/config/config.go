package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	mutex sync.RWMutex

	Binary Binary `json:"binary"`
	Source Source `json:"source"`
	Site   Site   `json:"site"`
}

func NewConfig() *Config {
	return &Config{
		Binary: Binary{
			Root:  "",
			Files: newBinaryFiles(),
		},
		Source: Source{
			Root: "",
			Ignore: []string{
				"tool",
				".git",
				".idea",
				".gitignore",
				"README.md",
			},
		},
		Site: Site{
			Omw: Source{
				Name: "omw",
				Root: "/home/dev/vue/om/omw",
				Ignore: []string{
					"node_modules",
					"dist",
					".git",
					".idea",
					".gitignore",
					"README.md",
				},
			},
		},
	}
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
