package config

import "github.com/spf13/viper"

type Config struct {
	Port                  string `mapstructure:"PORT"`
	AuthSvcUrl            string `mapstructure:"AUTH_SVC_URL"`
	PoolTransactionSvcUrl string `mapstructure:"POOL_TRANSACTION_SVC_URL"`
	UserComplianceSvcUrl  string `mapstructure:"USER_COMPLIANCE_SVC_URL"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
