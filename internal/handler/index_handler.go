package handler

import (
	"html/template"
	"net/http"

	"s3/internal/infrastructure/logger"
	"s3/internal/s3"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

type pageData struct {
	Files  []string
	Bucket string
}

func IndexHandler(s3Client *s3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := s3Client.ListObjects(r.Context())
		if err != nil {
			logger.Error("Ошибка получения списка файлов из S3", "error", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		data := pageData{
			Files:  files,
			Bucket: s3Client.Bucket,
		}

		if err := tmpl.Execute(w, data); err != nil {
			logger.Error("Ошибка рендеринга шаблона", "error", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
	}
}
