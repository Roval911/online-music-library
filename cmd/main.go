package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"online-music-library/internal/db"
	"online-music-library/internal/handlers"
)

func init() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Инициализация базы данных
	db.InitDB()

	// Запуск миграций
	db.RunMigrations()
}

func main() {
	// Настройка роутов
	router := gin.Default()

	// Роуты
	router.GET("/songs", handlers.GetSongs)             // Получение песен с фильтрацией и пагинацией
	router.GET("/songs/:id/text", handlers.GetSongText) // Получение текста песни с пагинацией по куплетам
	router.POST("/songs", handlers.AddSong)             // Добавление песни с запросом в API
	router.PUT("/songs/:id", handlers.UpdateSong)       // Обновление данных песни
	router.DELETE("/songs/:id", handlers.DeleteSong)    // Удаление песни

	// Запуск сервера
	router.Run(":8080")
}
