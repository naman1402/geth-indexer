package cli

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Run initializes the application configuration by reading the config.yaml file,
// unmarshaling the configuration into a Config struct, and parsing any command-line flags.
// It returns a pointer to the initialized Config struct.
func Run() *Config {

	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// println("hello world")

	// 1. tells viper where to find the configuration file and what format to expect
	//SetConfigFile explicitly defines the path, name and extension of the config file. Viper will use this and not check any of the config paths.
	viper.SetConfigFile("config.yaml")
	// AddConfigPath adds a path for Viper to search for the config file in. Can be called multiple times to define multiple search paths.
	viper.AddConfigPath(".")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("failed to read config file: %v\n", err)
	}
	// Map to struct Config
	// unmarshal config into Struct (config),
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Printf("failed to unmarshal config: %v\n", err)
	}

	config.Query = ParseFlags()
	return &config
}
