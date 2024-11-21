package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2" // импортируем resty
	"net/http"
	"online-music-library/internal/db"
	"online-music-library/internal/models"
	"os"
	"strconv"
	"strings"
)

// Получение песен с фильтрацией и пагинацией
func GetSongs(c *gin.Context) {
	group := c.DefaultQuery("group", "")
	song := c.DefaultQuery("song", "")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Преобразуем строковые значения в int
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset := (page - 1) * limit
	songs, err := db.GetFilteredSongs(group, song, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// Получение текста песни с пагинацией по куплетам
func GetSongText(c *gin.Context) {
	id := c.Param("id")

	songText, err := db.GetSongText(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch song text"})
		return
	}

	// Пагинация по куплетам
	verses := splitSongTextIntoVerses(songText)
	c.JSON(http.StatusOK, verses)
}

// Разбиение текста на куплеты
func splitSongTextIntoVerses(text string) []string {
	// Простой пример разбиения текста по строкам
	// Реализовать более сложное разбиение по куплетам, если нужно
	return strings.Split(text, "\n\n")
}

// Добавление песни и запрос к внешнему API
func AddSong(c *gin.Context) {
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Запрос к внешнему API для получения информации о песне
	apiUrl := os.Getenv("API_URL")
	client := resty.New()
	resp, err := client.R().
		SetQueryParam("group", song.Group).
		SetQueryParam("song", song.Song).
		Get(apiUrl)

	if err != nil || resp.StatusCode() != 200 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch song details"})
		return
	}

	// Обработка полученных данных
	var songDetails models.Song
	if err := json.Unmarshal(resp.Body(), &songDetails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshalling API response"})
		return
	}

	// Сохранение песни в базу данных
	err = db.SaveSong(&song, songDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving song to database"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// Обновление данных песни
func UpdateSong(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Обновление в базе данных
	err := db.UpdateSong(id, &song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating song"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// Удаление песни
func DeleteSong(c *gin.Context) {
	id := c.Param("id")

	// Удаление из базы данных
	err := db.DeleteSong(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
