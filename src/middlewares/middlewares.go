package middlewares

import (
	"archive/zip"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func SliceCommand(command string) []string {
	slicedCommand := strings.Split(strings.TrimSuffix(command, "\n"), " ")
	return slicedCommand
}

func B64Encode(resp string) string {
	b64 := base64.StdEncoding.EncodeToString([]byte(resp))
	return b64
}

func ZipFile(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/"
		}

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})

}
