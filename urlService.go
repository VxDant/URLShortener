package main

import (
	"URLShortener/model"
	"fmt"
)

type URLService struct {
	store *InMemoryStore
}

func NewURLService(store *InMemoryStore) *URLService {
	return &URLService{store}
}

func (service *URLService) ProcessAndAddLongURLtoMap(longUrl model.LongURL) {

	exists := service.store.ifLongURLExist(longUrl)
	if exists {
		fmt.Println("Long URL already exists")
		return
	}

	service.generateShortURL(longUrl)

	service.AddLongURL(longUrl)

}

func (service *URLService) AddLongURL(longUrl model.LongURL) {
	service.store.AddURL(longUrl)
}

func (service *URLService) generateShortURL(url model.LongURL) {
	dataset := model.DatasetNew().ValueSet

}
