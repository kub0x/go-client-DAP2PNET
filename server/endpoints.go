package server

import (
	"github.com/gin-gonic/gin"
)

func InitPeerEndpoints(router *gin.RouterGroup) {
	router.POST("/subscribe")
}
