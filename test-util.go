package gin_session

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
)

//TODO 따로 뽑아서 util 묶기

func GetRequest(r *gin.Engine, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func GetRequestWithCookie(r *gin.Engine, url string, cookie []string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", strings.Join(cookie, ";"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func GetCookie(w *httptest.ResponseRecorder) []string {
	cookie := w.HeaderMap["Set-Cookie"]
	return cookie
}
