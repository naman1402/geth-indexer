package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/naman1402/geth-indexer/cli"
	"github.com/naman1402/geth-indexer/indexer"
	"github.com/naman1402/geth-indexer/subsrciber"
)

const channelBufferSize = 1000

// main is the entry point of the application. It calls the exec function and exits the program with the returned status code.
func main() {
	os.Exit(exec())
}

func exec() int {
	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(2)

	// Returns Config (Query, Database, API) ✅
	options := cli.Run()
	fmt.Printf("Loaded configuration\nRPC Node URL (WS): %+v\nEtherscan API Key: %+v\n", options.API.EthNodeURL, options.API.EtherscanAPI)
	fmt.Printf("Database configuration: Host=%s, Port=%d, User=%s, DBName=%s\n", options.Database.DBHost, options.Database.DBPort, options.Database.DBUser, options.Database.DBName)
	fmt.Printf("Query configuration: Address=%s, From=%d, To=%d\n", options.Query.Address, options.Query.From, options.Query.To)

	// Reading non-flags arguments
	flag.Parse() // go run test.go Transfer
	events := flag.Args()
	if len(events) == 0 {
		log.Println("no events provided, please specify smart contract events")
		return 1
	}

	// Create a channel to receive events from the subscriber
	eventChannel := make(chan *subsrciber.Event, channelBufferSize)
	// Create quitChannel to know when to terminate the program
	quitChannel := make(chan bool)

	go stopSignal(quitChannel)
	go subsrciber.Subscribe(events, eventChannel, options, quitChannel)

	// Connect to Postgres database using provided configuration options ✅
	db, err := indexer.Connect(options.Database)
	if err != nil {
		log.Println(err)
		return 1
	}

	// Ensure database connection is closed when the function exits ✅
	defer func() int {
		if err := db.Close(); err != nil {
			log.Println(err)
			return 1
		}
		return 0
	}()

	// Indexes events from the eventChannel and stores them in the database
	go indexer.Index(eventChannel, db, quitChannel)

	// Wait for all goroutines to finish and then return 0 s
	wg.Wait()
	return 0
}

// stopSignal listens for user input on the console and sends a signal to the quitChannel
// when the user types "stop". This allows the main program to gracefully exit.
func stopSignal(quitChannel chan bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either OS signal or manual stop
	select {
	case <-sigChan:
		quitChannel <- true
	}

}
