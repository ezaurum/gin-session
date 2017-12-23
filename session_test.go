package session

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetNew(t *testing.T) {
	sm := Default()

	s := sm.GetNew()

	assert.NotNil(t, s)
	assert.True(t, 0 != s.ID(),"session ID cannot be 0")
}

func TestGetNewSerial(t *testing.T) {
	sm := Default()

	s := sm.GetNew()
	s0 := sm.GetNew()

	assert.NotEqual(t, s0.ID, s.ID)

}


