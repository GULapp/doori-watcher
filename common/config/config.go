package config

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	LOG "watcher/common/log"
)

type To struct {
	Ip   string `toml:"ip"`
	Port string `toml:"port"`
}

type Info struct {
	To    To `toml:"to"`
}

type Config struct {
	Site   string `toml:"site"`
	Domain string `toml:"domain"`
	Agent  Info   `toml:"agent"`
	Ui     Info   `toml:"ui"`
}

func InitConfig(path string) (Config, error) {
	var config Config
	var err error
	var file *os.File
	file, err = os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		LOG.Error("failed to call os.OpenFile()")
		return config, err
	}

	tomolData, err := ioutil.ReadAll(file)
	if err != nil {
		LOG.Error("failed to call ioutil.ReadAll()")
		return config, err
	}

	err = toml.Unmarshal(tomolData, &config)
	if err != nil {
		LOG.Error("failed to call ioutil.ReadAll()")
		return config, err
	}
	return config, nil
}
