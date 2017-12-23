package session

import (
	"github.com/bwmarrin/snowflake"
	"github.com/patrickmn/go-cache"
	"time"
)

type Manager interface {
	GetNew() Session
	Get(id string) (Session, bool)
}

type Session interface {
	ID() string
}

var _ Manager = DefaultSessionManager{}
var _ Session = DefaultSession{}

func Default() Manager {
	node, err := snowflake.NewNode(0)
	if nil != err {
		panic(err)
	}
	c := cache.New(5*time.Minute, 10*time.Minute)

	return DefaultSessionManager{
		generator:node,
		cache:c,
	}
}

type DefaultSessionManager struct {
	generator *snowflake.Node
	cache     *cache.Cache
}

func (sm DefaultSessionManager) GetNew() Session {
	id := sm.generator.Generate().String()
	s := DefaultSession{
		id:id,
	}

	sm.cache.Set(id, s, cache.DefaultExpiration)

	return s
}

func (sm DefaultSessionManager) Get(id string) (Session, bool) {
	s, e := sm.cache.Get(id)
	return s.(Session), e
}

type DefaultSession struct {
	id string
}

func (s DefaultSession) ID() string {
	return s.id
}
