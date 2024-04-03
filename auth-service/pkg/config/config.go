package config

import "github.com/spf13/viper"

type Config struct {
	Port              string `mapstructure:"PORT"`
	DBUrl             string `mapstructure:"DB_URL"`
	JWTSecretKey      string `mapstructure:"JWT_SECRET_KEY"`
	GalileoUrl        string `mapstructure:"GALILEO_URL"`
	GalileoLogin      string `mapstructure:"GALILEO_LOGIN"`
	GalileoTranskey   string `mapstructure:"GALILEO_TRANSKEY"`
	GalileoProviderId string `mapstructure:"GALILEO_PROVIDER_ID"`
	GalileoProductId  string `mapstructure:"GALILEO_PRODUCT_ID"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
