package main

import (
	"log"
	"context"
	"github.com/satori/go.uuid"
	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	log.Println(event)
	switch event.RequestType {
		case cfn.RequestCreate: {
			log.Println("create")
			physicalResourceID = uuid.NewV4().String()
		}
		default: {
			log.Println("update or delete")
			physicalResourceID = event.PhysicalResourceID
		}
	}
	log.Println("sending : " + physicalResourceID)
	data = map[string]interface{} {}
	return physicalResourceID, data, nil
}

func main() {
	lambda.Start(cfn.LambdaWrap(handler))
}