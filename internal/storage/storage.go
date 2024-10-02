package storage

import (
	shortlinkgenerator "github.com/WeisseNacht18/url-shortener/internal/shortLinkGenerator"
)

type Storage struct {
	shortUrls map[string]string
}

var (
	storage Storage
)

func Init() {
	storage = Storage{}
	storage.shortUrls = map[string]string{}
}

func InitWithMap(shortUrls map[string]string) { //переименовать на New
	storage = Storage{}
	storage.shortUrls = shortUrls
}

func AddURLToStorage(url string) (result string) {
	shortLink := shortlinkgenerator.GenerateShortLink()
	storage.shortUrls[shortLink] = url
	return shortLink
}

func GetURLFromStorage(shortURL string) (result string, ok bool) {
	result, ok = storage.shortUrls[shortURL]
	return
}
