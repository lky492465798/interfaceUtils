package interfaceUtils

import (
	"sync"
)

type WebAppStatValue struct {
	Path                                      string
	Method                                    string
	rw                                        sync.RWMutex
	runningcount                              int
	requestIntervalHistogram_0_1              int
	requestIntervalHistogram_1_10             int
	requestIntervalHistogram_10_100           int
	requestIntervalHistogram_100_1000         int
	requestIntervalHistogram_1000_10000       int
	requestIntervalHistogram_10000_100000     int
	requestIntervalHistogram_100000_1000000   int
	requestIntervalHistogram_1000000_10000000 int
	requestIntervalHistogram_10000000_more    int
}
