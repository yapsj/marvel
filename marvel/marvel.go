package marvel

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/yapsj/marvel/util"
)

const GET = "GET"
const API_URL = "https://gateway.marvel.com:443/v1/public"
const LIMIT = 20

type Client struct {
	httpClient *http.Client
	publicKey  string
	privateKey string
	logger     *log.Logger
}

func NewClient(client *http.Client, publicKey string, privateKey string, w io.Writer) *Client {
	return &Client{
		httpClient: client,
		publicKey:  publicKey,
		privateKey: privateKey,
		logger:     util.NewInfoLogger(w),
	}
}

func (m *Client) GetCharacter(ctx context.Context, id int) (*Character, error) {
	path := fmt.Sprintf("/characters/%s", strconv.Itoa(id))

	request, err := http.NewRequestWithContext(ctx, GET, getUrl(path), nil)
	if err != nil {
		return nil, err
	}

	m.formRequest(request, nil)
	response, err := m.execute(request)
	if err != nil {
		return nil, err
	}

	var result Result
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}

	if result.Data.Count == 0 {
		return nil, fmt.Errorf("no character with id %v is found", id)
	}

	return &result.Data.Results[0], nil
}

func (m *Client) GetCharacters(ctx context.Context) (*[]Character, error) {
	offset := 0
	result, err := m.getCharacters(ctx, offset)
	if err != nil {
		return nil, err
	}

	characters := &result.Data.Results
	if result.Data.Count == 0 {
		return characters, nil
	}

	total := result.Data.Total

	m.logger.Printf("total: %v, limit: %v\n", total, LIMIT)
	var n = total / LIMIT
	if total%LIMIT == 0 {
		n -= 1
	}

	offset += LIMIT

	var c []Character
	c = append(c, *characters...)
	for i := 0; i < n; i++ {
		r, err := m.getCharacters(ctx, offset)
		if err != nil {
			return nil, err
		}

		c = append(c, r.Data.Results...)
		m.logger.Printf("MARVEL.GetCharacters, Getting offset: %v \n", offset)
		offset += LIMIT

	}

	return &c, nil
}

func (m *Client) getCharacters(ctx context.Context, offset int) (*Result, error) {

	request, err := http.NewRequestWithContext(ctx, GET, getUrl("/characters"), nil)
	if err != nil {
		return nil, err
	}
	queryParams := map[string]string{
		"offset": strconv.Itoa(offset),
		"limit":  strconv.Itoa(LIMIT),
	}

	m.formRequest(request, queryParams)

	responseBody, err := m.execute(request)
	if err != nil {
		return nil, err
	}

	var result Result
	if unmarshalErr := json.Unmarshal(responseBody, &result); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &result, nil
}

func (m *Client) Test() {
	fmt.Println("test")
	/*
		aa, _ := m.GetCharacter(context.Background(), 1011334)
		aa1, _ := json.Marshal(aa)

		fmt.Printf("%s", string(aa1))
	*/

	aa, _ := m.GetCharacters(context.Background())
	count := len(*aa)

	fmt.Printf("%v", count)
}

func hash(s1, s2, s3 string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s%s", s1, s2, s3))))
}

func getUrl(format string) string {
	return fmt.Sprintf(API_URL + format)
}

func (m *Client) formRequest(request *http.Request, props map[string]string) {
	var ts = fmt.Sprintf("%v", time.Now().Unix())
	var hash = hash(ts, m.privateKey, m.publicKey)

	var query = request.URL.Query()
	query.Add("ts", ts)
	query.Add("apikey", m.publicKey)
	query.Add("hash", hash)

	for key, value := range props {
		query.Add(key, value)
	}

	request.URL.RawQuery = query.Encode()
}

func (m *Client) execute(request *http.Request) ([]byte, error) {
	response, err := m.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

type Result struct {
	Code            int    `json:"code"`
	Status          string `json:"status"`
	Copyright       string `json:"copyright"`
	AttributionText string `json:"attributionText"`
	AttributionHTML string `json:"attributionHTML"`
	Etag            string `json:"etag"`
	Data            struct {
		Offset  int         `json:"offset"`
		Limit   int         `json:"limit"`
		Total   int         `json:"total"`
		Count   int         `json:"count"`
		Results []Character `json:"results"`
	} `json:"data"`
}

type Character struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Modified    string `json:"modified"`
	Thumbnail   struct {
		Path      string `json:"path"`
		Extension string `json:"extension"`
	} `json:"thumbnail"`
	ResourceURI string `json:"resourceURI"`
	Comics      struct {
		Available     int    `json:"available"`
		CollectionURI string `json:"collectionURI"`
		Items         []struct {
			ResourceURI string `json:"resourceURI"`
			Name        string `json:"name"`
		} `json:"items"`
		Returned int `json:"returned"`
	} `json:"comics"`
	Series struct {
		Available     int    `json:"available"`
		CollectionURI string `json:"collectionURI"`
		Items         []struct {
			ResourceURI string `json:"resourceURI"`
			Name        string `json:"name"`
		} `json:"items"`
		Returned int `json:"returned"`
	} `json:"series"`
	Stories struct {
		Available     int    `json:"available"`
		CollectionURI string `json:"collectionURI"`
		Items         []struct {
			ResourceURI string `json:"resourceURI"`
			Name        string `json:"name"`
			Type        string `json:"type"`
		} `json:"items"`
		Returned int `json:"returned"`
	} `json:"stories"`
	Events struct {
		Available     int    `json:"available"`
		CollectionURI string `json:"collectionURI"`
		Items         []struct {
			ResourceURI string `json:"resourceURI"`
			Name        string `json:"name"`
		} `json:"items"`
		Returned int `json:"returned"`
	} `json:"events"`
	Urls []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"urls"`
}
