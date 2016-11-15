package goller

type Configuration struct {
	waitTimeSeconds int64
	visibilityTimeout int64
	maxNumberOfMessages int64
	region string
	queueUrl string
	accessKeyId string
	secretKey string
}

var defaultConfig = Configuration{
	waitTimeSeconds: 20,
	visibilityTimeout: 10,
	maxNumberOfMessages: 10,
	region: "us-east-1",
}

func mergeWithDefaultConfig(c *Configuration) {
	if c.region == "" {
		c.region = defaultConfig.region
	}
	if c.maxNumberOfMessages == 0 {
		c.maxNumberOfMessages = defaultConfig.maxNumberOfMessages
	}
}
