package define

const (
	//  0/1
	OP_HANDSHAKE byte = byte(iota)
	OP_HANDSHAKE_REPLY

	// ping  2/3
	OP_PING
	OP_PONG

	// connect client -> client 4/5
	OP_CONE
	OP_CONE_REPLY

	// message client -> client 6/7
	OP_MSG
	OP_MSG_REPLY
)
