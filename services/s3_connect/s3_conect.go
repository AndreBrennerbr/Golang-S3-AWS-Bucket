package s3conect

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioClient *minio.Client
	err         error
)

var endpoint string = "192.168.0.52:9000"
var accessKeyID string = "ROOTNAME"
var secretAccessKey string = "CHANGEME123"
var useSSL bool = false

func Create() {

	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

}
