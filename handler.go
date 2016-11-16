package goller

//Handler is the interface that will handle all SQS messages
type Handler interface {
	Handle(message *string)
}
