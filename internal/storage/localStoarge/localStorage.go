package localstorage

type Container struct {
	ShortURLs    map[string]string
	OriginalURLs map[string]string
}

type LocalStorage struct {
	Users map[string]Container
}

func NewLocalStorage() *LocalStorage {
	storage := LocalStorage{
		Users: map[string]Container{},
	}

	return &storage
}

func (storage *LocalStorage) AddURL(userID string, originalURL string, shortURL string) (ok bool) {
	ok = true
	storage.Users[userID].ShortURLs[shortURL] = originalURL
	storage.Users[userID].OriginalURLs[originalURL] = shortURL
	return
}

func (storage *LocalStorage) GetURL(userID string, shortURL string) (originalURL string, ok bool) {
	originalURL, ok = storage.Users[userID].ShortURLs[shortURL]
	return
}

func (storage *LocalStorage) GetAllURLs(userID string) map[string]string {
	return storage.Users[userID].ShortURLs
}

func (storage *LocalStorage) CheckStorage() error {
	return nil
}

func (storage *LocalStorage) CheckURL(userID, originalURL string) (string, bool) {
	val, ok := storage.Users[userID].OriginalURLs[originalURL]
	return val, ok
}

func (storage *LocalStorage) Close() {
	//supporting interface Storage
}
