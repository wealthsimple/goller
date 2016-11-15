package goller

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"os"
	"fmt"
)


type testStruct struct {
	Handler
}

func (t testStruct) Handle(message *string) {
	fmt.Printf(*message + "\n")
}

func TestGoller_NewSqsPoller(t *testing.T) {
	l :=  log.New(os.Stdout, "goller: ", log.Lshortfile|log.LstdFlags)
	ts := testStruct{}
	config := Configuration{
		waitTimeSeconds: 20,
		visibilityTimeout: 10,
		queueUrl: "https://sqs.us-east-1.amazonaws.com/526316940316/vishals-test-queue",

	}
	res := NewSqsPoller(config, ts, l)
	assert.NotEmpty(t, res, "it should not be empty")
	res.Poll()
}
