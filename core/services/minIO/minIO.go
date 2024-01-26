package minio_service

import (
	"context"
	"file_upload_project/core/config"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioClient     *minio.Client
	err             error
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	BucketName      string
)

func Create() {
	BucketName = config.EnvBucketName()
	endpoint = config.EnvEndPoint()
	accessKeyID = config.EnvAccessKey()
	secretAccessKey = config.EnvSecretAccesKey()
	useSSL = false

	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

}

func ExistBucket(background context.Context) (bool, error) {

	existBucket, err := MinioClient.BucketExists(background, BucketName)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	if !existBucket {

		err := MinioClient.MakeBucket(background, BucketName, minio.MakeBucketOptions{})

		if err != nil {
			fmt.Println(err)
			return false, err
		}

	}

	return true, nil

}

func Start(background context.Context) {
	Create()
	exist, err := ExistBucket(background)
	if !exist {
		log.Fatal(err)
		return
	}
}
