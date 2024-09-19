package config

type MongoDB struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Db       string `mapstructure:"db"`
	URI      string `mapstructure:"uri"`
}

type Redis struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Db       string `mapstructure:"db"`
	Prefix   string `mapstructure:"prefix"`
	Secure   bool   `mapstructure:"secure"`
	URI      string `mapstructure:"uri"`
}
