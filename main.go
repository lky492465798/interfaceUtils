package interfaceUtils

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	rand.Seed(time.Now().Unix())
	r.Use(UrlFilter)
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

	InitUrlFilter(r, []string{".txt", ".html"})
	r.Run(":9000")

}
