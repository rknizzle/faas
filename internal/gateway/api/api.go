package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rknizzle/faas/internal/gateway/deployer"
	"github.com/rknizzle/faas/internal/models"
)

type GatewayHandler struct {
	Deployer deployer.Deployer
}

func NewGatewayHandler(r *gin.Engine, deploy deployer.Deployer) {
	handler := &GatewayHandler{Deployer: deploy}
	r.GET("/ping", ping)

	r.POST("/functions", handler.addFunctionHandler)
	r.POST("/functions/:fn", handler.invokeHandler)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (gw GatewayHandler) invokeHandler(c *gin.Context) {
	c.JSON(400, gin.H{
		"message": "Function invocation not implemented yet",
	})
}

func (gw GatewayHandler) addFunctionHandler(c *gin.Context) {
	fnData, err := fnDataFromReq(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get function data from request body",
		})
	}
	err = gw.Deployer.Deploy(fnData)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to deploy function",
			"info":    err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"message": "ok",
	})
}
func fnDataFromReq(c *gin.Context) (models.FnData, error) {
	rawData, err := c.GetRawData()
	if err != nil {
		return models.FnData{}, err
	}

	var fnData models.FnData
	err = json.Unmarshal(rawData, &fnData)
	if err != nil {
		return models.FnData{}, err
	}
	return fnData, nil
}
