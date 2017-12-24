package cookie

import (
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/http"
	"github.com/ezaurum/gin-session"
	"strings"
)

func TestCookie(t *testing.T) {
	r := gin.Default()

	authenticator := Default()
	r.Use(authenticator.Handler())

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(session.DefaultSessionContextKey).(session.Session)
		c.String(http.StatusOK, s.ID())
	})

	// first request
	w := getRequest(r, "/")
	cookie := w.HeaderMap["Set-Cookie"]

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String() , strings.Split(strings.Split(cookie[0],";")[0], "=")[1])

	// first request with cookie
	w0 := getRequestWithCookie(r, "/", cookie)
	_, isExist := w0.HeaderMap["Set-Cookie"]

	assert.Equal(t, http.StatusOK, w0.Code)
	assert.True(t, !isExist)
	assert.Equal(t, w0.Body.String() , w.Body.String())
}

func getRequest(r *gin.Engine, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func getRequestWithCookie(r *gin.Engine, url string, cookie []string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", strings.Join(cookie,";"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
