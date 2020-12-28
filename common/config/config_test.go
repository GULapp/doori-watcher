package config

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	config, err := InitConfig("./config.toml")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(config.Site)
	t.Log(config.Domain)
	t.Log(config.Ui.To.Ip)
	t.Log(config.Ui.To.Port)
}