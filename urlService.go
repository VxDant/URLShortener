package main

import (
	"URLShortener/model"
	"fmt"
	"math/rand/v2"
)

type URLService struct {
	store *InMemoryStore
}

func NewURLService(store *InMemoryStore) *URLService {
	return &URLService{store}
}

func (service *URLService) ProcessAndAddLongURLtoMap(longUrl model.LongURL) string {

	exists := service.store.ifLongURLExist(longUrl)
	if exists {
		fmt.Println("Long URL already exists")
		return ""
	}

	shortUrl := service.generateShortURL(longUrl)

	ShortUrlFinal := "http://localhost:8080/" + shortUrl.Url

	service.AddLongURL(longUrl, shortUrl)

	return ShortUrlFinal
}

func (service *URLService) AddLongURL(longUrl model.LongURL, shortUrl model.ShortURL) {
	service.store.AddURL(longUrl, shortUrl)
}

func (service *URLService) generateShortURL(url model.LongURL) model.ShortURL {
	dataset := model.DatasetNew().ValueSet

	s := make([]byte, 4)

	for i := 0; i < 4; i++ {
		s[i] = byte(rand.UintN(8))
	}

	result := ""

	for _, value := range s {
		result += string(dataset[value])
	}

	return model.ShortURL{result}

}

func (service *URLService) getAllURLs() map[model.LongURL]model.ShortURL {
	return service.store.getAllURL()
}
