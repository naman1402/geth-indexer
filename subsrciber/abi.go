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

const etherscanURL = "https://api.etherscan.io/api?module=contract&action=getabi&address=0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2&apikey=%s"

type EtherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

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
	log.Printf("raw etherscan response: %s\n", string(data))

	// var res EtherscanResponse
	// if err := json.Unmarshal(data, &res); err != nil {
	// 	log.Printf("failed to unmarshal response: %v\n", err)
	// }

	// Convert ABI JSON (from resp) to Go-ethereum ABI type
	// a, err := abi.JSON(strings.NewReader(unmarshalToMapping(data)))
	result := unmarshalToMapping(data)
	if result == "" {
		return abi.ABI{}
	}
	parsedABI, err := abi.JSON(strings.NewReader(result))
	if err != nil {
		log.Println(err)
		return abi.ABI{}
	}

	return parsedABI
}

// unmarshalToMapping unmarshals the provided data (assumed to be JSON) into a map[string]string
// and returns the value of the "result" key from the map.
func unmarshalToMapping(data []byte) string {

	// abiMap := make(map[string]string)
	var response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}
	// converts JSON response to Go map
	// Unmarshal: JSON string -> Go data
	err := json.Unmarshal(data, &response)
	if err != nil {
		log.Printf("failed to unmarshal data into map: %v\n", err)
	}

	// Returns only the "result" field containing actual ABI
	return response.Result
}
