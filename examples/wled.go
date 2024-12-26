package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/jwetzell/ddp-go"
)

// receive data from WLED virtual output
func main() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:4048")
	if err != nil {
		slog.Error("error making UDP address", "err", err)
		return
	}

	client, err := net.ListenUDP("udp", addr)
	if err != nil {
		slog.Error("error listening to UDP", "err", err)
		return
	}
	defer client.Close()

	for {
		buffer := make([]byte, 2048)

		length, _, err := client.ReadFromUDP(buffer)
		if err != nil {
			slog.Error("error reading from UDP", "err", err)
		} else if length > 0 {
			message, err := ddp.Decode(buffer[0:length])
			if err != nil {
				slog.Error("error decoding", "err", err)
			}
			ledCount := uint32(message.Header.DataLength) / uint32(message.Header.DataType.BitsPerPixel/8) / 3

			for i := 0; i < int(ledCount); i++ {
				ledOffset := i * 3
				ledR := message.Data[ledOffset+0]
				ledG := message.Data[ledOffset+1]
				ledB := message.Data[ledOffset+2]

				fmt.Printf("led %d:\tr=%d\tg=%d\tb=%d\n", i, ledR, ledG, ledB)
			}
		}
	}
}
