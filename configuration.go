package goller

type Configuration struct {
	waitTimeSeconds int64
	visibilityTimeout int64
	maxNumberOfMessages int64
	region string
	queueUrl string
}
