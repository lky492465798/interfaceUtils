package interfaceUtils

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 1. r.Use(UrlFilter) 注册中间件
// 2. 注册GetStatHandler的路径, 根据需要给FlushUrlWebStats绑定一个路径
// 3. r.InitUrlFilter(r *gin.Engine, suffixs []string) 初始化监控路由,过滤路由列表
// *ps: 第3步(在注册完接口后调用)防止不能正确加载需要监控的路由

// TODO: 1. 拿到nano后时间处理过程,
//  	 2. 封装接口信息

// 是否使用默认配置(注册路由后失效)
var useDefault bool = true

// 忽略地址(目前支持.xxx格式)
var ignorePathMap map[string]struct{}

var WebAppStats map[string]*WebAppStat = make(map[string]*WebAppStat)

var methodTries map[string]*node

var lock sync.RWMutex

func UrlFilter(c *gin.Context) {
	path := c.Request.URL.Path
	method := c.Request.Method
	if !useDefault {
		root := methodTries[method]
		node := root.search(parsePattern(path), 0)
		if node == nil || isIgnore(path) {
			return
		}
		path = node.pattern
	}
	_, ok := WebAppStats[path+"-"+method]
	if !ok {
		lock.Lock()
		if _, ok := WebAppStats[path+"-"+method]; !ok {
			WebAppStats[path+"-"+method] = &WebAppStat{Path: path, Method: method}
		}
		lock.Unlock()
	}
	concurrentCal(WebAppStats[path+"-"+method])
	start := time.Now().UnixNano()
	c.Next()
	diff := time.Now().UnixNano() - start
	go addInterInfoASYNC(WebAppStats[path+"-"+method], diff)
}

func addInterInfoASYNC(w *WebAppStat, nanos int64) {
	// 时间处理逻辑
	w.DecreRunningCount()
	w.SetRequestTimeNano(nanos)
	w.SetRequestTimeNanoMax(nanos)
	w.histogramRecord(nanos)
}

func concurrentCal(w *WebAppStat) {
	running := w.IncreRunningCount()
	w.SetConcurrentMax(running)
}

// 获取Stats结果集
func GetStatHandler(c *gin.Context) {
	stats := make(Webstats, len(GetUrlWebStats()))
	currIdx := 0
	for _, v := range GetUrlWebStats() {
		stats[currIdx] = v.GetValue()
		currIdx++
	}
	sort.Sort(stats)
	c.JSON(200, gin.H{" 接口信息: ": stats})
}

// 清空Stats结果集
// func FlushUrlWebStats(c *gin.Context) {
// 	lock.Lock()
// 	defer lock.Unlock()
// 	urlWebStats = make(map[string]*urlWebStat)
// }

func InitUrlFilter(r *gin.Engine, suffixs []string) {
	r.GET("/urlstat/info", GetStatHandler)
	// r.GET("/urlstat/del", FlushUrlWebStats)
	registryPath(r)
	registryIgnorePath(suffixs)
}

func registryPath(r *gin.Engine) {
	if r == nil {
		return
	}

	useDefault = false
	methodTries = make(map[string]*node)
	paths := r.Routes()
	for _, v := range paths {
		parts := parsePattern(v.Path)
		method := v.Method
		_, ok := methodTries[method]
		if !ok {
			methodTries[method] = &node{}
		}
		methodTries[method].insert(v.Path, parts, 0)
	}
}

func isIgnore(path string) bool {
	index := strings.LastIndex(path, ".")
	if index == -1 {
		return false
	}

	_, ok := ignorePathMap[path[index:]]
	return ok
}

func registryIgnorePath(suffixs []string) {
	if suffixs == nil {
		return
	}

	ignorePathMap = make(map[string]struct{})
	for _, v := range suffixs {
		ignorePathMap[v] = struct{}{}
	}
}

func GetUrlWebStats() map[string]*WebAppStat {
	return WebAppStats
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
