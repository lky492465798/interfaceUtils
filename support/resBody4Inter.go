package interfaceUtils

type ResBody4Inter struct {
	Path    string `json:" 请求路径: "`
	Max     string `json:" 最长响应时间: "`
	TP50    string `json:" TP50: "`
	TP99    string `json:" TP99: "`
	Average string `json:" 平均响应时间: "`
	Times   int    `json:" 统计次数: "`
}

type ResBody4Inters []ResBody4Inter

