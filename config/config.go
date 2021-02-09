package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

//Config wrapper for Viper (exposed to all project's packages)
type Config struct {
	*viper.Viper
	log *zap.Logger
}

//InitConfig initializes the configuration for the service
func InitConfig(zlog *zap.Logger) *Config {
	configFilePath := getConfigFilePath()
	serviceConfig, err := readConfig(configFilePath)
	serviceConfig.log = zlog
	if err != nil {
		panic(fmt.Errorf("failed to load configuration from folder %s", configFilePath))
	}
	serviceConfig.log.Info("Configuration loaded successfully")
	return serviceConfig
}

func getConfigFilePath() string {
	path := "config/"
	fs, err := os.Stat(path)
	if err != nil || !fs.IsDir() {
		panic(fmt.Sprintf("Error fetching conf file, %s is not a valid path", path))
	}
	return path
}

func readConfig(path string) (c *Config, err error) {
	c = &Config{Viper: viper.New()}
	c.AddConfigPath(path)
	if err = c.ReadInConfig(); err != nil {
		return c, err
	}
	return c, err
}

//GetStringValue returns the string value of the given configuration key
func (c Config) GetStringValue(key string) string {
	return c.GetString(key)
}
