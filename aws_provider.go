package goller

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type awsProvider interface {
	getSession(region string) (*session.Session, error)
	getSessionWithCredentials(region string, accessKey string, secretKey string) (*session.Session, error)
	getQueue(*session.Session) (*sqs.SQS)
	receiveMessages(*sqs.ReceiveMessageInput, *sqs.SQS) (*sqs.ReceiveMessageOutput, error)
	deleteMessage(*sqs.DeleteMessageInput, *sqs.SQS) (error)
}

type defaultProvider struct {
	awsProvider
}

func (d *defaultProvider) getSession(region string) (*session.Session, error) {
	return session.NewSession(&aws.Config{Region: aws.String(region)})
}

func (d *defaultProvider) getSessionWithCredentials(region string, accessKey string, secretKey string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
}

func (d *defaultProvider) getQueue(sess *session.Session) (*sqs.SQS) {
	return sqs.New(sess)
}

func (d *defaultProvider) receiveMessages(params *sqs.ReceiveMessageInput, client *sqs.SQS) (*sqs.ReceiveMessageOutput, error) {
	return client.ReceiveMessage(params)
}

func (d*defaultProvider) deleteMessage(params *sqs.DeleteMessageInput, client *sqs.SQS) (error) {
	_, err := client.DeleteMessage(params)
	return err
}
