package float_heap

import "container/heap"

type FloatHeap interface {
	heap.Interface
	Init()
	DoPush(float64)
	DoPop() float64
}

type minHeap []float64
type maxHeap []float64

const (
	HEAP_MIN = 1.
	HEAP_MAX = -1.
)

func NewFloatHeap(direction int) FloatHeap {
	var fh FloatHeap
	if direction == HEAP_MAX {
		fh = &maxHeap{}
	} else {
		fh = &minHeap{}
	}
	heap.Init(fh)
	return fh
}

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i]-h[j] > 0 }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(float64))
}

func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *maxHeap) DoPush(f float64) {
	heap.Push(h, f)
}

func (h *maxHeap) DoPop() float64 {
	return heap.Pop(h).(float64)
}

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i]-h[j] < 0 }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(float64))
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *minHeap) DoPush(f float64) {
	heap.Push(h, f)
}

func (h *minHeap) DoPop() float64 {
	return heap.Pop(h).(float64)
}

func (fh *maxHeap) Init() {
	heap.Init(fh)
}

func (fh *minHeap) Init() {
	heap.Init(fh)
}
