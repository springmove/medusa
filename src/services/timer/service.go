package timer

import (
	"fmt"
	"sync"
	"time"

	"github.com/linshenqi/medusa/src/services/base"
	"github.com/linshenqi/sptty"
)

type Service struct {
	sptty.BaseService
	base.ITimerService

	mtx    sync.Mutex
	timers map[string]*Timer
}

func (s *Service) ServiceName() string {
	return base.ServiceTimer
}

func (s *Service) Init(app sptty.ISptty) error {
	s.mtx = sync.Mutex{}
	s.timers = map[string]*Timer{}

	return nil
}

func (s *Service) Release() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, v := range s.timers {
		v.Release()
	}
}

func (s *Service) getTimerByName(dispatcherName string) (*Timer, error) {

	timer, exist := s.timers[dispatcherName]
	if !exist {
		return nil, fmt.Errorf("Not Found")
	}

	return timer, nil
}

func (s *Service) CreateTimer(timerName string, handler base.TimerHander, itv ...time.Duration) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, err := s.getTimerByName(timerName)
	if err == nil {
		return err
	}

	s.timers[timerName] = CreateTimer(handler, itv...)
	return nil
}

func (s *Service) RemoveTimer(timerName string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	timer, err := s.getTimerByName(timerName)
	if err != nil {
		return err
	}

	timer.Release()
	delete(s.timers, timerName)

	return nil
}
