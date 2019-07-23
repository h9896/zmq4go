package jsunzmq

import (
	"errors"
	"fmt"
	"strings"

	zmq "github.com/pebbe/zmq4"
)

const (
	recThread  = 1
	sendThread = 2
)

// JsSocket includes socket, receive chan, and send chan
type JsSocket struct {
	socket       *zmq.Socket
	closeChannel bool
	RecQueue     chan [][]byte
	SendQueue    chan [][]byte
}

//var recQueue = make([]chan [][]byte,recThread)
//var sendQueue = make([]chan [][]byte,sendThread)

// SetZmq is used to set the socket information.
func SetZmq(typ string, endpoint string, bind bool, subscrib string) (JsSocket, error) {
	var zmqType zmq.Type
	send := make(chan [][]byte)
	rec := make(chan [][]byte)
	switch t := strings.ToUpper(typ); t {
	case "PUB":
		zmqType = zmq.PUB
		rec = nil
	case "SUB":
		zmqType = zmq.SUB
		send = nil
	default:
		info := JsSocket{socket: nil, RecQueue: nil, SendQueue: nil}
		return info, errors.New("Not Pub or Sub type")
	}
	sock, err := zmq.NewSocket(zmqType)
	if err != nil {
		fmt.Println(err)
	}
	if zmqType == zmq.SUB {
		sock.SetSubscribe(subscrib)
	}
	if bind {
		err = sock.Bind(endpoint)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err = sock.Connect(endpoint)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = sock.SetLinger(-1)
	if err != nil {
		fmt.Println(err)
	}
	info := JsSocket{socket: sock, RecQueue: rec, SendQueue: send, closeChannel: false}
	return info, err
}

// Start receive or send data
func (z JsSocket) Start() {
	if z.RecQueue != nil {
		go func() {
			defer z.socket.Close()
			//defer func() {
			//	if e := recover(); e != nil {
			//		fmt.Println(e)
			//	}
			//	z.socket.Close()
			//}()
			z.recData()
		}()
	}
	if z.SendQueue != nil {
		go func() {
			defer z.socket.Close()
			//defer func() {
			//	if e := recover(); e != nil {
			//		fmt.Println(e)
			//	}
			//	z.socket.Close()
			//}()
			z.sendData()
		}()
	}
}

// PreData is used to prepare send data
func (z JsSocket) PreData(data [][]byte) {
	z.SendQueue <- data
}

// CloseSocket is used to close socket
func (z JsSocket) CloseSocket() {
	z.closeChannel = true
}

func (z JsSocket) sendData() {
	for {
		if z.closeChannel {
			break
		}
		select {
		case val := <-z.SendQueue:
			_, err := z.socket.SendMessage(val)
			if err != nil {
				fmt.Println(err)
				//panic(err)
				break
			}
		}
	}
}

func (z JsSocket) recData() {
	for {
		if z.closeChannel {
			break
		}
		item, err := z.socket.RecvMessageBytes(0)
		if err != nil {
			fmt.Println(err)
			//panic(err)
			break
		}
		z.RecQueue <- item
	}
}
