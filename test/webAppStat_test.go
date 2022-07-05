package test

import (
	"fmt"
	"sync"
	"testing"

	is "github.com/lky492465798/interfaceUtils/support"
)

func TestRequestIntervalHistogramRecord(t *testing.T) {
	w := new(is.WebAppStat)

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 1; i <= 1000; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				w.SetRequestTimeNano(20)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("结果: ", w.GetValue())
}
