package localstorage

type Container struct {
	ShortURLs    map[string]string
	OriginalURLs map[string]string
}

type LocalStorage struct {
	Users        map[string]Container
	ShortURLs    map[string]string
	OriginalURLs map[string]string
}

func NewLocalStorage() *LocalStorage {
	storage := LocalStorage{
		Users: map[string]Container{},
	}

	return &storage
}

func (storage *LocalStorage) AddURL(userID string, originalURL string, shortURL string) (ok bool) {
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
	return
}

func (storage *LocalStorage) GetURL(userID string, shortURL string) (originalURL string, ok bool) {
	if userID != "" {
		originalURL, ok = storage.Users[userID].ShortURLs[shortURL]
	} else {
		originalURL, ok = storage.ShortURLs[shortURL]
	}

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
