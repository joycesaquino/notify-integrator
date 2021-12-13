package reader

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"notify-integrator/internal/types"
	"os"
	"path/filepath"
)

type S3Reader struct {
	dwl *s3manager.Downloader
}

func (sr S3Reader) NewUser(bytes []byte) *types.User {
	var user *types.User
	err := json.Unmarshal(bytes, &user)
	if err != nil {
		return nil
	}
	return user
}

func (sr S3Reader) Read(ctx context.Context, key string, bucketName string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(bucketName),
	}

	temp, err := os.Create(filepath.Join(os.TempDir(), "temp.json"))
	if err != nil {
		return nil, err
	}

	_, err = sr.dwl.DownloadWithContext(ctx, temp, input)
	if err != nil {
		return nil, err

	}
	bytes, err := ioutil.ReadAll(temp)
	if err != nil {
		return nil, err
	}

	return bytes, nil

}

func NewS3Reader(session *session.Session) *S3Reader {

	dwl := s3manager.NewDownloaderWithClient(s3.New(session), func(downloader *s3manager.Downloader) {
		downloader.PartSize = 10 * 1024 * 1024
		downloader.Concurrency = 4
	})

	return &S3Reader{
		dwl: dwl,
	}
}
