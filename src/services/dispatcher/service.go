package dispatcher

import (
	"fmt"
	"sync"

	"github.com/linshenqi/medusa/src/services/base"
	"github.com/linshenqi/sptty"
)

const (
	ServiceName = "dispatcher"
)

type Service struct {
	sptty.BaseService
	base.IDispatcher

	mtx         sync.Mutex
	dispatchers map[string]*base.Dispatcher
}

func (s *Service) ServiceName() string {
	s.mtx = sync.Mutex{}
	s.dispatchers = map[string]*base.Dispatcher{}

	return ServiceName
}

func (s *Service) Init(app sptty.Sptty) error {
	return nil
}

func (s *Service) Release() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, v := range s.dispatchers {
		v.Release()
	}
}

func (s *Service) getDispatcherByName(dispatcherName string) (*base.Dispatcher, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	dispatcher, exist := s.dispatchers[dispatcherName]
	if !exist {
		return nil, fmt.Errorf("Not Found")
	}

	return dispatcher, nil
}

func (s *Service) CreateDispatcher(dispatcherName string, params ...uint) error {
	_, err := s.getDispatcherByName(dispatcherName)
	if err == nil {
		return nil
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.dispatchers[dispatcherName] = base.CreateDispatcher(params...)

	return nil
}

func (s *Service) RemoveDispatcher(dispatcherName string) {
	dispatch, err := s.getDispatcherByName(dispatcherName)
	if err != nil {
		return
	}

	dispatch.Release()

	s.mtx.Lock()
	defer s.mtx.Unlock()
	delete(s.dispatchers, dispatcherName)
}

func (s *Service) Dispatch(dispatcherName string, data interface{}) error {
	dispatch, err := s.getDispatcherByName(dispatcherName)
	if err != nil {
		return err
	}

	dispatch.Dispatch(data)
	return nil
}

func (s *Service) AddHandler(dispatcherName string, handler interface{}, handlerName ...string) error {
	dispatch, err := s.getDispatcherByName(dispatcherName)
	if err != nil {
		return err
	}

	dispatch.AddHandler(handler, handlerName...)
	return nil
}
