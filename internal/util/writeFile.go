package util

import (
	"fmt"
	"log/slog"
	"mime/multipart"
	"os"
)

func WriteFile(input *multipart.FileHeader, output *os.File) error {
	f, err := input.Open()
	if err != nil {
		slog.Error("error open file", "error", err, "file name", input.Filename)
		return fmt.Errorf("error open file: %w", err)
	}
	defer f.Close()

	fileContent := make([]byte, input.Size)
	_, err = f.Read(fileContent)
	if err != nil {
		slog.Error("error read file", "error", err, "file name", input.Filename)
		return fmt.Errorf("error read file: %w", err)
	}

	_, err = output.Write(fileContent)
	if err != nil {
		slog.Error("error write file", "error", err, "file name", input.Filename)
		return fmt.Errorf("error write file: %w", err)
	}

	return nil
}
