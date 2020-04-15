package zmq4go

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	zmq "github.com/pebbe/zmq4/draft"
)

type Igomq interface {
}

var disconnectTopic = "MsgDisconnect"

type goSocket struct {
	socket       *zmq.Socket
	recQueue     chan [][]byte
	sendQueue    chan [][]byte
	closeChannel chan string
	status       bool
}

func Set(zmqName string, endPoint string, zmqType string, flag bool, identity string) (goSocket, error) {
	var form zmq.Type
	switch typ := strings.ToUpper(zmqType); typ {
	case "DEALER":
		form = zmq.DEALER
	case "ROUTER":
		form = zmq.ROUTER
	case "PUB":
		form = zmq.PUB
	case "SUB":
		form = zmq.SUB
	default:
		return goSocket{socket: nil, recQueue: nil, sendQueue: nil, closeChannel: make(chan string), status: false}, errors.New("Set Error: Not support this type of zmq")
	}
	sock, err := zmq.NewSocket(form)
	if err != nil {
		Record(fmt.Sprintf("NewSocket Error: %s", err))
	}
	err = sock.SetIdentity(identity)
	if err != nil {
		Record(fmt.Sprintf("SetIdentity Error: %s", err))
	}
	if !flag {
		err = sock.Connect(endPoint)
		if err != nil {
			Record(fmt.Sprintf("Connect Error: %s", err))
		}
	} else {
		err = sock.Bind(endPoint)
		if err != nil {
			Record(fmt.Sprintf("Bind Error: %s", err))
		}
	}
	err = sock.SetLinger(-1)
	if err != nil {
		Record(fmt.Sprintf("SetLinger Error: %s", err))
	}
	return goSocket{socket: sock, recQueue: make(chan [][]byte), sendQueue: make(chan [][]byte), closeChannel: make(chan string), status: true}, err
}
func (z goSocket) sendData() error {
	var mutex sync.Mutex
	for {
		select {
		case val := <-z.sendQueue:
			_, err := z.socket.SendMessage(val)
			if err != nil {
				Record(fmt.Sprintf("Sanding Error: %s; Data: %s", err, val))
				return fmt.Errorf("Sanding Error: %s; Data: %s", err, val)
			}
			for i, b := range val {
				Record(fmt.Sprintf("Send Len: %d, Data: %s", i, string(b)))
			}
		case val := <-z.closeChannel:
			if strings.Split(val, ".")[2] == "MsgDisconnect" {
				mutex.Lock()
				defer mutex.Unlock()
				_, err := z.socket.SendMessage(val)
				if err != nil {
					Record(fmt.Sprintf("Sanding Error: %s; Data: %s", err, val))
					return fmt.Errorf("Sanding Error: %s; Data: %s", err, val)
				}
				close(z.sendQueue)
				z.status = false
				return nil
			}
		}
	}
}

func (z goSocket) recData() error {
	var mutex sync.Mutex
	for {
		item, err := z.socket.RecvMessageBytes(0)
		if len(item) >= 1 {
			if (strings.Split(string(item[0]), ".")[2]) == "MsgDisconnect" {
				mutex.Lock()
				defer mutex.Unlock()
				return nil
			}
		}
		if err != nil {
			Record(fmt.Sprintf("Receive Error: %s; Data Missing", err))
			return fmt.Errorf("Receive Error: %s; Data Missing", err)
		}
		z.recQueue <- item
		for i, b := range item {
			Record(fmt.Sprintf("Receive Len: %d, Data: %s", i, string(b)))
		}
	}
}
func (z goSocket) Start() {
	go func() {
		err := z.recData()
		if err != nil {
			z.socket.Close()
			Record(fmt.Sprintf("recData Error: %s", err))
		}
		close(z.recQueue)
		z.socket.Close()
		z.status = false
	}()
	go func() {
		err := z.sendData()
		if err != nil {
			z.socket.Close()
			Record(fmt.Sprintf("sendData Error: %s", err))
		}
		close(z.closeChannel)
		z.status = false
	}()
}
func (z goSocket) PreData(data [][]byte) {
	defer func() {
		if recover() != nil {
			z.status = false
			id, _ := z.socket.GetIdentity()
			Record(fmt.Sprintf("SendQueue channel %s already close", id))
			z.closeSocket(disconnectTopic)
		}
	}()
	z.sendQueue <- data
	z.status = true
}
func (z goSocket) closeSocket(val string) {
	defer func() {
		if recover() != nil {
			z.status = false
		}
	}()
	z.closeChannel <- val
	z.status = true
}
