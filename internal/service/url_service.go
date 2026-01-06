package service

import (
	"URLShortener/internal/models"
	"URLShortener/internal/repository"
	"fmt"
	"math/rand/v2"
)

const dataSet = "ABCDEFGHIJKLMNOPQRSTUVXYZabcdefghijklmnopqrstuvwxyz"

type URLService struct {
	repository *repository.URLRepository
}

func NewURLService(repository *repository.URLRepository) *URLService {
	return &URLService{repository: repository}
}

func (service *URLService) CreateAndProcessShortURL(longURL string) (*models.URL, error) {
	shortURLCode, error := service.GenerateShortURLCode(longURL)

	if error != nil {
		return nil, error
	}

	url, error := service.AddShortURLCode(shortURLCode, longURL)

	if error != nil {
		return nil, error
	}

	return url, nil

}

func (service *URLService) AddShortURLCode(shortURLCode string, longURL string) (*models.URL, error) {
	url, err := service.repository.Create(shortURLCode, longURL)

	if err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	return url, nil
}

func (service *URLService) GetAllURL() ([]*models.URL, error) {
	urls, error := service.repository.GetAll()

	if error != nil {
		return nil, error
	}

	return urls, nil
}

//func (service *URLService) GenerateShortURL(shortURLCode string) (string, error) {
//
//	shortURL := fmt.Sprintf("https://shortly%v.com",shortURLCode)
//
//	return shortURL, nil
//
//}

func (service *URLService) GenerateShortURLCode(longURL string) (string, error) {
	s := make([]byte, 6)

	for i := 0; i < 6; i++ {
		s[i] = byte(rand.UintN(52))
	}

	result := ""

	for _, value := range s {
		result += string(dataSet[value])
	}

	return result, nil
}
