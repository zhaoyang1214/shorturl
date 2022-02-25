package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"sort"
	"time"
)

type Manager struct {
	app      contract.Application
	channels map[string]contract.Logger
}

var _ contract.LoggerManager = (*Manager)(nil)

func NewManager(app contract.Application) *Manager {
	return &Manager{
		app:      app,
		channels: make(map[string]contract.Logger),
	}
}

func (m *Manager) Channel(channel string) contract.Logger {
	return m.Driver(channel)
}

func (m *Manager) Driver(driver string) contract.Logger {
	if driver == "" {
		driver = m.getConfig().GetString("logger.default")
	}
	return m.get(driver)
}

func (m *Manager) get(name string) contract.Logger {
	if log, ok := m.channels[name]; ok {
		return log
	}
	m.channels[name] = m.resolve(name)
	return m.channels[name]
}

func (m *Manager) configName(name string) string {
	return "logger.channels." + name
}

func (m *Manager) resolve(name string) contract.Logger {
	core := m.createCore(name)
	opts := m.buildOptions(name)
	return &Logger{zap.New(core, opts...)}
}

func (m *Manager) createCore(name string) zapcore.Core {
	configName := m.configName(name)
	c := m.getConfig()
	conf := c.GetStringMap(configName)
	if conf == nil {
		panic("Logger config [" + configName + "] is not defined")
	}
	driver := c.GetString(configName + ".driver")
	var core zapcore.Core
	switch driver {
	case "stderr":
		core = m.createStderrCore(name)
	case "single":
		core = m.createSingleCore(name)
	case "stack":
		core = m.createStackCore(name)
	case "rotation":
		core = m.createRotationCore(name)
	default:
		panic("Logger driver [" + driver + "] is not supported")
	}
	return core
}

func (m *Manager) createStderrCore(name string) zapcore.Core {
	configName := m.configName(name)
	enc := m.buildEncoder(configName)
	sink, _, err := zap.Open("stderr")
	if err != nil {
		panic(err)
	}

	level := m.getConfigLevel(configName)
	return zapcore.NewCore(enc, sink, level)
}

func (m *Manager) createSingleCore(name string) zapcore.Core {
	configName := m.configName(name)
	enc := m.buildEncoder(configName)
	c := m.getConfig()
	p := c.GetString(configName + ".path")
	if p == "" {
		panic("Driver [" + c.GetString(configName+".driver") + "] path is not defined")
	}

	p = path.Join(m.app.RuntimePath(), p)
	sink, err := os.OpenFile(p, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	level := m.getConfigLevel(configName)
	return zapcore.NewCore(enc, sink, level)
}

func (m *Manager) createStackCore(name string) zapcore.Core {
	configName := m.configName(name)
	c := m.getConfig().Sub(configName)
	channels := c.GetStringSlice("channels")
	if len(channels) == 0 {
		panic("Driver stack channels is empty")
	}
	var cores []zapcore.Core
	for _, ch := range channels {
		cores = append(cores, m.createCore(ch))
	}
	return zapcore.NewTee(cores...)
}

func (m *Manager) createRotationCore(name string) zapcore.Core {
	configName := m.configName(name)
	enc := m.buildEncoder(configName)
	c := m.getConfig().Sub(configName)
	p := c.GetString("path")
	if p == "" {
		panic("Driver [" + c.GetString(configName+".driver") + "] path is not defined")
	}
	p = path.Join(m.app.RuntimePath(), p)

	var opts []rotatelogs.Option
	if c.Has("maxAge") {
		v := c.GetDuration("maxAge") * time.Hour
		opts = append(opts, rotatelogs.WithMaxAge(v))
	}

	if c.Has("rotationTime") {
		v := c.GetDuration("rotationTime") * time.Hour
		opts = append(opts, rotatelogs.WithRotationTime(v))
	}

	if c.Has("rotationCount") {
		opts = append(opts, rotatelogs.WithRotationCount(c.GetUint("rotationCount")))
	}

	if c.Has("rotationSize") {
		opts = append(opts, rotatelogs.WithRotationSize(c.GetInt64("rotationSize")))
	}

	hook, err := rotatelogs.New(p, opts...)
	if err != nil {
		panic(err)
	}

	return zapcore.NewCore(enc, zapcore.AddSync(hook), m.getConfigLevel(configName))
}

func (m *Manager) buildEncoder(configName string) zapcore.Encoder {
	c := m.getConfig().Sub(configName)
	cfg := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "ts",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	timeFormat := "2006-01-02 15:04:05.000"

	if c.Has("encoderConfig") {
		ec := c.Sub("encoderConfig")
		if v := ec.GetString("messageKey"); v != "" {
			cfg.MessageKey = v
		}

		if v := ec.GetString("levelKey"); v != "" {
			cfg.LevelKey = v
		}

		if v := ec.GetString("timeKey"); v != "" {
			cfg.TimeKey = v
		}

		if v := ec.GetString("nameKey"); v != "" {
			cfg.NameKey = v
		}

		if v := ec.GetString("callerKey"); v != "" {
			cfg.CallerKey = v
		}

		if v := ec.GetString("functionKey"); v != "" {
			cfg.FunctionKey = v
		}

		if v := ec.GetString("stacktraceKey"); v != "" {
			cfg.StacktraceKey = v
		}

		if v := ec.GetString("lineEnding"); v != "" {
			cfg.LineEnding = v
		}

		if v := ec.GetString("timeEncoder"); v != "" {
			timeFormat = v
		}

		if v := ec.GetString("callerEncoder"); v != "" {
			var callerEncoder zapcore.CallerEncoder
			if err := callerEncoder.UnmarshalText([]byte(v)); err != nil {
				panic(err)
			}
			cfg.EncodeCaller = callerEncoder
		}

		if v := ec.GetString("consoleSeparator"); v != "" {
			cfg.ConsoleSeparator = v
		}
	}

	cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(timeFormat))
	}

	if v := c.GetString("encoding"); v == "json" {
		return zapcore.NewJSONEncoder(cfg)
	} else {
		return zapcore.NewConsoleEncoder(cfg)
	}
}

func (m *Manager) buildOptions(name string) []zap.Option {
	configName := m.configName(name)
	c := m.getConfig().Sub(configName)

	disableCaller := c.GetBool("disableCaller")
	opts := []zap.Option{zap.WithCaller(!disableCaller)}
	if !disableCaller {
		callerSkip := 2
		if v := c.GetInt("callerSkip"); v > 0 {
			callerSkip = v
		}
		opts = append(opts, zap.AddCallerSkip(callerSkip))
	}

	stackLevel := ErrorLevel
	if v := c.GetBool("development"); v {
		opts = append(opts, zap.Development())
		stackLevel = WarnLevel
	}

	if v := c.GetBool("disableStacktrace"); v {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if v := c.GetStringMap("sampling"); len(v) > 0 {
		initial := c.GetInt("sampling.initial")
		thereafter := c.GetInt("sampling.thereafter")
		if initial > 0 && thereafter > 0 {
			opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewSamplerWithOptions(core, time.Second, initial, thereafter)
			}))
		}
	}

	if v := c.GetStringMap("initialFields"); len(v) > 0 {
		fs := make([]zap.Field, 0, len(v))
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, v[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	return opts
}

func (m *Manager) getConfig() contract.Config {
	configServer, err := m.app.Get("config")
	if err != nil {
		panic(err)
	}
	return configServer.(contract.Config)
}

func (m *Manager) getConfigLevel(configName string) zap.AtomicLevel {
	c := m.getConfig()
	level := c.GetString(configName + ".level")
	l := zap.NewAtomicLevel()
	if err := l.UnmarshalText([]byte(level)); err != nil {
		panic(err)
	}
	return l
}

func (m *Manager) Debug(msg string, context ...interface{}) {
	m.Driver("").Debug(msg, context...)
}

func (m *Manager) Info(msg string, context ...interface{}) {
	m.Driver("").Info(msg, context...)
}

func (m *Manager) Warn(msg string, context ...interface{}) {
	m.Driver("").Warn(msg, context...)
}

func (m *Manager) Error(msg string, context ...interface{}) {
	m.Driver("").Error(msg, context...)
}

func (m *Manager) Panic(msg string, context ...interface{}) {
	m.Driver("").Panic(msg, context...)
}

func (m *Manager) Fatal(msg string, context ...interface{}) {
	m.Driver("").Fatal(msg, context...)
}

func (m *Manager) Log(level interface{}, msg string, context ...interface{}) {
	m.Driver("").Log(level, msg, context...)
}

func (m *Manager) Sync() error {
	return m.Driver("").Sync()
}
