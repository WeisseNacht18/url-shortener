package filestorage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/WeisseNacht18/url-shortener/internal/logger"
)

type Container struct {
	ShortURLs    map[string]string
	OriginalURLs map[string]string
}

type FileStorage struct {
	Users        map[string]Container
	ShortURLs    map[string]string
	OriginalURLs map[string]string
	Path         string
	LastID       uint
}

func NewFileStorage(path string) *FileStorage {
	shortURLs, originalURLs, usersContainer := GetConfigFromFile(path)

	storage := FileStorage{
		Users:        usersContainer,
		ShortURLs:    shortURLs,
		OriginalURLs: originalURLs,
		Path:         path,
	}

	return &storage
}

func (storage *FileStorage) AddURL(userID string, originalURL string, shortURL string) (ok bool) {
	_, hasURL := storage.Users[userID].OriginalURLs[originalURL]
	if hasURL {
		ok = false
		return
	}
	ok = true
	storage.Users[userID] = Container{
		ShortURLs:    map[string]string{},
		OriginalURLs: map[string]string{},
	}
	storage.Users[userID].ShortURLs[shortURL] = originalURL
	storage.Users[userID].OriginalURLs[originalURL] = shortURL
	storage.ShortURLs[shortURL] = originalURL
	storage.OriginalURLs[originalURL] = shortURL
	storage.SaveLineToFile(userID, shortURL, originalURL)
	storage.LastID++
	return
}

func (storage *FileStorage) GetURL(userID string, shortURL string) (originalURL string, ok bool) {
	if userID != "" {
		originalURL, ok = storage.Users[userID].ShortURLs[shortURL]
	} else {
		originalURL, ok = storage.ShortURLs[shortURL]
	}
	return
}

func (storage *FileStorage) GetAllURLs(userID string) map[string]string {
	if userID != "" {
		return storage.Users[userID].ShortURLs
	} else {
		return storage.ShortURLs
	}
}

func (storage *FileStorage) CheckStorage() error {
	return nil
}

func (storage *FileStorage) CheckURL(userID string, originalURL string) (val string, ok bool) {
	if userID != "" {
		val, ok = storage.Users[userID].OriginalURLs[originalURL]
	} else {
		val, ok = storage.OriginalURLs[originalURL]
	}
	return
}

func (storage *FileStorage) GetUsers() map[string]string {
	result := map[string]string{}

	for userID := range storage.Users {
		result[userID] = ""
	}

	return result
}

func (storage *FileStorage) Close() {
	//supporting interface Storage
}

type URLStorageData struct {
	UUID        uint   `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
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

func GetConfigFromFile(filename string) (map[string]string, map[string]string, map[string]Container) {
	shortURLs := map[string]string{}
	originalURLs := map[string]string{}
	container := map[string]Container{}

	if _, err := os.Stat(filename); err != nil {
		return shortURLs, originalURLs, container
	}

	consumer, err := NewConsumer(filename)
	if err != nil {
		logger.Logger.Fatalln(err)
	}

	for {
		data, err := consumer.ReadStorageLine()
		if err != nil {
			break
		}
		shortURLs[data.ShortURL] = data.OriginalURL
		originalURLs[data.OriginalURL] = data.ShortURL
		container[data.UserID] = Container{
			ShortURLs:    map[string]string{},
			OriginalURLs: map[string]string{},
		}
		container[data.UserID].ShortURLs[data.ShortURL] = data.OriginalURL
		container[data.UserID].OriginalURLs[data.OriginalURL] = data.ShortURL
	}

	return shortURLs, originalURLs, container
}

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

func (storage *FileStorage) SaveLineToFile(userID string, shortURL string, url string) {
	lineData := URLStorageData{
		UUID:        storage.LastID,
		ShortURL:    shortURL,
		OriginalURL: url,
		UserID:      userID,
	}

	producer, err := NewProducer(storage.Path)
	if err != nil {
		logger.Logger.Fatalln(err)
	}

	producer.WriteStorageLine(&lineData)
}
