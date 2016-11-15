package goller

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"reflect"
)

type testStruct struct {
	Handler
	field string
}

func (t testStruct) Handle(message Handler) {
	fmt.Printf("%+v\n", message)
}

func TestGoller_NewSqsPoller(t *testing.T) {
	l := log.New(os.Stdout, "goller: ", log.Lshortfile|log.LstdFlags)
	ts := new(testStruct)
	config := Configuration{
		waitTimeSeconds:   20,
		visibilityTimeout: 10,
		queueUrl:          "https://sqs.us-east-1.amazonaws.com/526316940316/vishals-test-queue",
	}
	res := NewSqsPoller(config, ts, reflect.TypeOf(testStruct{}), l)
	assert.NotEmpty(t, res, "it should not be empty")
	res.Poll()
}
