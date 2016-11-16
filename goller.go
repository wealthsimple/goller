package goller

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//The internal structure that contains config and session information for a particular poller
type sqsQueue struct {
	client  *sqs.SQS
	logger  *log.Logger
	config  Configuration
	handler Handler
}

//NewSqsPoller returns a new sqs poller for a given configuration and handler
func NewSqsPoller(c Configuration, h Handler, l *log.Logger) *sqsQueue {
	mergeWithDefaultConfig(&c)

	sess := getSession(&c, l)

	return &sqsQueue{client: c.provider.getQueue(sess), config: c, handler: h, logger: l}
}

//Poll long polls the sqs queue (provided that the WaitTimeSeonds is set in the config and > 0)
func (s *sqsQueue) Poll() {
	if s.handler == nil {
		panic("A message handler needs to be registered first!")
	}

	s.logger.Printf("Long polling on %s\n", s.config.QueueUrl)

	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.config.QueueUrl),
		WaitTimeSeconds:     aws.Int64(s.config.WaitTimeSeconds),
		VisibilityTimeout:   aws.Int64(s.config.VisibilityTimeout),
		MaxNumberOfMessages: aws.Int64(s.config.MaxNumberOfMessages),
	}

	result, err := s.config.provider.receiveMessages(params, s.client)
	checkErr(err, s.logger)

	messages := result.Messages
	for _, v := range messages {
		receipt := v.ReceiptHandle
		s.handler.Handle(v.Body)
		s.deleteMessage(receipt)
	}

	s.logger.Printf("Finished long polling after %d seconds", s.config.WaitTimeSeconds)
}

//Deletes the message after long polling
func (s *sqsQueue) deleteMessage(receipt *string) {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.config.QueueUrl),
		ReceiptHandle: receipt,
	}
	err := s.config.provider.deleteMessage(params, s.client)

	checkErr(err, s.logger)
}

//Gets the session based on the configuration: checks if credentials are set, otherwise, uses aws provider chain
func getSession(c *Configuration, l *log.Logger) *session.Session {
	var sess *session.Session
	var err error

	if c.AccessKeyId != "" && c.SecretKey != "" {
		sess, err = c.provider.getSessionWithCredentials(c.Region, c.AccessKeyId, c.SecretKey)
	} else {
		sess, err = c.provider.getSession(c.Region)
	}

	checkErr(err, l)
	return sess
}
