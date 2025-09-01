package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/naman1402/geth-indexer/cli"
	"github.com/naman1402/geth-indexer/indexer"
	"github.com/naman1402/geth-indexer/subsrciber"
)

func exec_test() int {
	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(1)

	options := cli.Run()
	fmt.Printf("Loaded configuration\nRPC Node URL (WS): %+v\nEtherscan API Key: %+v\n", options.API.EthNodeURL, options.API.EtherscanAPI)
	fmt.Printf("Database configuration: Host=%s, Port=%d, User=%s, DBName=%s\n", options.Database.DBHost, options.Database.DBPort, options.Database.DBUser, options.Database.DBName)
	fmt.Printf("Query configuration: Address=%s, From=%d, To=%d\n", options.Query.Address, options.Query.From, options.Query.To)

	flag.Parse()
	events := flag.Args()
	if len(events) == 0 {
		log.Println("no events provided, please specify smart contract events")
		return 1
	}
	fmt.Printf("Events to subscribe: %+v\n", events)
	//go run test.go Transfer

	const channelBufferSize = 1000
	eventChannel := make(chan *subsrciber.Event, channelBufferSize)
	quitChannel := make(chan bool)
	fmt.Printf("Created channels\neventChannel: %v\nquitChannel: %v\n", eventChannel, quitChannel)
	// go stopSignal(quitChannel)
	// go subsrciber.Subscribe(events, eventChannel, options, quitChannel)

	// Start postgres container
	// 	docker compose pull postgres
	// docker compose up -d --no-deps --no-build postgres
	_, err := indexer.Connect(options.Database)
	if err != nil {
		log.Println(err)
		return 1
	}

	etherscanAPI := options.API.EtherscanAPI
	if etherscanAPI == "" {
		log.Fatal("ETHERSCAN_API_KEY environment variable is not set")
	}

	contractAddr := options.Query.Address
	if contractAddr == "" {
		log.Fatal("CONTRACT_ADDRESS environment variable is not set")
	}
	// const etherscanURLTemplate = "https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s"
	// // url := fmt.Sprintf(etherscanURLTemplate, contractAddr, etherscanAPI)
	// const etherscanURLGetSourceCode = "https://api.etherscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	// // fmt.Printf("Calling etherscan for ABI, URL: %s\n", url)
	// // abi := subsrciber.FetchABI(options)
	// // fmt.Println(abi)
	// getSourceCodeURL := fmt.Sprintf(etherscanURLGetSourceCode, "0xB8c77482e45F1F44dE1745F52C74426C631bDD52", etherscanAPI)
	// resp, _ := http.Get(getSourceCodeURL)
	// // if err != nil {
	// // 	log.Printf("failed to fetch source code from etherscan: %v\n", err)
	// // }
	// // defer resp.Body.Close()
	// data, _ := io.ReadAll(resp.Body)
	// // parse proxy info and ABI from etherscan response
	// isProxy, implAddr, abiStr, err := parseEtherscanProxyInfo(data)
	// if err != nil {
	// 	log.Printf("failed to parse etherscan response: %v", err)
	// } else {
	// 	fmt.Printf("isProxy=%v implementation=%s\n", isProxy, implAddr)
	// 	if abiStr != "" {
	// 		parsedABI, err := abi.JSON(strings.NewReader(abiStr))
	// 		if err != nil {
	// 			log.Printf("failed to parse ABI: %v", err)
	// 		} else {
	// 			fmt.Printf("Parsed ABI: constructor=%v methods=%d events=%d\n", parsedABI.Constructor, len(parsedABI.Methods), len(parsedABI.Events))
	// 		}
	// 	}
	// }

	// _ = subsrciber.FetchABI(options)
	// fmt.Println(abi)
	go subsrciber.Subscribe(events, eventChannel, options, quitChannel)
	// go indexer.Index(eventChannel, db, quitChannel)

	wg.Wait()
	return 0
}

// ... helper removed; parsing is handled by parseEtherscanProxyInfo

// parseEtherscanProxyInfo extracts proxy flag, implementation address and ABI string from an Etherscan
// 'getsourcecode' response payload. It returns (isProxy, implementationAddress, abiString, error).
// func parseEtherscanProxyInfo(data []byte) (bool, string, string, error) {
// 	var resp struct {
// 		Status  string `json:"status"`
// 		Message string `json:"message"`
// 		Result  []struct {
// 			Proxy          string `json:"Proxy"`
// 			Implementation string `json:"Implementation"`
// 			ABI            string `json:"ABI"`
// 		} `json:"result"`
// 	}

// 	if err := json.Unmarshal(data, &resp); err != nil {
// 		return false, "", "", err
// 	}
// 	if len(resp.Result) == 0 {
// 		return false, "", "", fmt.Errorf("no result in etherscan response")
// 	}
// 	r := resp.Result[0]
// 	isProxy := false
// 	if r.Proxy == "1" || strings.EqualFold(r.Proxy, "true") {
// 		isProxy = true
// 	}
// 	return isProxy, r.Implementation, r.ABI, nil
// }

// func main() {
// 	os.Exit(exec_test())
// }
