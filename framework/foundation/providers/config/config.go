package config

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/zhaoyang1214/ginco/framework/config"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"github.com/zhaoyang1214/ginco/framework/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
}

var _ contract.Provider = (*Config)(nil)

func (c *Config) Build(container contract.Container, params ...interface{}) (interface{}, error) {
	conf := config.NewConfig()
	conf.SetTypeByDefaultValue(true)
	conf.AutomaticEnv()
	appServer, err := container.Get("app")
	if err != nil {
		return nil, err
	}

	app := appServer.(contract.Application)
	envFile := app.BasePath(".env")
	envFileExists := util.PathExists(envFile)
	if envFileExists {
		conf.SetConfigFile(envFile)
		_ = conf.MergeInConfig()
	}

	configPath := app.BasePath("config")
	files, err := ioutil.ReadDir(configPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		if len(ext) > 1 {
			ext = ext[1:]
		}

		if _, ok := util.FindInSlice(viper.SupportedExts, ext); !ok {
			continue
		}

		conf.SetConfigFile(configPath + string(os.PathSeparator) + file.Name())
		if err := conf.MergeInConfig(); err != nil {
			return nil, errors.New(err.Error() + "  Using config file: " + conf.ConfigFileUsed())
		}
	}

	if envFileExists {
		for _, key := range conf.AllKeys() {
			value := conf.GetString(key)
			if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
				envKey := strings.Trim(value, "${}")
				/*v := conf.Get(envKey)
				if v == nil {
					return nil, errors.New("Env Key not found:" + envKey)
				}*/
				newCfg := generateMap(strings.Split(key, "."), conf.GetString(envKey))
				if err = conf.MergeConfigMap(newCfg); err != nil {
					return nil, err
				}
			}
		}
	}

	return conf, nil
}

func generateMap(keys []string, value interface{}) map[string]interface{} {
	if len(keys) == 1 {
		return map[string]interface{}{keys[0]: value}
	}
	v := generateMap(keys[1:], value)
	return map[string]interface{}{keys[0]: v}
}
