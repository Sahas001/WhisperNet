package main

import (
	"fmt"

	"github.com/Sahas001/whispernet/internal/p2p"
)

func main() {
	host, err := p2p.NewNode()
	defer host.Close()
	if err != nil {
		fmt.Println("Error creating node:", err)
		return
	}
	fmt.Println("Starting whisperd daemon...")

	select {} // Keep the main function running
}
