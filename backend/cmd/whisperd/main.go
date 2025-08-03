package main

import (
	"crypto/rand"
	"flag"
	"fmt"

	"github.com/Sahas001/whispernet/internal/chat"
	"github.com/Sahas001/whispernet/internal/p2p"
)

func main() {
	sourcePort := flag.Int("s", 0, "source port number")
	dest := flag.String("d", "", "destination multiaddr to connect to")
	flag.Parse()

	fmt.Println("Starting whisperd daemon on port:", sourcePort)
	host, err := p2p.NewNode(*sourcePort, rand.Reader)
	defer host.Close()
	if err != nil {
		fmt.Println("Error creating node:", err)
		return
	}
	if *dest == "" {
		fmt.Println("No destination provided, starting discovery service...")
		p2p.StartPeer(host, chat.HandleStream)

	} else {
		rw, err := p2p.PeerConnect(host, *dest)
		if err != nil {
			fmt.Println("Error connecting to peer:", err)
			return
		}
		go chat.WriteData(rw)
		go chat.ReadData(rw)
	}

	select {} // Keep the main function running
}
