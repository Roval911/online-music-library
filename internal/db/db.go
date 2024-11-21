package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"online-music-library/internal/models"
)

var db *sql.DB

// Инициализация базы данных
func InitDB() {
	var err error
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
}

// Закрытие соединения с базой данных
func CloseDB() {
	db.Close()
}

// Запуск миграций
func RunMigrations() {
	// Пример SQL миграций (создание таблиц)
	createTableQuery := `CREATE TABLE IF NOT EXISTS songs (
		id SERIAL PRIMARY KEY,
		group_name TEXT NOT NULL,
		song_name TEXT NOT NULL,
		release_date TEXT,
		text TEXT,
		link TEXT
	);`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Error running migrations: ", err)
	}
}

// Сохранение песни
func SaveSong(song *models.Song, details models.Song) error {
	query := `INSERT INTO songs (group_name, song_name, release_date, text, link) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.QueryRow(query, song.Group, song.Song, details.ReleaseDate, details.Text, details.Link).Scan(&song.ID)
	return err
}

// Получение фильтрованных песен
func GetFilteredSongs(group, song string, limit, offset int) ([]models.Song, error) {
	query := `SELECT id, group_name, song_name, release_date, text, link 
		FROM songs WHERE group_name LIKE $1 AND song_name LIKE $2 LIMIT $3 OFFSET $4`
	rows, err := db.Query(query, "%"+group+"%", "%"+song+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var s models.Song
		err := rows.Scan(&s.ID, &s.Group, &s.Song, &s.ReleaseDate, &s.Text, &s.Link)
		if err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}

	return songs, nil
}

// Получение текста песни
func GetSongText(id string) (string, error) {
	query := `SELECT text FROM songs WHERE id = $1`
	var text string
	err := db.QueryRow(query, id).Scan(&text)
	if err != nil {
		return "", err
	}
	return text, nil
}

// Обновление данных песни
func UpdateSong(id string, song *models.Song) error {
	query := `UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6`
	_, err := db.Exec(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, id)
	return err
}

// Удаление песни
func DeleteSong(id string) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
