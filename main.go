package main

import (
	"log"
	"sync"
	"time"

	"github.com/IamStubborN/filmtracker/youtube"

	"github.com/IamStubborN/filmtracker/gsrv"

	"github.com/IamStubborN/filmtracker/scrapper"

	_ "github.com/IamStubborN/filmtracker/docs"
)

// @title FilmTracker API
// @version 1.0
// @description DITS test FilmTracker project.
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host filmtracker-api.com

// @securityDefinitions.apiKey Token
// @in cookies
// @name Token

// @securityDefinitions.apiKey Refresh
// @in cookies
// @name Refresh

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

func main() {
	go updateFilmsDatabase(false)
	go youtube.StartSearchTrailers()
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
