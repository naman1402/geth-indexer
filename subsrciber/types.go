package subsrciber

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

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

type EtherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}
