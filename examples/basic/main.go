package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	ponstrade "github.com/0xfnzero/pons-trade-sdk/ponstrade"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	rpcURL := os.Getenv("ROBINHOOD_RPC_URL")
	if rpcURL == "" {
		fmt.Println("set ROBINHOOD_RPC_URL to run the example")
		return
	}

	ctx := context.Background()
	client, eth, err := ponstrade.Dial(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer eth.Close()

	curveValue := os.Getenv("PONS_CURVE")
	if !common.IsHexAddress(curveValue) {
		fmt.Println("set PONS_CURVE to a valid non-zero address")
		return
	}
	curve := common.HexToAddress(curveValue)
	if curve == (common.Address{}) {
		fmt.Println("set PONS_CURVE to a valid non-zero address")
		return
	}
	recipientValue := os.Getenv("PONS_RECIPIENT")
	if !common.IsHexAddress(recipientValue) {
		fmt.Println("set PONS_RECIPIENT to a valid non-zero address")
		return
	}
	recipient := common.HexToAddress(recipientValue)
	if recipient == (common.Address{}) {
		fmt.Println("set PONS_RECIPIENT to a valid non-zero address")
		return
	}
	quoteIn := new(big.Int).SetUint64(1_000_000_000_000_000)

	quote, err := client.QuoteBuy(ctx, curve, quoteIn, recipient, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("tokens out:", quote.TokensOut)
	fmt.Println("spent:", quote.Spent)
	fmt.Println("refund:", quote.Refund)
}
