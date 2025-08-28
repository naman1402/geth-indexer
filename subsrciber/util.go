package subsrciber

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naman1402/geth-indexer/cli"
)

// major functions: filter and listen

var (
	from *big.Int
	to   *big.Int
)

// ethClient.Client defines typed wrappers for the Ethereum RPC API.
// Log represents a contract log event
func filter(client *ethclient.Client, opts *cli.Config) []types.Log {

	// if from/to field in query is not zero, assign it to from/to variables
	// or else set them as nil (latest block)
	if opts.Query.From != 0 {
		from = big.NewInt(int64(opts.Query.From))
	} else {
		from = nil
	}
	if opts.Query.To != 0 {
		to = big.NewInt(int64(opts.Query.To))
	} else {
		to = nil
	}

	// FilterQuery contains options for contract log filtering.
	// Defines the filter criteria for retrieving logs from the Ethereum blockchain.
	query := ethereum.FilterQuery{
		FromBlock: from,
		ToBlock:   to,
		Addresses: []common.Address{
			common.HexToAddress(opts.Query.Address),
		},
	}

	// FilterLogs executes a filter query.
	// executes filter query on client with current context, retrieves logs that match the query and assign them to logs
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	return logs
}

// ethereum.Subscription represents an event subscription where events are delivered on a data channel.
func listen(client *ethclient.Client, opts *cli.Config) ethereum.Subscription {
	// make a channel of type types.Log
	logs := make(chan types.Log)
	// Creates a query that sets Addresses field to a slice containing the address specified in opts, converts to common.Address
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(opts.Query.Address)},
	}

	// SubscribeFilterLogs subscribes to the results of a streaming filter query.
	// This sets up a subscription to continuously receive logs from the Ethereum blockchain based on the specified query.
	// If there's an issue with the query or the connection to the blockchain, it logs the error and stops execution.
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	// returns ethereum.Subscription
	return sub
}
