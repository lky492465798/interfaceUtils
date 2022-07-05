package interfaceUtils

import (
	"sync/atomic"
	"time"
)

type WebAppStat struct {
	Path                                      string
	Method                                    string
	runningCount                              int32 // before
	concurrentMax                             int32 // before
	requestTimeNanoMax                        int64
	requestTimeNanoMaxOccurTime               int64
	requestTimeNano                           int64
	requestIntervalHistogram_0_1              int64
	requestIntervalHistogram_1_10             int64
	requestIntervalHistogram_10_100           int64
	requestIntervalHistogram_100_1000         int64
	requestIntervalHistogram_1000_10000       int64
	requestIntervalHistogram_10000_100000     int64
	requestIntervalHistogram_100000_1000000   int64
	requestIntervalHistogram_1000000_10000000 int64
	requestIntervalHistogram_10000000_more    int64
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
	if millis < 1 {
		atomic.AddInt64(&(w.requestIntervalHistogram_0_1), 1)
	} else if millis < 10 {
		atomic.AddInt64(&(w.requestIntervalHistogram_1_10), 1)
	} else if millis < 50 {
		atomic.AddInt64(&(w.requestIntervalHistogram_10_100), 1)
	} else if millis < 100 {
		atomic.AddInt64(&(w.requestIntervalHistogram_100_1000), 1)
	} else if millis < 300 {
		atomic.AddInt64(&(w.requestIntervalHistogram_1000_10000), 1)
	} else if millis < 500 {
		atomic.AddInt64(&(w.requestIntervalHistogram_10000_100000), 1)
	} else if millis < 1000 {
		atomic.AddInt64(&(w.requestIntervalHistogram_100000_1000000), 1)
	} else if millis < 5000 {
		atomic.AddInt64(&(w.requestIntervalHistogram_1000000_10000000), 1)
	} else {
		atomic.AddInt64(&(w.requestIntervalHistogram_10000000_more), 1)
	}
}

func (w *WebAppStat) GetValue() WebAppStatValue {
	res := WebAppStatValue{Path: w.Path, Method: w.Method}
	res.ConcurrentMax = get32(&w.concurrentMax)
	res.RequestTimeNano = get64(&w.requestTimeNano) / 1000 / 1000
	res.RequestTimeNanoMax = get64(&w.requestTimeNanoMax) / 1000 / 1000
	res.RequestTimeNanoMaxOccurTime = time.Unix(0, get64(&w.requestTimeNanoMaxOccurTime)).Format("2006-01-02 15:04:05")
	res.RunningCount = get32(&w.runningCount)
	res.RequestIntervalHistogram_0_1 = get64(&w.requestIntervalHistogram_0_1)
	res.RequestIntervalHistogram_1_10 = get64(&w.requestIntervalHistogram_1_10)
	res.RequestIntervalHistogram_1_10 = get64(&w.requestIntervalHistogram_1_10)
	res.RequestIntervalHistogram_10_100 = get64(&w.requestIntervalHistogram_10_100)
	res.RequestIntervalHistogram_100_1000 = get64(&w.requestIntervalHistogram_100_1000)
	res.RequestIntervalHistogram_1000_10000 = get64(&w.requestIntervalHistogram_1000_10000)
	res.RequestIntervalHistogram_10000_100000 = get64(&w.requestIntervalHistogram_10000_100000)
	res.RequestIntervalHistogram_100000_1000000 = get64(&w.requestIntervalHistogram_100000_1000000)
	res.RequestIntervalHistogram_1000000_10000000 = get64(&w.requestIntervalHistogram_1000000_10000000)
	res.RequestIntervalHistogram_10000000_more = get64(&w.requestIntervalHistogram_10000000_more)
	res.calTP()
	return res
}

func get64(t *int64) int64 {
	return atomic.LoadInt64(t)
}

func get32(t *int32) int32 {
	return atomic.LoadInt32(t)
}
