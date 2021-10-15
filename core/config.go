package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/jinzhu/configor"
	"github.com/takemo101/dc-scheduler/core/contract"
)

const (
	Production  = "production"
	Development = "development"
	Local       = "local"
)

// App config
type App struct {
	Name    string
	Host    string
	Port    int
	URL     string
	Version string
	Env     string
	Secret  string
	Debug   bool
	Config  string
}

// DB config
type DB struct {
	Type       string
	Host       string
	Protocol   string
	Port       int
	Name       string
	User       string
	Pass       string
	Charset    string
	Collation  string
	Connection struct {
		Max int
	}
}

// Server fiber config
type Server struct {
	Prefork     bool
	Strict      bool
	Case        bool
	Etag        bool
	BodyLimit   int
	Concurrency int
	Timeout     struct {
		Read  time.Duration
		Write time.Duration
		Idel  time.Duration
	}
	Buffer struct {
		Read  int
		Write int
	}
}

// Log config
type Log struct {
	Server string
}

// File config
type File struct {
	Storage string
	Public  string
	Current string
}

// SMTP config
type SMTP struct {
	Host       string
	Port       int
	Identity   string
	User       string
	Pass       string
	Encryption string
	From       struct {
		Address string
		Name    string
	}
}

// Static config
type Static struct {
	Prefix string
	Root   string
	Index  string
}

// Template config
type Template struct {
	Path   string
	Suffix string
	Reload bool
}

// Cache config
type Cache struct {
	Expiration time.Duration
	Control    bool
}

// Session is config
type Session struct {
	Expiration time.Duration
	Name       string
	Domain     string
	Path       string
	Secure     bool
	HTTPOnly   bool
}

// Cors is config
type Cors struct {
	Origins []string
	MaxAge  time.Duration
}

// Config full config
type Config struct {
	App
	DB
	Server
	Log
	File
	SMTP
	Static
	Template
	Cache
	Session
	Cors
	GoVersion     string
	ConfigMapData map[string]interface{}
}

func (app App) SecretKey() []byte {
	return []byte(app.Secret)
}

// Conf is static
var Conf = Config{}

// NewConfig create configure
func NewConfig(
	path contract.ConfigPath,
) Config {
	// config.yml
	err := configor.Load(&Conf, string(path))
	if err != nil {
		log.Fatalf("fail to load config.yml : %v", err)
	}

	Conf.GoVersion = runtime.Version()

	if Conf.App.Env == "" {
		Conf.App.Env = Local
	}

	if Conf.File.Current == "" {
		current, _ := os.Getwd()
		Conf.File.Current = current
	}

	Conf.ConfigMapData = make(map[string]interface{})
	return Conf
}

// Load config json data
func (c *Config) Load(key string) (map[string]interface{}, error) {
	if mapValue, ok := c.ConfigMapData[key]; ok {
		return mapValue.(map[string]interface{}), nil
	}

	path := path.Join(c.App.Config, key+".json")
	jsonString, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(jsonString, &v)
	if err != nil {
		return nil, err
	}

	c.ConfigMapData[key] = v

	return v.(map[string]interface{}), nil
}

// LoadToValue load config and to value
func (c *Config) LoadToValue(key string, keys string, def interface{}) interface{} {
	data, err := c.Load(key)
	if err == nil {
		arr := strings.Split(keys, ".")
		length := len(arr) - 1
		for i, k := range arr {
			if v, ok := data[k]; ok {
				if length == i {
					return v
				} else {
					data = data[k].(map[string]interface{})
				}
			}
		}
	}

	return def
}

func (c *Config) LoadToValueInt(key string, keys string, def interface{}) int {
	value := c.LoadToValue(key, keys, def).(float64)
	return int(value)
}

func (c *Config) LoadToValueUint(key string, keys string, def interface{}) uint {
	value := c.LoadToValue(key, keys, def).(float64)
	return uint(value)
}

func (c *Config) LoadToValueString(key string, keys string, def interface{}) string {
	return c.LoadToValue(key, keys, def).(string)
}

func (c *Config) LoadToValueMap(key string, keys string, def interface{}) map[string]interface{} {
	return c.LoadToValue(key, keys, def).(map[string]interface{})
}

func (c *Config) LoadToValueArray(key string, keys string, def interface{}) []interface{} {
	return c.LoadToValue(key, keys, def).([]interface{})
}
