package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/WeisseNacht18/url-shortener/internal/database"
	"github.com/WeisseNacht18/url-shortener/internal/logger"
	shortlinkgenerator "github.com/WeisseNacht18/url-shortener/internal/shortLinkGenerator"
)

const (
	DatabaseStorage = "database_storage"
	FileStorage     = "file_storage"
	LocalStorage    = "local_storage"
)

type URLStorage struct {
	shortUrls       map[string]string
	FileStoragePath string
	lastID          int
	Type            string
}

type URLStorageData struct {
	UUID        uint   `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var (
	storage URLStorage
)

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteStorageLine(storageLine *URLStorageData) error {
	data, err := json.Marshal(&storageLine)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

type Consumer struct {
	file   *os.File
	reader *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

func (c *Consumer) ReadStorageLine() (*URLStorageData, error) {
	data, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	line := URLStorageData{}
	err = json.Unmarshal(data, &line)
	if err != nil {
		return nil, err
	}

	return &line, nil
}

func NewURLStorage(fileStoragePath string, databaseDSN string) URLStorage {
	storage = URLStorage{}

	storage.Type = LocalStorage
	storage.shortUrls = map[string]string{}

	if fileStoragePath != "" {
		storage.Type = FileStorage
		storage.FileStoragePath = fileStoragePath
		logger.Logger.Infoln(fileStoragePath)
	}

	if databaseDSN != "" {
		storage.Type = DatabaseStorage
		logger.Logger.Infoln(databaseDSN)
	}

	if storage.Type == FileStorage {
		storage.shortUrls = GetConfigFromFile(fileStoragePath)
	}

	storage.lastID = len(storage.shortUrls)
	return storage
}

func NewEmptyURLStorage() {
	storage = URLStorage{}
	storage.Type = LocalStorage
	storage.FileStoragePath = "storage.txt"
	storage.shortUrls = map[string]string{}
	storage.lastID = 0
}

func NewURLStorageWithMap(shortUrls map[string]string) {
	storage = URLStorage{}
	storage.Type = LocalStorage
	storage.shortUrls = shortUrls
	storage.FileStoragePath = "storage.txt"
	storage.lastID = len(shortUrls)
}

func AddURLToStorage(url string) (result string) {
	shortLink := shortlinkgenerator.GenerateShortLink()

	if storage.Type == FileStorage {
		SaveLineToFile(shortLink, url)
	}

	if storage.Type == LocalStorage || storage.Type == FileStorage {
		storage.shortUrls[shortLink] = url
		storage.lastID += 1
	}

	if storage.Type == DatabaseStorage {
		err := SaveURLToDatabase(shortLink, url)
		if err != nil {
			logger.Logger.Infoln(err)
		}
	}

	return shortLink
}

func AddArrayOfURLToStorage(originalURLs map[string]string) (result map[string]string) {
	urls := map[string]string{}
	result = map[string]string{}

	for correlationID, originalURL := range originalURLs {
		shortLink := shortlinkgenerator.GenerateShortLink()
		result[correlationID] = shortLink
		urls[originalURL] = shortLink
	}

	if storage.Type == FileStorage {
		for originalURL, shortURL := range urls {
			SaveLineToFile(shortURL, originalURL)
		}
	}

	if storage.Type == LocalStorage || storage.Type == FileStorage {
		for originalURL, shortURL := range urls {
			storage.shortUrls[shortURL] = originalURL
			storage.lastID += 1
		}
	}

	if storage.Type == DatabaseStorage {
		for originalURL, shortURL := range urls {
			err := SaveURLToDatabase(shortURL, originalURL)
			if err != nil {
				logger.Logger.Infoln(err)
			}
		}
	}

	return
}

func GetURLFromStorage(shortURL string) (result string, ok bool) {
	if storage.Type == LocalStorage || storage.Type == FileStorage {
		result, ok = storage.shortUrls[shortURL]
	}

	if storage.Type == DatabaseStorage {
		result, ok = GetURLFromDatabase(shortURL)
		if !ok {
			result, ok = storage.shortUrls[shortURL]
		}
	}

	return
}

func SaveLineToFile(shortURL string, url string) {

	lineData := URLStorageData{
		UUID:        uint(storage.lastID),
		ShortURL:    shortURL,
		OriginalURL: url,
	}

	producer, err := NewProducer(storage.FileStoragePath)
	if err != nil {
		logger.Logger.Fatalln(err)
	}

	producer.WriteStorageLine(&lineData)
}

func GetConfigFromFile(filename string) map[string]string {
	result := map[string]string{}

	if _, err := os.Stat(storage.FileStoragePath); err != nil {
		return result
	}

	consumer, err := NewConsumer(storage.FileStoragePath)
	if err != nil {
		logger.Logger.Fatalln(err)
	}

	for {
		data, err := consumer.ReadStorageLine()
		if err != nil {
			break
		}
		result[data.ShortURL] = data.OriginalURL
	}

	return result
}

func SaveURLToDatabase(shortURL string, originalURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := database.Database.ExecContext(ctx, "INSERT INTO url (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL)

	return err
}

func GetURLFromDatabase(shortURL string) (result string, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var originalURL string

	row := database.Database.QueryRowContext(ctx, "SELECT original_url FROM url WHERE short_url = $1 LIMIT 1", shortURL)

	err := row.Scan(&originalURL)
	if err != nil {
		return originalURL, false
	}

	return originalURL, true
}
