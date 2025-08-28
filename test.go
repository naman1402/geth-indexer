package main

import (
	"fmt"

	"github.com/naman1402/geth-indexer/cli"
)

func main() {
	options := cli.Run()
	fmt.Printf("Loaded configuration\nRPC Node URL (WS): %+v\nEtherscan API Key: %+v\n", options.API.EthNodeURL, options.API.EtherscanAPI)
	fmt.Printf("Database configuration: Host=%s, Port=%d, User=%s, DBName=%s\n", options.Database.DBHost, options.Database.DBPort, options.Database.DBUser, options.Database.DBName)
	fmt.Printf("Query configuration: Address=%s, From=%d, To=%d\n", options.Query.Address, options.Query.From, options.Query.To)
}
