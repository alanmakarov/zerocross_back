
#переменные 

BUILD_DIR=bin
# Makefile для Go проекта

.PHONY: all build clean

# Задача для сборки релизной версии
build:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/app-linux-amd64 ./
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/app-windows-amd64.exe ./
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/app-darwin-arm64 ./

# Задача для чистки
clean:
	@echo "cleaning"
	rm -rf $(BUILD_DIR)/*
	@echo "Очистка завершена."
