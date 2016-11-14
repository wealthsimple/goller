package goller

type Handler interface {
	Handle(message *string)
}
