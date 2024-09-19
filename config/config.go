package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/spf13/viper"
)

var Cfg Config
var singleton sync.Once

type Config struct {
	App struct {
		Name              string `mapstructure:"name"`
		Env               string `mapstructure:"env"`
		Debug             bool   `mapstructure:"debug"`
		Host              string `mapstructure:"host"`
		Port              string `mapstructure:"port"`
		ApiUrl            string `mapstructure:"api_url"`
		AdminDashboardUrl string `mapstructure:"admin_dashboard_url"`
		PortalUrl         string `mapstructure:"portal_url"`
		SsoUrl            string `mapstructure:"sso_url"`
	} `mapstructure:"app"`

	Auth struct {
		JwtSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"auth"`

	DB struct {
		MongoDB MongoDB `mapstructure:"mongodb"`
		Redis   Redis   `mapstructure:"redis"`
	} `mapstructure:"db"`

	Messaging struct {
		RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
	} `mapstructure:"messaging"`

	// Indibiz struct {
	// 	FABDApiUrl       string `mapstructure:"fabd_api_url"`
	// 	FABDAccessKey    string `mapstructure:"fabd_access_key"`
	// 	FABDAccessSecret string `mapstructure:"fabd_access_secret"`
	// } `mapstructure:"indibiz"`
}

func InitConfig(cfgPath string) error {
	var oerr error
	singleton.Do(func() {
		defaultPath := "config.yml"

		if cfgPath == "" {
			cfgPath = defaultPath
		}

		// Silently make config file if it doesn't exist
		_, err := os.Stat(cfgPath)
		if err != nil {
			if os.IsNotExist(err) {
				cfgFile, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_WRONLY, 0644)
				defer cfgFile.Close()
				if err == nil {
					expCfgFile, err := os.Open(defaultPath + ".example")
					defer expCfgFile.Close()
					if err == nil {
						_, _ = io.Copy(cfgFile, expCfgFile)
					}
				}
			}
		}

		viper.SetConfigFile(cfgPath)
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		err = viper.ReadInConfig()
		if err != nil {
			oerr = err
			return
		}
		fmt.Println("default path ", cfgPath)
		viper.SetDefault("app.env", "development")
		viper.SetDefault("app.debug", "false")
		viper.SetDefault("app.host", "")
		viper.SetDefault("app.port", "8000")
		viper.SetDefault("app.name", "")
		viper.SetDefault("app.api_url", "")
		viper.SetDefault("app.admin_dashboard_url", "")
		viper.SetDefault("app.portal_url", "")

		viper.SetDefault("db.mongodb.db", "rivality")
		viper.SetDefault("db.mongodb.host", "127.0.0.1")
		viper.SetDefault("db.mongodb.port", "27017")
		viper.SetDefault("db.mongodb.username", "")
		viper.SetDefault("db.mongodb.password", "")
		viper.SetDefault("db.mongodb.uri", "mongodb://localhost:27017/")

		viper.SetDefault("db.redis.db", "0")
		viper.SetDefault("db.redis.host", "127.0.0.1")
		viper.SetDefault("db.redis.port", "6379")
		viper.SetDefault("db.redis.username", "")
		viper.SetDefault("db.redis.password", "")
		viper.SetDefault("db.redis.secure", "false")
		viper.SetDefault("db.redis.prefix", "")
		viper.SetDefault("db.redis.uri", "")

		viper.SetDefault("messaging.rabbitmq.secure", "false")
		viper.SetDefault("messaging.rabbitmq.host", "127.0.0.1")
		viper.SetDefault("messaging.rabbitmq.port", "5671")
		viper.SetDefault("messaging.rabbitmq.username", "")
		viper.SetDefault("messaging.rabbitmq.password", "")
		viper.SetDefault("messaging.rabbitmq.uri", "")

		err = viper.Unmarshal(&Cfg)
		if err != nil {
			log.Fatalln("cannot unmarshaling config")
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
