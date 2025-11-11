package node

import (
	"context"
	"log"
	"time"

	proto "node/node/src/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Request struct {
	senderID  int
	timestamp int
}

// gRPC server metoder
func (n *Node) RequestAccess(ctx context.Context, req *proto.RequestMessage) (*proto.Ack, error) {
	log.Printf("Node %d received request from Node %d with timestamp %d", n.Id, req.Id, req.Timestamp)
	receiveRequest(n, int(req.Id), int(req.Timestamp))
	return &proto.Ack{}, nil
}

func (n *Node) ReplyAccess(ctx context.Context, reply *proto.ReplyMessage) (*proto.Ack, error) {
	log.Printf("Node %d received reply from Node %d", n.Id, reply.Id)
	receiveReply(n, int(reply.Id), 0)
	return &proto.Ack{}, nil
}

// Ricart Agrawala algoritme
func SendRequest(n *Node) {
	n.mu.Lock()
	n.state = "WANTED"
	n.timestamp = int(time.Now().Unix()-baseTimestamp)*10 + n.Id
	n.pendingReplies = len(n.peers)
	myTimestamp := n.timestamp
	n.mu.Unlock()

	log.Printf("None %d want to enter CS with timestamp %d", n.Id, myTimestamp)

	for _, peer := range n.peers {
		go func(peer *Node) {
			conn, _ := grpc.Dial(peer.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			defer conn.Close()

			client := proto.NewMutualExclusionClient(conn)
			client.RequestAccess(context.Background(), &proto.RequestMessage{
				Id:        int32(n.Id),
				Timestamp: int32(myTimestamp),
			})
			log.Printf("Node %d sent request to Node %d with timestamp %d", n.Id, peer.Id, peer.timestamp)
		}(peer)
	}

	for {
		n.mu.Lock()
		if n.pendingReplies == 0 {
			n.state = "HELD"
			n.mu.Unlock()
			enterCS(n)
			exitCS(n)
			break
		}
		n.mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func receiveRequest(n *Node, senderID int, senderTS int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.state == "HELD" {
		log.Printf("Node %d sets Node %d in queue", n.Id, senderID)
		n.queue = append(n.queue, senderID)
		return
	}

	if n.state == "WANTED" {
		if n.timestamp < senderTS || (n.timestamp == senderTS && n.Id < senderID) {
			log.Printf("Node %d has the priority and sets Node %d in queue", n.Id, senderID)
			n.queue = append(n.queue, senderID)
			return
		}
	}
	log.Printf("Node %d send reply to Node %d", n.Id, senderID)
	go reply(n, senderID)
}

func receiveReply(n *Node, senderID int, senderTimestamp int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.pendingReplies--
}

func enterCS(n *Node) {
	log.Printf("Node %d entering CS\n", n.Id)
	time.Sleep(2 * time.Second)
	log.Printf("Node %d exit CS\n", n.Id)
}

func exitCS(n *Node) {
	n.mu.Lock()
	n.state = "RELEASED"
	n.mu.Unlock()

	for _, waiterId := range n.queue {
		reply(n, waiterId)
	}
	n.mu.Lock()
	n.queue = n.queue[:0]
	n.mu.Unlock()
}

func reply(n *Node, recieverID int) {
	var target *Node
	for _, peer := range n.peers {
		if peer.Id == recieverID {
			target = peer
			break
		}
	}

	if target == nil {
		return
	}

	// gPRC connection
	conn, _ := grpc.NewClient(target.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := proto.NewMutualExclusionClient(conn)
	client.ReplyAccess(context.Background(), &proto.ReplyMessage{
		Id: int32(n.Id),
	})
	log.Printf("Node %d sent reply to Node %d", n.Id, recieverID)
}
