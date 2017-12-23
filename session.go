package session

type Store interface {
	GetNew() Session
	Get(id string) (Session, bool)
}

type Session interface {
	ID() string
}
