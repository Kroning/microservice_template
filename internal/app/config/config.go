package config

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"{{index .App "git"}}/pkg/logger"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"

	configPathDefault = "config/config.json"
)

const (
	configFilePathEnv  = "CONFIG_FILE_PATH"
)

type Options struct {
	Path          string
	Type          string
	EnvVarsPrefix string
}

func GetConfig(ctx context.Context) *Config {
	return NewConfig(ctx, Options{})
}

// NewConfig loads a config from the sources
func NewConfig(ctx context.Context, opts Options) *Config {
	cfg := &Config{}

	configPath := opts.Path
	if configPath == "" {
		configPath = configPathDefault
	}

	err := Parse(ctx, cfg, Options{Path: configPath}, logger.Logger())
	if err != nil {
		logger.Fatal(ctx, "load config error: ", zap.Error(err))
	}

	if cfg.App.Env == EnvLocal {
		fmt.Printf("%+v\n", cfg)
	} else {
		fmt.Println("Not printing config because cfg.App.Env: ", cfg.App.Env)
	}

	return cfg
}

func Parse(ctx context.Context, configStruct interface{}, opts Options, log *zap.Logger) error {
	err := opts.fill()
	if err != nil {
		return err
	}

	loaded, err := loadFromFile(ctx, opts)
	if err != nil {
		return err
	}
	{{if index .Modules "vault"}}
	fromVault, err := loadFromVault(ctx)
	if err != nil {
		return errors.Wrap(err, "loadFromVault error")
	}
	loaded = loaded || fromVault
	{{end}}
	if !loaded {
		return errors.New("cannot load config from any source")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(true)
	viper.SetEnvPrefix(opts.EnvVarsPrefix)
	viper.AutomaticEnv()

	err = viper.Unmarshal(configStruct)
	if err != nil {
		return err
	}

	return nil
}

func (o *Options) fill() error {
	if o.Type == "" {
		o.Type = "json"
	}
	if os.Getenv(configFilePathEnv) != "" {
		o.Path = os.Getenv(configFilePathEnv)
		o.Type = strings.ToLower(path.Ext(os.Getenv(configFilePathEnv)))
		o.Type = strings.ReplaceAll(o.Type, ".", "")
	}

	return nil
}

func loadFromFile(ctx context.Context, opts Options) (loaded bool, err error) {
	if err := fileExists(opts.Path); err != nil {
		logger.Debug(ctx, "file not exists or empty", zap.String("path", opts.Path), zap.Error(err))
		return false, nil
	}

	logger.Info(ctx, "load config from file", zap.String("path", opts.Path))
	opts.Path, err = filepath.Abs(opts.Path)
	if err != nil {
		return false, err
	}
	viper.SetConfigFile(opts.Path)
	viper.SetConfigType(opts.Type)

	err = viper.ReadInConfig()
	if err != nil {
		return false, err
	}
	logger.Info(ctx, "done loading config from file")
	return true, nil
}

func fileExists(path string) (err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file not exists: " + absPath)
		}
		return
	}
	if info.IsDir() {
		return errors.New("must be file: " + absPath)
	}
	if info.Size() == 0 {
		return errors.New("file is empty: " + absPath)
	}
	return
}
