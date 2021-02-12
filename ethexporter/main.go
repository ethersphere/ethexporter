package main

import (
	"bufio"
	"context"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	allWatching []*Watching
	data        string
	port        string
	updates     string
	prefix      string
	interval    int
	loadSeconds float64
	totalLoaded int64
	eth         *ethclient.Client
)

type Watching struct {
	Name    string
	Address string
	Balance string
}

//
// Connect to geth server
func ConnectionToGeth(url string) error {
	var err error
	eth, err = ethclient.Dial(url)
	return err
}

//
// Fetch ETH balance from Geth server
func GetEthBalance(address string) *big.Float {
	balance, err := eth.BalanceAt(context.TODO(), common.HexToAddress(address), nil)
	if err != nil {
		fmt.Printf("Error fetching ETH Balance for address: %v\n", address)
	}
	return ToEther(balance)
}

//
// Fetch ETH balance from Geth server
func CurrentBlock() uint64 {
	block, err := eth.BlockByNumber(context.TODO(), nil)
	if err != nil {
		fmt.Printf("Error fetching current block height: %v\n", err)
		return 0
	}
	return block.NumberU64()
}

//
// CONVERTS WEI TO ETH
func ToEther(o *big.Int) *big.Float {
	pul, int := big.NewFloat(0), big.NewFloat(0)
	int.SetInt(o)
	pul.Mul(big.NewFloat(0.000000000000000001), int)
	return pul
}

//
// HTTP response handler for /metrics
func MetricsHttp(w http.ResponseWriter, r *http.Request) {
	var allOut []string
	total := big.NewFloat(0)
	for _, v := range allWatching {
		if v.Balance == "" {
			v.Balance = "0"
		}
		bal := big.NewFloat(0)
		bal.SetString(v.Balance)
		total.Add(total, bal)
		allOut = append(allOut, fmt.Sprintf("%veth_balance{name=\"%v\",address=\"%v\"} %v", prefix, v.Name, v.Address, v.Balance))
	}
	allOut = append(allOut, fmt.Sprintf("%veth_balance_total %0.18f", prefix, total))
	allOut = append(allOut, fmt.Sprintf("%veth_load_seconds %0.2f", prefix, loadSeconds))
	allOut = append(allOut, fmt.Sprintf("%veth_loaded_addresses %v", prefix, totalLoaded))
	allOut = append(allOut, fmt.Sprintf("%veth_total_addresses %v", prefix, len(allWatching)))
	fmt.Fprintln(w, strings.Join(allOut, "\n"))
}

//
// Open the addresses.txt file (name:address)
func OpenAddresses(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		object := strings.Split(scanner.Text(), ":")
		if common.IsHexAddress(object[1]) {
			w := &Watching{
				Name:    object[0],
				Address: object[1],
			}
			allWatching = append(allWatching, w)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return err
}

func main() {
	gethUrl := os.Getenv("GETH")
	port = os.Getenv("PORT")
	if port == "" {
		port = "9890"
	}
	data = os.Getenv("DATA")
	if data == "" {
		data = "."
	}
	interval, _ = strconv.Atoi(os.Getenv("INTERVAL"))
	if interval == 0 {
		interval = 15
	}
	prefix = os.Getenv("PREFIX")

	err := OpenAddresses(data + "/addresses.txt")
	if err != nil {
		panic(err)
	}

	err = ConnectionToGeth(gethUrl)
	if err != nil {
		panic(err)
	}

	// check address balances
	go func() {
		for {
			totalLoaded = 0
			t1 := time.Now()
			fmt.Printf("Checking %v wallets...\n", len(allWatching))
			for _, v := range allWatching {
				v.Balance = GetEthBalance(v.Address).String()
				totalLoaded++
			}
			t2 := time.Now()
			loadSeconds = t2.Sub(t1).Seconds()
			fmt.Printf("Finished checking %v wallets in %0.0f seconds, sleeping for %v seconds.\n", len(allWatching), loadSeconds, 15)
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()

	block := CurrentBlock()

	fmt.Printf("ETHexporter has started on port %v using Geth server: %v at block #%v\n", port, gethUrl, block)
	http.HandleFunc("/metrics", MetricsHttp)
	panic(http.ListenAndServe("0.0.0.0:"+port, nil))
}
