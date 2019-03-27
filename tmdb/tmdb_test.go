package tmdb

import (
	"testing"
)

var testingFilm = &Film{
	Name: "Алита: Боевой Ангел",
	Overview: "XXVI век. На планете Земля богачи живут в Небесном городе, " +
		"беднота довольствуется жизнью в Нижнем городе, куда сбрасываются " +
		"все отходы и мусор. Однажды в куче металлолома Нижнего города учёный " +
		"находит части женщины-киборга и возвращает её к жизни. Придя в сознание, " +
		"киборг обнаруживает, что из ее памяти стерто все, кроме боевых приемов. " +
		"Теперь она должна обрести утерянные воспоминания и выяснить, " +
		"кто её отправил на свалку.",
}

func TestCreateMovieDB(t *testing.T) {
	mdb := GetMovieDB()
	if mdb == nil {
		t.Error("Creating TestCreateMovieDB not Passed")
	}
}

func TestMovieDB_CreateMovieFromName(t *testing.T) {
	mdb := GetMovieDB()
	film, err := mdb.CreateMovieFromName("Алита")
	if err != nil {
		t.Error(err)
	}
	if film.Name != testingFilm.Name || film.Overview != testingFilm.Overview {
		t.Error("Creating TestMovieDB_CreateMovieFromName not Passed")
	}
}

func TestMovieDB_fillGenres(t *testing.T) {
	mdb := GetMovieDB()
	mdb.fillGenres()
	if len(mdb.Genres) != 38 {
		t.Error("Filling TestMovieDB_FillGenres not Passed")
	}
}
