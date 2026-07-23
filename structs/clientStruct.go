package extras

import (
	"bufio"
	"encoding/json"
)

type SessionState struct {
	PacketListener            chan Packet
	AuthenticationProcessDone bool
	ConnectionWrapper         Connection
	Encoder                   *json.Encoder
	Decoder                   *json.Decoder
	InputScanner              *bufio.Scanner
}
