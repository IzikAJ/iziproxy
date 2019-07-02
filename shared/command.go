package shared

// Command Types for transfered messages
const (
	// - SERVICE COMMANDS:
	// CommandPing - ping request
	CommandPing = iota
	// CommandPong - ping response
	CommandPong

	// - CONTEXT COMMANDS:
	// CommandReady - if server connection ready
	CommandReady
	// CommandAuth - ask/send credentials
	CommandAuth
	// CommandScope - allow/request scope
	CommandScope
	// CommandFailed - response if previous command invalid/failed
	CommandFailed

	// - TRANSFER COMMANDS:
	// CommandRequest - request pack
	CommandRequest
	// CommandResponse - response pack
	CommandResponse
)
