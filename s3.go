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
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

const (
	messageClientNotDefined = "Aws Client not defined"
)

type Bucket struct {
	s3Client s3iface.S3API
	Name     string
	LocalDir string
}

func NewBucket(client s3iface.S3API, name, localDir string) Bucket {
	return Bucket{s3Client: client, Name: name, LocalDir: localDir}
}

//DownloadBucket ...
func (b *Bucket) DownloadBucket(excludePatten *string) error {
	var wg sync.WaitGroup

	if b.s3Client == nil {
		return fmt.Errorf(messageClientNotDefined)
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(b.Name),
	}

	result, err := b.s3Client.ListObjectsV2(input)
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
		go saveObjectToS3(b.Name, b.LocalDir, *s3Obj.Key, b.s3Client, &wg)
	}
	wg.Wait()
	return nil
}
func saveObjectToS3(bucket, baseDir, key string, s3Client s3iface.S3API, wg *sync.WaitGroup) {
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

//UploadBucket ...
func (b *Bucket) UploadBucket() error {
	var wg sync.WaitGroup

	if b.s3Client == nil {
		return fmt.Errorf(messageClientNotDefined)
	}

	for _, file := range getFiles(b.LocalDir) {
		wg.Add(1)
		go saveObjectFromS3(b.Name, b.LocalDir, file, b.s3Client, &wg)
	}
	wg.Wait()
	return nil
}
func saveObjectFromS3(bucket, baseDir, fileName string, s3Client s3iface.S3API, wg *sync.WaitGroup) {
	defer wg.Done()

	key := toKey(baseDir, fileName)

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(strings.NewReader(fileName)),
	}
	if _, err := s3Client.PutObject(input); err != nil {
		log.Println("Unable to upload file: " + err.Error())
		return
	}
	return

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
