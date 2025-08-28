package subsrciber

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const etherscanURLTemplate = "https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s"

type EtherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// fetchABI fetches the ABI (Application Binary Interface) from the Etherscan API
// using the provided etherscanAPI string. It returns the parsed ABI.
func fetchABI() abi.ABI {
	etherscanAPI := os.Getenv("ETHERSCAN_API_KEY")
	if etherscanAPI == "" {
		log.Fatal("ETHERSCAN_API_KEY environment variable is not set")
	}

	contractAddr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddr == "" {
		log.Fatal("CONTRACT_ADDRESS environment variable is not set")
	}

	url := fmt.Sprintf(etherscanURLTemplate, contractAddr, etherscanAPI)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch ABI from etherscan: %v\n", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read from http response body: %v\n", err)
	}

	// Remove verbose logging
	result := unmarshalToMapping(data)
	if result == "" {
		return abi.ABI{}
	}

	parsedABI, err := abi.JSON(strings.NewReader(result))
	if err != nil {
		log.Println(err)
		return abi.ABI{}
	}

	var response EtherscanResponse
	json.Unmarshal(data, &response)

	log.Printf("\n"+
		"┌─────────────── ABI Fetched ───────────────┐\n"+
		"│ Status: %s\n"+
		"│ Message: %s\n"+
		"│ Events Found: %d\n"+
		"└──────────────────────────────────────────┘\n",
		response.Status,
		response.Message,
		len(parsedABI.Events))

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
