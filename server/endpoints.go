package server

import (
	"dap2pnet/client/handlers"
	"dap2pnet/client/kademlia/buckets"
	"dap2pnet/client/middlewares"

	"github.com/gin-gonic/gin"
)

func InitPeerEndpoints(router *gin.RouterGroup, bucks *buckets.Buckets) {
	router.GET("/key/:keyid", middlewares.CheckKey(), handlers.OnFindNode(bucks))
}
