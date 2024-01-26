package getfile

import (
	"context"
	"fmt"
	"net/http"

	minio_service "file_upload_project/core/services/minIO"

	"github.com/minio/minio-go/v7"
)

func GetObjects(w http.ResponseWriter, r *http.Request) {
	listOfObj := minio_service.MinioClient.ListObjects(context.Background(), minio_service.BucketName, minio.ListObjectsOptions{})

	for object := range listOfObj {
		if object.Err != nil {
			fmt.Println(object.Err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Erro ao listar objetos")
			return
		}
		fmt.Fprintln(w, object.Key)
	}
	w.WriteHeader(http.StatusOK)
}
