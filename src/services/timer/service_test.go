package timer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func initSrv() *Service {
	srv := Service{}
	srv.Init(nil)

	return &srv
}

func timerFunc() {
	fmt.Println(time.Now().UTC().Format(time.RFC3339Nano))
}

func TestCreateTimer(t *testing.T) {
	srv := initSrv()
	defer srv.Release()
	var err error

	timerName := "t1"
	err = srv.CreateTimer(timerName, timerFunc)
	assert.Nil(t, err)
	assert.Equal(t, len(srv.timers), 1)

	timer, err := srv.getTimerByName(timerName)
	assert.Nil(t, err)
	assert.Equal(t, timer.itv, 1*time.Second)

	timer2Name := "t2"
	err = srv.CreateTimer(timer2Name, timerFunc, 200*time.Millisecond)
	assert.Nil(t, err)
	assert.Equal(t, len(srv.timers), 2)

	time.Sleep(1 * time.Second)
}

func TestRemoveTimer(t *testing.T) {
	srv := initSrv()
	defer srv.Release()
	var err error

	timerName := "t1"
	err = srv.RemoveTimer(timerName)
	assert.NotNil(t, err)

	err = srv.CreateTimer(timerName, timerFunc)
	assert.Nil(t, err)
	assert.Equal(t, len(srv.timers), 1)

	err = srv.RemoveTimer(timerName)
	assert.Nil(t, err)
	assert.Equal(t, len(srv.timers), 0)
}
