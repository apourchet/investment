package float_heap_test

import (
	"testing"

	"math"

	"github.com/apourchet/investment/lib/float-heap"
)

const EPSILON = 0.000001

func assertDoubleEqual(t *testing.T, a, b float64) {
	if math.Abs(a-b) > EPSILON {
		t.Logf("%f != %f", a, b)
		t.Fail()
	}
}

func TestBasic(t *testing.T) {
	h := float_heap.NewFloatHeap(float_heap.HEAP_MAX)
	h.DoPush(1.)
	h.DoPush(2.)
	h.DoPush(0.)

	assertDoubleEqual(t, h.DoPop(), 2.)
	assertDoubleEqual(t, h.DoPop(), 1.)
	assertDoubleEqual(t, h.DoPop(), 0.)

	h = float_heap.NewFloatHeap(float_heap.HEAP_MIN)
	h.DoPush(1.)
	h.DoPush(2.)
	h.DoPush(0.)

	assertDoubleEqual(t, h.DoPop(), 0.)
	assertDoubleEqual(t, h.DoPop(), 1.)
	assertDoubleEqual(t, h.DoPop(), 2.)
}
