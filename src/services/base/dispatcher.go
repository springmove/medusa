package base

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/linshenqi/sptty"
)

type IDispatcher interface {
	CreateDispatcher(dispatcherName string, params ...uint) error
	RemoveDispatcher(dispatcherName string)
	Dispatch(dispatcherName string, data interface{}) error
	AddHandler(dispatcherName string, handler interface{}, handlerName ...string) error
}

// 参数1： 缓冲长度
// 参数2： worker数量
func CreateDispatcher(params ...uint) *Dispatcher {

	vals := []uint{DefaultBufSize, 1}

	for i := 0; i < len(params); i++ {
		vals[i] = params[i]
	}

	dispatcher := Dispatcher{
		buf:      make(chan interface{}, vals[0]),
		handlers: map[string]interface{}{},
		mtx:      sync.Mutex{},
	}

	for i := 0; i < int(vals[1]); i++ {
		dispatcher.workers = append(dispatcher.workers, createDispatcherWorker(&dispatcher))
	}

	return &dispatcher
}

func createDispatcherWorker(dispatcher *Dispatcher) *dispatcherWorker {
	worker := dispatcherWorker{
		dispatcher: dispatcher,
		done:       make(chan bool, 1),
	}

	go worker.asyncWorker()
	return &worker
}

type dispatcherWorker struct {
	dispatcher *Dispatcher
	done       chan bool
}

func (s *dispatcherWorker) release() {
	s.done <- true
}

func (s *dispatcherWorker) asyncWorker() {
	for {
		select {
		case data := <-s.dispatcher.buf:
			s.dispatcher.doDispatch(data)

		case <-s.done:
			return
		}
	}
}

type Dispatcher struct {
	mtx      sync.Mutex
	buf      chan interface{}
	workers  []*dispatcherWorker
	handlers map[string]interface{}
}

func (s *Dispatcher) doDispatch(data interface{}) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	fmt.Println(data)
	for _, v := range s.handlers {
		f := reflect.ValueOf(v)
		go f.Call([]reflect.Value{
			reflect.ValueOf(data),
		})
	}
}

func (s *Dispatcher) Release() {
	for _, v := range s.workers {
		v.release()
	}
}

func (s *Dispatcher) Dispatch(data interface{}) {
	s.buf <- data
}

func (s *Dispatcher) AddHandler(handler interface{}, handlerName ...string) {
	name := sptty.GenerateUID()
	if len(handlerName) > 0 {
		name = handlerName[0]
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.handlers[name] = handler
}

func (s *Dispatcher) RemoveHandler(handlerName string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, exist := s.handlers[handlerName]
	if exist {
		delete(s.handlers, handlerName)
	}
}
