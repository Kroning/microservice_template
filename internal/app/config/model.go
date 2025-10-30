package config

import (
	"time"
)

type Config struct {
	App        App        `mapstructure:"app"`{{if index .Modules "http_chi"}}
	HTTPServer HTTPServer `mapstructure:"http_server"`{{end}}
	DB         DB         `mapstructure:"db"`{{if index .Modules "vault"}}
	Vault      Vault      `mapstructure:"vault"`{{end}}
}

type App struct {
	Env      string `mapstructure:"env"`
	Name     string `mapstructure:"name"`
	LogLevel string `mapstructure:"log_level"`
}

type HTTPServer struct {
	Port                    uint          `mapstructure:"port"`
	ReadHeaderTimeout       time.Duration `mapstructure:"read_header_timeout"`
	ReadTimeout             time.Duration `mapstructure:"read_timeout"`
	WriteTimeout            time.Duration `mapstructure:"write_timeout"`
	IdleTimeout             time.Duration `mapstructure:"idle_timeout"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
}

type DB struct {
	Master     DBConfig `mapstructure:"master"`
	Slave      DBConfig `mapstructure:"slave"`
	Metrics    bool     `mapstructure:"metrics"`
	Migrations bool     `mapstructure:"migrations"`
}

type DBConfig struct {
	Host     string        `mapstructure:"host"`
	Port     string        `mapstructure:"port"`
	User     string        `mapstructure:"user"`
	Password SecretString  `mapstructure:"password"`
	Database string        `mapstructure:"database"`
	MaxOpen  uint          `mapstructure:"max_open"`
	Timeout  time.Duration `mapstructure:"timeout"`
}
{{if index .Modules "vault"}}
type Vault struct {
	Healthy string `mapstructure:"healthy"`
}
{{end}}
type SecretString string

func (SecretString) String() string {
	return "*********"
}

func (s SecretString) Value() string {
	return string(s)
}
