package cookie

import (
	gs "github.com/ezaurum/gin-session"
	"github.com/ezaurum/session/stores/memstore"
	"github.com/gin-gonic/gin"
	"github.com/ezaurum/session"
)

var _ gs.Authenticator = cookieAuthenticator{}

const cookieName = "ca-default-name"

func Default() gs.Authenticator {
	return New(memstore.Default())
}

func New(store session.Store) gs.Authenticator {
	return cookieAuthenticator{
		store : store,
	}
}

type cookieAuthenticator struct {
	store session.Store
}

func (ca cookieAuthenticator) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var session session.Session
		cookie, e := c.Cookie(cookieName)
		if nil != e {
			// set cookie
			session = ca.store.GetNew()
			c.SetCookie(cookieName, session.ID(), 100000, "/", "*", false, true )
		} else {
			// authenticate
			s, is := ca.store.Get(cookie)

			if !is {
				//TODO error
				return
			}

			session = s
		}

		c.Set(gs.DefaultSessionContextKey, session)
	}
}
