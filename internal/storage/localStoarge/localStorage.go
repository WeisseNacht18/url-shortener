package localstorage

type LocalStorage struct {
	ShortURLs    map[string]string
	originalURLs map[string]string
}

func NewLocalStorage() *LocalStorage {
	storage := LocalStorage{
		ShortURLs:    make(map[string]string),
		originalURLs: make(map[string]string),
	}

	return &storage
}

func (storage *LocalStorage) AddURL(userID string, originalURL string, shortURL string) (ok bool) {
	ok = true
	storage.ShortURLs[shortURL] = originalURL
	storage.originalURLs[originalURL] = shortURL
	return
}

func (storage *LocalStorage) GetURL(userID string, shortURL string) (originalURL string, ok bool) {
	originalURL, ok = storage.ShortURLs[shortURL]
	return
}

func (storage *LocalStorage) GetAllURLs(userID int) map[string]string {
	return storage.ShortURLs
}

func (storage *LocalStorage) CheckStorage() error {
	return nil
}

func (storage *LocalStorage) CheckURL(originalURL string) (string, bool) {
	val, ok := storage.originalURLs[originalURL]
	return val, ok
}

func (storage *LocalStorage) Close() {
	//supporting interface Storage
}
