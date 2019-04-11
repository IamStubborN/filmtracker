package youtube

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IamStubborN/filmtracker/database"
	"github.com/IamStubborN/filmtracker/tmdb"

	"google.golang.org/api/option"

	"google.golang.org/api/youtube/v3"
)

var service *youtube.Service
var db *database.Database

func init() {
	ctx := context.Background()
	s, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_API")))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	service = s
	db = database.GetDB()
}

func StartSearchTrailers() {
	for range time.NewTicker(time.Duration(2 * time.Hour)).C {
		films, err := db.GetAllFilms()
		if err != nil {
			log.Println(err)
		}
		for _, film := range films {
			if film.YoutubeID == "" {
				getTrailerIDWithRetry(film)
				fmt.Println(film.YoutubeID)
				if err := db.UpsertFilm(film); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func getTrailerIDWithRetry(film *tmdb.Film) {
	fmt.Println(123)
	retryCount := 5
	youtubeID, err := GetTrailerID(
		film.Name + " " + strings.Split(film.ReleaseDate, "-")[0])
	if err != nil {
		count := 1
		fmt.Println("Try", count, err)
		time.Sleep(24 * time.Hour)
		for count <= retryCount {
			count++
			youtubeID, err = GetTrailerID(
				film.Name + " " + strings.Split(film.ReleaseDate, "-")[0])
			if err != nil {
				fmt.Println("Try", count, err)
			} else {
				break
			}
		}
	}
	film.YoutubeID = youtubeID
}

func GetTrailerID(query string) (string, error) {
	fmt.Println("youtube", query)
	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(query + " трейлер").
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		return "", err
	}
	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			return item.Id.VideoId, nil
		}
	}
	return "", fmt.Errorf(`can't find video on youtube by query: %s`, query)
}
