package mysqlkv

import (
	"fmt"
	"time"
)

type GC struct {
	executor Executor
	ticker   *time.Ticker
	done     chan bool
}

func newGC(executor Executor, interval time.Duration, done chan bool) *GC {
	return &GC{
		executor: executor,
		ticker:   time.NewTicker(interval),
		done:     done,
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
				err := gc.executor.batchDeleteExpiredKVs(1000)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	return gc.ticker, gc.done
}
