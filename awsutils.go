// Package awsutils provides some helper function for common aws task.
package awsutils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//DownloadBucket ...
func DownloadBucket(s3Client s3iface.S3API, baseDir, bucket string, excludePatten *string) error {
	var wg sync.WaitGroup

	if s3Client == nil {
		return fmt.Errorf(messageClientNotDefined)
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	result, err := s3Client.ListObjectsV2(input)
	if err != nil {
		return err
	}

	for _, s3Obj := range result.Contents {
		if excludePatten != nil {
			matched, err := regexp.Match(*excludePatten, []byte(*s3Obj.Key))
			if err != nil || matched {
				continue
			}
		}

		wg.Add(1)
		go saveObject(bucket, baseDir, *s3Obj.Key, s3Client, &wg)
	}
	wg.Wait()
	return nil
}
func saveObject(bucket, baseDir, key string, s3Client s3iface.S3API, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := mkDirIfNeeded(baseDir, key); err != nil {
		log.Println("Unable to create dir: " + err.Error())
		return
	}

	fileName := path.Join(baseDir, key)
	file, err := os.Create(fileName)

	if err != nil {
		log.Println("Unable to create file: " + err.Error())
		return
	}
	defer file.Close()

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	results, err := s3Client.GetObject(input)
	if err != nil {
		log.Println("Unable to download item: " + err.Error())
		return
	}
	defer results.Body.Close()

	if _, err := io.Copy(file, results.Body); err != nil {
		log.Println("Unable to copy item: " + err.Error())
		return
	}
}
func mkDirIfNeeded(baseDir string, key string) (err error) {
	err = nil
	if lastIdx := strings.LastIndex(key, "/"); lastIdx != -1 {
		prefix := key[:lastIdx]
		dirPath := path.Join(baseDir, prefix)
		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return
		}
	}
	return
}

//Upload
func UploadBucket(baseDir, bucket, region string) error {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	iter := createIterator(baseDir, bucket)
	uploader := s3manager.NewUploader(sess)

	if err := uploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
		return err
	}
	return nil
}
func getFiles(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return files
}
func toKey(baseDir, fileName string) string {
	dir := filepath.ToSlash(fileName)
	key := dir[len(baseDir+"/"):]
	return key
}

type directoryIterator struct {
	filePaths []string
	bucket    string
	baseDir   string
	next      struct {
		path string
		key  string
		f    *os.File
	}
	err error
}

func createIterator(baseDir, bucket string) s3manager.BatchUploadIterator {
	paths := getFiles(baseDir)
	return &directoryIterator{
		filePaths: paths,
		bucket:    bucket,
		baseDir:   baseDir,
	}
}

func (iter *directoryIterator) Next() bool {
	if len(iter.filePaths) == 0 {
		iter.next.f = nil
		return false
	}

	f, err := os.Open(iter.filePaths[0])
	iter.err = err

	// Iterate next
	iter.next.f = f
	iter.next.path = iter.filePaths[0]
	iter.next.key = toKey(iter.baseDir, iter.filePaths[0])

	iter.filePaths = iter.filePaths[1:]
	return true && iter.Err() == nil
}

func (iter *directoryIterator) Err() error {
	return iter.err
}

func (iter *directoryIterator) UploadObject() s3manager.BatchUploadObject {
	f := iter.next.f
	return s3manager.BatchUploadObject{
		Object: &s3manager.UploadInput{
			Bucket: &iter.bucket,
			Key:    &iter.next.key,
			Body:   f,
		},
		After: func() error {
			return f.Close()
		},
	}
}
