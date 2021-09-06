package server

import (
	"dap2pnet/client/kademlia/buckets"
	"dap2pnet/client/middlewares"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	TLSCertPath string
	TLSKeytPath string
	Port        uint16
	bucks       *buckets.Buckets
}

func Run(port uint16, bucks *buckets.Buckets) error {
	servConfig := &ServerConfig{
		TLSCertPath: "./certs/rendezvous.dap2p.net.pem",
		TLSKeytPath: "./certs/rendezvous.dap2p.net.key",
		Port:        port,
		bucks:       bucks,
	}

	return initializeEndpoints(servConfig)

}

func initializeEndpoints(servConfig *ServerConfig) error {
	gin.ForceConsoleColor()
	router := gin.New()
	router.Use(gin.Recovery(), gin.LoggerWithFormatter(middlewares.Logger))

	peersGroup := router.Group("/peers/")
	//peersGroup.Use(middlewares.SetPeerIdentity())

	InitPeerEndpoints(peersGroup, servConfig.bucks)
	return router.Run(":" + fmt.Sprint(servConfig.Port))
	//return router.RunTLS(":6667", servConfig.TLSCertPath, servConfig.TLSKeytPath)
}
