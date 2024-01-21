package main

import (
	"go-eth-activity-tracker/internal/activity"
	"go-eth-activity-tracker/internal/ethclient"
	"log"
	"os"
)

func main() {
	accessToken := os.Getenv("ETH_ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatal("Access Token not found. Please set the ETH_ACCESS_TOKEN environment variable.")
	}

	client := ethclient.NewClient(accessToken)
	defer client.Stop() // to release ticker
	blocks, err := client.GetLatestBlocks(100)
	if err != nil {
		log.Fatal(err)
	}

	tracker := activity.NewTracker()
	for _, block := range blocks {
		tracker.UpdateActivity(block)
	}

	topAddresses := tracker.GetTopAddresses(5)
	log.Println("Top 5 active addresses for last 100 blocks:")
	for i, address := range topAddresses {
		log.Printf("%d. Address: %s, Activity: %d\n", i+1, address.Address, address.Activity)
	}
}
