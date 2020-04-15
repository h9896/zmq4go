package zmq4go

import (
	"fmt"
	"time"
)

func Record(msg string) {
	t := time.Now()
	fmt.Printf("%s->%s", t.Format("2006/01/02-15:04:05.000"), msg)
}
