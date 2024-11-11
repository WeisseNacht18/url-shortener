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
		Users:        map[string]Container{},
		ShortURLs:    map[string]string{},
		OriginalURLs: map[string]string{},
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
	_, hasUser := storage.Users[userID]
	if !hasUser {
		storage.Users[userID] = Container{
			ShortURLs:    map[string]string{},
			OriginalURLs: map[string]string{},
		}
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
	if userID != "" {
		return storage.Users[userID].ShortURLs
	} else {
		return storage.ShortURLs
	}
}

func (storage *LocalStorage) CheckStorage() error {
	return nil
}

func (storage *LocalStorage) CheckURL(userID, originalURL string) (val string, ok bool) {
	if userID != "" {
		val, ok = storage.Users[userID].OriginalURLs[originalURL]
	} else {
		val, ok = storage.OriginalURLs[originalURL]
	}
	return
}

func (storage *LocalStorage) GetUsers() map[string]string {
	return map[string]string{}
}

func (storage *LocalStorage) Close() {
	//supporting interface Storage
}
