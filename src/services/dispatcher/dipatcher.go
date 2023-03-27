package dispatcher

import (
	"reflect"
	"sync"

	"github.com/springmove/medusa/src/base"
	"github.com/springmove/sptty"
)

// param1： buffer size
// param2： worker num
func CreateDispatcher(params ...uint) *Dispatcher {

	vals := []uint{base.DefaultBufSize, 1}

	for i := 0; i < len(params); i++ {
		vals[i] = params[i]
	}

	dispatcher := Dispatcher{
		buf:      make(chan interface{}, vals[0]),
		handlers: map[string]base.DispatcherHander{},
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
	handlers map[string]base.DispatcherHander
}

func (s *Dispatcher) doDispatch(data interface{}) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

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

func (s *Dispatcher) AddHandler(handler base.DispatcherHander, handlerName ...string) {
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
