package config

import (
	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	mySQL MySQLConfig `mapstructure:"mysql"`
	jwt   JwtConfig   `mapstructure:"jwt"`
	echo  EchoConfig  `mapstructure:"echo"`
}

type MySQLConfig struct {
	Name    string `mapstructure:"name"`
	Address string `mapstructure:"address"`
	Net     string `mapstructure:"net"`
	User    string `mapstructure:"user"`
	Pass    string `mapstructure:"pass"`
}
type JwtConfig struct {
	Secret        string `mapstructure:"secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
}
type EchoConfig struct {
	Mode         string `mapstructure:"mode"`
	HttpPort     string `mapstructure:"http_port"`
	HttpsPort    string `mapstructure:"https_port"`
	LoggerFormat string `mapstructure:"logger_format"`
}

func (c Config) GetMySQlConfig() MySQLConfig { return c.mySQL }
func (c Config) GetJWTConfig() JwtConfig     { return c.jwt }
func (c Config) GetEchoConfig() EchoConfig   { return c.echo }

func (c *Config) Read() error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.UnmarshalKey("mysql", &c.mySQL); err != nil {
		return err
	}
	if err := v.UnmarshalKey("jwt", &c.jwt); err != nil {
		return err
	}
	if err := v.UnmarshalKey("echo", &c.echo); err != nil {
		return err
	}

	return nil
}
