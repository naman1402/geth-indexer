package subsrciber

import (
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naman1402/geth-indexer/cli"
)

const buffer = 100

func Subscribe(events []string, eventCh chan<- *Event, opts *cli.Config, quit chan bool) {

	fmt.Println("\nSubscribing to events...")
	fmt.Printf("\nContract Address: %s\nBlock range: %d to %d\nEvents: %s\n", opts.Query.Address, opts.Query.From, opts.Query.To, strings.Join(events, ", "))

	// 1. Connecting to EVM using RPC URL
	client, err := ethclient.Dial(opts.API.EthNodeURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// fmt.Printf("Subscribing to these events on contract %s ... %s\n", opts.Query.Address, strings.Join(events, " "))
	fmt.Println("\nConnected to RPC URL:", opts.API.EthNodeURL)

	// 2. Initialize a Contract struct with the provided address and ABI ✅
	c := &Contract{
		Address: common.HexToAddress(opts.Query.Address),
		ABI:     fetchABI(opts),
		// Initially this will be an empty mapping, populated using ABI events
		events: make(map[common.Hash]string),
	}

	for _, e := range c.ABI.Events {
		c.events[e.ID] = e.Name
	}
	// fmt.Printf("Contract Events Mapping: %+v\n", c.events)
	// fmt.Printf("Contract ABI fetched: %+v\n", c.ABI)✅
	logCh := make(chan types.Log, buffer)

	////////////////////////////////////////////////////////////////////////////
	// 3. Filter Historical Logs ✅ ////////////////////////////////////////////
	// starts goroutine that filters historical logs and sends them to logCh ///
	// Ensures that historical logs are processed and sent to the log channel //
	////////////////////////////////////////////////////////////////////////////
	var topicList []common.Hash
	for h := range c.events {
		topicList = append(topicList, h)
	}
	var topics [][]common.Hash
	if len(topicList) > 0 {
		topics = append(topics, topicList)
	}
	go func() {
		for _, l := range filter(client, opts, topics) {
			logCh <- l
		}
	}()

	///////////////////////////////////////////////////////////////////////////
	// 4. Subscribe to Real-Time Logs /////////////////////////////////////////
	// Sets up a subscription to real-time logs from the Ethereum blockchain //
	///////////////////////////////////////////////////////////////////////////
	sub := listen(client, opts)

	// 5. Process Logs
	for {
		select {
		case err := <-sub.Err():
			log.Println(err)
		case l := <-logCh:
			if data := parseEvents(events, l, c); data != nil {
				// Pretty print event data
				log.Printf("\n"+
					"╔═══════════════ New Event ═══════════════\n"+
					"║ Type: %s\n"+
					"║ Block: %d\n"+
					"║ Contract: %s\n"+
					"║ Data: %+v\n"+
					"╚══════════════════════════════════════════\n",
					data.Name,
					data.BlockNumber,
					data.Contract.Hex(),
					data.Data)
				eventCh <- data
			}
		case stop := <-quit:
			if stop {
				return
			}
		}
	}

}

func parseEvents(events []string, log types.Log, c *Contract) *Event {
	// fmt.Println("\nparseEvents called")
	name, ok := c.events[log.Topics[0]]
	// fmt.Printf("name from topics: %s, error: %s", name, ok)
	if !ok {
		return nil
	}
	event := ""

	// iterates over the events provided the param
	// if the event name (from events param) matches the name retrieved from the contract's events map, then assign the name to event and break the loop
	for _, e := range events {
		if e == name {
			event = name
			break
		}
	}

	// if no event matches
	if event == "" {
		return nil
	}
	fmt.Println("\nEvent matched:", event)
	// decoded the log data using event name, abi
	data, err := unpackLog(event, log.Data, c.ABI)
	if err != nil || data == nil {
		return nil
	}

	// create new object and populate it with event name, block number, block hash, contract address and unpacked data
	e := &Event{
		Name:        event,
		BlockNumber: log.BlockNumber,
		BlockHash:   log.BlockHash,
		Contract:    log.Address,
		Data:        data,
	}
	fmt.Print("\nEvent parsed: ", e)
	return e
}

// unpackLog unpacks the event data from the provided log data and ABI, returning a map of the event parameters.
// If the event cannot be unpacked, an error is returned.
func unpackLog(event string, data []byte, abi abi.ABI) (map[string]interface{}, error) {
	// mapping of string keys to interface{} values
	logMap := make(map[string]interface{})

	// UnpackIntoMap unpacks a log into the provided map[string]interface{}.
	// decodes the log data into a structured format based on the ABI
	err := abi.UnpackIntoMap(logMap, event, data)
	if err != nil {
		return nil, err
	}

	return logMap, nil
}
