package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"s3/internal/image"
	"s3/internal/infrastructure/logger"
	"s3/internal/s3"
)

func UploadHandler(s3Client *s3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := "static/"
		entries, err := os.ReadDir(dir)
		if err != nil {
			logger.Error("Ошибка чтения папки static", "error", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		uploaded := 0
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filename := entry.Name()
			if !strings.HasSuffix(strings.ToLower(filename), ".jpg") && !strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
				continue
			}

			path := filepath.Join(dir, filename)

			originalFile, err := os.Open(path)
			if err != nil {
				logger.Error("Ошибка открытия оригинального файла", "file", filename, "error", err)
				continue
			}
			stat, err := originalFile.Stat()
			if err != nil {
				logger.Error("Ошибка stat файла", "file", filename, "error", err)
				originalFile.Close()
				continue
			}
			originalSize := stat.Size()
			originalFile.Close()

			compressed, err := image.CompressJPEG(path, 75)
			if err != nil {
				logger.Error("Ошибка сжатия", "file", filename, "error", err)
				continue
			}

			compressedSize := int64(len(compressed))

			reduction := 100.0 - (float64(compressedSize) / float64(originalSize) * 100.0)

			if err := s3Client.UploadFile(r.Context(), filename, compressed); err != nil {
				logger.Error("Ошибка загрузки в S3", "file", filename, "error", err)
				continue
			}

			logger.Info("Картинка загружена в S3",
				"file", filename,
				"original_kb", originalSize/1024,
				"compressed_kb", compressedSize/1024,
				"reduction_percent", fmt.Sprintf("%.1f%%", reduction),
			)

			uploaded++
		}

		logger.Info("Загрузка завершена", "uploaded", uploaded)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
