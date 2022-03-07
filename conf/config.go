/*
-- @Time : 2022/3/7 14:35
-- @Author : raoxiaoya
-- @Desc :
*/
package conf

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Apps Apps `yaml:"Apps"`
}

type Item struct {
	Type     string `yaml:"Type"`
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Database string `yaml:"Database"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

type Apps struct {
	Database Item `yaml:"Database"`
}

func ReadYamlConfig(path string) (*Config, error) {
	conf := &Config{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(conf)
		if err != nil {
			return nil, err
		}
	}
	return conf, nil
}
