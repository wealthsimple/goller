package goller

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

type testStruct struct {
	Handler
}

func (t testStruct) Handle(message *string) {
	fmt.Printf("%+v\n", *message)
}

func TestGoller_NewSqsPoller(t *testing.T) {
	l := log.New(os.Stdout, "goller: ", log.Lshortfile|log.LstdFlags)
	ts := new(testStruct)
	config := Configuration{
		waitTimeSeconds:   20,
		visibilityTimeout: 10,
		queueUrl:/*YOUR QUEUE URL HERE */ "test_url.com",
	}
	res := NewSqsPoller(config, ts, l)
	assert.NotEmpty(t, res, "it should not be empty")
	res.Poll()
}
