package core

import (
	"encoding/json"
	"fmt"
)

type BlockchainInfo struct {
	Blocks int `json:"blocks"`
}

func CallBlockchainInfo(btcClient *Client) (BlockchainInfo, error) {
	result, err := btcClient.Call("getblockchaininfo", nil)
	if err != nil {
		return BlockchainInfo{}, fmt.Errorf("error in call blockchain info: %w", err)
	}

	var blockchainInfo BlockchainInfo

	if err := json.Unmarshal(result, &blockchainInfo); err != nil {
		return BlockchainInfo{}, fmt.Errorf("error in call blockchain info: %w", err)
	}

	return blockchainInfo, nil
}
