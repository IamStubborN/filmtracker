package main

import (
	"log"
	"sync"
	"time"

	"github.com/IamStubborN/filmtracker/gsrv"

	"github.com/IamStubborN/filmtracker/scrapper"

	_ "github.com/IamStubborN/filmtracker/docs"
)

var newTorrent = &scrapper.FilmTracker{
	URL:            "http://newtorrent.org/",
	URLCategory:    "films/",
	PostfixURL:     "page/",
	ContainerClass: ".entry",
	EntryDetails:   ".entry__info-download",
	FilmName:       ".inner-entry__title",
	PageCount:      ".pages",
}

var toreentsClub = &scrapper.FilmTracker{
	URL:            "http://toreents.club/",
	URLCategory:    "katalog-torrent-films/",
	PostfixURL:     "page/",
	ContainerClass: ".dpad",
	EntryDetails:   ".argmore a[href]",
	FilmName:       ".btl",
	PageCount:      ".navigation",
}

var wg = &sync.WaitGroup{}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	//go updateFilmsDatabase(true)
	//go youtube.StartSearchTrailers()
	server := gsrv.CreateServer()
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func updateFilmsDatabase(isFullUpdate bool) {
	listTrackers := []*scrapper.FilmTracker{
		newTorrent,
		toreentsClub,
	}
	for _, tracker := range listTrackers {
		wg.Add(1)
		sc, err := scrapper.CreateScrapper([]string{
			"toreents.club", "newtorrent.org",
			"www.proxy-list.download", "www.ua-tracker.com"}, false)
		if err != nil {
			log.Fatal(err)
		}
		go sc.GetAllMovies(tracker, isFullUpdate)
	}
	wg.Wait()
	for range time.NewTicker(time.Duration(3 * time.Hour)).C {
		updateFilmsDatabase(false)
	}
}
