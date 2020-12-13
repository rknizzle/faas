package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rknizzle/faas/internal/gateway/datastore"
	"github.com/rknizzle/faas/internal/gateway/deployer"
	"github.com/rknizzle/faas/internal/gateway/loadbalancer"
	"github.com/rknizzle/faas/internal/models"
)

type gatewayHandler struct {
	d  deployer.Deployer
	lb loadbalancer.LoadBalancer
	ds datastore.Datastore
}

func NewGatewayHandler(r *gin.Engine, d deployer.Deployer, lb loadbalancer.LoadBalancer, ds datastore.Datastore) {
	handler := &gatewayHandler{d: d, lb: lb, ds: ds}
	r.GET("/ping", ping)

	r.POST("/functions", handler.addFunctionHandler)
	r.POST("/functions/:fn", handler.invokeHandler)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (gw gatewayHandler) invokeHandler(c *gin.Context) {
	// get the function name from the path param
	fn := c.Param("fn")
	// TODO: get fn input from request body

	// invoke the function on a runner machine
	_, err := gw.lb.SendToRunner(fn, "inputPlaceholder")
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

func (gw gatewayHandler) addFunctionHandler(c *gin.Context) {
	fnData, err := fnDataFromReq(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get function data from request body",
		})
		return
	}
	err = gw.d.Deploy(fnData)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to deploy function",
			"info":    err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"invoke": c.Request.Host + "/functions/" + fnData.Name,
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
