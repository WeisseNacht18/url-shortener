package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/WeisseNacht18/url-shortener/internal/logger"
	shortlinkgenerator "github.com/WeisseNacht18/url-shortener/internal/shortLinkGenerator"
)

type Storage struct {
	shortUrls       map[string]string
	FileStoragePath string
	lastId          int
}

type URLStorageData struct {
	UUID        uint   `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var (
	storage Storage
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

func New(fileStoragePath string) {
	storage = Storage{}
	storage.FileStoragePath = fileStoragePath
	storage.shortUrls = GetConfigFromFile(fileStoragePath)
	storage.lastId = len(storage.shortUrls)

}

func NewEmpty() {
	storage = Storage{}
	storage.FileStoragePath = "storage.txt"
	storage.shortUrls = map[string]string{}
	storage.lastId = 0
}

func NewWithMap(shortUrls map[string]string) { //переименовать на New
	storage = Storage{}
	storage.shortUrls = shortUrls
	storage.FileStoragePath = "storage.txt"
	storage.lastId = len(shortUrls)
}

func AddURLToStorage(url string) (result string) {
	shortLink := shortlinkgenerator.GenerateShortLink()
	storage.shortUrls[shortLink] = url
	SaveLineToFile(shortLink, url)
	storage.lastId += 1
	return shortLink
}

func GetURLFromStorage(shortURL string) (result string, ok bool) {
	result, ok = storage.shortUrls[shortURL]
	return
}

func SaveLineToFile(shortURL string, url string) {

	lineData := URLStorageData{
		UUID:        uint(storage.lastId),
		ShortURL:    shortURL,
		OriginalURL: url,
	}

	logger.Logger.Infoln(storage.FileStoragePath)
	producer, err := NewProducer(storage.FileStoragePath)
	if err != nil {
		logger.Logger.Fatalln(err)
	}

	producer.WriteStorageLine(&lineData)
}

func GetConfigFromFile(filename string) map[string]string {
	//проверить существует ли файл конфига. Если существует, то строчка за строчкой вычитать все хранилище

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
