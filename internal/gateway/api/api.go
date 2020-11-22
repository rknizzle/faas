package api

import (
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	r.GET("/ping", ping)

	r.POST("/functions", addFunctionHandler)
	r.POST("/functions/:fn", invokeHandler)

	// Listen and serve on localhost
	r.Run()
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func invokeHandler(c *gin.Context) {
	c.JSON(400, gin.H{
		"message": "Function invocation not implemented yet",
	})
}

func addFunctionHandler(c *gin.Context) {
	c.JSON(400, gin.H{
		"message": "Adding functions not implemented yet",
	})
}
