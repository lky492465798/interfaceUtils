package interfaceUtils

import (
	"sync/atomic"
)

type WebSessionStat struct {
	path                                      string
	method                                    string
	runningcount                              uint64
	requestIntervalHistogram_0_1              uint64
	requestIntervalHistogram_1_10             uint64
	requestIntervalHistogram_10_100           uint64
	requestIntervalHistogram_100_1000         uint64
	requestIntervalHistogram_1000_10000       uint64
	requestIntervalHistogram_10000_100000     uint64
	requestIntervalHistogram_100000_1000000   uint64
	requestIntervalHistogram_1000000_10000000 uint64
	requestIntervalHistogram_10000000_more    uint64
}

//
func (w *WebSessionStat) requestIntervalHistogramRecord(nanoSpan uint64) {
	millis := nanoSpan / 1000 / 1000
	if millis < 1 {
		atomic.AddUint64(&w.requestIntervalHistogram_0_1, millis)
	} else if millis < 10 {
		atomic.AddUint64(&w.requestIntervalHistogram_1_10, millis)
	} else if millis < 100 {
		atomic.AddUint64(&w.requestIntervalHistogram_10_100, millis)
	} else if millis < 1000 {
		atomic.AddUint64(&w.requestIntervalHistogram_100_1000, millis)
	} else if millis < 10000 {
		atomic.AddUint64(&w.requestIntervalHistogram_1000_10000, millis)
	} else if millis < 100000 {
		atomic.AddUint64(&w.requestIntervalHistogram_10000_100000, millis)
	} else if millis < 1000000 {
		atomic.AddUint64(&w.requestIntervalHistogram_100000_1000000, millis)
	} else if millis < 10000000 {
		atomic.AddUint64(&w.requestIntervalHistogram_1000000_10000000, millis)
	} else {
		atomic.AddUint64(&w.requestIntervalHistogram_10000000_more, millis)
	}
}
