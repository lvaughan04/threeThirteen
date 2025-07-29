package interfaces

type Event interface {
	GetType() string
	GetGameID() string
}