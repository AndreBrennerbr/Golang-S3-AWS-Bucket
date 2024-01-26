package http_server

import (
	"context"
	s3conect "file_upload_project/services/s3_connect"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

var bucketName string

func Start() {
	bucketName = os.Getenv("BUCKETNAME")
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
	getHandle := http.HandlerFunc(getObjects)
	downloadHandle := http.HandlerFunc(downloadObject)
	deleteHandle := http.HandlerFunc(deletObject)

	http.HandleFunc("/upload", isPostMethodMiddleware(uploadHandle))

	http.HandleFunc("/get_objects", isGetMethodMiddleware(getHandle))

	http.HandleFunc("/download", isGetMethodMiddleware(downloadHandle))

	http.HandleFunc("/delete/", isDeleteMethodMiddleware(deleteHandle))

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

func isGetMethodMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func isDeleteMethodMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
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

func getObjects(w http.ResponseWriter, r *http.Request) {
	listOfObj := s3conect.MinioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})

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

func downloadObject(w http.ResponseWriter, r *http.Request) {
	objectKey := r.URL.Query().Get("objectKey")
	fmt.Println(objectKey)

	object, err := s3conect.MinioClient.GetObject(context.Background(), bucketName, objectKey, minio.GetObjectOptions{})

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

func deletObject(w http.ResponseWriter, r *http.Request) {

	caminho := strings.Split(r.URL.Path, "/")

	if len(caminho) <= 1 {
		http.Error(w, "Erro ao ler o path", http.StatusInternalServerError)
		return
	}

	objectName := caminho[2]

	err := s3conect.MinioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintln(w, "Deletado com sucesso")
	w.WriteHeader(http.StatusOK)
}
