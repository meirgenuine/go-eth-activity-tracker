package ethclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-eth-activity-tracker/internal/model"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"
)

const (
	maxRetries   = 3
	retryTimeout = 2 * time.Second
)

type Client struct {
	accessToken string
	baseURL     string
	rateLimiter *time.Ticker
}

func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		baseURL:     "https://go.getblock.io/" + accessToken + "/",
		rateLimiter: time.NewTicker(time.Second / 60), // 60 requests per second to not exceed getblock.io account limits
	}
}

func (c *Client) Stop() {
	c.rateLimiter.Stop()
}

func (c *Client) GetLatestBlocks(count int) ([]model.Block, error) {
	latestBlockNumber, err := c.getLatestBlockNumber()
	if err != nil {
		return nil, err
	}

	blocks := make([]model.Block, count)
	errChan := make(chan error, count)
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		i := i
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			retryCount := 0
			for {
				blockNumber := new(big.Int).Sub(latestBlockNumber, big.NewInt(int64(i)))
				block, err := c.getBlockByNumber(blockNumber)
				if err != nil {
					if retryCount < maxRetries {
						retryCount++
						time.Sleep(retryTimeout)
						continue
					}
					errChan <- err
					return
				}
				blocks[i] = block
				break
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return blocks, nil
}

func (c *Client) getLatestBlockNumber() (*big.Int, error) {
	response, err := c.sendRequest("eth_blockNumber", nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	blockNumberHex, ok := result["result"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid response format for block number")
	}

	blockNumber := new(big.Int)
	blockNumber.SetString(blockNumberHex[2:], 16) // Convert hex to big.Int
	return blockNumber, nil
}

func (c *Client) getBlockByNumber(blockNumber *big.Int) (model.Block, error) {
	hexBlockNumber := fmt.Sprintf("0x%x", blockNumber)
	response, err := c.sendRequest("eth_getBlockByNumber", []interface{}{hexBlockNumber, true})
	if err != nil {
		return model.Block{}, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return model.Block{}, err
	}

	blockData, ok := result["result"].(map[string]interface{})
	if !ok {
		return model.Block{}, fmt.Errorf("invalid data for block number %s", hexBlockNumber)
	}

	var block model.Block
	blockJSON, err := json.Marshal(blockData)
	if err != nil {
		return model.Block{}, err
	}
	if err := json.Unmarshal(blockJSON, &block); err != nil {
		return model.Block{}, err
	}

	return block, nil
}

func (c *Client) sendRequest(method string, params []interface{}) ([]byte, error) {
	<-c.rateLimiter.C // Wait for the next tick before proceeding

	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      "getblock.io",
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.baseURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
