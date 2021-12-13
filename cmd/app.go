package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"notify-integrator/internal/reader"
	"notify-integrator/internal/types"
	"sync"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, event events.S3Event) {
	newSession := session.Must(session.NewSession())
	s3Reader := reader.NewS3Reader(newSession)
	var users []*types.User

	var wg sync.WaitGroup
	wg.Add(len(event.Records))

	for _, record := range event.Records {
		go func(eventRecord events.S3EventRecord) {
			key := eventRecord.S3.Object.Key
			bucket := eventRecord.S3.Bucket.Name
			bytes, err := s3Reader.Read(ctx, key, bucket)
			if err != nil {
				return
			}

			user := s3Reader.NewUser(bytes)
			users = append(users, user)
			wg.Done()
		}(record)
	}
}
