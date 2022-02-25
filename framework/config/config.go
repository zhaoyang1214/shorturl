package config

import (
	"github.com/spf13/viper"
	"github.com/zhaoyang1214/ginco/framework/contract"
)

type Config struct {
	viper.Viper
}

var _ contract.Config = (*Config)(nil)

type DecoderConfigOption = viper.DecoderConfigOption

func NewConfig() *Config {
	return &Config{
		*viper.New(),
	}
}

func (c *Config) Has(key string) bool {
	return c.IsSet(key)
}

func (c *Config) convOpts(opts ...interface{}) []DecoderConfigOption {
	dco := make([]DecoderConfigOption, len(opts))
	for i, opt := range opts {
		switch opt.(type) {
		case DecoderConfigOption:
			dco[i] = opt.(DecoderConfigOption)
		}
	}
	return dco
}

func (c *Config) Sub(key string) contract.Config {
	return &Config{
		*c.Viper.Sub(key),
	}
}

func (c *Config) UnmarshalKey(key string, rawVal interface{}, opts ...interface{}) error {
	return c.Viper.UnmarshalKey(key, rawVal, c.convOpts(opts...)...)
}

func (c *Config) Unmarshal(rawVal interface{}, opts ...interface{}) error {
	return c.Viper.Unmarshal(rawVal, c.convOpts(opts...)...)
}

func (c *Config) UnmarshalExact(rawVal interface{}, opts ...interface{}) error {
	return c.Viper.UnmarshalExact(rawVal, c.convOpts(opts...)...)
}
