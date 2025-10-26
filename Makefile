.PHONY: run build clean clear-port

# Запуск сервера
run:
	go run ./cmd/server

# Сборка бинарного файла
build:
	go build -o bin/server ./cmd/server

# Очистка и перезапуск
clean: clear-port
	rm -f blog.db blog.db-shm blog.db-wal
	go run ./cmd/server

# Очистить порт
clear-port:
	@if lsof -t -i :3000 > /dev/null 2>&1; then \
		kill -9 $$(lsof -t -i :3000); \
		echo "Порт 3000 освобожден"; \
	else \
		echo "Порт 3000 уже свободен"; \
	fi
