package interfaceUtils

import (
	"sync/atomic"
	"time"
)

// 时间区间(ms)
var flag []int64 = []int64{1, 10, 50, 100, 300, 500, 1000, 5000}

type WebAppStat struct {
	Path                        string
	Method                      string
	runningCount                int32 // before
	concurrentMax               int32 // before
	requestTimeNanoMax          int64
	requestTimeNanoMaxOccurTime int64
	requestTimeNano             int64
	requestIntervalHistogram1   int64
	requestIntervalHistogram2   int64
	requestIntervalHistogram3   int64
	requestIntervalHistogram4   int64
	requestIntervalHistogram5   int64
	requestIntervalHistogram6   int64
	requestIntervalHistogram7   int64
	requestIntervalHistogram8   int64
	requestIntervalHistogram9   int64
}

func (w *WebAppStat) SetRequestTimeNano(nanos int64) {
	atomic.AddInt64(&(w.requestTimeNano), nanos)
}

func (w *WebAppStat) SetRequestTimeNanoMax(nanos int64) {
	for {
		current := atomic.LoadInt64(&(w.requestTimeNanoMax))
		if current < nanos {
			if atomic.CompareAndSwapInt64(&(w.requestTimeNanoMax), current, nanos) {
				w.requestTimeNanoMaxOccurTime = time.Now().UnixNano()
				break
			} else {
				continue
			}
		} else {
			break
		}
	}

}

func (w *WebAppStat) SetConcurrentMax(running int32) {
	for {
		max := atomic.LoadInt32(&(w.concurrentMax))
		if running > max {
			if atomic.CompareAndSwapInt32(&(w.concurrentMax), w.concurrentMax, running) {
				break
			}
		} else {
			break
		}
	}
}

func (w *WebAppStat) IncreRunningCount() int32 {
	return atomic.AddInt32(&w.runningCount, 1)
}

func (w *WebAppStat) DecreRunningCount() {
	atomic.AddInt32(&(w.runningCount), -1)
}

func (w *WebAppStat) histogramRecord(nanoSpan int64) {
	millis := nanoSpan / 1000 / 1000
	if millis < flag[0] {
		atomic.AddInt64(&(w.requestIntervalHistogram1), 1)
	} else if millis < flag[1] {
		atomic.AddInt64(&(w.requestIntervalHistogram2), 1)
	} else if millis < flag[2] {
		atomic.AddInt64(&(w.requestIntervalHistogram3), 1)
	} else if millis < flag[3] {
		atomic.AddInt64(&(w.requestIntervalHistogram4), 1)
	} else if millis < flag[4] {
		atomic.AddInt64(&(w.requestIntervalHistogram5), 1)
	} else if millis < flag[5] {
		atomic.AddInt64(&(w.requestIntervalHistogram6), 1)
	} else if millis < flag[6] {
		atomic.AddInt64(&(w.requestIntervalHistogram7), 1)
	} else if millis < flag[7] {
		atomic.AddInt64(&(w.requestIntervalHistogram8), 1)
	} else {
		atomic.AddInt64(&(w.requestIntervalHistogram9), 1)
	}
}

func (w *WebAppStat) GetValue() WebAppStatValue {
	res := WebAppStatValue{Path: w.Path, Method: w.Method}
	res.ConcurrentMax = get32(&w.concurrentMax)
	res.RequestTimeNano = get64(&w.requestTimeNano) / 1000 / 1000
	res.RequestTimeNanoMax = get64(&w.requestTimeNanoMax) / 1000 / 1000
	res.RequestTimeNanoMaxOccurTime = time.Unix(0, get64(&w.requestTimeNanoMaxOccurTime)).Format("2006-01-02 15:04:05")
	res.RunningCount = get32(&w.runningCount)
	res.RequestIntervalHistogram1 = get64(&w.requestIntervalHistogram1)
	res.RequestIntervalHistogram2 = get64(&w.requestIntervalHistogram2)
	res.RequestIntervalHistogram3 = get64(&w.requestIntervalHistogram3)
	res.RequestIntervalHistogram4 = get64(&w.requestIntervalHistogram4)
	res.RequestIntervalHistogram5 = get64(&w.requestIntervalHistogram5)
	res.RequestIntervalHistogram6 = get64(&w.requestIntervalHistogram6)
	res.RequestIntervalHistogram7 = get64(&w.requestIntervalHistogram7)
	res.RequestIntervalHistogram8 = get64(&w.requestIntervalHistogram8)
	res.RequestIntervalHistogram9 = get64(&w.requestIntervalHistogram9)
	res.calTP()
	return res
}

func get64(t *int64) int64 {
	return atomic.LoadInt64(t)
}

func get32(t *int32) int32 {
	return atomic.LoadInt32(t)
}
