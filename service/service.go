package service

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/yapsj/marvel/marvel"
	"github.com/yapsj/marvel/redis"
	"github.com/yapsj/marvel/util"
)

const CHARACTERS_KEY = "characters"

type MarvelService struct {
	marvel *marvel.Client
	redis  *redis.Client
	logger *log.Logger
	ttl    int
}

func NewService(marvel *marvel.Client, redis *redis.Client, w io.Writer, ttl int) *MarvelService {
	return &MarvelService{
		marvel: marvel,
		redis:  redis,
		logger: util.NewInfoLogger(w),
		ttl:    ttl,
	}
}

func (m *MarvelService) ClearCharactersCache() error {
	return m.redis.Delete(CHARACTERS_KEY)
}

func (m *MarvelService) GetCharacters(ctx context.Context) (*[]int, error) {
	logger := m.logger

	//check in redis for cached records
	logger.Printf("Checking REDIS cache server for key: %s\n", CHARACTERS_KEY)
	redisResult, _ := m.redis.Get(CHARACTERS_KEY)

	var characters *[]marvel.Character
	// if no cached result in redis
	if redisResult == "" {
		//get characters from marvel api
		logger.Printf("REDIS has no result, proceeding to query from marvel server.\n")
		var err error
		characters, err = m.marvel.GetCharacters(ctx)

		logger.Printf("RESULT from marvel server, length %v\n", len(*characters))
		if err != nil {
			return nil, err
		}

		raw, err := json.Marshal(characters)
		if err != nil {
			return nil, err
		}

		//set result into redis
		logger.Printf("Setting to REDIS cache server, key: %s\n", CHARACTERS_KEY)
		err = m.redis.Set(CHARACTERS_KEY, string(raw), m.ttl)
		if err != nil {
			return nil, err
		}

	} else {
		err := json.Unmarshal([]byte(redisResult), &characters)
		if err != nil {
			return nil, err
		}
	}

	var idList []int
	for _, c := range *characters {
		idList = append(idList, c.ID)
	}

	return &idList, nil
}

func (m *MarvelService) GetCharacter(ctx context.Context, id int) (*marvel.Character, error) {
	logger := m.logger
	logger.Printf("Getting character with id: %v\n", id)
	character, err := m.marvel.GetCharacter(ctx, id)
	if err != nil {
		return nil, err
	}
	return character, nil
}
