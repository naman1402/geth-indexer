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

// Contract represents an Ethereum contract with its address, ABI, and a mapping of event IDs to event names.
type Contract struct {
	Address common.Address
	ABI     abi.ABI
	events  map[common.Hash]string
}

// Event represents an Ethereum event with its name, block number, block hash, contract address, and event data.
type Event struct {
	Name        string
	BlockNumber uint64
	BlockHash   common.Hash
	Contract    common.Address
	Data        map[string]interface{}
}

func Subscribe(events []string, eventCh chan<- *Event, opts *cli.Config, quit chan bool) {

	// 1. Connect to the Ethereum node
	// Dial the Ethereum node using the provided URL
	client, err := ethclient.Dial(opts.API.EthNodeURL)
	if err != nil {
		log.Fatal(err)
	}
	// Close the client connection when the function returns
	defer client.Close()
	fmt.Printf("Subscribing to these events on contract %s ... %s\n", opts.Query.Address, strings.Join(events, " "))

	// 2. Initialize a Contract struct with the provided address and ABI
	c := &Contract{
		Address: common.HexToAddress(opts.Query.Address),
		ABI:     fetchABI(opts.API.EtherscanAPI),
		events:  make(map[common.Hash]string),
	}
	// Polpulates the event map in Contract with event ID and their corresponding names from the ABI
	for _, e := range c.ABI.Events {
		c.events[e.ID] = e.Name
	}

	// create channel of type types.Log and buffer size of buffer
	logCh := make(chan types.Log, buffer)

	// 3. Filter Historical Logs
	// starts goroutine that filters historical logs and sends them to logCh
	// Ensures that historical logs are processed and sent to the log channel
	go func() {
		for _, l := range filter(client, opts) {
			logCh <- l
		}
	}()

	// 4. Subscribe to Real-Time Logs
	// Sets up a subscription to real-time logs from the Ethereum blockchain
	sub := listen(client, opts)

	// 5. Process Logs
	for {
		select {
		// errors occur in the subscription, it logs the error
		case err := <-sub.Err():
			log.Println(err)
			// l log is received from the logCh channel, parses the log and if the event is not nil, sends the event to the eventCh channel
		case l := <-logCh:
			if data := parseEvents(events, l, c); data != nil {
				eventCh <- data
			}
			// for stop signal, exit the function
		case stop := <-quit:
			if stop {
				return
			}
		}
	}

}

func parseEvents(events []string, log types.Log, c *Contract) *Event {
	// retrieves the event's name associated with the first topic in the log's topics
	// if event name is not found, return nil
	name, ok := c.events[log.Topics[0]]
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
