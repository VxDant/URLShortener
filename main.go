package main

import (
	"fmt"
	"math/rand/v2"
)

var longurl string
var longToShort map[string]string

func main() {

	longToShort = make(map[string]string)
	dataSet := "abcdABCD"

	fmt.Println("Welcome to the url shortener app")

	for {
		fmt.Println("Enter the URL: ")
		fmt.Scanln(&longurl)
		fmt.Println(shorten(dataSet))

	}

}

func shorten(dataSet string) string {
	s := make([]byte, 4)

	for i := 0; i < 4; i++ {
		s[i] = byte(rand.UintN(8))
	}

	return getShortenedURL(dataSet, s)

}

func getShortenedURL(dataSet string, s []byte) string {
	//s = [0 1 2 3]

	result := ""

	for _, value := range s {
		result += string(dataSet[value])
	}

	return result
}
