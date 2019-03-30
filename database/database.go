package database

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	"github.com/IamStubborN/filmtracker/jwtmanager"

	"github.com/IamStubborN/filmtracker/tmdb"
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

type (
	Database struct {
		Films *mgo.Collection
		Users *mgo.Collection
	}

	User struct {
		UserID       string `bson:"user_id" json:"user_id"`
		Login        string `bson:"login" json:"login" validate:"min=3, max=10"`
		Password     string `bson:"password" json:"password" validate:"min=2, max=40"`
		Role         string `bson:"role" json:"role"`
		RefreshToken string `bson:"refresh_token" json:"refresh_token"`
	}

	filter struct {
		Name  string
		Genre string
		Year  string
		Page  string
	}
)

var database *Database
var jmg *jwtmanager.JwtManager
var mdb *tmdb.MovieDB

func init() {
	database = &Database{}
	jmg = jwtmanager.GetJWTManager()
	mdb = tmdb.GetMovieDB()
	session, err := mgo.Dial(os.Getenv("MONGO_HOST"))
	if err != nil {
		panic(err)
	}
	db := session.DB(os.Getenv("DB_NAME"))
	database.Films = db.C("films")
	database.Users = db.C("users")
	session.SetMode(mgo.Monotonic, true)
}

func GetDB() *Database {
	return database
}

func (database *Database) UpsertFilm(film *tmdb.Film) error {
	_, err := database.Films.Upsert(bson.M{"id": film.ID}, film)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) InsertFilm(film *tmdb.Film) error {
	err := database.Films.Insert(film)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) IsExistFilm(idStr string) (bool, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return false, err
	}
	count, err := database.Films.Find(bson.M{"id": id}).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (database *Database) GetCountFromCollection(nameCol string) (count int, err error) {
	switch nameCol {
	case "films":
		count, err = database.Films.Count()
	case "users":
		count, err = database.Users.Count()
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (database *Database) GetAllFilms() (films []*tmdb.Film, err error) {
	var film *tmdb.Film
	if err = database.Films.Find(bson.M{}).For(&film, func() error {
		films = append(films, film)
		return nil
	}); err != nil {
		return nil, err
	}
	return films, nil
}

func (database *Database) GetFilmByID(idStr string) (film *tmdb.Film, err error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	err = database.Films.Find(bson.M{"id": id}).One(&film)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err.Error())
	}
	return
}

func (database *Database) DeleteFilmByID(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	err = database.Films.Remove(bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("database error: %s", err.Error())
	}
	return nil
}

func (database *Database) getFilmsByFilter(filter *filter) ([]*tmdb.Film, error) {
	var films []*tmdb.Film
	var pageID int
	if filter.Page == "" {
		pageID = 1
	} else {
		page, err := strconv.Atoi(filter.Page)
		if err != nil {
			return nil, err
		}
		pageID = page
	}
	fmt.Println(filter.Genre)
	q := database.Films.Find(
		bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": filter.Name, "$options": "i"}},
				{"original_name": bson.M{"$regex": filter.Name, "$options": "i"}},
			},
			"release_date": bson.M{"$regex": filter.Year},
			"genres":       bson.M{"$elemMatch": bson.M{"russian_name": bson.M{"$regex": filter.Genre, "$options": "i"}}}}).
		Sort("-release_date").
		Limit(20)
	q = q.Skip((pageID - 1) * 20)

	if err := q.All(&films); err != nil {
		return nil, err
	}
	fmt.Println(films)
	if len(films) <= 0 {
		return nil, fmt.Errorf("no films matches")
	}
	return films, nil
}

func (database *Database) IsExistUser(login string) (bool, error) {
	count, err := database.Users.Find(bson.M{"login": login}).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (database *Database) GetUser(login string, password string) (*User, error) {
	var user User
	err := database.Users.Find(bson.M{"login": login}).One(&user)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err.Error())
	}
	return &user, nil
}

func (database *Database) CreateUser(login, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := User{
		UserID:   uuid.New().String(),
		Login:    login,
		Password: string(hashedPassword),
		Role:     "user",
	}
	if err := database.UpsertUser(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (database *Database) GetUserByUserID(userID string) (*User, error) {
	var user User
	err := database.Users.Find(bson.M{"user_id": userID}).One(&user)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err.Error())
	}
	return &user, nil
}

func (database *Database) UpdateUser(user *User) error {
	return database.Users.Update(bson.M{"user_id": user.UserID}, user)
}

func (database *Database) DeleteRefreshToken(refreshToken string) error {
	return database.Users.Update(bson.M{"refresh_token": refreshToken},
		bson.M{"$set": bson.M{"refresh_token": ""}})
}

func validateQuery(reg, query string) (string, error) {
	if query != "." {
		regex, err := regexp.Compile(reg)
		if err != nil {
			return "", fmt.Errorf("compile regex error: %s", regex)
		}
		if regex.MatchString(query) {
			return query, nil
		} else {
			return "", fmt.Errorf("bad request: %s", query)
		}
	}
	return "", nil
}

func createFilter(queries ...string) (*filter, error) {
	regexSlice := []string{
		`^[a-zA-ZА-Яа-я\s]{0,30}$`, //name
		`^[a-z]{3,12}$`,            //genre
		`^\d{4}$`,                  //year
		`^[1-9][0-9]{0,4}$`,        //page
	}
	name, err := validateQuery(regexSlice[0], queries[0])
	if err != nil {
		return nil, err
	}
	genre, err := validateQuery(regexSlice[1], queries[1])
	if err != nil {
		return nil, err
	}
	year, err := validateQuery(regexSlice[2], queries[2])
	if err != nil {
		return nil, err
	}
	page, err := validateQuery(regexSlice[3], queries[3])
	if err != nil {
		return nil, err
	}

	filter := &filter{
		Name:  name,
		Year:  year,
		Page:  page,
		Genre: genre,
	}
	for _, g := range mdb.Genres {
		if g.EnglishName == genre {
			filter.Genre = g.RussianName
		}
	}

	return filter, nil

}

func (database *Database) GetFilmsByQuery(queries ...string) ([]*tmdb.Film, error) {
	filter, err := createFilter(queries...)
	if err != nil {
		return nil, err
	}
	films, err := database.getFilmsByFilter(filter)
	if err != nil {
		return nil, err
	}
	return films, nil
}

func (database *Database) UpdateRefreshTokenForUser(userID, refreshToken string) error {
	if err := database.Users.Update(bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"refresh_token": refreshToken}}); err != nil {
		return err
	}
	return nil
}

func (database *Database) UpsertUser(user *User) error {
	_, err := database.Users.Upsert(bson.M{"user_id": user.UserID}, user)
	if err != nil {
		return err
	}
	return nil
}
