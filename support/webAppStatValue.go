package interfaceUtils

import "fmt"

type Webstats []WebAppStatValue

type WebAppStatValue struct {
	Path                                      string `json:" 请求路径: "`
	Method                                    string `json:" 请求方式: "`
	RunningCount                              int32  `json:" 当前运行数量: "`
	ConcurrentMax                             int32  `json:" 最大并发数量: "`
	RequestTimeNanoMax                        int64  `json:" 最大耗时时间(ms): "`
	RequestTimeNanoMaxOccurTime               string `json:" 最大耗时发生时间: "`
	RequestTimeNano                           int64  `json:" 请求总时间(ms): "`
	AvgRequestTimeNano                        int64  `json:" 请求平均时间(ms): "`
	RequestIntervalHistogram_0_1              int64  `json:" 第一区间: "`
	RequestIntervalHistogram_1_10             int64  `json:" 第二区间: "`
	RequestIntervalHistogram_10_100           int64  `json:" 第三区间: "`
	RequestIntervalHistogram_100_1000         int64  `json:" 第四区间: "`
	RequestIntervalHistogram_1000_10000       int64  `json:" 第五区间: "`
	RequestIntervalHistogram_10000_100000     int64  `json:" 第六区间: "`
	RequestIntervalHistogram_100000_1000000   int64  `json:" 第七区间: "`
	RequestIntervalHistogram_1000000_10000000 int64  `json:" 第八区间: "`
	RequestIntervalHistogram_10000000_more    int64  `json:" 第九区间: "`
	TP99                                      string `json:" TP99: "`
	TP50                                      string `json:" TP50: "`
}

func (w *WebAppStatValue) calTP() {
	arr := []int64{
		w.RequestIntervalHistogram_0_1,
		w.RequestIntervalHistogram_1_10,
		w.RequestIntervalHistogram_10_100,
		w.RequestIntervalHistogram_100_1000,
		w.RequestIntervalHistogram_1000_10000,
		w.RequestIntervalHistogram_10000_100000,
		w.RequestIntervalHistogram_100000_1000000,
		w.RequestIntervalHistogram_1000000_10000000,
		w.RequestIntervalHistogram_10000000_more,
	}
	sum := count(arr)
	if sum != 0 {
		w.AvgRequestTimeNano = w.RequestTimeNano / int64(sum)
	}
	itp99 := sum * 0.99
	itp50 := sum * 0.50
	var t int64 = 0
	for i, v := range arr {
		t += v
		if t >= int64(itp99) {
			w.TP99 = fmt.Sprintf("位于第 %d 区间!", i+1)
			break
		}
		if w.TP50 == "" && t >= int64(itp50) {
			w.TP50 = fmt.Sprintf("位于第 %d 区间!", i+1)
		}
	}
}

func (s Webstats) Len() int {
	return len(s)
}
func (s Webstats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Webstats) Less(i, j int) bool {
	return len(s[i].Path) < len(s[j].Path)
}

func count(arr []int64) float64 {
	var res int64 = 0
	for _, v := range arr {
		res += v
	}
	return float64(res)
}
