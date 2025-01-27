package core

import (
	"sync"
)

type Info struct {
	blocks int
	client *Client
	mu     *sync.Mutex
}

func NewInfo(rpcURL string, username string, password string) *Info {
	infoClient := NewClient(rpcURL, username, password)
	var infoMu sync.Mutex
	info := Info{
		blocks: 0,
		client: infoClient,
		mu:     &infoMu,
	}
	return &info
}

func (info *Info) GetBlocks(forceUpdate bool) (int, error) {
	info.mu.Lock()
	defer info.mu.Unlock()

	if info.blocks != 0 && !forceUpdate {
		return info.blocks, nil
	}

	result, err := CallBlockchainInfo(info.client)
	if err != nil {
		return info.blocks, err
	}

	info.blocks = result.Blocks
	return info.blocks, nil
}
