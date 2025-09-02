package common

type MessageType int

const (
	BetsMessage MessageType = iota
	DoneMessage
	WinnersMessage
)

