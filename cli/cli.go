package cli

import (
	"log"

	"github.com/spf13/viper"
)

// Run initializes the application configuration by reading the config.yaml file,
// unmarshaling the configuration into a Config struct, and parsing any command-line flags.
// It returns a pointer to the initialized Config struct.
func Run() *Config {
	// println("hello world")

	//SetConfigFile explicitly defines the path, name and extension of the config file. Viper will use this and not check any of the config paths.
	viper.SetConfigFile("config.yaml")
	// AddConfigPath adds a path for Viper to search for the config file in. Can be called multiple times to define multiple search paths.
	viper.AddConfigPath(".")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("failed to read config file: %v\n", err)
	}
	// unmarshal config into Struct (config),
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Printf("failed to unmarshal config: %v\n", err)
	}

	config.Query = ParseFlags()
	return &config
}
