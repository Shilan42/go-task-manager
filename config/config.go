package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Глобальные переменные для хранения значений из окружения.
var (
	Port        string // Порт приложения (из TODO_PORT)
	DatabaseURL string // Путь к БД (из TODO_DBFILE)
	Password    string // Мастер‑пароль (из TODO_PASSWORD)
	JWTSecret   string // Секрет для подписи JWT (из TODO_JWT_SECRET)
)

// LoadEnv загружает переменные окружения из .env‑файла.
// Если файл не найден, использует системные переменные окружения.
// При критических ошибках (не связанных с отсутствием файла) возвращает ошибку.
//
// Возвращает:
//   - nil, если переменные загружены успешно или .env не найден (используются системные переменные);
//   - ошибку, если возникла проблема при чтении .env (кроме отсутствия файла).
func LoadEnv() error {
	// Пытаемся загрузить .env‑файл с переменными окружения
	err := godotenv.Load()
	if err != nil {
		// Если файл не найден - это не критичная ошибка: продолжаем, используя системные переменные
		if os.IsNotExist(err) {
			log.Println(".env file not found, using system environment variables")
			return nil
		}
		// Любая другая ошибка (например, проблемы с правами, синтаксис .env) - критична
		return err
	}

	// Загружаем значения из окружения (после загрузки .env они доступны через os.Getenv)
	Port = os.Getenv("TODO_PORT")
	DatabaseURL = os.Getenv("TODO_DBFILE")
	Password = os.Getenv("TODO_PASSWORD")
	JWTSecret = os.Getenv("TODO_JWT_SECRET")

	return nil
}
