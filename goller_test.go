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

func (t *testStruct) handle(message *string) {
	fmt.Printf("%+v\n", message)
}

func TestGoller_NewSqsPoller(t *testing.T) {
	l :=  log.New(os.Stdout, "goller: ", log.Lshortfile|log.LstdFlags)
	res := NewSqsPoller("https://sqs.us-east-1.amazonaws.com/526316940316/vishals-test-queue", l)
	assert.NotEmpty(t, res, "it should not be empty")
	ts := testStruct{}
	res.RegisterHandler(ts)
	res.Poll()
}
