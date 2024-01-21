package activity

import (
	"encoding/hex"
	"fmt"
	"go-eth-activity-tracker/internal/model"
	"sort"
	"strings"
)

const (
	addrCountEst          = 200 * 100    // Estimated avg number of addresses in the top 100 blocks
	erc20TransferMethodID = "0xa9059cbb" // Method ID for ERC20 transfers
)

type AddressActivity struct {
	Address  string
	Activity int
}

type Tracker struct {
	activities map[string]int // address -> activity score (number of transactions)
}

func NewTracker() *Tracker {
	return &Tracker{
		activities: make(map[string]int, addrCountEst),
	}
}

func (t *Tracker) UpdateActivity(block model.Block) {
	for _, tx := range block.Transactions {
		if len(tx.Input) >= 138 && isERC20Transfer(tx.Input) {
			recipient, err := decodeERC20Recipient(tx.Input)
			if err != nil {
				continue
			}
			t.activities[tx.From]++
			t.activities[recipient]++
		}
	}
}

func (t *Tracker) GetTopAddresses(count int) []AddressActivity {
	var addrActivities []AddressActivity
	for addr, activity := range t.activities {
		addrActivities = append(addrActivities, AddressActivity{Address: addr, Activity: activity})
	}

	sort.Slice(addrActivities, func(i, j int) bool {
		return addrActivities[i].Activity > addrActivities[j].Activity
	})

	if count > len(addrActivities) {
		count = len(addrActivities)
	}

	return addrActivities[:count]
}

func isERC20Transfer(inputData string) bool {
	if len(inputData) < len(erc20TransferMethodID) {
		return false
	}
	return strings.HasPrefix(inputData, erc20TransferMethodID)
}

func decodeERC20Recipient(inputData string) (string, error) {
	if len(inputData) < 138 {
		return "", fmt.Errorf("input data too short to be an ERC20 transfer")
	}

	// Extracting the address part of the input data
	addrData := inputData[2+8+24 : 2+8+24+40] // Skip "0x", method ID, and first 12 bytes of address padding
	recipientBytes, err := hex.DecodeString(addrData)
	if err != nil {
		return "", fmt.Errorf("failed to decode recipient address: %v", err)
	}

	recipientAddress := "0x" + hex.EncodeToString(recipientBytes)
	if len(recipientAddress) != 42 {
		return "", fmt.Errorf("decoded address is invalid: %s", recipientAddress)
	}

	return recipientAddress, nil
}
