package workpool

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/blueturbo-ad/go-utils/config_manage"
)

func TestWorkPool_NonblockingDropAfterCapacity(t *testing.T) {
	wp := NewAntsPool()
	cfg := []config_manage.WorkPoolConfig{{
		Name:     "stress",
		PoolSize: 1000,
	}}
	mgr := &config_manage.WorkPoolConfigManager{Config: &cfg}
	if err := wp.BuildWorkPool(mgr); err != nil {
		t.Fatalf("BuildWorkPool failed: %v", err)
	}
	pool, err := wp.GetGinCtxPool("stress")
	if err != nil {
		t.Fatalf("GetGinCtxPool failed: %v", err)
	}
	defer wp.Release()

	var accepted int32
	var rejected int32
	block := make(chan struct{})

	total := 1200
	for i := 0; i < total; i++ {
		err := pool.Submit(func() {
			<-block
		})
		if err != nil {
			atomic.AddInt32(&rejected, 1)
		} else {
			atomic.AddInt32(&accepted, 1)
		}
	}

	if accepted != 1000 {
		t.Fatalf("accepted=%d, want 1000", accepted)
	}
	if rejected != 200 {
		t.Fatalf("rejected=%d, want 200", rejected)
	}

	close(block)
	for pool.Running() > 0 {
		time.Sleep(10 * time.Millisecond)
	}
}
