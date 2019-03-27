package scrapper

import (
	"testing"
)

func TestCreateScrapper(t *testing.T) {
	sc, err := CreateScrapper([]string{"toreents.club", "newtorrent.org",
		"www.proxy-list.download", "www.ua-tracker.com"}, false)
	if err != nil {
		t.Error(err)
	}
	if sc.mainCollector == nil || len(sc.listUA) == 0 {
		t.Error("Creating TestCreateScrapper not Passed")
	}
}
