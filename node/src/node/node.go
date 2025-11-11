package node

import (
	"node/node/src/proto"
	"sync"
)

type Node struct {
	proto.UnimplementedMutualExclusionServer
	mu             sync.Mutex
	Id             int
	Address        string
	state          string // (released, wanted, held)
	timestamp      int
	queue          []int
	peers          []*Node // pointer to a node
	pendingReplies int
}

type State string

const (
	RELEASED State = "RELEASED"
	WANTED   State = "WANTED"
	HELD     State = "HELD"
)

func NewNode(id int, address string) *Node {
	n := &Node{}
	n.Id = id
	n.Address = address
	n.state = "RELEASED"
	n.timestamp = 0
	n.queue = make([]int, 0)
	n.peers = make([]*Node, 0)
	n.pendingReplies = 0
	return n
}

func (n *Node) setState(newState State) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.state = string(newState)
}

func (n *Node) AddPeer(peer *Node) {
	n.peers = append(n.peers, peer)
}
