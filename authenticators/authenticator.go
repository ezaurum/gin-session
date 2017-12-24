package authenticators

import "github.com/gin-gonic/gin"

type Authenticator interface {
	Handler() gin.HandlerFunc
}


