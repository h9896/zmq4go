package zmq4go

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	zmq "github.com/pebbe/zmq4/draft"
)

var disconnectTopic = "MsgDisconnect"
var log Ilog

// GoSocket - information about socket.
type GoSocket struct {
	socket       *zmq.Socket
	RecQueue     chan [][]byte
	sendQueue    chan [][]byte
	closeChannel chan string
	Status       bool
}

// Set - create a Gosocket.
func Set(zmqName string, endPoint string, zmqType string, flag bool, identity string) (GoSocket, error) {
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
		return GoSocket{socket: nil, RecQueue: nil, sendQueue: nil, closeChannel: make(chan string), Status: false}, errors.New("Set Error: Not support this type of zmq")
	}
	sock, err := zmq.NewSocket(form)
	if err != nil {
		return GoSocket{socket: nil, RecQueue: nil, sendQueue: nil, closeChannel: make(chan string), Status: false}, fmt.Errorf("NewSocket Error: %s", err)
	}
	err = sock.SetIdentity(identity)
	if err != nil {
		return GoSocket{socket: nil, RecQueue: nil, sendQueue: nil, closeChannel: make(chan string), Status: false}, fmt.Errorf("SetIdentity Error: %s", err)
	}
	if !flag {
		err = sock.Connect(endPoint)
		if err != nil {
			return GoSocket{socket: nil, RecQueue: nil, sendQueue: nil, closeChannel: make(chan string), Status: false}, fmt.Errorf("Connect Error: %s", err)
		}
	} else {
		err = sock.Bind(endPoint)
		if err != nil {
			return GoSocket{socket: nil, RecQueue: nil, sendQueue: nil, closeChannel: make(chan string), Status: false}, fmt.Errorf("Bind Error: %s", err)
		}
	}
	err = sock.SetLinger(-1)
	if err != nil {
		return GoSocket{socket: nil, RecQueue: nil, sendQueue: nil, closeChannel: make(chan string), Status: false}, fmt.Errorf("SetLinger Error: %s", err)
	}
	return GoSocket{socket: sock, RecQueue: make(chan [][]byte), sendQueue: make(chan [][]byte), closeChannel: make(chan string), Status: true}, err
}
func (z GoSocket) sendData() error {
	var mutex sync.Mutex
	for {
		select {
		case val := <-z.sendQueue:
			_, err := z.socket.SendMessage(val)
			if err != nil {
				log.Record(fmt.Sprintf("Sanding Error: %s; Data: %s", err, val))
				return fmt.Errorf("Sanding Error: %s; Data: %s", err, val)
			}
			for i, b := range val {
				log.Record(fmt.Sprintf("Send Len: %d, Data: %s", i, string(b)))
			}
		case val := <-z.closeChannel:
			if strings.Split(val, ".")[2] == "MsgDisconnect" {
				mutex.Lock()
				defer mutex.Unlock()
				_, err := z.socket.SendMessage(val)
				if err != nil {
					log.Record(fmt.Sprintf("Sanding Error: %s; Data: %s", err, val))
					return fmt.Errorf("Sanding Error: %s; Data: %s", err, val)
				}
				close(z.sendQueue)
				z.Status = false
				return nil
			}
		}
	}
}

func (z GoSocket) recData() error {
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
			log.Record(fmt.Sprintf("Receive Error: %s; Data Missing", err))
			return fmt.Errorf("Receive Error: %s; Data Missing", err)
		}
		z.RecQueue <- item
		for i, b := range item {
			log.Record(fmt.Sprintf("Receive Len: %d, Data: %s", i, string(b)))
		}
	}
}

// Start - start receive or send data
func (z GoSocket) Start() {
	go func() {
		err := z.recData()
		if err != nil {
			z.socket.Close()
			log.Record(fmt.Sprintf("recData Error: %s", err))
		}
		close(z.RecQueue)
		z.socket.Close()
		z.Status = false
	}()
	go func() {
		err := z.sendData()
		if err != nil {
			z.socket.Close()
			log.Record(fmt.Sprintf("sendData Error: %s", err))
		}
		close(z.closeChannel)
		z.Status = false
	}()
}

// PreData - prepare data to send
func (z GoSocket) PreData(data [][]byte) {
	defer func() {
		if recover() != nil {
			z.Status = false
			id, _ := z.socket.GetIdentity()
			log.Record(fmt.Sprintf("SendQueue channel %s already close", id))
			z.closeSocket(disconnectTopic)
		}
	}()
	z.sendQueue <- data
	z.Status = true
}
func (z GoSocket) closeSocket(val string) {
	defer func() {
		if recover() != nil {
			z.Status = false
		}
	}()
	z.closeChannel <- val
	z.Status = true
}
