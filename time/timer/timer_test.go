package timer_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gopherd/core/time/timer"
)

func TestMain(m *testing.M) {
	const N = 1000000
	const d = time.Second * 1
	var scheduler = timer.NewMemoryScheduler()
	scheduler.Start()
	defer scheduler.Shutdown()

	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		timer.SetTimeoutFunc(scheduler, d, func(id timer.ID) {
			fmt.Println("SetTimeoutFunc called")
			wg.Done()
		})
	}
	wg.Wait()
}
