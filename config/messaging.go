package config

type RabbitMQ struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Secure   bool   `mapstructure:"secure"`
	URI      string `mapstructure:"uri"`
}
