package activity

import (
	"go-eth-activity-tracker/internal/model"
	"sort"
)

const addrCountEst = 200 * 100 // Estimated avg number of addresses in the top 100 blocks

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
		t.activities[tx.From]++
		if tx.To != "" {
			t.activities[tx.To]++
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
