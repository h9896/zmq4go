package main

import (
	"fmt"

	gomq "../gomq"
)

type socketInfo struct {
	name       string
	endpoint   string
	socketType string
	flag       bool
	identity   string
}

func main() {
	pub := socketInfo{
		name:       "TryPub",
		endpoint:   "tcp://127.0.0.1:20886",
		socketType: "pub",
		flag:       true,
		identity:   "123456",
	}
	sub := socketInfo{
		name:       "TrySub",
		endpoint:   "tcp://127.0.0.1:20886",
		socketType: "Sub",
		flag:       false,
		identity:   "223456",
	}
	pubMQ, err := gomq.Set(pub.name, pub.endpoint, pub.socketType, pub.flag, pub.identity)
	if err != nil {
		fmt.Println(err)
	}
	subMQ, err := gomq.Set(sub.name, sub.endpoint, sub.socketType, sub.flag, sub.identity)
	if err != nil {
		fmt.Println(err)
	}
	b := [][]byte{[]byte("Try Send"), []byte("Try Send1")}
	pubMQ.PreData(b)
	for {
		if subMQ.Status == false {
			return
		}
		select {
		case val := <-subMQ.RecQueue:
			message := fmt.Sprintf("%s", val)
			fmt.Println(message)
		}
	}
}
