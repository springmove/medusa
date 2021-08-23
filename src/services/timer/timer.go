package timer

import (
	"time"

	"github.com/linshenqi/medusa/src/services/base"
)

// param1： handler func
// param2： timer duration
func CreateTimer(handler base.TimerHander, itv ...time.Duration) *Timer {
	timer := Timer{
		done:    make(chan bool, 1),
		handler: handler,
		itv:     1 * time.Second,
	}

	if len(itv) > 0 {
		timer.itv = itv[0]
	}

	go timer.asyncWorker()
	return &timer
}

type Timer struct {
	done    chan bool
	handler base.TimerHander
	itv     time.Duration
}

func (s *Timer) Release() {
	s.done <- true
}

func (s *Timer) asyncWorker() {
	ticker := time.NewTicker(s.itv)

	for {
		go s.handler()

		select {
		case <-ticker.C:
		case <-s.done:
			return
		}
	}
}
