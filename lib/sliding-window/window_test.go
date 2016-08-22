package slidwin_test

import (
	. "github.com/apourchet/investment/lib/sliding-window"
	"testing"
)

func TestBasicSlidingWindow(t *testing.T) {
	s := NewSlidingWindow(5)
	s.Push(0)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	for i, d := range s {
		if i != 4-d.(int) {
			t.Fatal()
		}
	}
	s.Push(5)
	if s[0] != 5 {
		t.Fatal()
	}
	if s[4] != 1 {
		t.Fatal()
	}
}
