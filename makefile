# Makefile для S3 Files Manager

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

help:
	@echo ''
	@echo '${YELLOW}S3 Files Manager — Makefile команды${RESET}'
	@echo ''
	@echo 'make <команда>'
	@echo ''
	@echo 'Команды:'
	@echo '  ${GREEN}run${RESET}          Запуск сервиса'
	@echo '  ${GREEN}minio-up${RESET}     Запуск MinIO'
	@echo '  ${GREEN}minio-down${RESET}   Остановка MinIO'
	@echo '  ${GREEN}minio-logs${RESET}   Логи MinIO'
	@echo '  ${GREEN}tidy${RESET}         go mod tidy'
	@echo ''

run:
	@echo '${GREEN}Запуск сервиса...${RESET}'
	go run cmd/main.go

minio-up:
	@echo '${GREEN}Запуск MinIO...${RESET}'
	docker-compose up -d

minio-down:
	@echo '${GREEN}Остановка MinIO...${RESET}'
	docker-compose down

minio-logs:
	@echo '${GREEN}Логи MinIO${RESET}'
	docker-compose logs -f minio

tidy:
	@echo '${GREEN}go mod tidy${RESET}'
	go mod tidy

.PHONY: help run minio-up minio-down minio-logs tidy