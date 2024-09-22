package config

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/spf13/viper"
)

var singleton sync.Once
var Cfg Config

type Config struct {
	App struct {
		Name  string `mapstructure:"name"`
		Env   string `mapstructure:"env"`
		Debug bool   `mapstructure:"debug"`
		Host  string `mapstructure:"host"`
		Port  string `mapstructure:"port"`
	} `mapstructure:"app"`

	Auth struct {
		JwtSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"auth"`

	DB struct {
		MongoDB  MongoDB `mapstructure:"mongodb"`
		Redis    Redis   `mapstructure:"redis"`
		MysqlUri string  `mapstructure:"mysql_uri"`
	} `mapstructure:"db"`

	Messaging struct {
		RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
	} `mapstructure:"messaging"`
}

func InitConfig(cfgPath string) error {
	var oerr error
	singleton.Do(func() {
		defaultPath := "config.yml"

		if cfgPath == "" {
			cfgPath = defaultPath
		}

		// Silently create config file if it doesn't exist
		_, err := os.Stat(cfgPath)
		if os.IsNotExist(err) {
			cfgFile, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				defer cfgFile.Close()
				expCfgFile, err := os.Open(defaultPath + ".example")
				if err == nil {
					defer expCfgFile.Close()
					_, _ = io.Copy(cfgFile, expCfgFile)
				}
			}
		}

		viper.SetConfigFile(cfgPath)
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// Membaca konfigurasi
		err = viper.ReadInConfig()
		if err != nil {
			oerr = err
			return
		}

		// Set default configurations
		viper.SetDefault("app.name", "xyz-multifinance")
		viper.SetDefault("app.env", "development")
		viper.SetDefault("app.debug", true)
		viper.SetDefault("app.host", "localhost")
		viper.SetDefault("app.port", "8000")

		viper.SetDefault("auth.jwt_secret", "ajh73&39h2j3b(*31)")

		viper.SetDefault("db.mongodb_uri", "mongodb://localhost:27017/xyz_multifinance")
		viper.SetDefault("db.mysql_uri", "root@tcp(127.0.0.1:3306)/xyz_multifinance?charset=utf8mb4&parseTime=True&loc=Local")
		viper.SetDefault("db.redis.addr", "localhost:6379")
		viper.SetDefault("db.redis.password", "")
		viper.SetDefault("db.redis.db", 0)

		viper.SetDefault("messaging.rabbitmq.url", "amqp://localhost:5672")

		// Unmarshal ke struct Config
		err = viper.Unmarshal(&Cfg)
		if err != nil {
			log.Fatalln("cannot unmarshal config:", err)
		}
	})
	return oerr
}

func InitTestConfig() error {
	p, err := utils.PathFromRoot("config.test.yml")
	if err != nil {
		panic(err)
	}
	gin.SetMode(gin.TestMode)
	return InitConfig(p)
}
