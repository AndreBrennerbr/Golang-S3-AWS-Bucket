package routes

import (
	"file_upload_project/core/middlewares"
	delete_handler "file_upload_project/http_server/handlers/delete"
	"file_upload_project/http_server/handlers/download"
	getfile "file_upload_project/http_server/handlers/getFile"
	"file_upload_project/http_server/handlers/upload"
	"log"
	"net/http"
)

func Routes() {

	uploadHandle := http.HandlerFunc(upload.UploadFileMinIO)
	getHandle := http.HandlerFunc(getfile.GetObjects)
	downloadHandle := http.HandlerFunc(download.DownloadObject)
	deleteHandle := http.HandlerFunc(delete_handler.DeletObject)

	http.HandleFunc("/upload", middlewares.IsPostMethodMiddleware(uploadHandle))

	http.HandleFunc("/get_objects", middlewares.IsGetMethodMiddleware(getHandle))

	http.HandleFunc("/download", middlewares.IsGetMethodMiddleware(downloadHandle))

	http.HandleFunc("/delete/", middlewares.IsDeleteMethodMiddleware(deleteHandle))

	log.Fatal(http.ListenAndServe(":8080", nil))

}
