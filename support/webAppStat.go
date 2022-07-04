package interfaceUtils

import (
	"sync/atomic"
	"time"
)

type WebAppStat struct {
	Path                                      string
	Method                                    string
	runningCount                              int32
	concurrentMax                             int32
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

func (w *WebAppStat) GetRequestTimeNano() int64 {
	return w.requestTimeNano
}

func (w *WebAppStat) GetRequestTimeNanoMaxOccurTime() int64 {
	return w.requestTimeNanoMaxOccurTime
}

func (w *WebAppStat) GetRequestTimeNanoMax() int64 {
	return w.requestTimeNanoMax
}

func (w *WebAppStat) SetRequestTimeNanoMax(nanos int64) {
	for {
		current := atomic.LoadInt64(&(w.requestTimeNanoMax))
		if current < nanos {
			if atomic.CompareAndSwapInt64(&(w.requestTimeNanoMax), current, nanos) {
				w.requestTimeNanoMaxOccurTime = time.Now().UnixMilli()
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

func (w *WebAppStat) GetConcurrentMax() int32 {
	return w.concurrentMax
}

func (w *WebAppStat) IncreRunningCount() {
	atomic.AddInt32(&w.runningCount, 1)
}

func (w *WebAppStat) DecreRunningCount() {
	atomic.AddInt32(&(w.runningCount), -1)
}

func (w *WebAppStat) GetRunningCount() int32 {
	return w.runningCount
}

func (w *WebAppStat) histogramRecord(nanoSpan int64) {
	millis := nanoSpan / 1000 / 1000
	if millis < 1 {
		atomic.AddInt64(&(w.requestIntervalHistogram_0_1), millis)
	} else if millis < 10 {
		atomic.AddInt64(&(w.requestIntervalHistogram_1_10), millis)
	} else if millis < 100 {
		atomic.AddInt64(&(w.requestIntervalHistogram_10_100), millis)
	} else if millis < 1000 {
		atomic.AddInt64(&(w.requestIntervalHistogram_100_1000), millis)
	} else if millis < 10000 {
		atomic.AddInt64(&(w.requestIntervalHistogram_1000_10000), millis)
	} else if millis < 100000 {
		atomic.AddInt64(&(w.requestIntervalHistogram_10000_100000), millis)
	} else if millis < 1000000 {
		atomic.AddInt64(&(w.requestIntervalHistogram_100000_1000000), millis)
	} else if millis < 10000000 {
		atomic.AddInt64(&(w.requestIntervalHistogram_1000000_10000000), millis)
	} else {
		atomic.AddInt64(&(w.requestIntervalHistogram_10000000_more), millis)
	}
}
