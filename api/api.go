package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	Init()
}

func Init() {
	fmt.Println("Starting...")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Listen and serve on localhost
	r.Run()
}
