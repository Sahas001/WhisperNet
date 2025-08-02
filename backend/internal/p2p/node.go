package p2p

import (
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type DiscoveryNotifee struct{}

func (d *DiscoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Println("Discovered new peer:", pi.ID)
}

func NewNode() (host.Host, error) {
	node, err := libp2p.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p node: %w", err)
	}
	fmt.Println("Node created with ID:", node.ID())
	fmt.Println("Listening on addresses:")
	for _, addr := range node.Addrs() {
		fmt.Println(" -", addr)
	}

	s := discovery.NewMdnsService(node, "whisperd-mDNS", &DiscoveryNotifee{})

	if err := s.Start(); err != nil {
		log.Printf("Failed to start mDNS discovery service: %v", err)
	}
	return node, nil
}
