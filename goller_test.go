package goller

import (
	"testing"
	"log"
	"os"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	 "github.com/stretchr/testify/assert"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type dummyProvider struct {
	awsProvider
	calledGetSession bool
	calledGetWithCreds bool
	calledGetQueue bool
	calledReceivedMessages bool
	calledDelete bool
	output *sqs.ReceiveMessageOutput
	sess *session.Session
}

func (d *dummyProvider) getSession(region string) (*session.Session, error) {
	d.calledGetSession = true
	return d.sess, nil
}

func (d *dummyProvider) getSessionWithCredentials(region string, accessKey string, secretKey string) (*session.Session, error) {
	d.calledGetWithCreds = true
	return d.sess, nil
}

func (d *dummyProvider) getQueue(*session.Session) (*sqs.SQS) {
	d.calledGetQueue = true
	return nil
}

func (d *dummyProvider) receiveMessages(params *sqs.ReceiveMessageInput, client *sqs.SQS) (*sqs.ReceiveMessageOutput, error) {
	d.calledReceivedMessages = true
	return d.output, nil
}

func (d *dummyProvider) deleteMessage(params *sqs.DeleteMessageInput, client *sqs.SQS) (error) {
	d.calledDelete = true
	return nil
}

type dummyHandler struct {
	Handler
	called bool
}

func (d *dummyHandler) Handle(message *string) {
	d.called = true
	fmt.Printf("%+v\n", *message)
}


var l *log.Logger

func init() {
	l = log.New(os.Stdout, "goller: ", log.Lshortfile|log.LstdFlags)
}

func TestNewSqsPoller_WithoutCredentials(t *testing.T) {
	p := &dummyProvider{}

	config := Configuration{
		WaitTimeSeconds:   20,
		VisibilityTimeout: 10,
		QueueUrl:          "my_url.com",
		provider:    	p,
	}

	h := new(dummyHandler)

	poller := NewSqsPoller(config, h, l)

	assert.NotNil(t, poller, "The queue should not be nil")
	assert.True(t, p.calledGetSession, "Should have called get session")
	assert.False(t, p.calledGetWithCreds, "Should not have called get session with creds")
	assert.True(t, p.calledGetQueue, "Should have called get queue")
}

func TestNewSqsPoller_WithCredentials(t *testing.T) {
	p := &dummyProvider{}

	config := Configuration{
		WaitTimeSeconds:   20,
		VisibilityTimeout: 10,
		QueueUrl:          "my_url.com",
		provider:    	p,
		AccessKeyId: "akid",
		SecretKey: "secretKey",
	}

	h := new(dummyHandler)

	poller := NewSqsPoller(config, h, l)

	assert.NotNil(t, poller, "The queue should not be nil")
	assert.False(t, p.calledGetSession, "Should not have called get session")
	assert.True(t, p.calledGetWithCreds, "Should have called get session with creds")
	assert.True(t, p.calledGetQueue, "Should have called get queue")
}

func TestSqsQueue_Poll_WithoutHandler(t *testing.T) {
	p := &dummyProvider{}

	config := Configuration{
		WaitTimeSeconds:   20,
		VisibilityTimeout: 10,
		QueueUrl:          "my_url.com",
		provider:    	p,
		AccessKeyId: "akid",
		SecretKey: "secretKey",
	}

	assert.Panics(t, func() {
		NewSqsPoller(config, nil, l).Poll()
	}, "No handler should panic")
}


func TestSqsQueue_PollWithNoMessages(t *testing.T) {
	p := &dummyProvider{output: new(sqs.ReceiveMessageOutput)}

	config := Configuration{
		WaitTimeSeconds:   20,
		VisibilityTimeout: 10,
		QueueUrl:          "my_url.com",
		provider:    	p,
		AccessKeyId: "akid",
		SecretKey: "secretKey",
	}

	h := new(dummyHandler)

	poller := NewSqsPoller(config, h, l)

	poller.Poll()

	assert.False(t, p.calledDelete, "It should not have called delete")
	assert.False(t, h.called, "It should not have called the handler")
}

func TestSqsQueue_Poll(t *testing.T) {
	messages := make([]*sqs.Message, 1)
	receiptHandle := "myReceiptHandle"
	receiptBody := "myMessageBody"
	messages[0] = &sqs.Message{
		ReceiptHandle: &receiptHandle,
		Body: &receiptBody,
	}
	output := &sqs.ReceiveMessageOutput{
		Messages: messages,
	}

	p := &dummyProvider{output: output}

	config := Configuration{
		WaitTimeSeconds:   20,
		VisibilityTimeout: 10,
		QueueUrl:          "my_url.com",
		provider:    	p,
		AccessKeyId: "akid",
		SecretKey: "secretKey",
	}

	h := new(dummyHandler)

	poller := NewSqsPoller(config, h, l)

	poller.Poll()

	assert.True(t, p.calledDelete, "It should have called delete")
	assert.True(t, h.called, "It should have called the handler")
}
