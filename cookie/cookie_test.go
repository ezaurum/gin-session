package cookie

import (
	gs "github.com/ezaurum/gin-session"
	"github.com/ezaurum/session"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestPanic(t *testing.T) {
	authenticator := Default()
	assert.Panics(t, func() { authenticator.Handler() }, "Handler before init must panic.")
}

func TestCookie(t *testing.T) {
	r := getDefault()

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(DefaultSessionContextKey).(session.Session)
		c.String(http.StatusOK, s.ID())
	})

	// first request
	w := gs.GetRequest(r, "/")
	cookie := w.HeaderMap["Set-Cookie"]

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), strings.Split(strings.Split(cookie[0], ";")[0], "=")[1])

	// first request with cookie
	w0 := gs.GetRequestWithCookie(r, "/", cookie)
	_, isExist := w0.HeaderMap["Set-Cookie"]

	assert.Equal(t, http.StatusOK, w0.Code)
	assert.True(t, !isExist)
	assert.Equal(t, w0.Body.String(), w.Body.String())
}

func TestCallbacks(t *testing.T) {

	var activated, created, loadedAuth, createAuth, invalidCookie bool

	r := gin.Default()
	authenticator := Default().Callbacks(
		func(c *gin.Context, session session.Session) {
			created = true
		},
		func(c *gin.Context, session session.Session) {
			activated = true
		},
		func(context *gin.Context, i session.Session, s string) {
			invalidCookie = true
		},
		func(c *gin.Context, session session.Session) {
			createAuth = true
		},
		func(context *gin.Context, i session.Session, s string) {
			loadedAuth = true
		},
	).Init()

	r.Use(authenticator.Handler())

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(DefaultSessionContextKey).(session.Session)
		c.String(http.StatusOK, s.ID())
	})

	w := gs.GetRequest(r, "/")
	cookie := gs.GetCookie(w)

	assert.True(t, created)
	assert.True(t, activated)

	created = false
	activated = false
	gs.GetRequestWithCookie(r, "/", cookie)

	assert.True(t, !created)
	assert.True(t, activated)
}

//TODO load auth - 여기서 valid, invalid를 체크해야 한다. has auth info 정도가 좋을까?
//TODO create auth - 실제로 auth 를 생성하는 건 아니다. 그러면 여기서는 create 가 아니라 first? no auth info?
//TODO invalid cookie
//TODO no session 이지 invalid cookie 가 아니라

func getDefault() *gin.Engine {
	r := gin.Default()
	authenticator := Default()
	authenticator.Init()
	r.Use(authenticator.Handler())
	return r
}
