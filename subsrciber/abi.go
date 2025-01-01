package subsrciber

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const etherscanURL = "https://api-goerli.etherscan.io/api?module=contract&action=getabi&address=0xdBFC942264f5CebF8C59f4065af2EFfB92D12475&apikey=%s"

// fetchABI fetches the ABI (Application Binary Interface) from the Etherscan API
// using the provided etherscanAPI string. It returns the parsed ABI.
func fetchABI(etherscanAPI string) abi.ABI {

	url := fmt.Sprintf(etherscanURL, etherscanAPI)
	// HTTP GET request to Etherscan
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch ABI from etherscan: %v\n", err)
	}

	// ensures the HTTP response body is properly closed after the function completes.
	// Defer: schedules the anonymous function to be executed when the surrounding function (fetchABI) returns
	// cleanup guaranteed to run even if an error occurs
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read from http response body: %v\n", err)
	}
	// Convert ABI JSON (from resp) to Go-ethereum ABI type
	a, err := abi.JSON(strings.NewReader(unmarshalToMapping(data)))
	if err != nil {
		log.Println(err)
	}

	return a
}

// unmarshalToMapping unmarshals the provided data (assumed to be JSON) into a map[string]string
// and returns the value of the "result" key from the map.
func unmarshalToMapping(data []byte) string {

	abiMap := make(map[string]string)
	// converts JSON response to Go map
	// Unmarshal: JSON string -> Go data
	err := json.Unmarshal(data, &abiMap)
	if err != nil {
		log.Printf("failed to unmarshal data into map: %v\n", err)
	}
	// Returns only the "result" field containing actual ABI
	return abiMap["result"]
}
