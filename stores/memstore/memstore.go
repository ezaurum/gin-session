package memstore

import (
	"github.com/patrickmn/go-cache"
	"time"
	"github.com/ezaurum/gin-session"
	"github.com/ezaurum/gin-session/generators/snowflake"
	"github.com/ezaurum/gin-session/generators"
)

var _ session.Store = memorySessionStore{}

func Default() session.Store {
	k := snowflake.New(0)
	return New(k)
}

func New(k generators.IDGenerator) session.Store {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return memorySessionStore{
		keyGenerator:k,
		cache:c,
	}
}


type memorySessionStore struct {
	cache     *cache.Cache
	keyGenerator generators.IDGenerator
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
	if e {
		return s.(session.Session),true
	}

	return nil, false
}

func (sm memorySessionStore) Count() int {
	return sm.cache.ItemCount()
}

func (sm memorySessionStore) Sessions() session.SessionMap {
	 m := make(session.SessionMap)
	for k, v := range sm.cache.Items() {
		m[k] = v.Object.(session.Session)
	}
	return m
}

type DefaultSession struct {
	id string
}

func (s DefaultSession) ID() string {
	return s.id
}

