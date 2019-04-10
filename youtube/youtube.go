package youtube

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/option"

	"google.golang.org/api/youtube/v3"
)

var service *youtube.Service

func init() {
	ctx := context.Background()
	s, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_API")))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	service = s
}

func GetTrailerID(query string) (string, error) {
	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(query + " трейлер").
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Println(err)
	}

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			return item.Id.VideoId, nil
		}
	}
	return "", fmt.Errorf(`can't find video on youtube by query: %s`, query)
}
