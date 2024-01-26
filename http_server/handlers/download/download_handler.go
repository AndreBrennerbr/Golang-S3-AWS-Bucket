package download

import (
	"context"
	minio_service "file_upload_project/core/services/minIO"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
)

func DownloadObject(w http.ResponseWriter, r *http.Request) {
	objectKey := r.URL.Query().Get("objectKey")
	fmt.Println(objectKey)

	object, err := minio_service.MinioClient.GetObject(context.Background(), minio_service.BucketName, objectKey, minio.GetObjectOptions{})

	if err != nil {
		http.Error(w, "Erro ao obter o objeto do MinIO", http.StatusInternalServerError)
		return
	}

	defer object.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", objectKey))
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, object)
	if err != nil {
		log.Println("Erro ao copiar o conteúdo do objeto para a resposta HTTP:", err)
		http.Error(w, "Erro ao servir o objeto", http.StatusInternalServerError)
		return
	}
}
