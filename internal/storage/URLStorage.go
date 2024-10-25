package storage

import (
	"github.com/WeisseNacht18/url-shortener/internal/logger"
	shortlinkgenerator "github.com/WeisseNacht18/url-shortener/internal/shortLinkGenerator"
	databasestorage "github.com/WeisseNacht18/url-shortener/internal/storage/databaseStorage"
	filestorage "github.com/WeisseNacht18/url-shortener/internal/storage/fileStorage"
	localstorage "github.com/WeisseNacht18/url-shortener/internal/storage/localStoarge"
)

type Storage interface {
	AddURL(string, string) bool
	GetURL(string) (string, bool)
	CheckStorage() error
	CheckURL(string) (string, bool)
	Close()
}

var storage Storage

func NewURLStorage(fileStoragePath string, databaseDSN string) {
	if databaseDSN != "" {
		storage = databasestorage.NewDatabaseStorage(databaseDSN)

	} else if fileStoragePath != "" {
		storage = filestorage.NewFileStorage(fileStoragePath)
	} else {
		storage = localstorage.NewLocalStorage()
	}
}

func NewEmptyURLStorage() {
	storage = localstorage.NewLocalStorage()
}

func NewURLStorageWithMap(shortUrls map[string]string) {
	storage = localstorage.NewLocalStorage()
	for shortURL, originalURL := range shortUrls {
		storage.AddURL(originalURL, shortURL)
	}
}

func AddURLToStorage(url string) (shortURL string, hasURL bool) {
	shortLink := shortlinkgenerator.GenerateShortLink()

	shortURL, hasURL = storage.CheckURL(url)

	if !hasURL {
		ok := storage.AddURL(url, shortLink)
		if !ok {
			logger.Logger.Fatalln("error: don't add url to storage")
		}
		shortURL = shortLink
	}

	return
}

func AddArrayOfURLToStorage(originalURLs map[string]string) (result map[string]string) {
	result = map[string]string{}

	for correlationID, originalURL := range originalURLs {
		shortLink := shortlinkgenerator.GenerateShortLink()
		result[correlationID] = shortLink
		ok := storage.AddURL(shortLink, originalURL)
		if !ok {
			logger.Logger.Fatalln("error: don't add url to storage")
		}
	}

	return
}

func GetURLFromStorage(shortURL string) (result string, ok bool) {
	return storage.GetURL(shortURL)
}

func CheckConnection() error {
	return storage.CheckStorage()
}

func Close() {
	storage.Close()
}
