package config

import "github.com/spf13/viper"

type Config struct {
	Port                 string `mapstructure:"PORT"`
	DBUrl                string `mapstructure:"DB_URL"`
	GalileoUrl           string `mapstructure:"GALILEO_URL"`
	GalileoLogin         string `mapstructure:"GALILEO_LOGIN"`
	GalileoTranskey      string `mapstructure:"GALILEO_TRANSKEY"`
	GalileoProviderId    int32  `mapstructure:"GALILEO_PROVIDER_ID"`
	GalileoProductId     int32  `mapstructure:"GALILEO_PRODUCT_ID"`
	UserComplianceSvcUrl string `mapstructure:"USER_COMPLIANCE_SVC_URL"`
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
