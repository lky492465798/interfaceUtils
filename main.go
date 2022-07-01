package main

import (
	"math/rand"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	stat "github.com/lky492465798/interfaceUtils/support"
)

func main() {
	r := gin.Default()
	rand.Seed(time.Now().Unix())
	r.Use(stat.UrlFilter)
	r.GET("/index", func(c *gin.Context) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
	})

	r.GET("/params/a/:t", func(c *gin.Context) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(80)))
	})
	r.GET("/index/a", func(c *gin.Context) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	})
	r.GET("/", func(c *gin.Context) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(30)))
	})

	r.GET("/info", func(c *gin.Context) {
		infos := make(stat.ResBody4Inters, len(stat.GetUrlWebStats()))
		currIdx := 0
		for _, v := range stat.GetUrlWebStats() {
			infos[currIdx] = v.ShowInfo()
			currIdx++
		}
		sort.Sort(infos)
		c.JSON(200, gin.H{" 接口信息: ": infos})

	})
	stat.InitUrlFilter(r, []string{".txt", ".html"})
	r.Run(":9000")

}
