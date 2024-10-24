package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/WeisseNacht18/url-shortener/internal/logger"
	shortlinkgenerator "github.com/WeisseNacht18/url-shortener/internal/shortLinkGenerator"
)

type URLStorage struct {
	shortUrls       map[string]string
	FileStoragePath string
	lastID          int
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

func NewURLStorage(fileStoragePath string) URLStorage {
	storage = URLStorage{}
	storage.FileStoragePath = fileStoragePath
	storage.shortUrls = GetConfigFromFile(fileStoragePath)
	storage.lastID = len(storage.shortUrls)
	return storage
}

func NewEmptyURLStorage() {
	storage = URLStorage{}
	storage.FileStoragePath = "storage.txt"
	storage.shortUrls = map[string]string{}
	storage.lastID = 0
}

func NewURLStorageWithMap(shortUrls map[string]string) {
	storage = URLStorage{}
	storage.shortUrls = shortUrls
	storage.FileStoragePath = "storage.txt"
	storage.lastID = len(shortUrls)
}

func AddURLToStorage(url string) (result string) {
	shortLink := shortlinkgenerator.GenerateShortLink()
	storage.shortUrls[shortLink] = url
	SaveLineToFile(shortLink, url)
	storage.lastID += 1
	return shortLink
}

func GetURLFromStorage(shortURL string) (result string, ok bool) {
	result, ok = storage.shortUrls[shortURL]
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
