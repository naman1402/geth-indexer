package subsrciber

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/naman1402/geth-indexer/cli"
)

const etherscanURLTemplate = "https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s"

// fetchABI fetches the ABI (Application Binary Interface) from the Etherscan API
// using the provided etherscanAPI string. It returns the parsed ABI. ✅
func fetchABI(opts *cli.Config) abi.ABI {
	etherscanAPI := opts.API.EtherscanAPI
	if etherscanAPI == "" {
		log.Fatal("ETHERSCAN_API_KEY environment variable is not set")
	}

	contractAddr := opts.Query.Address
	if contractAddr == "" {
		log.Fatal("CONTRACT_ADDRESS environment variable is not set")
	}

	proxyResult, ActualImplementationAddress, _ := getProxyInfoAndImplementation(contractAddr, etherscanAPI)

	if proxyResult {
		fmt.Printf("Address: %s is a proxy contract, using implementation address: %s to get the ABI\n ", contractAddr, ActualImplementationAddress)
		contractAddr = ActualImplementationAddress
	} else {
		fmt.Printf("Address: %s is not a proxy contract, using it to get the ABI\n", contractAddr)
	}

	url := fmt.Sprintf(etherscanURLTemplate, contractAddr, etherscanAPI)
	fmt.Printf("Calling etherscan for ABI, URL: %s\n", url)
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

	log.Printf("ABI fetched: status=%s message=%s events=%d\n",
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

// func fetchByteCode(opts *cli.Config) string {
// 	return ""
// }

func getProxyInfoAndImplementation(contractAddress, etherScanAPI string) (bool, string, error) {
	const etherscanURLGetSourceCode = "https://api.etherscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	getSourceCodeURL := fmt.Sprintf(etherscanURLGetSourceCode, contractAddress, etherScanAPI)
	resp, _ := http.Get(getSourceCodeURL)
	data, _ := io.ReadAll(resp.Body)

	var responseStruct struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  []struct {
			Proxy          string `json:"Proxy"`
			Implementation string `json:"Implementation"`
			ABI            string `json:"ABI"`
		} `json:"result"`
	}

	if err := json.Unmarshal(data, &responseStruct); err != nil {
		return false, "", err
	}
	if len(responseStruct.Result) == 0 {
		return false, "", fmt.Errorf("no result in etherscan response")
	}
	r := responseStruct.Result[0]
	isProxy := false
	if r.Proxy == "1" || strings.EqualFold(r.Proxy, "true") {
		isProxy = true
	}
	return isProxy, r.Implementation, nil
}
