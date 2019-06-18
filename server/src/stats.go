package server

// Stats - simple stats
type Stats struct {
	Connected int
	Pending   int
	Done      int
	Failed    int
	Timeout   int
}

func (stats *Stats) connected() {
	(*stats).Connected++
}

func (stats *Stats) disconnected() {
	(*stats).Connected--
}

func (stats *Stats) complete() {
	(*stats).Done++
	(*stats).Pending--
}

func (stats *Stats) start() {
	(*stats).Pending++
}

func (stats *Stats) fail() {
	(*stats).Pending--
	(*stats).Failed++
}

func (stats *Stats) timeout() {
	(*stats).Pending--
	(*stats).Timeout++
}
