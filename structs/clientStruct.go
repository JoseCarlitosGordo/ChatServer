package extras

type SessionState struct {
	PacketListener            chan Packet[any]
	AuthenticationProcessDone bool
	ConnectionWrapper         Connection
}
