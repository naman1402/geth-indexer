package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/naman1402/geth-indexer/cli"
	"github.com/naman1402/geth-indexer/indexer"
	"github.com/naman1402/geth-indexer/subsrciber"
)

const channelBufferSize = 1000

func main() {
	os.Exit(exec())
}

func exec() int {

	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(2)

	options := cli.Run()
	events := flag.Args()
	if len(events) == 0 {
		log.Println("no events provided, please specify smart contract events")
		return 1
	}

	eventChannel := make(chan *subsrciber.Event, channelBufferSize)
	quitChannel := make(chan bool)

	go stopSignal(quitChannel)
	go subsrciber.Subscribe(events, eventChannel, options, quitChannel)

	db, err := indexer.Connect(options.Database)
	if err != nil {
		log.Println(err)
		return 1
	}

	defer func() int {
		if err := db.Close(); err != nil {
			log.Println(err)
			return 1
		}
		return 0
	}()

	go indexer.Index(eventChannel, db, quitChannel)

	wg.Wait()
	return 0
}

func stopSignal(quitChannel chan bool) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("input stop to exit from application")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error reading input, try again or press ctrl+c to exit")
	}

	if strings.TrimSpace(text) == "stop" {
		quitChannel <- true
	}
}
