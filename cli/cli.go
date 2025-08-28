package cli

import (
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// Run initializes the application configuration by reading the config.yaml file,
// unmarshaling the configuration into a Config struct, and parsing any command-line flags.
// It returns a pointer to the initialized Config struct.
func Run() *Config {

	_ = gotenv.Load()
	dbConfig := DatabaseConfig{
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvAsIntOrDefault("DB_PORT", 5432),
		DBUser:     getEnvOrDefault("DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:     getEnvOrDefault("DB_NAME", "geth_indexer"),
	}

	apiConfig := APIConfig{
		EtherscanAPI: os.Getenv("ETHERSCAN_API_KEY"),
		EthNodeURL:   os.Getenv("RPC_URL"),
	}

	queryConfig := QueryFlagOptions{
		Address: os.Getenv("CONTRACT_ADDRESS"),
		From:    getEnvAsIntOrDefault("START_BLOCK", 0),
		To:      getEnvAsIntOrDefault("END_BLOCK", 0),
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	// var config Config
	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Printf("failed to read config file: %v\n", err)
	// }
	// err := viper.Unmarshal(&config)
	// if err != nil {
	// 	log.Printf("failed to unmarshal config: %v\n", err)
	// }

	// config.Query = ParseFlags()
	return &Config{
		Query:    queryConfig,
		Database: dbConfig,
		API:      apiConfig,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
