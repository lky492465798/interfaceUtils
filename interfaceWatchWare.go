package interfaceUtils

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var lock sync.RWMutex

func UrlFilter(c *gin.Context) {
	fmt.Println("被执行!!!!")
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

func InitUrlFilter(r *gin.Engine, suffixs []string) {
	r.GET("/urlwebstat/list", func(c *gin.Context) {
		infos := make(ResBody4Inters, len(GetUrlWebStats()))
		currIdx := 0
		for _, v := range GetUrlWebStats() {
			infos[currIdx] = v.ShowInfo()
			currIdx++
		}
		sort.Sort(infos)
		c.JSON(200, gin.H{" 接口信息: ": infos})

	})
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

	ignorePathMap = make(map[string]string)
	for _, v := range suffixs {
		ignorePathMap[v] = ""
	}
}

func GetUrlWebStats() map[string]*urlWebStat {
	return urlWebStats
}

var useDefault bool = true

var ignorePathMap map[string]string

var urlWebStats map[string]*urlWebStat = make(map[string]*urlWebStat)

var methodTries map[string]*node

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
