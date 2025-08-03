package p2p

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
)

const ProtocolID = "/whispernet/1.0.0"

type DiscoveryNotifee struct{}

func (d *DiscoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Println("Discovered new peer:", pi.ID)
}

func NewNode(port int, random io.Reader) (host.Host, error) {
	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, random)
	if err != nil {
		log.Printf("Failed to generate private key: %v", err)
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	if err != nil {
		log.Printf("Failed to create multiaddr: %v", err)
		return nil, fmt.Errorf("failed to create multiaddr: %w", err)
	}
	node, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(privKey),
	)
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

func PeerConnect(node host.Host, dest string) (*bufio.ReadWriter, error) {
	log.Println("Connecting to peer:", dest)
	addr, err := multiaddr.NewMultiaddr(dest)
	if err != nil {
		log.Printf("Failed to parse multiaddr: %v", err)
		return nil, fmt.Errorf("failed to parse multiaddr: %w", err)
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		log.Printf("Failed to get peer info from address: %v", err)
		return nil, fmt.Errorf("failed to get peer info from address: %w", err)
	}

	if peerInfo.ID == node.ID() {
		log.Println("Attempting to dial self — skipping")
		return nil, nil
	}

	node.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)
	s, err := node.NewStream(context.Background(), peerInfo.ID, ProtocolID)
	if err != nil {
		log.Printf("Failed to create new stream: %v", err)
		return nil, fmt.Errorf("failed to create new stream: %w", err)
	}
	log.Println("Connected to peer:", peerInfo.ID)

	for _, addr := range peerInfo.Addrs {
		fmt.Printf(" - Address: %s\n", addr.String())
	}

	fmt.Printf("Self ID: %s\n", node.ID())

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	return rw, nil
}

func StartPeer(node host.Host, streamHandler network.StreamHandler) {
	node.SetStreamHandler(ProtocolID, streamHandler)

	var port string
	for _, la := range node.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		log.Println("was not able to find actual local port")
		return
	}
}
