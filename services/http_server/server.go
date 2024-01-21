package http_server

import (
	"context"
	s3conect "file_upload_project/services/s3_connect"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

var bucketName string = "brenner"

func Start() {
	s3conect.Create()

	existBucket, _ := s3conect.MinioClient.BucketExists(context.Background(), bucketName)

	if !existBucket {
		err := s3conect.MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	uploadHandle := http.HandlerFunc(uploadFileMinIO)

	http.HandleFunc("/upload", isPostMethodMiddleware(uploadHandle))

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func isPostMethodMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func uploadFileMinIO(w http.ResponseWriter, r *http.Request) {

	validadeTypeOfFile(w, r)

	file, handler, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Erro ao obter o arquivo do formulário", http.StatusBadRequest)
		return
	}

	objectName := handler.Filename

	info, err := s3conect.MinioClient.PutObject(context.Background(), bucketName, objectName, file, -1, minio.PutObjectOptions{})

	if err != nil {
		http.Error(w, "Erro ao salvar arquivo", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Upload concluido com sucesso "+info.ETag)
	return
}

func validadeTypeOfFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao obter o arquivo do formulário", http.StatusBadRequest)
		return
	}

	defer file.Close()

	acceptedExtensions := []string{".pdf", ".doc", ".docx"}
	extension := filepath.Ext(handler.Filename)

	if !contains(acceptedExtensions, extension) {
		http.Error(w, "Tipo de arquivo não permitido", http.StatusBadRequest)
	}
}

func contains(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
}
