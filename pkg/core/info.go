package core

import (
	"fmt"
	"sync"
)

type Info struct {
	blocks int
	client *Client
	mu     *sync.RWMutex
}

func NewInfo(rpcURL string, username string, password string) *Info {
	infoClient := NewClient(rpcURL, username, password)
	var infoMu sync.RWMutex
	info := Info{
		blocks: 0,
		client: infoClient,
		mu:     &infoMu,
	}
	return &info
}

func (info *Info) GetBlocks(forceUpdate bool) (int, error) {
	if info.blocks != 0 && !forceUpdate {
		info.mu.RLock()
		defer info.mu.RUnlock()
		return info.blocks, nil
	}

	info.mu.Lock()
	defer info.mu.Unlock()

	result, err := CallBlockchainInfo(info.client)
	if err != nil {
		return 0, fmt.Errorf("error in get blocks: %w", err)
	}

	info.blocks = result.Blocks
	return info.blocks, nil
}
