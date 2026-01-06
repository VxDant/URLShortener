package service

import (
	"fmt"
	"testing"
)

func TestURLService_GenerateShortURLCode(t *testing.T) {
	testLongURL := "https://google.com"

	service := NewURLService(nil)

	shortCode, _ := service.GenerateShortURLCode(testLongURL)

	fmt.Println(shortCode)

}
