package foxyclient

import (
	"github.com/spf13/viper"
	"log"
)

type FoxyConfig struct {
	ClientID     string `mapstructure:"clientid"`
	ClientSecret string `mapstructure:"clientsecret"`
	RefreshToken string `mapstructure:"refreshtoken"`
	BaseUrl      string `mapstructure:"baseurl"`
}

func readConfig() FoxyConfig {
	viper.SetEnvPrefix("FOXY") // So the env variable "FOXY_CLIENTSECRET" can be used to set the client secret
	viper.AutomaticEnv()
	viper.SetConfigFile("config.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Unable to read configuration file: %v", err)
	}

	// Reading the config manually, rather than via unmarshalling above, so the environment variable override works
	config := FoxyConfig{
		ClientID:     viper.GetString("clientid"),
		ClientSecret: viper.GetString("clientsecret"),
		RefreshToken: viper.GetString("refreshtoken"),
		BaseUrl:      viper.GetString("baseurl"),
	}

	return config
}
