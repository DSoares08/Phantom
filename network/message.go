package network

type GetStatusMessage struct {}

type StatusMessage struct {
	// Server id
	ID string
	Version uint32
	CurrentHeight uint32
}