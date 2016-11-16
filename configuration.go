package goller

type Configuration struct {
	WaitTimeSeconds     int64
	VisibilityTimeout   int64
	MaxNumberOfMessages int64
	Region              string
	QueueUrl            string
	AccessKeyId         string
	SecretKey           string
	provider            awsProvider
}

var defaultConfig = Configuration{
	WaitTimeSeconds:     20,
	VisibilityTimeout:   10,
	MaxNumberOfMessages: 10,
	Region:              "us-east-1",
	provider:            &defaultProvider{},
}

func mergeWithDefaultConfig(c *Configuration) {
	if c.Region == "" {
		c.Region = defaultConfig.Region
	}
	if c.MaxNumberOfMessages == 0 {
		c.MaxNumberOfMessages = defaultConfig.MaxNumberOfMessages
	}
	if c.provider == nil {
		c.provider = defaultConfig.provider
	}
}
