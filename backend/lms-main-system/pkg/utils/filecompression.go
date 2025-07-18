package utils

import (
	"github.com/klauspost/compress/gzip"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func CompressBackupFile(filename string) error {
	original, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(original *os.File) {
		err := original.Close()
		if err != nil {
			logrus.Errorf("Failed to close gzip writer: %v", err)
		}
	}(original)

	compressed, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer func(compressed *os.File) {
		err := compressed.Close()
		if err != nil {
			logrus.Errorf("Failed to close gzip writer: %v", err)
		}
	}(compressed)

	gz := gzip.NewWriter(compressed)
	defer func(gz *gzip.Writer) {
		err := gz.Close()
		if err != nil {
			logrus.Errorf("Failed to close gzip writer: %v", err)
		}
	}(gz)

	if _, err := io.Copy(gz, original); err != nil {
		return err
	}

	err = os.Remove(filename)
	if err != nil {
		logrus.Errorf("Failed to remove backups file: %v", err)
		return err
	}
	return nil
}
