package ethclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-eth-activity-tracker/internal/model"
	"io"
	"math/big"
	"net/http"
)

type Client struct {
	accessToken string
	baseURL     string
}

func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		baseURL:     "https://go.getblock.io/" + accessToken + "/",
	}
}

func (c *Client) GetLatestBlocks(count int) ([]model.Block, error) {
	latestBlockNumber, err := c.getLatestBlockNumber()
	if err != nil {
		return nil, err
	}

	var blocks []model.Block
	for i := 0; i < count; i++ {
		blockNumber := new(big.Int).Sub(latestBlockNumber, big.NewInt(int64(i)))
		block, err := c.getBlockByNumber(blockNumber)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
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
