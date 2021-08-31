package server

import (
	"dap2pnet/client/middlewares"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	TLSCertPath string
	TLSKeytPath string
	Port        uint16
}

func Run(port uint16) error {
	servConfig := &ServerConfig{
		TLSCertPath: "./certs/rendezvous.dap2p.net.pem",
		TLSKeytPath: "./certs/rendezvous.dap2p.net.key",
		Port:        port,
	}

	return InitializeEndpoints(servConfig)

}

func InitializeEndpoints(servConfig *ServerConfig) error {
	gin.ForceConsoleColor()
	router := gin.New()
	router.Use(gin.Recovery(), gin.LoggerWithFormatter(middlewares.Logger))

	peersGroup := router.Group("/peers/")
	peersGroup.Use(middlewares.SetPeerIdentity())

	InitPeerEndpoints(peersGroup)
	return router.Run(":" + fmt.Sprint(servConfig.Port))
	//return router.RunTLS(":6667", servConfig.TLSCertPath, servConfig.TLSKeytPath)
}
