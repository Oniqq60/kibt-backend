package middleware

import (
	"net/http"
	"github.com/go-chi/cors"
)

func CORS() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		// Разрешаем запросы с фронтенда
		AllowedOrigins: []string{
			"http://localhost:3000", // Next.js dev
			"http://localhost:5173", // Vite dev
			"https://kibt.ru",       // Продакшн домен
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,  // Время кэширования preflight-запроса
	}).Handler
}