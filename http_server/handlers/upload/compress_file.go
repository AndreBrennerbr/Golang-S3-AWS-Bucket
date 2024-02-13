package upload

import (
	"archive/tar"
	"compress/gzip"
	"file_upload_project/core/config"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

type Sizer interface {
	Size() int64
}

func compress(w http.ResponseWriter, r *http.Request, file multipart.File, fileName string) (*os.File, error) {
	dir, err := config.GetUniversalDir()
	if err != nil {
		log.Println("Erro ao obter o diretório raiz:", err)
		return nil, err
	}

	tmpDir := filepath.Join(dir, "tmp")
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		if err := os.Mkdir(tmpDir, 0777); err != nil {
			log.Println("Erro ao criar diretório:", err)
			return nil, err
		}
	}

	// Substituindo a extensão para .tar.gz
	fileNameFormatado := regexp.MustCompile(`\.(.*)$`).ReplaceAllString(fileName, ".tar.gz")
	tarGzFilePath := filepath.Join(tmpDir, fileNameFormatado)

	// Criando um novo arquivo .tar.gz para a compressão
	tarGzFile, err := os.Create(tarGzFilePath)
	if err != nil {
		log.Println("Erro ao criar arquivo .tar.gz para compressão:", err)
		return nil, err
	}
	defer tarGzFile.Close()

	gzipWriter := gzip.NewWriter(tarGzFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Criar o cabeçalho do tar para o arquivo
	header := &tar.Header{
		Name: filepath.Base(fileName), // Preserva o nome original do arquivo no tar
		Size: file.(Sizer).Size(),
	}

	// Escrever o cabeçalho no arquivo tar
	err = tarWriter.WriteHeader(header)
	if err != nil {
		log.Println("Erro ao escrever cabeçalho no tar:", err)
		return nil, err
	}

	// Copiando e comprimindo o conteúdo do arquivo para o tar.gz
	if _, err := io.Copy(tarWriter, file); err != nil {
		log.Println("Erro ao copiar e comprimir arquivo para o tar.gz:", err)
		return nil, err
	}

	// Garantindo que todos os dados sejam flushados para o gzip e tar writer
	if err := tarWriter.Close(); err != nil {
		log.Println("Erro ao fechar o tar writer:", err)
		return nil, err
	}
	if err := gzipWriter.Close(); err != nil {
		log.Println("Erro ao fechar o gzip writer:", err)
		return nil, err
	}

	// Reabrindo o arquivo .tar.gz para retorno
	tarGzFile, err = os.Open(tarGzFilePath)
	if err != nil {
		log.Println("Erro ao abrir o arquivo .tar.gz:", err)
		return nil, err
	}

	return tarGzFile, nil
}
