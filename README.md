# goller 
[![CircleCI](https://circleci.com/gh/wealthsimple/goller.svg?style=svg)](https://circleci.com/gh/wealthsimple/goller) [![GoDoc](https://godoc.org/github.com/wealthsimple/goller?status.svg)](https://godoc.org/github.com/wealthsimple/goller) 

A Golang AWS SQS long poller.

## Installation
```
go get github.com/wealthsimple/goller
```

## Usage
```go
import (
	"fmt"
	"github.com/wealthsimple/goller"
	"log"
	"os"
)

type testStruct struct {
	goller.Handler
}

func (t testStruct) Handle(message *string) {
	fmt.Printf("%+v\n", *message)
}

func main() {
	l := log.New(os.Stdout, "goller: ", log.Lshortfile|log.LstdFlags)
	config := goller.Configuration{
		WaitTimeSeconds:     20,
		VisibilityTimeout:   10,
		MaxNumberOfMessages: 1,
		Region:              "us-east-1",
		QueueURL:            os.Getenv("QUEUE_URL"),
		AccessKeyID:         os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:           os.Getenv("AWS_SECRET_KEY"),
	}

	t := new(testStruct)
	sqs := goller.NewSqsPoller(config, t, l)
	sqs.Poll()
}
```

* See goller_integration_test.go for more details
