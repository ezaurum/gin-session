package cookie

import (
	"github.com/ezaurum/gin-session/authenticators"
	gs "github.com/ezaurum/gin-session"
	"github.com/ezaurum/gin-session/stores/memstore"
	"github.com/gin-gonic/gin"
)

var _ authenticators.Authenticator = cookieAuthenticator{}

const cookieName = "ca-default-name"

func Default() authenticators.Authenticator {
	return New(memstore.Default())
}

func New(store gs.Store) authenticators.Authenticator {
	return cookieAuthenticator{
		store : store,
	}
}

type cookieAuthenticator struct {
	store gs.Store
}

func (ca cookieAuthenticator) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var session gs.Session
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
