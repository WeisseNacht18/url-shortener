package storage

import (
	"log"

	"github.com/WeisseNacht18/url-shortener/internal/generator"
	"github.com/WeisseNacht18/url-shortener/internal/logger"
	databasestorage "github.com/WeisseNacht18/url-shortener/internal/storage/databaseStorage"
	filestorage "github.com/WeisseNacht18/url-shortener/internal/storage/fileStorage"
	localstorage "github.com/WeisseNacht18/url-shortener/internal/storage/localStoarge"
)

type Storage interface {
	AddURL(string, string, string) bool
	GetURL(string, string) (string, bool)
	GetAllURLs(string) map[string]string
	CheckStorage() error
	CheckURL(string, string) (string, bool)
	Close()
}

var storage Storage
var userTokens map[string]string

func NewURLStorage(fileStoragePath string, databaseDSN string) {
	if databaseDSN != "" {
		storage = databasestorage.NewDatabaseStorage(databaseDSN)

	} else if fileStoragePath != "" {
		storage = filestorage.NewFileStorage(fileStoragePath)
	} else {
		storage = localstorage.NewLocalStorage()
	}
	userTokens = map[string]string{}
}

func NewEmptyURLStorage() {
	storage = localstorage.NewLocalStorage()
	userTokens = map[string]string{}
}

func NewURLStorageWithMap(shortUrls map[string]string) {
	storage = localstorage.NewLocalStorage()
	for shortURL, originalURL := range shortUrls {
		storage.AddURL("0", originalURL, shortURL)
	}
	userTokens = map[string]string{}
}

func AddURLToStorage(userID string, url string) (shortURL string, hasURL bool) {
	shortLink := generator.GenerateShortLink()

	shortURL, hasURL = storage.CheckURL(userID, url)

	if !hasURL {
		ok := storage.AddURL(userID, url, shortLink)
		if !ok {
			logger.Logger.Fatalln("error: don't add url to storage")
		}
		shortURL = shortLink
	}

	return
}

func AddArrayOfURLToStorage(userID string, originalURLs map[string]string) (result map[string]string) {
	result = map[string]string{}

	for correlationID, originalURL := range originalURLs {
		shortLink := generator.GenerateShortLink()
		result[correlationID] = shortLink
		ok := storage.AddURL(userID, originalURL, shortLink)
		if !ok {
			logger.Logger.Fatalln("error: don't add url to storage")
		}
	}

	return
}

func GetURLFromStorage(userID string, shortURL string) (result string, ok bool) {
	return storage.GetURL(userID, shortURL)
}

func GetAllURLsFromStorage(userID string) map[string]string {
	return storage.GetAllURLs(userID)
}

func CheckConnection() error {
	return storage.CheckStorage()
}

func CheckUserID(userID string) bool {
	_, ok := userTokens[userID]
	return ok
}

func CheckUserIDWithToken(userID string, token string) bool {
	return userTokens[userID] == token
}

func AddUserID(userID string) bool {
	if !CheckUserID(userID) {
		userTokens[userID] = ""
	}
	return false
}

func AddUserIDWithToken(userID string, token string) bool {
	log.Println("has user:", CheckUserID(userID))
	if !CheckUserID(userID) {
		userTokens[userID] = token
		return true
	}
	return false
}

func Close() {
	storage.Close()
}
