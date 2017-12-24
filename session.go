package session

const DefaultSessionContextKey = "session context key - session"

type Store interface {
	GetNew() Session
	Get(id string) (Session, bool)
	Sessions() SessionMap
	Count() int
}

type Session interface {
	ID() string
}

type SessionMap map[string]Session
