package manager

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	m := newLimiter(100, time.Second)
	z := time.Second / 100
	for i := 0; i < 100; i++ {
		time.Sleep(z)
		if !m.addNewTime() {
			t.Fatalf("%d", i)
		}
	}
	if m.addNewTime() {
		t.Fatalf("%d", 101)
	}
	m.check()
	t.Fail()
}
