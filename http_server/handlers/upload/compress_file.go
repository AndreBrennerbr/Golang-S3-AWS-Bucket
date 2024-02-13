package upload

import (
	"bufio"
	"compress/gzip"
	"file_upload_project/core/config"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func compress(w http.ResponseWriter, r *http.Request, file multipart.File, fileName string) (*os.File, error) {
	dir, err := config.GetUniversalDir()

	if err != nil {
		log.Println("Erro ao obter o diretório raiz:", err)
		return nil, err
	}

	tmpDir := filepath.Join(dir, "tmp")

	// Verifica se o caminho existe.
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		// Se não existir, cria.
		if err := os.Mkdir(tmpDir, 0777); err != nil {
			log.Println("Erro ao criar diretório:", err)
			return nil, err
		}
	}

	// Criando arquivo temporário localmente.
	createdFileEmpty, err := os.Create(filepath.Join(tmpDir, fileName))
	if err != nil {
		log.Println("Erro ao criar arquivo temporário:", err)
		return nil, err
	}

	defer func() {
		createdFileEmpty.Close()
		os.Remove(createdFileEmpty.Name()) // Limpar arquivo temporário após o uso.
	}()

	// Copiando o conteúdo do arquivo original para o arquivo temporário.
	_, err = io.Copy(createdFileEmpty, file)
	if err != nil {
		log.Println("Erro ao copiar conteúdo do arquivo original:", err)
		return nil, err
	}

	// Lendo o arquivo temporário.
	createdFileEmpty.Seek(0, 0) // Volta ao início do arquivo.
	reader := bufio.NewReader(createdFileEmpty)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("Erro ao ler arquivo temporário:", err)
		return nil, err
	}

	// Substituindo a extensão.
	fileNameFormatado := regexp.MustCompile(`\.(.*)`).ReplaceAllString(fileName, ".gz")

	// Criando um novo arquivo.
	f, err := os.Create(filepath.Join(tmpDir, fileNameFormatado))
	if err != nil {
		log.Println("Erro ao criar novo arquivo:", err)
		return nil, err
	}
	defer f.Close()

	// Fazendo compressão do arquivo.
	gzipFile := gzip.NewWriter(f)
	defer gzipFile.Close()
	_, err = gzipFile.Write(content)
	if err != nil {
		log.Println("Erro ao escrever conteúdo comprimido:", err)
		return nil, err
	}

	// Abrindo o arquivo comprimido para retornar à função principal.
	zipedFile, err := os.Open(filepath.Join(tmpDir, fileNameFormatado))
	if err != nil {
		log.Println("Erro ao abrir arquivo comprimido:", err)
		return nil, err
	}

	fmt.Println(zipedFile.Name())

	return zipedFile, nil
}
