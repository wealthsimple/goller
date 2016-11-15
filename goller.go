package goller

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"reflect"
	"encoding/json"
)

//The internal structure that contains config and session information for a particular poller
type sqsQueue struct {
	client  *sqs.SQS
	logger  *log.Logger
	config  Configuration
	handler Handler
	responseObject reflect.Type
}

//Returns a new sqs poller for a given configuration and handler
func NewSqsPoller(c Configuration, h Handler, responseObject reflect.Type, l *log.Logger) *sqsQueue {
	mergeWithDefaultConfig(&c)

	sess := getSession(&c, l)

	return &sqsQueue{client: sqs.New(sess), config: c, handler: h, responseObject: responseObject, logger: l}
}

//Long polls the sqs queue (provided that the waitTimeSeconds is set in the config and > 0)
func (s *sqsQueue) Poll() {
	if s.handler == nil {
		panic("A message handler needs to be registered first!")
	}

	s.logger.Printf("Long polling on %s", s.config.queueUrl)

	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.config.queueUrl),
		WaitTimeSeconds:     aws.Int64(s.config.waitTimeSeconds),
		VisibilityTimeout:   aws.Int64(s.config.visibilityTimeout),
		MaxNumberOfMessages: aws.Int64(s.config.maxNumberOfMessages),
	}

	result, err := s.client.ReceiveMessage(params)
	checkErr(err, s.logger)

	messages := result.Messages
	for _, v := range messages {
		receipt := v.ReceiptHandle
		obj := reflect.New(s.responseObject).Interface().(Handler)
		json.Unmarshal([]byte(*v.Body), &obj)
		s.handler.Handle(obj)
		s.deleteMessage(receipt)
	}
}

//Deletes the message after long polling
func (s *sqsQueue) deleteMessage(receipt *string) {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.config.queueUrl),
		ReceiptHandle: receipt,
	}
	_, err := s.client.DeleteMessage(params)

	checkErr(err, s.logger)
}

//Gets the session based on the configuration: checks if credentials are set, otherwise, uses aws provider chain
func getSession(c *Configuration, l *log.Logger) *session.Session {
	var sess *session.Session
	var err error

	if c.accessKeyId != "" && c.secretKey != "" {
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(c.region),
			Credentials: credentials.NewStaticCredentials(c.accessKeyId, c.secretKey, ""),
		})
	} else {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(c.region)})
	}
	checkErr(err, l)
	return sess
}
