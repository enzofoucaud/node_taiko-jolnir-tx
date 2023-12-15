package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"taiko-tx/config"
	"taiko-tx/internal/etherscan"
	"taiko-tx/internal/taiko"
	"taiko-tx/internal/utils"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	c, err := config.Config()
	if err != nil {
		log.Fatal(err)
	}

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Hours().Do(checkNode, c) // nolint
	scheduler.StartBlocking()
}

func checkNode(c *config.Configuration) {
	for idWallet, wallet := range c.Wallets {
		_, address, err := utils.GenerateECDSAKeys(wallet.PrivateKey)
		if err != nil {
			log.Fatal(err)
		}

		txList, err := etherscan.GetAccountTxList(c.EtherScan.Url, c.EtherScan.ApiKey, address.Hex())
		if err != nil {
			log.Fatal(err)
		}

		t := taiko.Init(c)

		fmt.Println("Address", address.Hex(), "is checking...")
		for _, tx := range txList.Result {
			if tx.To == strings.ToLower(t.DepositSmartContract.Hex()) {
				if tx.TimeStamp == "" {
					fmt.Println("No timestamp, go to next tx")
					continue
				}

				timestamp, err := strconv.ParseInt(tx.TimeStamp, 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				ts := time.Unix(timestamp, 0)

				if time.Since(ts).Hours() > 3 {
					err = t.BondTTKO(idWallet)
					if err != nil {
						log.Fatal(err)
					}
					break
				} else {
					fmt.Println(ts.Format("02/01/2006 15:04:05") + " is less than 3 hours, do nothing")
					break
				}
			}
		}
	}
}
