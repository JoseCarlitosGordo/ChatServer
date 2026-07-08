package extras

import (
	"bufio"
	"encoding/gob"
)

type SessionState struct {
	PacketListener            chan Packet
	AuthenticationProcessDone bool
	ConnectionWrapper         Connection
	Encoder                   *gob.Encoder
	Decoder                   *gob.Decoder
	InputScanner              *bufio.Scanner
}
