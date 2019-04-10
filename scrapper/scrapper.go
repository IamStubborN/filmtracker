package scrapper

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/IamStubborN/filmtracker/database"

	"github.com/IamStubborN/filmtracker/tmdb"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

// updateProxyAndUATimeout - duration timeout.
const updateProxyAndUATimeout = 1

type (
	FilmTracker struct {
		URL            string
		URLCategory    string
		PostfixURL     string
		ContainerClass string
		EntryDetails   string
		FilmName       string
		PageCount      string
	}
	Scrapper struct {
		mainCollector    *colly.Collector
		tmdbCollector    *colly.Collector
		torrentCollector *colly.Collector
		MovieDB          *tmdb.MovieDB
		listUA           []string
		isUseProxy       bool
		sync.RWMutex
	}
)

var db *database.Database

// createScrapper func
func CreateScrapper(allowedDomains []string, isUseProxy bool) (*Scrapper, error) {
	scrapper := Scrapper{
		isUseProxy: isUseProxy,
		MovieDB:    tmdb.GetMovieDB(),
	}
	scrapper.mainCollector = colly.NewCollector(
		colly.DetectCharset(),
		colly.AllowURLRevisit(),
		colly.AllowedDomains(allowedDomains...),
	)

	scrapper.torrentCollector = colly.NewCollector(
		colly.DetectCharset(),
		colly.AllowURLRevisit(),
		colly.AllowedDomains("rutor.info"),
	)

	//supportCollector := colly.NewCollector(
	//	colly.AllowedDomains("www.proxy-list.download", "www.ua-tracker.com"),
	//	colly.DetectCharset(),
	//	colly.AllowURLRevisit(),
	//)

	scrapper.mainCollector.SetRequestTimeout(30 * time.Second)
	scrapper.tmdbCollector = scrapper.mainCollector.Clone()

	proxyFunc, err := scrapper.getProxyFunc(scrapper.mainCollector)
	if scrapper.isUseProxy {
		if err != nil {
			return nil, err
		}
		scrapper.mainCollector.SetProxyFunc(proxyFunc)
		scrapper.tmdbCollector.SetProxyFunc(proxyFunc)
	}
	// setProxy to torrentCollector is necessary, because rutor can ban you.
	scrapper.torrentCollector.SetProxyFunc(proxyFunc)
	list, err := scrapper.getUserAgentsList(scrapper.mainCollector)
	if err != nil {
		return nil, err
	}
	scrapper.listUA = list
	db = database.GetDB()
	go scrapper.ChangeUAWithTimeout(1)
	//go scrapper.UpdateProxyAndUAListWithTimeout(supportCollector)
	return &scrapper, nil
}

// ChangeUAWithTimeout - updating UA of all collectors.
func (scrapper *Scrapper) ChangeUAWithTimeout(changingTimeout time.Duration) {
	rand.Seed(time.Now().Unix())
	for range time.NewTicker(time.Duration(changingTimeout * time.Minute)).C {
		scrapper.mainCollector.UserAgent = scrapper.listUA[getRandomItem(len(scrapper.listUA))]
		scrapper.tmdbCollector.UserAgent = scrapper.listUA[getRandomItem(len(scrapper.listUA))]
		scrapper.torrentCollector.UserAgent = scrapper.listUA[getRandomItem(len(scrapper.listUA))]
	}
}

// UpdateProxyAndUAListWithTimeout - updating UA and Proxy's with timeout.
// This is necessary to avoid a ban because of the many requests.
func (scrapper *Scrapper) UpdateProxyAndUAListWithTimeout(supportCollector *colly.Collector) {
	fmt.Println("In function")
	for range time.NewTicker(time.Duration(updateProxyAndUATimeout * time.Minute)).C {
		fmt.Println("in for range")
		listUA, err := scrapper.getUserAgentsList(supportCollector)
		if err != nil {
			log.Printf("can't update useragent list")
		}
		scrapper.RLock()
		scrapper.listUA = listUA
		scrapper.RUnlock()

		proxyFunc, err := scrapper.getProxyFunc(supportCollector)
		if scrapper.isUseProxy {
			if err != nil {
				log.Printf("can't update proxy list")
			}
			scrapper.mainCollector.SetProxyFunc(proxyFunc)
			scrapper.tmdbCollector.SetProxyFunc(proxyFunc)
		}
		log.Println("Proxy's changed")
		scrapper.torrentCollector.SetProxyFunc(proxyFunc)
	}
}

// getRandomItem return random item from slice by length.
func getRandomItem(len int) int {
	return int(math.Abs(float64(rand.Intn(len - 1))))
}

// getPagesCount from torrent tracker.
func (scrapper *Scrapper) getPagesCount(tor *FilmTracker) (count int, err error) {
	scrapper.mainCollector.OnHTML(tor.PageCount, func(e *colly.HTMLElement) {
		count, err = strconv.Atoi(e.DOM.Children().Last().Text())
	})
	scrapper.visitWithRetry(scrapper.mainCollector, tor.URL+tor.URLCategory, 30)
	return
}

// getProxyFunc generate proxy.RoundRobinSwitcher and return it.
func (scrapper *Scrapper) getProxyFunc(collector *colly.Collector) (colly.ProxyFunc, error) {
	proxyList, err := scrapper.getProxyList(collector)
	if err != nil {
		return nil, err
	}
	proxyFunc, err := proxy.RoundRobinProxySwitcher(proxyList...)
	if err != nil {
		return nil, err
	}
	return proxyFunc, nil
}

// getProxyList get proxy list from https://www.proxy-list.download/api/v1/get?type=http
func (scrapper *Scrapper) getProxyList(collector *colly.Collector) ([]string, error) {
	const api = "https://www.proxy-list.download/api/v1/get?type=http"
	var proxyList []string
	collector.OnResponse(func(response *colly.Response) {
		if response.StatusCode == 200 {
			for _, pr := range strings.Split(string(response.Body), "\r\n") {
				proxyList = append(proxyList, fmt.Sprintf(`http://%s`, pr))
			}
		}
	})
	scrapper.visitWithRetry(collector, api, 30)
	return proxyList, nil
}

// getUserAgentsList get UA list from http://www.ua-tracker.com/user_agents.txt
func (scrapper *Scrapper) getUserAgentsList(collector *colly.Collector) ([]string, error) {
	const api = "http://www.ua-tracker.com/user_agents.txt"
	var uaList []string
	collector.OnResponse(func(response *colly.Response) {
		if response.StatusCode == 200 {
			for _, ua := range strings.Split(string(response.Body), "\n") {
				uaList = append(uaList, ua)
			}
		}
	})
	scrapper.visitWithRetry(collector, api, 30)
	return uaList, nil
}

// GetAllMovies from torrent tracker.
func (scrapper *Scrapper) GetAllMovies(tor *FilmTracker, isFullUpdate bool) {
	baseURL := tor.URL
	pagesURLs := make([]string, 0, 20)
	pagesURLs = append(pagesURLs, baseURL+tor.URLCategory)
	if isFullUpdate {
		pagesCount, err := scrapper.getPagesCount(tor)
		if err != nil {
			log.Println(err)
		}
		for idx := 1; idx <= pagesCount; idx++ {
			pageURL := baseURL + tor.URLCategory + tor.PostfixURL + strconv.Itoa(idx) + "/"
			pagesURLs = append(pagesURLs, pageURL)
		}
	}
	filmURLs := make([]string, 0, 20)
	scrapper.mainCollector.OnHTML(tor.EntryDetails, func(e *colly.HTMLElement) {
		detailsURL := e.Attr("href")
		filmURLs = append(filmURLs, detailsURL)
	})
	scrapper.mainCollector.OnRequest(func(request *colly.Request) {
		fmt.Println("main ", request.AbsoluteURL(request.URL.Path))
	})

	for _, URL := range pagesURLs {
		scrapper.visitWithRetry(scrapper.mainCollector, URL, 30)
	}
	if err := scrapper.getDetailsAboutMovies(filmURLs, tor.FilmName); err != nil {
		log.Println(err)
	}
}

// visitWithRetry try visit with count.
func (scrapper *Scrapper) visitWithRetry(collector *colly.Collector, URL string, retryCount int) {
	if err := collector.Visit(URL); err != nil {
		count := 1
		fmt.Println("Try", count, err)
		for count <= retryCount {
			count++
			if err := collector.Visit(URL); err != nil {
				fmt.Println("Try", count, err)
			} else {
				break
			}
		}
	}
}

// getDetailsAboutMovies
func (scrapper *Scrapper) getDetailsAboutMovies(URLs []string, filmNameSelector string) (err error) {
	var film *tmdb.Film
	var releaseYear string

	scrapper.tmdbCollector.OnHTML(filmNameSelector, func(e *colly.HTMLElement) {
		name := strings.TrimSpace(strings.Split(e.Text, "(")[0])
		film, err = scrapper.MovieDB.CreateMovieFromName(name)
		if err != nil {
			return
		}
		releaseYear = strings.Split(film.ReleaseDate, "-")[0]
		torURL := "http://rutor.info/search/0/0/000/0/" + film.Name
		scrapper.visitWithRetry(scrapper.torrentCollector, torURL, 30)
	})

	scrapper.tmdbCollector.OnRequest(func(request *colly.Request) {
		fmt.Println("tmdb ", request.AbsoluteURL(request.URL.Path))
	})

	re := regexp.MustCompile("TS|BDRip|HDRip|BDRemux|UHD|DVD|WEB")
	reWebTorrent := regexp.MustCompile("WEB|HDRip")
	scrapper.torrentCollector.OnHTML("#index tbody tr", func(e *colly.HTMLElement) {
		name := e.Text
		if strings.Contains(strings.ToLower(name), strings.ToLower(film.Name)) &&
			strings.Contains(name, releaseYear) &&
			re.MatchString(name) {
			a := e.DOM.Find("a")
			name := strings.TrimSpace(a.Text())
			link := a.Get(1).Attr[0].Val
			film.MagnetLinks[name] = link
			fmt.Println(reWebTorrent.MatchString(name))
			fmt.Println(film.WebTorrentMagnet)
			fmt.Println(name)
			if film.WebTorrentMagnet == "" && reWebTorrent.MatchString(name) {
				fmt.Println("12312414", link)
				film.WebTorrentMagnet = link
			}
		}
	})
	scrapper.torrentCollector.OnScraped(func(response *colly.Response) {
		if len(film.MagnetLinks) > 0 {
			fmt.Println(film)
			if err = db.UpsertFilm(film); err != nil {
				return
			}
		}
	})

	scrapper.torrentCollector.OnRequest(func(request *colly.Request) {
		fmt.Println("torrents ", request.AbsoluteURL(request.URL.Path))
	})

	for idx := range URLs {
		scrapper.visitWithRetry(scrapper.tmdbCollector, URLs[idx], 30)
	}
	return nil
}
