package goller

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type sqsQueue struct {
	client   *sqs.SQS
	logger   *log.Logger
	config	 Configuration
	handler  Handler
}

var defaultConfig = Configuration{
	waitTimeSeconds: 20,
	visibilityTimeout: 10,
	maxNumberOfMessages: 10,
	region: "us-east-1",
}

func NewSqsPoller(c Configuration, h Handler, l *log.Logger) *sqsQueue {
	mergeWithDefaults(&c)

	sess, err := session.NewSession(&aws.Config{Region: aws.String(c.region)})
	checkErr(err, l)

	return &sqsQueue{client: sqs.New(sess), config: c, handler: h, logger: l}
}

func (s *sqsQueue) Poll() {
	if s.handler == nil {
		panic("A message handler needs to be registered first!")
	}

	s.logger.Printf("Long polling on %s", s.config.queueUrl)

	params := &sqs.ReceiveMessageInput{
		QueueUrl:        aws.String(s.config.queueUrl),
		WaitTimeSeconds: aws.Int64(s.config.waitTimeSeconds),
		VisibilityTimeout: aws.Int64(s.config.visibilityTimeout),
		MaxNumberOfMessages: aws.Int64(s.config.maxNumberOfMessages),
	}

	result, err := s.client.ReceiveMessage(params)
	checkErr(err, s.logger)

	messages := result.Messages
	for _, v := range messages {
		receipt := v.ReceiptHandle
		s.handler.Handle(v.Body)
		s.deleteMessage(receipt)
	}
}

func (s *sqsQueue) deleteMessage(receipt *string) {
	params := &sqs.DeleteMessageInput{
		QueueUrl: aws.String(s.config.queueUrl),
		ReceiptHandle: receipt,
	}
	_, err := s.client.DeleteMessage(params)

	checkErr(err, s.logger)
}

func mergeWithDefaults(c *Configuration) {
	if c.region == "" {
		c.region = defaultConfig.region
	}
	if c.maxNumberOfMessages == 0 {
		c.maxNumberOfMessages = defaultConfig.maxNumberOfMessages
	}
}
