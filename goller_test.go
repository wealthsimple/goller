package goller

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)


type testStruct struct {
	Handler
}
func TestGoller_NewSqsPoller(t *testing.T) {
	l := &log.Logger{}
	res := NewSqsPoller("https://sqs.us-east-1.amazonaws.com/526316940316/vishals-test-queue", l)
	assert.NotEmpty(t, res, "it should not be empty")
	ts := testStruct{}
	res.RegisterHandler(ts)
	res.Poll()
}
