package main

import (
	"URLShortener/model"
	"fmt"
)

type InMemoryStore struct {
	mapLongUrlToShortUrl map[model.LongURL]model.ShortURL
	mapShortUrlToLongUrl map[model.ShortURL]model.LongURL
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		mapLongUrlToShortUrl: make(map[model.LongURL]model.ShortURL),
		mapShortUrlToLongUrl: make(map[model.ShortURL]model.LongURL),
	}
}

/*AddURL Adds www.facebook.com/vxDant/profile -> tinyurl.com/abc123
 */
func (s *InMemoryStore) AddURL(longUrl model.LongURL, shortUrl model.ShortURL) {
	s.mapLongUrlToShortUrl[longUrl] = shortUrl
	s.mapShortUrlToLongUrl[shortUrl] = longUrl

	fmt.Println(s.mapLongUrlToShortUrl)
}

/*GetShortURL returns a short url from the map */
func (s *InMemoryStore) GetShortURL(longUrl model.LongURL) model.ShortURL {
	return s.mapLongUrlToShortUrl[longUrl]
}

/*GetLongURL returns a long url from the map */
func (s *InMemoryStore) GetLongURL(shortUrl model.ShortURL) model.LongURL {
	return s.mapShortUrlToLongUrl[shortUrl]
}

func (s *InMemoryStore) ifLongURLExist(longUrl model.LongURL) bool {
	_, exists := s.mapLongUrlToShortUrl[longUrl]
	return exists
}

func (s *InMemoryStore) getAllURL() map[model.LongURL]model.ShortURL {
	return s.mapLongUrlToShortUrl
}
