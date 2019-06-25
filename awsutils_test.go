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
	return &s3.ListObjectsV2Output{Contents: []*s3.Object{}}, nil
}

func TestDownloadEmptyBucket(t *testing.T) {

	err := DownloadBucket(nil, "", "", nil)

	if err.Error() != messageClientNotDefined {
		t.Errorf("Expected error :%s, and got %s", messageClientNotDefined, err.Error())
	}

	mock := &mockedS3Client{}
	err = DownloadBucket(mock, "Dir", "Bucket", nil)
	if err != nil {
		t.Errorf(err.Error())
	}

}
