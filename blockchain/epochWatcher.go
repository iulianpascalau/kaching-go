package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type blockInfo struct {
	Hash  string `json:"hash"`
	Epoch int    `json:"epoch"`
	Shard int    `json:"shard"`
	Round int    `json:"round"`
}

//https://testnet-api.elrond.com/blocks?fields=epoch,shard
const endpoint = "/blocks?fields=epoch,shard,round"
const epochNotRead = -1
const httpTimeout = time.Second * 5
const metachainShardId = 0xFFFFFFFF

type epochWatcher struct {
	address         string
	poolingInterval time.Duration
	currentEpoch    int
	cancelFunc      func()
	chPlaySound     chan struct{}
	httpClient      *http.Client
}

func NewEpochWatcher(
	address string,
	poolingInterval time.Duration,
	chPlaySound chan struct{},
) *epochWatcher {
	ew := &epochWatcher{
		address:         address,
		poolingInterval: poolingInterval,
		currentEpoch:    epochNotRead,
		chPlaySound:     chPlaySound,
		httpClient:      http.DefaultClient,
	}

	ctx, cancel := context.WithCancel(context.Background())
	ew.cancelFunc = cancel
	ew.httpClient.Timeout = httpTimeout

	go ew.poll(ctx)

	return ew
}

func (ew *epochWatcher) poll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(ew.poolingInterval):
			ew.checkEpoch()
		}
	}
}

func (ew *epochWatcher) checkEpoch() {
	bi, err := ew.getBlockchainEpoch()
	if err != nil {
		fmt.Println(err)
		return
	}

	oldEpoch := ew.currentEpoch
	ew.currentEpoch = bi.Epoch
	if oldEpoch == epochNotRead {
		fmt.Printf("epoch initialized to %d in round %d\n", bi.Epoch, bi.Round)
		return
	}

	if oldEpoch < bi.Epoch {
		fmt.Printf("new epoch %d in round %d\n", bi.Epoch, bi.Round)
		ew.chPlaySound <- struct{}{}
	} else {
		fmt.Printf("read epoch %d and round %d\n", bi.Epoch, bi.Round)
	}
}

func (ew *epochWatcher) getBlockchainEpoch() (*blockInfo, error) {
	blocks, err := ew.getBlocksInfo()
	if err != nil {
		return nil, err
	}

	var highestMetablock *blockInfo
	var highestShardBlock *blockInfo
	for _, b := range blocks {
		if b.Shard == metachainShardId {
			if highestMetablock == nil {
				highestMetablock = b
			}

			if highestMetablock.Round < b.Round {
				highestMetablock = b
			}

			continue
		}

		if highestShardBlock == nil {
			highestShardBlock = b
		}

		if highestShardBlock.Round < b.Round {
			highestShardBlock = b
		}
	}

	if highestMetablock != nil {
		return highestMetablock, nil
	}

	return highestShardBlock, nil
}

func (ew *epochWatcher) getBlocksInfo() ([]*blockInfo, error) {
	req, err := http.NewRequest("GET", ew.address+endpoint, nil)
	if err != nil {
		return nil, err
	}

	userAgent := "Epoch-Kaching / 1.0.0 <Requesting data from gateway>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := ew.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	blocks := make([]*blockInfo, 0)
	err = json.NewDecoder(resp.Body).Decode(&blocks)
	if err != nil {
		return nil, err
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("empty response from gateway")
	}

	return blocks, nil
}

func (ew *epochWatcher) Close() {
	ew.cancelFunc()
}
