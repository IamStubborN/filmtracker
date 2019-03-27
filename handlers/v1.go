package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/IamStubborN/filmtracker/tmdb"

	"github.com/gin-gonic/gin"
	"github.com/kennygrant/sanitize"
)

func FetchFilmsFilter(c *gin.Context) {
	nameQuery := sanitize.HTML(c.Query("name"))
	genreQuery := sanitize.Name(c.Query("genre"))
	yearQuery := sanitize.Name(c.Query("year"))
	pageQuery := sanitize.Name(c.Query("page"))

	films, err := db.GetFilmsByQuery(nameQuery, genreQuery, yearQuery, pageQuery)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, films)
}

func FetchPage(c *gin.Context) {
	if pageQuery, ok := c.Params.Get("id"); ok {
		pageRegexp, err := regexp.Compile(`^[1-9][0-9]{0,4}$`)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "bad Page request")
			return
		}
		if pageRegexp.MatchString(pageQuery) {
			films, err := db.GetFilmsByQuery(pageQuery)
			if err != nil {
				log.Fatal(err)
			}
			c.JSON(http.StatusOK, films)
		} else {
			RespondWithError(c, http.StatusBadRequest, "invalid Page number")
		}
	}
}

func UpdateFilm(c *gin.Context) {
	film := tmdb.Film{}
	if err := c.BindJSON(&film); err != nil {
		log.Fatal(err)
	}
	isExist, err := db.IsExistFilm(strconv.Itoa(film.ID))
	if err != nil {
		log.Fatal(err)
	}
	if !isExist {
		c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("Film isn't exist with that id:%d", film.ID)})
	} else {
		if err := db.UpsertFilm(&film); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"ok": fmt.Sprintf("Film successful updated with that id:%d", film.ID)})
	}
}

func AddFilm(c *gin.Context) {
	film := tmdb.Film{}
	if err := c.BindJSON(&film); err != nil {
		log.Fatal(err)
	}
	isExist, err := db.IsExistFilm(strconv.Itoa(film.ID))
	if err != nil {
		log.Fatal(err)
	}
	if isExist {
		c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("Film is exist with that id:%d", film.ID)})
	} else {
		if err := db.InsertFilm(&film); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"ok": fmt.Sprintf("Film successful added with that id:%d", film.ID)})
	}
}

func DeleteFilm(c *gin.Context) {
	if id, ok := c.Params.Get("id"); ok {
		err := db.DeleteFilmByID(id)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "film id doesn't exist in database"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": fmt.Sprintf("film with id %s deleted", id)})
	}
}

func FetchSingleFilm(c *gin.Context) {
	if id, ok := c.Params.Get("id"); ok {
		film, err := db.GetFilmByID(id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "film id doesn't exist in database"})
			return
		}
		c.JSON(http.StatusOK, film)
	}
}

func FetchAllFilms(c *gin.Context) {
	films, err := db.GetAllFilms()
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, films)
}

func StartPage(c *gin.Context) {
	filmsCount, err := db.GetCountFromCollection("films")
	if err != nil {
		log.Fatal(filmsCount)
	}
	userCount, err := db.GetCountFromCollection("users")
	if err != nil {
		log.Fatal(filmsCount)
	}
	overviewDatabase := struct {
		Name       string
		FilmsCount int
		UsersCount int
	}{
		Name:       "FilmTracker",
		FilmsCount: filmsCount,
		UsersCount: userCount,
	}
	c.JSON(http.StatusOK, overviewDatabase)
}