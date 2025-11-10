package config

import (
	"meili_dog/models"

	"github.com/BurntSushi/toml"
)

// LoadConfig 从TOML文件加载配置
func LoadConfig(path string) (*models.AppConfig, error) {
	var config models.AppConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
