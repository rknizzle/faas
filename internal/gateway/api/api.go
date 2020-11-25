package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rknizzle/faas/internal/gateway/datastore"
	"github.com/rknizzle/faas/internal/gateway/deployer"
	"github.com/rknizzle/faas/internal/gateway/loadbalancer"
	"github.com/rknizzle/faas/internal/models"
)

type GatewayHandler struct {
	Deployer deployer.Deployer
	LB       loadbalancer.LoadBalancer
	DS       datastore.Datastore
}

func NewGatewayHandler(r *gin.Engine, deploy deployer.Deployer, lb loadbalancer.LoadBalancer, ds datastore.Datastore) {
	handler := &GatewayHandler{Deployer: deploy, LB: lb, DS: ds}
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
	// get the function name from the path param
	fn := c.Param("fn")

	// invoke the function on a runner machine
	err := gw.LB.SendToRunner("rkneills/" + fn)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func (gw GatewayHandler) addFunctionHandler(c *gin.Context) {
	fnData, err := fnDataFromReq(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get function data from request body",
		})
		return
	}
	err = gw.Deployer.Deploy(fnData)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to deploy function",
			"info":    err.Error(),
		})
		return
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
