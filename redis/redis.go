package redis

import (
	"github.com/go-redis/redis"
)

const UNLIMITED = -1

type Client struct {
	client *redis.Client
}

func NewClient(Addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: Addr,
		DB:   0, // use default DB
	})

	return &Client{
		client: rdb,
	}
}

func (c *Client) Delete(key string) error {
	_, err := c.client.Do("del", key).Result()
	return err
}

func (c *Client) Set(key string, value string, ttl int) error {
	_, err := c.client.Do("set", key, value).Result()
	if err != nil {
		return err
	}

	if ttl == UNLIMITED {
		return nil
	}

	_, err = c.client.Do("expire", key, ttl).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(key string) (string, error) {
	result, err := c.client.Do("get", key).Result()
	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func test() {

}
