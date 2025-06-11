package conf

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	// "github.com/kr/pretty"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Env      string
	Kitex    Kitex    `yaml:"kitex"`
	Hertz    Hertz    `yaml:"hertz"`
	Tidb     Tidb     `yaml:"tidb"`
	Alipay   Alipay   `yaml:"alipay"`
	Postgres Postgres `yaml:"postgres"`
	Redis    Redis    `yaml:"redis"`
	Registry Registry `yaml:"registry"`
	Auth     Auth     `yaml:"auth"`
}

type Tidb struct {
	DSN      string `yaml:"dsn"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type Postgres struct {
	DSN      string `yaml:"dsn"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type Alipay struct {
	APPID        string `yaml:"app_id"`
	PrivateKey   string `yaml:"private_key"`
	ServerDomain string `yaml:"server_domain"`
}

type Redis struct {
	DSN      string `yaml:"dsn"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
}

type Hertz struct {
	Address         string `yaml:"address"`
	MetricsPort     string `yaml:"metrics_port"`
	EnablePprof     bool   `yaml:"enable_pprof"`
	EnableGzip      bool   `yaml:"enable_gzip"`
	EnableAccessLog bool   `yaml:"enable_access_log"`
	LogLevel        string `yaml:"log_level"`
	LogFileName     string `yaml:"log_file_name"`
	LogMaxSize      int    `yaml:"log_max_size"`
	LogMaxBackups   int    `yaml:"log_max_backups"`
	LogMaxAge       int    `yaml:"log_max_age"`
	RegistryAddr    string `yaml:"registry_addr"`
}

type Kitex struct {
	Service         string `yaml:"service"`
	Address         string `yaml:"address"`
	MetricsPort     string `yaml:"metrics_port"`
	EnablePprof     bool   `yaml:"enable_pprof"`
	EnableGzip      bool   `yaml:"enable_gzip"`
	EnableAccessLog bool   `yaml:"enable_access_log"`
	LogLevel        string `yaml:"log_level"`
	LogFileName     string `yaml:"log_file_name"`
	LogMaxSize      int    `yaml:"log_max_size"`
	LogMaxBackups   int    `yaml:"log_max_backups"`
	LogMaxAge       int    `yaml:"log_max_age"`
}

type Registry struct {
	RegistryAddress []string `yaml:"registry_address"`
	ServiceId       string   `yaml:"service_id"`
}

type Auth struct {
	WhiteList []string `yaml:"white_list"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	prefix := "conf"
	confFileRelPath := filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yml"))
	content, err := os.ReadFile(confFileRelPath)
	if err != nil {
		panic(err)
	}
	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		klog.Error("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		klog.Error("validate config error - %v", err)
		panic(err)
	}
	conf.Env = GetEnv()
	// pretty.Printf("%+v\n", conf)
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "develop"
	}
	return e
}

func LogLevel() klog.Level {
	level := GetConf().Kitex.LogLevel
	switch level {
	case "trace":
		return klog.LevelTrace
	case "debug":
		return klog.LevelDebug
	case "info":
		return klog.LevelInfo
	case "notice":
		return klog.LevelNotice
	case "warn":
		return klog.LevelWarn
	case "error":
		return klog.LevelError
	case "fatal":
		return klog.LevelFatal
	default:
		return klog.LevelInfo
	}
}
