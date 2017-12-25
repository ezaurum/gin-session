package gin_session

import (
	"github.com/ezaurum/session"
	"github.com/gin-gonic/gin"
)

type Authenticator interface {
	Init() Authenticator
	Callbacks(create SessionHandlerFunc, activate SessionHandlerFunc,
		invalidCookie SessionCookieHandlerFunc,
		createAuth SessionHandlerFunc, loadAuth SessionCookieHandlerFunc) Authenticator
	Handler() gin.HandlerFunc
}
type SessionCookieHandlerFunc func(*gin.Context, session.Session, string)
type SessionHandlerFunc func(c *gin.Context, session session.Session)
