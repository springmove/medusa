package dispatcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initSrv() *Service {
	srv := Service{}
	_ = srv.Init(nil)

	return &srv
}

func TestCreateDispatcher(t *testing.T) {

	srv := initSrv()
	defer srv.Release()
	var err error

	err = srv.CreateDispatcher("d1")
	assert.Nil(t, err)
	assert.Equal(t, len(srv.dispatchers), 1)

	err = srv.CreateDispatcher("d1")
	assert.Nil(t, err)
	assert.Equal(t, len(srv.dispatchers), 1)

	err = srv.CreateDispatcher("d2")
	assert.Nil(t, err)
	assert.Equal(t, len(srv.dispatchers), 2)

	err = srv.CreateDispatcher("buffersize", 1024)
	assert.Nil(t, err)
	assert.Equal(t, len(srv.dispatchers), 3)
	dBuffersize, err := srv.getDispatcherByName("buffersize")
	assert.Nil(t, err)
	assert.Equal(t, cap(dBuffersize.buf), 1024)

	err = srv.CreateDispatcher("buffersize_worker", 2048, 4)
	assert.Nil(t, err)
	assert.Equal(t, len(srv.dispatchers), 4)
	dBuffersizeWorker, err := srv.getDispatcherByName("buffersize_worker")
	assert.Nil(t, err)
	assert.Equal(t, cap(dBuffersizeWorker.buf), 2048)
	assert.Equal(t, len(dBuffersizeWorker.workers), 4)
}

func TestRemoveDispatcher(t *testing.T) {
	srv := initSrv()
	defer srv.Release()
	var err error

	err = srv.CreateDispatcher("d1")
	assert.Nil(t, err)
	assert.Equal(t, len(srv.dispatchers), 1)

	err = srv.RemoveDispatcher("d2")
	assert.NotNil(t, err)
	assert.Equal(t, len(srv.dispatchers), 1)

	err = srv.RemoveDispatcher("d1")
	assert.Nil(t, err)
	assert.Equal(t, len(srv.dispatchers), 0)
}

func handler(i int) {
	fmt.Println(i)
}

func TestAddDispatcherHandler(t *testing.T) {
	srv := initSrv()
	defer srv.Release()
	var err error

	dispatcherName := "d1"
	err = srv.AddDispatcherHandler(dispatcherName, handler)
	assert.NotNil(t, err)

	err = srv.CreateDispatcher(dispatcherName)
	assert.Nil(t, err)
	assert.Equal(t, len(srv.dispatchers), 1)

	err = srv.AddDispatcherHandler(dispatcherName, handler)
	assert.Nil(t, err)
	dispatcher, err := srv.getDispatcherByName(dispatcherName)
	assert.Nil(t, err)
	assert.Equal(t, len(dispatcher.handlers), 1)

	err = srv.AddDispatcherHandler(dispatcherName, handler, "h2")
	assert.Nil(t, err)
	dispatcher, err = srv.getDispatcherByName(dispatcherName)
	assert.Nil(t, err)
	assert.Equal(t, len(dispatcher.handlers), 2)
}

func TestRemoveDispatcherHandler(t *testing.T) {
	srv := initSrv()
	defer srv.Release()
	var err error

	dispatcherName := "d1"
	handlerName := "h1"
	err = srv.RemoveDispatcherHandler(dispatcherName, handlerName)
	assert.NotNil(t, err)

	err = srv.CreateDispatcher(dispatcherName)
	assert.Nil(t, err)
	assert.Equal(t, len(srv.dispatchers), 1)

	err = srv.AddDispatcherHandler(dispatcherName, handler)
	assert.Nil(t, err)
	dispatcher, err := srv.getDispatcherByName(dispatcherName)
	assert.Nil(t, err)
	assert.Equal(t, len(dispatcher.handlers), 1)

	err = srv.RemoveDispatcherHandler(dispatcherName, handlerName)
	assert.Nil(t, err)
	dispatcher, err = srv.getDispatcherByName(dispatcherName)
	assert.Nil(t, err)
	assert.Equal(t, len(dispatcher.handlers), 1)

	err = srv.AddDispatcherHandler(dispatcherName, handler, handlerName)
	assert.Nil(t, err)
	dispatcher, err = srv.getDispatcherByName(dispatcherName)
	assert.Nil(t, err)
	assert.Equal(t, len(dispatcher.handlers), 2)

	err = srv.RemoveDispatcherHandler(dispatcherName, handlerName)
	assert.Nil(t, err)
	dispatcher, err = srv.getDispatcherByName(dispatcherName)
	assert.Nil(t, err)
	assert.Equal(t, len(dispatcher.handlers), 1)
}

func TestDispatch(t *testing.T) {
	srv := initSrv()
	defer srv.Release()
	var err error

	dispatcherName := "d1"

	err = srv.Dispatch(dispatcherName, 1)
	assert.NotNil(t, err)

	err = srv.CreateDispatcher(dispatcherName)
	assert.Nil(t, err)
	assert.Equal(t, len(srv.dispatchers), 1)

	err = srv.AddDispatcherHandler(dispatcherName, handler)
	assert.Nil(t, err)
	dispatcher, err := srv.getDispatcherByName(dispatcherName)
	assert.Nil(t, err)
	assert.Equal(t, len(dispatcher.handlers), 1)

	err = srv.Dispatch("d1", 1)
	assert.Nil(t, err)
}
