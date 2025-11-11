package main

import (
	"log"
	"net"
	node "node/node/src/node"
	"node/node/src/proto"
	"time"

	"google.golang.org/grpc"
)

// Starter en node - Server + Client

func main() {
	// Create notes
	n1 := node.NewNode(1, "localhost:5001")
	n2 := node.NewNode(2, "localhost:5002")
	n3 := node.NewNode(3, "localhost:5003")

	// Giving the nodes peers
	n1.AddPeer(n2)
	n1.AddPeer(n3)
	n2.AddPeer(n1)
	n2.AddPeer(n3)
	n3.AddPeer(n1)
	n3.AddPeer(n2)

	log.Printf("Node 1: %s", n1.Address)
	log.Printf("Node 2: %s", n2.Address)
	log.Printf("Node 3: %s", n3.Address)

	// start gRPC server
	go startgRPC(n1)
	go startgRPC(n2)
	go startgRPC(n3)

	time.Sleep(2 * time.Second)

	log.Printf("Tests:")
	log.Printf("===========================================================================")
	// Tests if a node can request cs
	log.Printf("Test 1: One node can request CS")
	go node.SendRequest(n1)

	time.Sleep(5 * time.Second)

	// Two nodes can request cs at the same time
	log.Printf("===========================================================================")
	log.Printf("Tests 2: Two nodes can request CS at the same time")
	go node.SendRequest(n2)
	time.Sleep(10 * time.Second)
	go node.SendRequest(n3)

	// 3 nodes can request cs
	log.Printf("===========================================================================")
	log.Printf("\n Test 3: All nodes can request CS")
	go node.SendRequest(n1)
	go node.SendRequest(n2)
	go node.SendRequest(n3)

	select {}
}

func startgRPC(n *node.Node) {
	lis, err := net.Listen("tcp", n.Address)
	if err != nil {
		log.Fatalf("Node %d could failed to listen: %v", n.Id, err)
	}

	s := grpc.NewServer()
	// registers the node as gRPC service
	proto.RegisterMutualExclusionServer(s, n)

	log.Printf("Node %d starting gRPC Server on %s", n.Id, n.Address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Node %d failed to serve: %v", n.Id, err)
	}
}
