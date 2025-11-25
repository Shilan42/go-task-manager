package middleware

import (
	"crypto/sha256"
	"fmt"
	"go-task-manager-final_project/config"
	"go-task-manager-final_project/internal/api"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Auth - middleware-функция для проверки авторизации пользователя через JWT-токен.
// Параметр:
// next - обработчик HTTP-запроса, который будет вызван при успешной авторизации.
// Возвращает:
// http.HandlerFunc - обернутый обработчик с логикой авторизации.
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Если пароль задан, выполняем проверку авторизации.
		if config.Password != "" {
			// Пытаемся получить cookie с именем "token" из запроса.
			cookie, err := r.Cookie("token")
			if err != nil {
				// Если cookie отсутствует или возникла ошибка - возвращаем статус 401 (Неавторизован).
				api.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "unauthorized",
				})
				return
			}

			// Если переменная не задана, возвращаем ошибку 500 (Internal Server Error).
			if config.JWTSecret == "" {
				api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "JWT secret not configured",
				})
				return
			}
			secret := []byte(config.JWTSecret)

			// Парсим JWT-токен из значения cookie.
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				// Проверяем, что алгоритм подписи токена соответствует ожидаемому (HMAC).
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method %q", token.Header["alg"])
				}
				return secret, nil
			})

			// Если при парсинге токена произошла ошибка или токен недействителен - возвращаем ошибку.
			if err != nil || !token.Valid {
				api.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "token expired or invalid",
				})
				return
			}

			// Извлекаем claims (данные) из токена.
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				// Если claims не соответствуют ожидаемому типу - возвращаем ошибку.
				api.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid token: malformed claims",
				})
				return
			}

			// Вычисляем SHA-256 хэш текущего пароля из окружения.
			currentHash := sha256.Sum256([]byte(config.Password))
			currentHashStr := fmt.Sprintf("%x", currentHash)

			// Сравниваем хэш пароля из токена с текущим хэшем пароля.
			// Если хэши не совпадают - токен недействителен.
			if claims["password_hash"] != currentHashStr {
				api.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid token: password changed",
				})
				return
			}

		}
		// Если все проверки прошли - передаём запрос дальше по цепочке обработчиков.
		next(w, r)
	})
}
