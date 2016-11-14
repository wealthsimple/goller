package goller

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

type sqsQueue struct {
	client   *sqs.SQS
	logger   *log.Logger
	handler  Handler
	queueUrl string
}

func NewSqsPollerWithRegion(queueUrl string, region string, l *log.Logger) *sqsQueue {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	checkErr(err, l)

	return &sqsQueue{client: sqs.New(sess), logger: l, queueUrl: queueUrl}
}

func NewSqsPoller(queueName string, l *log.Logger) *sqsQueue {
	return NewSqsPollerWithRegion(queueName, "us-east-1", l)
}

func (s *sqsQueue) RegisterHandler(h Handler) {
	if s.handler == nil {
		s.handler = h
	} else {
		panic("There is already a message handler registed to this class!")
	}
}

func (s *sqsQueue) Poll() {
	if s.handler == nil {
		panic("A message handler needs to be registered first!")
	}

	s.logger.Printf("Polling on %s", s.queueUrl)

	params := &sqs.ReceiveMessageInput{
		QueueUrl:        aws.String(s.queueUrl),
		WaitTimeSeconds: aws.Int64(20),
		VisibilityTimeout: aws.Int64(10),
		MaxNumberOfMessages: aws.Int64(10),
	}

	result, err := s.client.ReceiveMessage(params)
	checkErr(err, s.logger)

	fmt.Printf("%+v\n", result)

	messages := result.Messages
	for _, v := range messages {
		receipt := v.ReceiptHandle
		s.handler.Handle(v.Body)
		s.deleteMessage(receipt)
	}
}

func (s *sqsQueue)deleteMessage(receipt *string) {
	params := &sqs.DeleteMessageInput{
		QueueUrl: aws.String(s.queueUrl),
		ReceiptHandle: receipt,
	}
	_, err := s.client.DeleteMessage(params)

	checkErr(err, s.logger)
}
