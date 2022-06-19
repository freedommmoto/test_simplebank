package tool

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

//all config from .env
type ConfigObject struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOUECE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenConfigKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenLiftTimeConfig time.Duration `mapstructure:"ACCESS_TOKEN_DURATION_TIME"`
	//value is work only start with uppercase
}
type makeNewCustomer struct {
	CustomerName string `json:"customer_name" binding:"required"`
	Currency     string `json:"currency" binding:"required,oneof=USD EUR" `
}

//load config
func LoadConfig(part string) (config ConfigObject, err error) {
	viper.AddConfigPath(part)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into struct %v", err)
	}
	return
}
