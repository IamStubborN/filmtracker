package main

import (
	"log"
	"sync"

	"github.com/IamStubborN/filmtracker/scrapper"

	"github.com/IamStubborN/filmtracker/gsrv"

	"github.com/joho/godotenv"
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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("File .env not found, reading configuration from ENV")
	}
}

func main() {
	go updateFilmsDatabase(false)
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
}
