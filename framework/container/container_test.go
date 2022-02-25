package container

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"testing"
)

type testServer struct {
	Name string
}

type testProvider struct {
}

func (t *testProvider) Build(container contract.Container, params ...interface{}) (interface{}, error) {
	return &testServer{"test"}, nil
}

func TestContainer_Bind_MAKE(t *testing.T) {
	a := NewContainer()
	a.Bind("test", &testProvider{})
	server, _ := a.Get("test")
	server1, _ := a.Get("test")
	server2, _ := a.Make("test")

	assert.IsType(t, server, &testServer{})
	assert.IsType(t, server1, &testServer{})
	assert.IsType(t, server2, &testServer{})
	assert.Equal(t, true, server == server1)
	assert.Equal(t, false, server == server2)

	server3, err := a.Get("test2")
	assert.Equal(t, nil, server3)
	assert.Equal(t, errors.New("No provider for test2"), err)
}

func TestContainer_Set(t *testing.T) {
	a := NewContainer()
	ts := &testServer{}
	a.Set("test", ts)
	server, _ := a.Get("test")
	assert.Equal(t, true, ts == server)
}

func TestContainer_Alias(t *testing.T) {
	a := NewContainer()
	a.Bind("test", &testProvider{})
	a.Alias("test", "t")
	server, _ := a.Get("test")
	server1, _ := a.Get("t")
	assert.Equal(t, true, server == server1)
}

func TestContainer_GetNameByAlias(t *testing.T) {
	a := NewContainer()
	a.Alias("test", "t")
	assert.Equal(t, "test", a.GetNameByAlias("t"))
}
