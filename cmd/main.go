package main

import (
	"fmt"
	"go-eth-activity-tracker/internal/activity"
	"go-eth-activity-tracker/internal/ethclient"
	"log"
	"os"
)

const (
	blocksCount = 100
	topCount    = 5
)

func main() {
	accessToken := os.Getenv("ETH_ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatal("Access Token not found. Please set the ETH_ACCESS_TOKEN environment variable.")
	}

	client := ethclient.NewClient(accessToken)
	defer client.Stop() // to release ticker
	log.Printf("Retrieving latest %d blocks...\n", blocksCount)
	blocks, err := client.GetLatestBlocks(blocksCount)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Calculating activity metrics...")
	tracker := activity.NewTracker()
	for _, block := range blocks {
		tracker.UpdateActivity(block)
	}

	topAddresses := tracker.GetTopAddresses(topCount)
	fmt.Printf("Top %d active addresses for last %d blocks:\n", topCount, blocksCount)
	for i, address := range topAddresses {
		fmt.Printf("%d. Address: %s, Activity: %d\n", i+1, address.Address, address.Activity)
	}
}
