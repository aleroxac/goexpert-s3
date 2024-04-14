package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	TEMP_DIR              = ".temp"
	MAX_UPLOAD_CONCURRENT = 100
	MAX_RETRY_CONCURRENT  = 10
)

var (
	s3Client *s3.S3
	s3Bucket string
	wg       sync.WaitGroup
)

func init() {
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}

	s3Client = s3.New(session)
	s3Bucket = "goexpert-bucket-example-aleroxac-20240414"
}

func main() {
	dir, err := os.Open(TEMP_DIR)
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	uploadControl := make(chan struct{}, MAX_UPLOAD_CONCURRENT)
	errorFileUpload := make(chan string, MAX_RETRY_CONCURRENT)

	go func() {
		for {
			select {
			case filename := <-errorFileUpload:
				uploadControl <- struct{}{}
				wg.Add(1)
				go uploadFile(filename, uploadControl, errorFileUpload)
				fmt.Printf("------------> File %s failed to upload <------------\n", filename)
			}
		}
	}()

	for {
		files, err := dir.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading directory: %f\n", err)
			continue
		}

		wg.Add(1)
		uploadControl <- struct{}{}
		go uploadFile(files[0].Name(), uploadControl, errorFileUpload)
	}
	wg.Wait()
}

func uploadFile(filename string, uploadControl <-chan struct{}, errorFileUpload chan<- string) {
	defer wg.Done()
	completeFileName := fmt.Sprintf("%s/%s", TEMP_DIR, filename)

	file, err := os.Open(completeFileName)
	if err != nil {
		fmt.Printf("Error opening file %v\n", completeFileName)
		<-uploadControl // esvazia o canal
		errorFileUpload <- completeFileName
		return
	}
	defer file.Close()

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		fmt.Printf("Error uploading file %s: %v\n", completeFileName, err)
		<-uploadControl // esvazia o canal
		errorFileUpload <- completeFileName
		return
	}
	fmt.Printf("File %s uploaded successfully\n", completeFileName)
	<-uploadControl // esvazia o canal
}
