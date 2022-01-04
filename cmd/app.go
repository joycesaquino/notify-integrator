package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"notify-integrator/internal/client"
	"notify-integrator/internal/converter"
	"notify-integrator/internal/reader"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, event events.S3Event) error {

	newSession := session.Must(session.NewSession())
	s3Reader := reader.NewS3Reader(newSession)
	integrationClient := client.NewClient()

	for _, record := range event.Records {

		key := record.S3.Object.Key
		bucket := record.S3.Bucket.Name
		bytes, err := s3Reader.Download(ctx, key, bucket)
		if err != nil {
			return err
		}

		outputs, err := converter.ToObject(bytes)
		if err != nil {
			return err
		}
		for _, body := range outputs {
			if integrationClient.Post(ctx, body); err != nil {
				return err
			}
		}
	}

	return nil
}
