package middlewares

import (
	"errors"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	AuthorizationMiddlewareErrIdentityNotFound = errors.New("identity not found. Proxy didn't send CN field")
	AuthorizationMiddlewareErrUnvalidKey       = errors.New("key format is unvalid")
	AuthorizationMiddlewareErrKeyIDNotfound    = errors.New("please provide a key id")
)

func SetPeerIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		tlsHeader := c.GetHeader("X-Forwarded-Tls-Client-Cert-Info")
		if tlsHeader == "" {
			c.AbortWithError(http.StatusForbidden, AuthorizationMiddlewareErrIdentityNotFound)
		}

		q, err := url.QueryUnescape(tlsHeader)
		if err != nil {
			c.AbortWithError(http.StatusForbidden, AuthorizationMiddlewareErrIdentityNotFound)
		}

		q = strings.Split(q, ",")[0]
		identity := strings.ReplaceAll(strings.Split(q, "CN=")[1], "\"", "")
		c.Request.Header.Add("Peer-Identity", identity) // For logging purposes
		c.Set("Identity", identity)
	}
}

func CheckKey() gin.HandlerFunc {

	return func(c *gin.Context) {
		keyID := c.Param("keyid")
		if keyID == "" {
			c.AbortWithError(http.StatusForbidden, AuthorizationMiddlewareErrKeyIDNotfound)
			return
		}
		keyIntID := new(big.Int)
		keyIntID.SetString(keyID, 16)
		println(keyIntID.BitLen())

		if keyIntID.BitLen() > 256 {
			c.AbortWithError(http.StatusForbidden, AuthorizationMiddlewareErrUnvalidKey)
			return
		}

		c.Set("keyID", keyID)
	}

}
