package stat

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

const DEFAULTCAPICITY int = 1000

type urlWebStat struct {
	Path      string
	Method    string
	TimeArray []float64
	Head      int
	IsCircle  bool
	rw        sync.RWMutex
}

func (urlWebStat *urlWebStat) ShowInfo() ResBody4Inter {
	dst := make([]float64, len(urlWebStat.TimeArray))
	urlWebStat.rw.RLock()
	copy(dst, urlWebStat.TimeArray)
	urlWebStat.rw.RUnlock()
	sort.Float64s(dst)
	tp50 := getTP50(dst)
	tp99 := getTP99(dst)
	avg := getAverage(dst)
	max := getMax(dst)
	length := len(urlWebStat.TimeArray)
	path := urlWebStat.Path
	return ResBody4Inter{TP50: tp50, TP99: tp99, Path: path, Max: max, Average: avg, Times: length}
}

func (urlWebStat *urlWebStat) Add(time float64) {
	urlWebStat.rw.Lock()
	defer urlWebStat.rw.Unlock()
	if !urlWebStat.IsCircle {
		urlWebStat.TimeArray = append(urlWebStat.TimeArray, time)
		urlWebStat.Head++
		if urlWebStat.Head == DEFAULTCAPICITY {
			urlWebStat.IsCircle = true
			urlWebStat.Head = 0
		}
		return
	}
	urlWebStat.TimeArray[urlWebStat.Head] = time
	urlWebStat.Head++
	if urlWebStat.Head == DEFAULTCAPICITY {
		urlWebStat.Head = 0
	}
}

func getAverage(arr []float64) string {
	var sum float64
	for _, v := range arr {
		sum += v
	}
	return Float64ToTimeOfms(sum / float64(len(arr))).String()
}

func getMax(arr []float64) string {
	return Float64ToTimeOfms(arr[len(arr)-1]).String()
}

func getTP99(arr []float64) string {
	i := Ftoi(float64(len(arr))*0.99) - 1
	return Float64ToTimeOfms(arr[i]).String()
}

func getTP50(arr []float64) string {
	i := Ftoi(float64(len(arr)) * 0.5)
	return Float64ToTimeOfms(arr[i]).String()
}

func Ftoi(f float64) int {
	i, _ := strconv.Atoi(fmt.Sprintf("%1.0f", f))
	return i
}

func (urlWebStat *urlWebStat) Resize() {
	urlWebStat.rw.Lock()
	defer urlWebStat.rw.Unlock()
	urlWebStat.Head = 0
	urlWebStat.IsCircle = false
	urlWebStat.TimeArray = []float64{}
}

func TimeToFloatOfms(t time.Duration) float64 {
	return float64(t.Milliseconds())
}

func Float64ToTimeOfms(f float64) time.Duration {
	str := strconv.FormatFloat(f, 'f', -1, 64) + "ms"
	time, _ := time.ParseDuration(str)
	return time
}
