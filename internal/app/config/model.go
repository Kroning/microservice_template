package config

import (
	"time"
)

type Config struct {
	App    App    `mapstructure:"app"`
	Server Server `mapstructure:"server"`
	DB     DB     `mapstructure:"db"`{{if index .Modules "vault"}}
	Vault  Vault  `mapstructure:"vault"`{{end}}
}

type App struct {
	Env      string `mapstructure:"env"`
	Name     string `mapstructure:"name"`
	LogLevel string `mapstructure:"log_level"`
}

type Server struct {
	HTTPPort       *uint `mapstructure:"http_port"`
	GRPC           GRPC  `mapstructure:"grpc"`
	MonitoringPort uint  `mapstructure:"monitoring_port"`
	Logging        bool  `mapstructure:"logging"`
}

type GRPC struct {
	Port                  *uint         `mapstructure:"port"`
	AuthTokens            []string      `mapstructure:"auth_tokens"`
	BadAuthRandomSleepMin time.Duration `mapstructure:"bad_auth_random_sleep_min"`
	BadAuthRandomSleepMax time.Duration `mapstructure:"bad_auth_random_sleep_max"`
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
	Healthy string `json:"healthy"`
}
{{end}}
type SecretString string

func (SecretString) String() string {
	return "*********"
}

func (s SecretString) Value() string {
	return string(s)
}
