package session

import "github.com/bwmarrin/snowflake"

type Manager interface {
	GetNew() Session
}

type Session interface {
	ID() int64
}

var _ Manager = DefaultSessionManager{}
var _ Session = DefaultSession{}

type DefaultSessionManager struct {
	generator *snowflake.Node
}

func Default() Manager {
	node, err := snowflake.NewNode(0)
	if nil != err {
		panic(err)
	}
	return DefaultSessionManager{
		generator:node,
	}
}

func (manager DefaultSessionManager) GetNew() Session {
	id := manager.generator.Generate().Int64()
	return DefaultSession{
		id:id,
	}
}

type DefaultSession struct {
	id int64
}

func (s DefaultSession) ID() int64 {
	return s.id
}
