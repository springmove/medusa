package base

import (
	"fmt"
	"testing"
	"time"
)

func tp(a int, i ...int) {
	fmt.Println(i)

	if len(i) > 0 {
		fmt.Println(i[0])
	}
}

func TestBase(t *testing.T) {
	tp(23, 12)
}

type P struct {
	A int
	B string
}

func onData(p *P) {
	fmt.Println(p.A)
}

func TestDispatcher(t *testing.T) {
	d := CreateDispatcher()
	d.AddHandler(onData)

	for i := 0; i < 100; i++ {
		d.Dispatch(&P{
			A: i,
		})
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
