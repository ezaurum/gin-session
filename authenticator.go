package gin_session

import "github.com/gin-gonic/gin"

const DefaultSessionContextKey = "default session context key"

type Authenticator interface {
	Handler() gin.HandlerFunc
}


