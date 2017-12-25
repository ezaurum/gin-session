package cookie

import (
	gs "github.com/ezaurum/gin-session"
	ezs "github.com/ezaurum/session"
	"github.com/ezaurum/session/stores/memstore"
	"github.com/gin-gonic/gin"
	"log"
)

var _ gs.Authenticator = &cookieAuthenticator{}

const (
	DefaultSessionContextKey = "default session context key"
	sessionIDCookieName      = "ca-default-name"
	persistCookieName        = "ca-default-remember-me"
)

func Default() gs.Authenticator {
	return New(memstore.Default())
}

func New(store ezs.Store) gs.Authenticator {
	logF := func(c *gin.Context, s ezs.Session, cs string) {
		log.Printf("%v, %v", s, cs)
	}
	logF0 := func(c *gin.Context, s ezs.Session) {
		log.Printf("%v", s)
	}
	return &cookieAuthenticator{
		store: store,
		cookieInvalidCallback:        DefaultCookieInvalid,
		createAuthenticationCallback: logF0,
		loadAuthenticationCallback:   logF,
		activateCallback:             logF0,
		createCallback:               logF0,
	}
}

type cookieAuthenticator struct {
	store                        ezs.Store
	isInit                       bool
	createCallback               gs.SessionHandlerFunc
	activateCallback             gs.SessionHandlerFunc
	cookieInvalidCallback        gs.SessionCookieHandlerFunc
	createAuthenticationCallback gs.SessionHandlerFunc
	loadAuthenticationCallback   gs.SessionCookieHandlerFunc
}

func (ca *cookieAuthenticator) Callbacks(create gs.SessionHandlerFunc,
	activate gs.SessionHandlerFunc, invalid gs.SessionCookieHandlerFunc,
	createAuth gs.SessionHandlerFunc, load gs.SessionCookieHandlerFunc) gs.Authenticator {
	if ca.isInit {
		panic("Must before init")
	}
	if nil != create {
		ca.createCallback = create
	}
	if nil != activate {
		ca.activateCallback = activate
	}
	if nil != invalid {
		ca.cookieInvalidCallback = invalid
	}
	if nil != createAuth {
		ca.createAuthenticationCallback = createAuth
	}
	if nil != load {
		ca.loadAuthenticationCallback = load
	}
	return ca
}

func (ca cookieAuthenticator) Handler() gin.HandlerFunc {
	if !ca.isInit {
		panic("Not initialized")
	}

	return func(c *gin.Context) {

		session, needSession := ca.getSession(c)

		if needSession {
			// set cookie
			session = ca.store.GetNew()
			//TODO 여기 duration 과 세션 duration 맞추기
			c.SetCookie(sessionIDCookieName, session.ID(),
				100000,
				"/", "",
				false, true)

			// created
			ca.createCallback(c, session)
		}

		//TODO 여기서 authenticate가 일어나야 한다
		ca.authenticate(c, session)

		//activate
		//TODO gin에 종속된 부분은 여기다
		c.Set(DefaultSessionContextKey, session)

		ca.activateCallback(c, session)
	}
}

func (ca cookieAuthenticator) getSession(c *gin.Context) (ezs.Session, bool) {
	sessionIDCookie, e := c.Cookie(sessionIDCookieName)
	needSession := nil != e
	var session ezs.Session
	var noSession bool
	if nil == e {
		//TODO secure
		s, n := ca.store.Get(sessionIDCookie)
		session = s
		noSession = !n

		if noSession {
			// 세션 유효하지 않은 경우
			// 만료되었거나, 값 조작이거나
			ca.cookieInvalidCallback(c, nil, sessionIDCookie)
		}
	}
	needSession = needSession || noSession
	return session, needSession
}

func (ca cookieAuthenticator) authenticate(c *gin.Context, session ezs.Session) {
	//TODO 세션이 이미 authenticate 되었는지 체크해서 되었다면 안 건드려도 됨
	pCookie, e := c.Cookie(persistCookieName)
	if nil != e {
		if nil != ca.createAuthenticationCallback {
			ca.createAuthenticationCallback(c, session)
		}
	} else {
		// load user with cookie
		if nil != ca.loadAuthenticationCallback {
			ca.loadAuthenticationCallback(c, session, pCookie)
		}
	}
}

func DefaultCookieInvalid(c *gin.Context, s ezs.Session, cookie string) {
	//TODO error
	log.Printf("cookie invalid %v", cookie)
	//쿠키 있는데 세션이 없는 경우?\
	//보통 세션 만료시나, 쿠키 변조시?
	// 쿠키 삭제
	//TODO c.SetCookie(sessionIDCookieName, "", -1,"/","", false, true  )
	// 리다이렉트로 처리
	//TODO c.Redirect(http.StatusFound, c.Request.URL.String())
	//TODO c.Abort()
}

func (ca *cookieAuthenticator) Init() gs.Authenticator {
	ca.isInit = true
	return ca
}
