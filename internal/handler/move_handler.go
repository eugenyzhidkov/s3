package handler

import (
	"net/http"
	"strings"

	"s3/internal/infrastructure/logger"
	"s3/internal/s3"
)

func MoveHandler(s3Client *s3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := s3Client.ListObjects(r.Context())
		if err != nil {
			logger.Error("Ошибка получения списка файлов", "error", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		moved := 0
		for _, file := range files {
			if strings.HasPrefix(file, "archive/") {
				continue
			}
			newKey := "archive/" + file
			if err := s3Client.MoveFile(r.Context(), file, newKey); err != nil {
				logger.Error("Ошибка перемещения файла", "file", file, "error", err)
			} else {
				logger.Info("Файл перемещён", "from", file, "to", newKey)
				moved++
			}
		}

		logger.Info("Перемещение завершено", "moved", moved)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
