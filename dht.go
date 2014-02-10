package main

import (
	//"fmt"
	"net"
	"os"
	//"time"
)

const (
	BITS           = 160 // sha1
	WORKER_THREADS = 1
)

type DHT struct {
	self          *Node
	finger        []*Node
	predecessor   *Node
	successor     *Node
	globalInbound chan *Message
}

func NewDHT(self *Node) *DHT {
	dht := DHT{
		self:          self,
		finger:        make([]*Node, BITS),
		predecessor:   nil,
		successor:     self,
		globalInbound: make(chan *Message, 100),
	}

	for i, _ := range dht.finger {
		dht.finger[i] = self
	}

	// TODO: implement worker number
	// also: make threadsafe, which _isn't_
	for i := 0; i < WORKER_THREADS; i++ {
		go dht.Worker()
	}

	return &dht
}

func (d *DHT) Store(object []byte) error {
	return nil
}

func (d *DHT) Retrieve(id int64) ([]byte, error) {
	return nil, nil
}

func (d *DHT) Join(node *Node) {
	err := node.Connect(d.globalInbound)

	if err != nil {
		panic(err)
	}

  node.SendPing()

	successor, err := node.GetSuccessor(d.self) // blocks

	if err != nil {
		panic(err)
	}

	println(successor)

}

func (d *DHT) Listen() {
	sock, err := net.Listen("tcp", d.self.Address())

	if err != nil {
		println("Error listening:", err)
		os.Exit(1)
	}

	println("Listening on", d.self.Address())

	for {
		conn, err := sock.Accept()

		if err != nil {
			println("Error accepting!")
		}

		node := NewNode(conn.RemoteAddr().String())
		node.Accept(conn, d.globalInbound)
	}
}

func (d *DHT) Worker() {
	println("dht worker started")

	for {
		m := <-d.globalInbound

    println(m.String())

    switch m.Intent {
    case REQUEST_SUCCESSOR:
    case REQUEST_PING:
      println("LOL")
      m.Sender.ReplyPing()
    }
	}
}

//
//
//
//
//

func (d *DHT) findSuccessor(node *Node) (*Node, error) {
	if node.Id().elementOf(d.self.Id(), d.successor.Id()) { // this interval is (]
		println("oix")
		return d.successor, nil
	} else {
		println("oix2")
		queryNode := d.closestPrecedingNode(node)
		resultNode, err := queryNode.GetSuccessor(node)

		if err != nil {
			return nil, err
		}

		return resultNode, nil
	}
}

func (d *DHT) closestPrecedingNode(node *Node) *Node {
	for i := BITS; i > 0; i-- {
		if d.finger[i].Id().elementOf(d.self.Id(), node.Id()) { // this interval is ()
			return d.finger[i]
		}
	}

	return d.self
}

func (d *DHT) stabilize() {
}

func (d DHT) notify(node Node) {
}

func (d *DHT) fixFingers() {
}
