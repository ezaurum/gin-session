package memstore

import (
	"github.com/patrickmn/go-cache"
	"time"
	"github.com/ezaurum/gin-session"
	"github.com/ezaurum/gin-session/generators/snowflake"
)

var _ session.Store = memorySessionStore{}

func Default() session.Store {
	c := cache.New(5*time.Minute, 10*time.Minute)
	k := snowflake.New(0)

	return memorySessionStore{
		keyGenerator:k,
		cache:c,
	}
}

type memorySessionStore struct {
	cache     *cache.Cache
	keyGenerator session.IDGenerator
}

func (sm memorySessionStore) GetNew() session.Session {
	id := sm.keyGenerator.Generate()
	s := DefaultSession{
		id:id,
	}

	sm.cache.Set(id, s, cache.DefaultExpiration)

	return s
}

func (sm memorySessionStore) Get(id string) (session.Session, bool) {
	s, e := sm.cache.Get(id)
	return s.(session.Session), e
}

type DefaultSession struct {
	id string
}

func (s DefaultSession) ID() string {
	return s.id
}
