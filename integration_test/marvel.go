package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

func marvelTest(ctx context.Context, url string) {
	client := http.DefaultClient
	responseBody := request(ctx, client, fmt.Sprint(url+"/characters"))

	var charactersId []int
	err := json.Unmarshal(responseBody, &charactersId)
	if err != nil {
		log.Panic(err)
	}

	randomId := rand.Intn(len(charactersId))
	getCharactersById := fmt.Sprintf(url+"/characters/%v", charactersId[randomId])
	request(ctx, client, getCharactersById)

}

func request(ctx context.Context, client *http.Client, url string) []byte {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Panic(err)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Get Characters result: ", string(responseBody))
	return responseBody
}
