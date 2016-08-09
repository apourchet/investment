package broadcaster

import (
	"math/rand"
)

// TODO Sync safe
type Broadcaster struct {
	receivers map[ReceiverID]chan interface{}
}

type ReceiverID int

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{make(map[ReceiverID]chan interface{})}
}

func (b *Broadcaster) Register(cb chan interface{}) ReceiverID {
	rid := ReceiverID(rand.Intn(10000))
	b.receivers[rid] = cb
	return rid
}

func (b *Broadcaster) Deregister(rid ReceiverID) {
	delete(b.receivers, rid)
}

func (b *Broadcaster) Emit(data interface{}) {
	for _, rec := range b.receivers {
		go func(cb chan interface{}) { cb <- data }(rec)
	}
}
