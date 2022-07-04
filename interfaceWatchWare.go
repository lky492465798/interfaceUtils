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

// 是否使用默认配置(注册路由后失效)
var useDefault bool = true

var ignorePathMap map[string]struct{}

var urlWebStats map[string]*urlWebStat = make(map[string]*urlWebStat)

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
	start := time.Now()
	c.Next()
	diff := time.Since(start)
	_, ok := urlWebStats[path]
	if !ok {
		lock.Lock()
		if _, ok := urlWebStats[path]; !ok {
			urlWebStats[path] = &urlWebStat{Path: path, Method: method, Head: 0, IsCircle: false}
		}
		lock.Unlock()
	}
	go addInterInfoASYNC(urlWebStats[path], &diff)
}

func addInterInfoASYNC(t *urlWebStat, time *time.Duration) {
	t.Add(TimeToFloatOfms(*time))
}

// 获取Stats结果集
func GetStatHandler(c *gin.Context) {
	infos := make(ResBody4Inters, len(GetUrlWebStats()))
	currIdx := 0
	for _, v := range GetUrlWebStats() {
		infos[currIdx] = v.ShowInfo()
		currIdx++
	}
	sort.Sort(infos)
	c.JSON(200, gin.H{" 接口信息: ": infos})
}

// 清空Stats结果集
func FlushUrlWebStats(c *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	urlWebStats = make(map[string]*urlWebStat)
}

func InitUrlFilter(r *gin.Engine, suffixs []string) {
	r.GET("/urlstat/info", GetStatHandler)
	r.GET("/urlstat/del", FlushUrlWebStats)
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

func GetUrlWebStats() map[string]*urlWebStat {
	return urlWebStats
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
