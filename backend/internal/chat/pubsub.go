package chat

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/network"
)

func HandleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go ReadData(rw)
	go WriteData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func ReadData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green: \x1b[32m
			// Reset: \x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func WriteData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(rw, "%s\n", sendData)
		rw.Flush()
	}
}
