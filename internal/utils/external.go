package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// FetchSongDetails запрашивает данные песни из внешнего API.
func FetchSongDetails(group, song string) (SongDetail, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", os.Getenv("EXTERNAL_API_URL"), group, song)
	resp, err := http.Get(url)
	if err != nil {
		return SongDetail{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SongDetail{}, fmt.Errorf("external API error: %s", resp.Status)
	}

	var details SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return SongDetail{}, err
	}

	return details, nil
}
