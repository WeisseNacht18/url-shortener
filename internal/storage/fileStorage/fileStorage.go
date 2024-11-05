package filestorage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/WeisseNacht18/url-shortener/internal/logger"
)

type FileStorage struct {
	ShortURLs    map[string]string
	originalURLs map[string]string
	Path         string
	LastID       int
}

func NewFileStorage(path string) *FileStorage {
	fromFile := GetConfigFromFile(path)
	reverseFromFile := changeKV(fromFile)

	storage := FileStorage{
		ShortURLs:    fromFile,
		originalURLs: reverseFromFile,
		Path:         path,
	}

	return &storage
}

func (storage *FileStorage) AddURL(originalURL string, shortURL string) (ok bool) {
	_, hasURL := storage.originalURLs[originalURL]
	if hasURL {
		ok = false
		return
	}
	ok = true
	storage.ShortURLs[shortURL] = originalURL
	storage.originalURLs[originalURL] = shortURL
	storage.SaveLineToFile(shortURL, originalURL)
	return
}

func (storage *FileStorage) GetURL(shortURL string) (originalURL string, ok bool) {
	originalURL, ok = storage.ShortURLs[shortURL]
	return
}

func (storage *FileStorage) GetAllURLs(userID int) map[string]string {
	return storage.ShortURLs
}

func (storage *FileStorage) CheckStorage() error {
	return nil
}

func (storage *FileStorage) CheckURL(originalURL string) (string, bool) {
	val, ok := storage.originalURLs[originalURL]
	return val, ok
}

func (storage *FileStorage) Close() {
	//supporting interface Storage
}

type URLStorageData struct {
	UUID        uint   `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
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

func GetConfigFromFile(filename string) map[string]string {
	result := map[string]string{}

	if _, err := os.Stat(filename); err != nil {
		return result
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
		result[data.ShortURL] = data.OriginalURL
	}

	return result
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

func (storage *FileStorage) SaveLineToFile(shortURL string, url string) {
	lineData := URLStorageData{
		UUID:        uint(storage.LastID),
		ShortURL:    shortURL,
		OriginalURL: url,
	}

	producer, err := NewProducer(storage.Path)
	if err != nil {
		logger.Logger.Fatalln(err)
	}

	producer.WriteStorageLine(&lineData)
}

func changeKV(in map[string]string) (out map[string]string) {
	out = make(map[string]string)
	for k, v := range in {
		out[v] = k
	}
	return
}
