package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestInn(t *testing.T) {
	//go test -v -run TestInn
	t.Parallel()
	t.Run("inter1", aTestInter2)
	t.Run("infoFunc", aTestInterGet)
	t.Run("inter2", aTestInter1)
	t.Run("inter3", aTestInter3)
	t.Run("inter4", aTestInter4)
	t.Run("inter5", aTestInter5)
	// t.Run("inter5", gettest)

}

func sorttest(t *testing.T) {
	rand.Seed(time.Now().Unix())
	arr := make([]float64, 1000000)
	for i := 0; i < 1000000; i++ {
		arr[i] = float64(rand.Int())
	}
	s := time.Now()
	sort.Float64s(arr)
	end := time.Now().Sub(s)
	fmt.Println("时间:", end)
}

func aTestInterGet(*testing.T) {
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 1; i <= 20; i++ {
		go func() {
			time.Sleep(time.Millisecond * 30)
			for j := 0; j < 5; j++ {
				_, err := http.Get("http://127.0.0.1:9000/info")
				if err != nil {
					fmt.Println("发生错误:", err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func aTestInter1(*testing.T) {
	var wg sync.WaitGroup
	fail := 0
	wg.Add(100)
	for i := 1; i <= 100; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				_, err := http.Get("http://127.0.0.1:9000/index")
				if err != nil {
					fail++
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("/index", "失败:", fail)
}

func aTestInter2(*testing.T) {
	fail := 0
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 1; i <= 100; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				_, err := http.Get("http://127.0.0.1:9000/index/a")
				if err != nil {
					fail++
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("/index/a", "失败:", fail)
}

func aTestInter3(*testing.T) {
	var wg sync.WaitGroup
	fail := 0
	wg.Add(100)
	for i := 1; i <= 100; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				_, err := http.Get("http://127.0.0.1:9000/params/a/33")
				if err != nil {
					fmt.Println(err)
					fail++
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("/params/a/33", "失败:", fail)
}

func aTestInter4(*testing.T) {
	var wg sync.WaitGroup
	fail := 0
	wg.Add(100)
	for i := 1; i <= 100; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				_, err := http.Get("http://127.0.0.1:9000/params/a/22")
				if err != nil {
					fail++
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("/params/a/22", "失败:", fail)
}

func aTestInter5(*testing.T) {
	fail := 0
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 1; i <= 100; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				_, err := http.Get("http://127.0.0.1:9000/")
				if err != nil {
					fail++
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("/", "失败:", fail)
}
