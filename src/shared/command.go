package shared

// Command Types for transfered messages
const (
	// CommandPing - ping
	CommandPing = iota
	// CommandAuth - auth req
	CommandAuth
	// CommandSpace - subdoamin req
	CommandSpace
	// CommandRequest - request pack
	CommandRequest
	// CommandResponse - response pack
	CommandResponse
	// CommandOK - ok response
	CommandOK
	// CommandFail - failed response
	CommandFail
)
