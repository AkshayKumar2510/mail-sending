package main

import (
	"fmt"
	"os"

	mail_send "mail-sending"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func main() {
	region := os.Getenv("AWS_REGION")

	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	s3Instance := s3.New(sess)
	sessInstance := session.Must(session.NewSession())
	ssmInstance := ssm.New(sess, aws.NewConfig().WithRegion(region))
	mailHandler := mail_send.MailHandler{
		S3Instance:      s3Instance,
		SessionInstance: sessInstance,
		SSMInstance:     ssmInstance,
	}

	lambda.Start(mailHandler.Handler)
}
