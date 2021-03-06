package lambda

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Uploader interface {
	Upload(f *Function) error
}

type S3Uploader struct {
}

func (s3 S3Uploader) Upload(lambdaF *Function) error {
	log.Printf("Uploading function %s to S3 Bucket %s", lambdaF.Name, lambdaF.S3Bucket)
	zf, _ := os.Open(lambdaF.Target)
	// The session the S3 Uploader will use
	sess, _ := session.NewSession()

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	upParams := &s3manager.UploadInput{
		Bucket: &lambdaF.S3Bucket,
		Key:    &lambdaF.S3Key,
		Body:   zf,
	}

	_, err := uploader.Upload(upParams, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024 // 10MB part size
		u.LeavePartsOnError = true    // Don't delete the parts if the upload fails.
	})

	if err != nil {
		log.Fatal(err)
	}
	return err

}
