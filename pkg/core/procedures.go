package core

import "encoding/json"

type BlockchainInfo struct {
	Blocks int `json:"blocks"`
}

func CallBlockchainInfo(btcClient *Client) (BlockchainInfo, error) {
	result, err := btcClient.Call("getblockchaininfo", nil)
	if err != nil {
		return BlockchainInfo{}, err
	}

	var blockchainInfo BlockchainInfo

	if err := json.Unmarshal(result, &blockchainInfo); err != nil {
		return BlockchainInfo{}, err
	}

	return blockchainInfo, nil
}
