package gin_session

import "github.com/gin-gonic/gin"

type Authenticator interface {
	Handler() gin.HandlerFunc
}
