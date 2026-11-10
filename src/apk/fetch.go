package apk

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const indexFileName = "APKINDEX"

const requestTimeout = 60 * time.Second

func FetchIndex(url string) (*Index, error) {
	client := &http.Client{Timeout: requestTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к репозиторию: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("репозиторий ответил с ошибкой: %s", resp.Status)
	}
	return ReadArchive(resp.Body)
}

func ReadArchive(r io.Reader) (*Index, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("файл не является gzip-архивом: %w", err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("в архиве нет файла %s", indexFileName)
		}
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения tar-архива: %w", err)
		}
		if header.Name == indexFileName {
			return ParseIndex(tr)
		}
	}
}
