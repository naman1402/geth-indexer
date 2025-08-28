package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/naman1402/geth-indexer/cli"
	"github.com/naman1402/geth-indexer/indexer"
	"github.com/naman1402/geth-indexer/subsrciber"
)

func exec_test() int {
	options := cli.Run()
	fmt.Printf("Loaded configuration\nRPC Node URL (WS): %+v\nEtherscan API Key: %+v\n", options.API.EthNodeURL, options.API.EtherscanAPI)
	fmt.Printf("Database configuration: Host=%s, Port=%d, User=%s, DBName=%s\n", options.Database.DBHost, options.Database.DBPort, options.Database.DBUser, options.Database.DBName)
	fmt.Printf("Query configuration: Address=%s, From=%d, To=%d\n", options.Query.Address, options.Query.From, options.Query.To)

	flag.Parse()
	events := flag.Args()
	if len(events) == 0 {
		log.Println("no events provided, please specify smart contract events")
		return 1
	}
	fmt.Printf("Events to subscribe: %+v\n", events)
	//go run test.go Transfer

	const channelBufferSize = 1000
	eventChannel := make(chan *subsrciber.Event, channelBufferSize)
	quitChannel := make(chan bool)
	fmt.Printf("Created channels\neventChannel: %v\nquitChannel: %v\n", eventChannel, quitChannel)
	// go stopSignal(quitChannel)
	// go subsrciber.Subscribe(events, eventChannel, options, quitChannel)

	// Start postgres container
	// 	docker compose pull postgres
	// docker compose up -d --no-deps --no-build postgres
	_, err := indexer.Connect(options.Database)
	if err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(exec_test())
}
