package tmdb

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/ryanbradynd05/go-tmdb"
)

type (
	// MovieDB struct. For manipulating with tmdb package.
	MovieDB struct {
		tmdbApi *tmdb.TMDb
		options map[string]string
		Genres  []*Genre
	}
	// Film is response struct.
	Film struct {
		ID           int               `bson:"id" json:"id" example:"299534"`
		Name         string            `bson:"name" json:"name" example:"Мстители: Финал"`
		OriginalName string            `bson:"original_name" json:"original_name" example:"Avengers: Endgame"`
		Poster       string            `bson:"poster_path" json:"poster_path" example:"http://image.tmdb.org/t/p/w500/4fjlSFUoW6L9FSAE5knxjjUDd6v.jpg"`
		ReleaseDate  string            `bson:"release_date" json:"release_date" example:"2019-04-24"`
		Genres       []*Genre          `bson:"genres" json:"genres"`
		Overview     string            `bson:"overview" json:"overview" example:"Мстители и Стражи Галактики вступают в последнюю стадию войны с Таносом, владеющим всемогущей Перчаткой Бесконечности. Грядёт финальная битва между силами героев и Безумного Титана, которая раз и навсегда определит дальнейшую судьбу не только Земли, но и всей вселенной."`
		AddedDate    string            `bson:"added_date" json:"added_date" example:"2019-06-02 19:01:10"`
		YoutubeID    string            `bson:"youtube_id" json:"youtube_id" example:"https://www.youtube.com/watch?v=gbcVZgO4n4E"`
		MagnetLinks  map[string]string `bson:"magnet_links" json:"magnet_links" example:"[magnet:?xt=urn:btih:855dca1c13fbca4eca6b701555a13d32ea065dbb&dn=rutor.info&tr=udp://opentor.org:2710&tr=udp://opentor.org:2710&tr=http://retracker.local/announce]"`
	}
	// Genre struct is overview of all genres from tmdb.
	Genre struct {
		ID          int    `bson:"id" json:"id" example:"28"`
		EnglishName string `bson:"english_name" json:"english_name" example:"action"`
		RussianName string `bson:"russian_name" json:"russian_name" example:"боевик"`
	}
)

// posterBaseURL basic link to get poster from tmdb.
const posterBaseURL = "http://image.tmdb.org/t/p/w500"

var movieDB *MovieDB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("File .env not found, reading configuration from ENV")
	}

	mdb := new(MovieDB)
	config := tmdb.Config{
		APIKey: os.Getenv("TMDB_API"),
	}
	mdb.tmdbApi = tmdb.Init(config)
	mdb.options = make(map[string]string)
	mdb.options["language"] = "ru"
	mdb.fillGenres()
	movieDB = mdb
}

// GetMovieDB return *MovieDB pointer to struct.
func GetMovieDB() *MovieDB {
	return movieDB
}

// CreateMovieFromName get film name and find in tmdb.
// Returns *Film struct with filled fields.
func (mdb *MovieDB) CreateMovieFromName(name string) (*Film, error) {
	movie, err := mdb.tmdbApi.SearchMovie(name, mdb.options)
	if err != nil {
		return nil, fmt.Errorf("can't get movie from tmdb: "+
			"name - %s, err - %s", name, err)
	}
	if len(movie.Results) == 0 {
		return nil, fmt.Errorf("can't search movie from tmdb: "+
			"name - %s, err - %s", name, err)
	}
	// get only first result, other isn't necessary.
	searchedFilm := movie.Results[0]

	// for AddedDate field.
	t := time.Now()
	formattedTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(),
		t.Day(), t.Hour(),
		t.Minute(), t.Second())
	film := &Film{
		ID:           searchedFilm.ID,
		Name:         searchedFilm.Title,
		OriginalName: searchedFilm.OriginalTitle,
		Poster:       posterBaseURL + searchedFilm.PosterPath,
		ReleaseDate:  searchedFilm.ReleaseDate,
		Overview:     searchedFilm.Overview,
		AddedDate:    formattedTime,
		MagnetLinks:  make(map[string]string),
	}

	//add Genres to film
	for _, genre := range mdb.Genres {
		for _, searchedGenres := range searchedFilm.GenreIDs {
			if genre.ID == int(searchedGenres) {
				film.Genres = append(film.Genres, genre)
			}
		}
	}
	if len(film.Genres) == 0 {
		g := []*Genre{{ID: 0, EnglishName: "No Genre", RussianName: "Без жанра"}}
		film.Genres = g
	}

	return film, nil
}

// fillGenres create a Genres slice with ID, RussianName, EnglishName
func (mdb *MovieDB) fillGenres() {
	// get russian genres names
	dbGenre, err := mdb.tmdbApi.GetMovieGenres(mdb.options)
	if err != nil {
		fmt.Println(err.Error())
	}
	// adding ID and RussianName
	for _, dbGenre := range dbGenre.Genres {
		genre := &Genre{
			ID:          dbGenre.ID,
			RussianName: dbGenre.Name,
		}
		mdb.Genres = append(mdb.Genres, genre)
	}
	options := make(map[string]string)
	options["language"] = "en"

	// get english genres names
	dbGenre, err = mdb.tmdbApi.GetMovieGenres(options)
	if err != nil {
		fmt.Println(err.Error())
	}
	// adding EnglishName
	for idx, dbGenre := range dbGenre.Genres {
		if mdb.Genres[idx].ID == dbGenre.ID {
			mdb.Genres[idx].EnglishName = strings.ToLower(dbGenre.Name)
		}
	}
}
