package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"context"
	"crypto/ecdsa"

	"github.com/mazen160/go-random"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/martinlindhe/notify"
)

var letters = []rune("abcdef0123456789")

func main() {
	charset := "abcdef0123456789"
	length := 64

	// create file
	f, err := os.Create("privateKeys.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {
			data, err := random.Random(length, charset, true)
			if err != nil {
				fmt.Println(err)
				continue
			}

			privateKey, err := crypto.HexToECDSA(data)
			if err != nil {
				fmt.Println(err)
				continue
			}

			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				fmt.Println(err)
				continue
			}

			fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
			balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if balance.Cmp(big.NewInt(0)) != 0 {
				fmt.Println("Account Private Key: ", data)
				fmt.Println("Account Balance: ", balance)
				_, err := fmt.Fprintln(f, "*", data, "*")
				notify.Notify(fromAddress.String(), data, balance.String(), "")
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Print(".")
			}
		}
	}()
	var input string
	fmt.Scanln(&input)
}
