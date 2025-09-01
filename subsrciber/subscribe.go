package subsrciber

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
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
	sub, subLogs := listen(client, opts)
	// fmt.Print("listen function called, the output is (sub): ", sub)
	// fmt.Print("listen function called, the output is (subLogs): ", subLogs)

	// 5. Process Logs
	for {
		select {
		case err := <-sub.Err():
			log.Println(err)
		case l := <-logCh:
			// fmt.Sprintln(events, l, c)
			if data := parseEvents(events, l, c); data != nil {
				log.Printf("received historical log. txn hash: %s and event data: %+v", data.TxnHash, data.Data)
				// Send the event data to the event channel
				eventCh <- data
			}
		case liveLog := <-subLogs:
			// fmt.Println("\nReceived log from subscription:", liveLog)
			if data := parseEvents(events, liveLog, c); data != nil {
				log.Printf("received live log. txn hash: %s and event data: %+v", data.TxnHash, data.Data)
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
	// defensive: ensure topics exist
	if len(log.Topics) == 0 {
		return nil
	}

	name, ok := c.events[log.Topics[0]]
	if !ok {
		return nil
	}

	// ensure requested
	found := false
	for _, e := range events {
		if e == name {
			found = true
			break
		}
	}
	if !found {
		fmt.Println("event not found in requested events")
		return nil
	}

	data, err := unpackLog(name, log.Topics, log.Data, c.ABI)
	if err != nil || data == nil {
		return nil
	}

	ev := &Event{
		Name:        name,
		BlockNumber: log.BlockNumber,
		TxnHash:     log.TxHash,
		Contract:    log.Address,
		Data:        data,
	}
	// fmt.Println("events parsing done: ", *ev)
	return ev
}

func unpackLog(eventName string, topics []common.Hash, data []byte, contractABI abi.ABI) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	ev, ok := contractABI.Events[eventName]
	if !ok {
		return nil, fmt.Errorf("event %s not found in ABI", eventName)
	}

	if len(data) > 0 {
		if err := contractABI.UnpackIntoMap(out, eventName, data); err != nil {
			return nil, err
		}
	}

	// collect indexed params from topics (topics[0] is event id)
	topicIdx := 1
	for _, input := range ev.Inputs {
		if !input.Indexed {
			continue
		}
		if topicIdx >= len(topics) {
			return nil, fmt.Errorf("missing topic for indexed arg %s", input.Name)
		}
		tb := topics[topicIdx].Bytes()
		switch input.Type.T {
		case abi.AddressTy:
			out[input.Name] = common.BytesToAddress(tb[12:]).Hex()
		case abi.UintTy, abi.IntTy:
			out[input.Name] = new(big.Int).SetBytes(tb)
		case abi.BoolTy:
			out[input.Name] = tb[len(tb)-1] == 1
		case abi.BytesTy, abi.StringTy:
			out[input.Name] = "indexed-hash:" + hex.EncodeToString(tb)
		default:
			out[input.Name] = hex.EncodeToString(tb)
		}
		topicIdx++
	}
	// fmt.Println("unpacking log done, output: ", out)
	return out, nil
}
