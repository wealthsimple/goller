package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wealthsimple/goller"
)

//READ ME FIRST:
//
//This is an integration test, as such it requires a real AWS queue url
//and the appropriate credentials to access that queue. Goller gets its aws credentials
//from the Provider chain. As a last resort, you can provide the AWS access key id and
//secret key to the config object; however, this is not recommended in production

type testStruct struct {
	goller.Handler
}

func (t testStruct) Handle(message *string) {
	fmt.Printf("%+v\n", *message)
}

func main() {
	l := log.New(os.Stdout, "goller: ", log.Lshortfile|log.LstdFlags)
	ts := new(testStruct)
	config := goller.Configuration{
		WaitTimeSeconds:   20,
		VisibilityTimeout: 10,
		QueueUrl:          /* Provide your queue url here */ "PROVIDE_YOUR_URL.com",
	}
	res := goller.NewSqsPoller(config, ts, l)
	res.Poll()
}
