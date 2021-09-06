package handlers

import (
	"dap2pnet/client/kademlia/buckets"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnFindNode(bucks *buckets.Buckets) gin.HandlerFunc {

	return func(c *gin.Context) {
		keyID := c.GetString("keyID")
		if keyID == "" {
			c.AbortWithError(http.StatusForbidden, errors.New("keyid not in context"))
			return
		}

		peers := bucks.NearestToKey(keyID)

		if peers == nil {
			c.AbortWithError(http.StatusForbidden, errors.New("peer list not available"))
			return
		}

		c.JSON(http.StatusOK, peers)
	}

}
