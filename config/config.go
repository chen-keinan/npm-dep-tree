package config

import (
	"fmt"
	"github.com/chen-keinan/go-simple-config/simple"
	"go.uber.org/zap"
	"os"
)

//Config wrapper for Viper (exposed to all project's packages)
type Config struct {
	*simple.Config
}

//InitConfig initializes the configuration for the service
func InitConfig(zlog *zap.Logger) *Config {
	configFilePath := getConfigFilePath()
	serviceConfig, err := readConfig(configFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to load configuration from folder %s", configFilePath))
	}
	zlog.Info("Configuration loaded successfully")
	return serviceConfig
}

func getConfigFilePath() string {
	path := "./config/config.json"
	_, err := os.Stat(path)
	if err != nil {
		panic(fmt.Sprintf("Error fetching conf file, %s is not a valid path", path))
	}
	return path
}

func readConfig(path string) (c *Config, err error) {
	c = &Config{Config: simple.New()}
	if err = c.Load(path); err != nil {
		return c, err
	}
	return c, err
}

//GetStringValue returns the string value of the given configuration key
func (c Config) GetStringValue(key string) string {
	return c.Config.GetStringValue(key)
}
