package base

import "time"

const (
	ServiceTimer = "timer"
)

type TimerHander func()

// loop task based on duration trigger
type ITimerService interface {
	CreateTimer(timerName string, handler TimerHander, itv ...time.Duration) error
	RemoveTimer(timerName string) error
}
