// Package awsutils provides some helper function for common aws task.
package awsutils

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

/*Mock stuff*/
type mockedS3Client struct {
	s3iface.S3API
}

func (s *mockedS3Client) ListObjectsV2(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	key := "someKey"
	contents := []*s3.Object{&s3.Object{Key: &key}}
	return &s3.ListObjectsV2Output{Contents: contents}, nil
}

func TestDownloadEmptyBucket(t *testing.T) {

	b := Bucket{}
	err := b.DownloadBucket(nil)

	if err.Error() != messageClientNotDefined {
		t.Errorf("Expected error :%s, and got %s", messageClientNotDefined, err.Error())
	}
	b = NewBucket(&mockedS3Client{}, "Bucket", "temp")

	err = b.DownloadBucket(nil)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUploadEmptyBucket(t *testing.T) {

	b := Bucket{}
	err := b.UploadBucket()

	if err.Error() != messageClientNotDefined {
		t.Errorf("Expected error :%s, and got %s", messageClientNotDefined, err.Error())
	}
	b = NewBucket(&mockedS3Client{}, "Bucket", "NotADir")

	err = b.UploadBucket()
	if err != nil {
		t.Errorf(err.Error())
	}
}
