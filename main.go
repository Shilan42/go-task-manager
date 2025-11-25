package main

import (
	"go-task-manager-final_project/config"
	"go-task-manager-final_project/internal/db"
	"go-task-manager-final_project/internal/server"
	"log"
	"os"
)

// В main инициализируем соединение с базой данных, обеспечиваем его корректное закрытие и запускаем HTTP-сервер для обработки запросов.
func main() {

	// Загружаем переменные окружения из .env или системных переменных
	if err := config.LoadEnv(); err != nil {
		log.Printf("failed to load environment variables: %v", err)
		os.Exit(1) // Критическая ошибка — без конфига работа невозможна
	}

	// Открываем соединения с БД и, при необходимости, создаем схему
	db, err := db.Init(config.DatabaseURL)
	if err != nil {
		log.Printf("failed to initialize database: %v", err)
	}
	// Обеспечиваем закрытие соединения с БД при завершении работы программы (даже в случае паники или ошибки).
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("failed to close database connection: %v", closeErr)
		}
	}()

	// Запускаем сервер
	err = server.StartServer(db)
	if err != nil {
		log.Printf("failed to start server: %v", err)
		return
	}
}
