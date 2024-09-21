package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type EnvType string // Is an application environment name

const (
	EnvLocal EnvType = "local"   // is a name of local environment
	EnvDev   EnvType = "develop" // is a name develop environment
	EnvProd  EnvType = "prod"    // is a name production name environment
)

type GRPCConfig struct {
	Port int `yaml:"port" env-default:"443"` // Port for starting server
}

type DatabaseConfig struct {
	Name     string `yaml:"dbname" env-default:"postgres"`   // is a name of database
	Host     string `yaml:"host" env-default:"localhost"`    // is a host name of database
	Port     int    `yaml:"port" env-default:"5432"`         // is a port of database
	Username string `yaml:"username" env-default:"postgres"` // user for connect to database
	Password string `yaml:"password" env-default:"postgres"` // password for connect for database
}

type Config struct {
	Env  EnvType        `yaml:"env" env-default:"local"` // Name of application environment
	GRPC GRPCConfig     `yaml:"grpc"`                    // gRPC server configuration
	DB   DatabaseConfig `yaml:"database"`                // Database configuration
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath
// Priority: flag > env > default
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
