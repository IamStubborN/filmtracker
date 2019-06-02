package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/IamStubborN/filmtracker/database"

	"github.com/IamStubborN/filmtracker/tmdb"

	"github.com/gin-gonic/gin"
	"github.com/kennygrant/sanitize"
)

// FetchFilmsFilter godoc
// @Summary Filter films
// @Description Filter films by name or/and genre, year, page
// @Tags api
// @Accept  json
// @Produce  json
// @Param name query string false "Film name"
// @Param year query string false "Film year production"
// @Param genre query string false "Film genre"
// @Param page query string false "Film page of the results"
// @Success 200 {array} tmdb.Film "Return slice of films with filter"
// @Failure 400 "{"error": "no films matches"}"
// @Security ApiKeyAuth Token, Refresh
// @Router /api/v1/films/filter [get]
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

// FetchPage godoc
// @Summary Fetch by page films
// @Description Filter films page
// @Tags api
// @Accept  json
// @Produce  json
// @Param ID path string true "Film page of the results"
// @Success 200 {array} tmdb.Film "Return slice of films with filter"
// @Failure 400 "{"error": "no films matches"/"bad page request"/ "invalid page number"}"
// @Security ApiKeyAuth Token, Refresh
// @securityDefinitions.ApiKeyAuth[Refresh,Token]
// @Router /api/v1/films/page/{ID} [get]
func FetchPage(c *gin.Context) {
	if pageQuery, ok := c.Params.Get("id"); ok {
		pageRegexp, err := regexp.Compile(`^[1-9][0-9]{0,4}$`)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "bad page request")
			return
		}
		if pageRegexp.MatchString(pageQuery) {
			films, err := db.GetFilmsByQuery(".", ".", ".", pageQuery)
			if err != nil {
				RespondWithError(c, http.StatusBadRequest, err.Error())
				return
			}
			c.JSON(http.StatusOK, films)
		} else {
			RespondWithError(c, http.StatusBadRequest, "invalid page number")
		}
	}
}

// UpdateFilm godoc
// @Summary Update film by json body
// @Description Update film by ID
// @Tags api
// @Accept  json
// @Produce  json
// @Param Film body tmdb.Film true "Update the film by fields"
// @Description "Only with user role = 'admin'"
// @Success 200 "{"ok": "film successful updated with that ID ___"}"
// @Failure 400 "{"error": "film isn't exist with that ID ___"}" string
// @Security ApiKeyAuth Token
// @Security ApiKeyAuth Refresh
// @Router /api/v1/films/ [put]
func UpdateFilm(c *gin.Context) {
	film := tmdb.Film{}
	if err := c.BindJSON(&film); err != nil {
		RespondWithError(c, http.StatusBadRequest, err)
		return
	}
	isExist, err := db.IsExistFilm(strconv.Itoa(film.ID))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err)
		return
	}
	if !isExist {
		RespondWithError(c, http.StatusBadRequest, fmt.Sprintf("film isn't exist with that ID:%d", film.ID))
		return
	} else {
		if err := db.UpsertFilm(&film); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"ok": fmt.Sprintf("film successful updated with that ID:%d", film.ID)})
	}
}

// AddFilm godoc
// @Summary Add film by json body
// @Description Add film json body
// @Tags api
// @Accept json
// @Produce json
// @Param Film body tmdb.Film true "Added the film by fields"
// @Description "Only with user role = 'admin'"
// @Success 200 "{"ok": "film successful added with that ID ___"}"
// @Failure 400 "{"error": "film is exist with that ___"}" string
// @Security ApiKeyAuth Token
// @Security ApiKeyAuth Refresh
// @Router /api/v1/films/ [post]
func AddFilm(c *gin.Context) {
	film := tmdb.Film{}
	if err := c.BindJSON(&film); err != nil {
		RespondWithError(c, http.StatusBadRequest, err)
	}
	isExist, err := db.IsExistFilm(strconv.Itoa(film.ID))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err)
		return
	}
	if isExist {
		RespondWithError(c, http.StatusBadRequest, fmt.Sprintf("film is exist with that ID:%d", film.ID))
		return
	} else {
		if err := db.InsertFilm(&film); err != nil {
			RespondWithError(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": fmt.Sprintf("film successful added with that ID:%d", film.ID)})
	}
}

// DeleteFilm godoc
// @Summary Delete film by ID
// @Description Delete film by ID
// @Tags api
// @Accept json
// @Produce json
// @Param ID path int true "Delete film by ID"
// @Success 200 "{"ok": "film with ID ___ deleted"}"
// @Failure 400 "{"error": "film ID doesn't exist in database: ___"}" string
// @Description "Only with user role = 'admin'"
// @Security ApiKeyAuth Token
// @Security ApiKeyAuth Refresh
// @Router /api/v1/films/film/{ID} [delete]
func DeleteFilm(c *gin.Context) {
	if id, ok := c.Params.Get("id"); ok {
		err := db.DeleteFilmByID(id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, fmt.Sprintf("film ID doesn't exist in database:%s", id))
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": fmt.Sprintf("film with ID %s deleted", id)})
	}
}

// FetchSingleFilm godoc
// @Summary Get single film by ID
// @Description Get single film by ID
// @Tags api
// @Accept json
// @Produce json
// @Param ID path int true "Get single film by ID"
// @Success 200 {object} tmdb.Film "Return single film by ID"
// @Failure 400 "{"error": "film ID doesn't exist in database: ____"}" string
// @Security ApiKeyAuth Token
// @Security ApiKeyAuth Refresh
// @Router /api/v1/films/film/{ID} [get]
func FetchSingleFilm(c *gin.Context) {
	if id, ok := c.Params.Get("id"); ok {
		film, err := db.GetFilmByID(id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, fmt.Sprintf("film ID doesn't exist in database:%d", film.ID))
			return
		}
		c.JSON(http.StatusOK, film)
	}
}

// FetchAllFilms godoc
// @Summary Get all films
// @Description Get all films
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {array} tmdb.Film "Return all films"
// @Failure 400 "{"error": "can't fetch from db films"}" string
// @Security ApiKeyAuth Token
// @Security ApiKeyAuth Refresh
// @Router /api/v1/films/ [get]
func FetchAllFilms(c *gin.Context) {
	films, err := db.GetAllFilms()
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "can't fetch from db films")
		return
	}
	c.JSON(http.StatusOK, films)
}

// ApiOverview godoc
// @Summary Api overview
// @Description Api overview
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} database.Overview "Return films and user count"
// @Failure 400 "{"error": "can't fetch from db films"}" string
// @Security ApiKeyAuth Token
// @Security ApiKeyAuth Refresh
// @Router /api/v1/ [get]
func ApiOverview(c *gin.Context) {
	filmsCount, err := db.GetCountFromCollection("films")
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err)
		return
	}
	userCount, err := db.GetCountFromCollection("users")
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err)
		return
	}
	overviewDatabase := database.Overview{
		Name:       "FilmTracker",
		FilmsCount: filmsCount,
		UsersCount: userCount,
	}
	c.JSON(http.StatusOK, overviewDatabase)
}
