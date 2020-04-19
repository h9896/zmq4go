package zmq4go

// Ilog - a interface about record log.
type Ilog interface {
	Record(rec string)
}
