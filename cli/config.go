package cli

import "flag"

// Config holds the configuration options for the application.
type Config struct {
	// Query holds the options for querying the smart contract.
	Query QueryFlagOptions
	// Database holds the configuration for the database connection.
	Database DatabaseConfig
	// API holds the configuration for the API endpoints.
	API APIConfig
}

// QueryFlagOptions holds the options for querying the smart contract.
type QueryFlagOptions struct {
	// Address is the address of the smart contract.
	Address string
	// From is the starting block number for the query.
	From int
	// To is the ending block number for the query.
	To int
}

// DatabaseConfig holds the configuration for the database connection.
type DatabaseConfig struct {
	// DBHost is the hostname or IP address of the database server.
	DBHost string
	// DBPort is the port number of the database server.
	DBPort int
	// DBUser is the username for the database connection.
	DBUser string
	// DBPassword is the password for the database connection.
	DBPassword string
	// DBName is the name of the database.
	DBName string
}

// APIConfig holds the configuration for the API endpoints.
type APIConfig struct {
	// EtherscanAPI is the API key for the Etherscan API.
	EtherscanAPI string
	// EthNodeURL is the URL of the Ethereum node.
	EthNodeURL string
}

// ParseFlags parses the command-line flags and returns a QueryFlagOptions struct
// containing the parsed values.
func ParseFlags() QueryFlagOptions {

	// func flag.String(name string, value string, usage string) *string
	address := flag.String("address", "", "Address of smart contract")
	from := flag.Int("from", 0, "Block range, default value is genesis block")
	to := flag.Int("to", 0, "Block range, default options makes application listen for future events")
	// Parse the command-line flags
	flag.Parse()

	return QueryFlagOptions{
		Address: *address,
		From:    *from,
		To:      *to,
	}
}
