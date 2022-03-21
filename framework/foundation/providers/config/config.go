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
			if r, value := replace(conf.Get(key), conf); r {
				newCfg := generateMap(strings.Split(key, "."), value)
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

func replace(value interface{}, conf *config.Config) (r bool, v interface{}) {
	switch value.(type) {
	case string:
		s := value.(string)
		if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
			r = true
			s = conf.GetString(strings.Trim(s, "${}"))
		}
		v = s
	case []interface{}:
		s := value.([]interface{})
		var r1 bool
		for i, v1 := range s {
			r1, s[i] = replace(v1, conf)
			if r1 {
				r = r1
			}
		}
		v = s
	case map[interface{}]interface{}:
		s := value.(map[interface{}]interface{})
		var r1 bool
		for i, v1 := range s {
			r1, s[i] = replace(v1, conf)
			if r1 {
				r = r1
			}
		}
		v = s
	default:
		v = value
	}
	return
}
