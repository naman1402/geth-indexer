package main

import (
	"flag"
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

	// WaitGroup is a collection of goroutines
	// provides simple way to coordinate and manage the lifecycles of multiple goroutines in a concurrent program.
	var wg sync.WaitGroup
	defer wg.Done()
	// Adds 2 goroutine to the wait group and ensures that wg.Done() is called when the function exits
	wg.Add(2)

	// Returns Config (Query, Database, API)
	options := cli.Run()
	// Reading non-flags arguments
	events := flag.Args()
	if len(events) == 0 {
		log.Println("no events provided, please specify smart contract events")
		return 1
	}

	// Create a channel to receive events from the subscriber
	// Create quitChannel to know when to terminate the program
	eventChannel := make(chan *subsrciber.Event, channelBufferSize)
	quitChannel := make(chan bool)

	go stopSignal(quitChannel)
	go subsrciber.Subscribe(events, eventChannel, options, quitChannel)

	// Connect to Postgres database using provided configuration options
	db, err := indexer.Connect(options.Database)
	if err != nil {
		log.Println(err)
		return 1
	}

	// Ensure database connection is closed when the function exits
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
