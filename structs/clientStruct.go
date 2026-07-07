package extras

import (
	"bytes"
	"encoding/gob"
)

type SessionState struct {
	PacketListener            chan Packet
	AuthenticationProcessDone bool
	ConnectionWrapper         Connection
	Encoder                   *gob.Encoder
	Decoder                   *gob.Decoder
	Buffer                    *bytes.Buffer
}
