package dispatcher

import (
	"fmt"
	"sync"

	"github.com/linshenqi/medusa/src/services/base"
	"github.com/linshenqi/sptty"
)

type Service struct {
	sptty.BaseService
	base.IDispatcherService

	mtx         sync.Mutex
	dispatchers map[string]*Dispatcher
}

func (s *Service) ServiceName() string {
	return base.ServiceDispatcher
}

func (s *Service) Init(app sptty.ISptty) error {
	s.mtx = sync.Mutex{}
	s.dispatchers = map[string]*Dispatcher{}

	return nil
}

func (s *Service) Release() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, v := range s.dispatchers {
		v.Release()
	}
}

func (s *Service) getDispatcherByName(dispatcherName string) (*Dispatcher, error) {

	dispatcher, exist := s.dispatchers[dispatcherName]
	if !exist {
		return nil, fmt.Errorf("Not Found")
	}

	return dispatcher, nil
}

func (s *Service) CreateDispatcher(dispatcherName string, params ...uint) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, err := s.getDispatcherByName(dispatcherName)
	if err == nil {
		return nil
	}

	s.dispatchers[dispatcherName] = CreateDispatcher(params...)

	return nil
}

func (s *Service) RemoveDispatcher(dispatcherName string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	dispatcher, err := s.getDispatcherByName(dispatcherName)
	if err != nil {
		return err
	}

	dispatcher.Release()
	delete(s.dispatchers, dispatcherName)

	return nil
}

func (s *Service) Dispatch(dispatcherName string, data interface{}) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	dispatcher, err := s.getDispatcherByName(dispatcherName)
	if err != nil {
		return err
	}

	dispatcher.Dispatch(data)
	return nil
}

func (s *Service) AddDispatcherHandler(dispatcherName string, handler base.DispatcherHander, handlerName ...string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	dispatcher, err := s.getDispatcherByName(dispatcherName)
	if err != nil {
		return err
	}

	dispatcher.AddHandler(handler, handlerName...)
	return nil
}

func (s *Service) RemoveDispatcherHandler(dispatcherName string, handlerName string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	dispatcher, err := s.getDispatcherByName(dispatcherName)
	if err != nil {
		return err
	}

	dispatcher.RemoveHandler(handlerName)
	return nil
}
