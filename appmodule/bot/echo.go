package bot

// EchoRequest defines bot request.
type EchoRequest struct {
	Payload string `form:"payload"`
}
