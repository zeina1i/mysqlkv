package mysqlkv

import (
	"fmt"
	"time"
)

type GC struct {
	store  Store
	ticker *time.Ticker
	done   chan bool
}

func newGC(store Store, interval time.Duration, done chan bool) *GC {
	return &GC{
		store:  store,
		ticker: time.NewTicker(interval),
		done:   done,
	}
}

func (gc *GC) collect() (*time.Ticker, chan bool) {
	go func() {
		for {
			select {
			case <-gc.done:
				return
			case t := <-gc.ticker.C:
				fmt.Println(t)
				err := gc.store.batchDeleteExpiredKVs(1000)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	return gc.ticker, gc.done
}
