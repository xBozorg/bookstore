package config

import (
	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	mySQL    MySQLConfig    `mapstructure:"mysql"`
	jwt      JwtConfig      `mapstructure:"jwt"`
	echo     EchoConfig     `mapstructure:"echo"`
	zarinpal ZarinpalConfig `mapstructure:"zarinpal"`
	redis    RedisConfig    `mapstructure:"redis"`
}

type MySQLConfig struct {
	Name    string `mapstructure:"name"`
	Address string `mapstructure:"address"`
	Net     string `mapstructure:"net"`
	User    string `mapstructure:"user"`
	Pass    string `mapstructure:"pass"`
}
type JwtConfig struct {
	Secret string `mapstructure:"secret"`
}
type EchoConfig struct {
	Mode         string `mapstructure:"mode"`
	HttpPort     string `mapstructure:"http_port"`
	HttpsPort    string `mapstructure:"https_port"`
	LoggerFormat string `mapstructure:"logger_format"`
}
type ZarinpalConfig struct {
	MerchantID string `mapstructure:"merchant_id"`
	Sandbox    bool   `mapstructure:"sandbox"`
}
type RedisConfig struct {
	Address string `mapstructure:"address"`
	Pass    string `mapstructure:"pass"`
	DB      int    `mapstructure:"db"`
}

func (c *Config) GetMySQlConfig() *MySQLConfig       { return &c.mySQL }
func (c *Config) GetJWTConfig() *JwtConfig           { return &c.jwt }
func (c *Config) GetEchoConfig() *EchoConfig         { return &c.echo }
func (c *Config) GetZarinpalConfig() *ZarinpalConfig { return &c.zarinpal }
func (c *Config) GetRedisConfig() *RedisConfig       { return &c.redis }

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
	if err := v.UnmarshalKey("zarinpal", &c.zarinpal); err != nil {
		return err
	}
	if err := v.UnmarshalKey("redis", &c.redis); err != nil {
		return err
	}

	return nil
}
