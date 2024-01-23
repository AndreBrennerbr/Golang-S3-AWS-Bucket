package s3conect

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioClient *minio.Client
	err         error
)

var endpoint string
var accessKeyID string
var secretAccessKey string
var useSSL bool

func Create() {
	endpoint = os.Getenv("ENDPOINT")
	accessKeyID = os.Getenv("ACCESSKEYID")
	secretAccessKey = os.Getenv("SECRETACCESSKEY")
	useSSL = false

	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

}
