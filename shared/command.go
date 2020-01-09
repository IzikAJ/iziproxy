package shared

// Command Types for transfered messages
const (
	// - SERVICE COMMANDS:
	// CommandPing - ping request
	CommandPing = iota
	// CommandPong - ping response
	CommandPong

	// - CONTEXT COMMANDS:
	// CommandSetup - send credentials + other setup data
	CommandSetup
	// CommandReady - if server connection ready
	CommandReady
	// CommandFailed - response if previous command invalid/failed
	CommandFailed

	// - TRANSFER COMMANDS:
	// CommandRequest - request pack
	CommandRequest
	// CommandResponse - response pack
	CommandResponse
)
