package chain

import (
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/log15"
	"sync"
)

type NeedSnapshotCache struct {
	chain     *Chain
	subLedger map[types.Address][]*ledger.AccountBlock

	lock sync.Mutex
	log  log15.Logger
}

func NewNeedSnapshotContent(chain *Chain, unconfirmedSubLedger map[types.Address][]*ledger.AccountBlock) *NeedSnapshotCache {
	cache := &NeedSnapshotCache{
		chain:     chain,
		subLedger: unconfirmedSubLedger,
		log:       log15.New("module", "chain/NeedSnapshotCache"),
	}

	return cache
}

func (cache *NeedSnapshotCache) Get(addr *types.Address) []*ledger.AccountBlock {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	return cache.subLedger[*addr]
}

func (cache *NeedSnapshotCache) Add(addr *types.Address, accountBlock *ledger.AccountBlock) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cachedChain := cache.subLedger[*addr]; cachedChain != nil && cachedChain[len(cachedChain)-1].Height >= accountBlock.Height {
		cache.log.Error("Can't add", "method", "Add")
		return
	}

	cache.subLedger[*addr] = append(cache.subLedger[*addr], accountBlock)
}

func (cache *NeedSnapshotCache) Remove(addr *types.Address, height uint64) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cachedChain := cache.subLedger[*addr]
	if cachedChain == nil {
		cache.log.Error("Can't remove", "method", "Remove")
		return
	}

	var deletedIndex = -1
	for index, item := range cachedChain {
		if item.Height > height {
			break
		}
		deletedIndex = index
	}
	cachedChain = cachedChain[deletedIndex+1:]
	cache.subLedger[*addr] = cachedChain
}
