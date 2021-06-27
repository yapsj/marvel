package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpserver "github.com/yapsj/marvel/http"
	"github.com/yapsj/marvel/marvel"
	"github.com/yapsj/marvel/redis"
	"github.com/yapsj/marvel/service"
	"github.com/yapsj/marvel/util"
)

const USE_DEFAULT = false

func readConfigFile(filepath string) (*Config, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

type Config struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
	LogFilePath string `json:"logFilePath"`
	TTL         int    `json:"cacheTTL"`
}

func main() {
	//declaration

	var publicKey string
	var privateKey string
	var redisAddress string
	var httpAddress string
	var httpPort int
	var ttl int
	var logFilePath string

	if USE_DEFAULT {
		publicKey = ""
		privateKey = ""
		redisAddress = "localhost:6379"
		httpAddress = "localhost"
		httpPort = 8080
		ttl = 30000
		logFilePath = "logs.txt"
	} else {
		var args = os.Args

		if len(args) < 4 {
			fmt.Println("Arguments is not complete, should contains config file location, publickey and privatekey")
			return
		}

		configFile := args[1]
		config, err := readConfigFile(configFile)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			return
		}

		publicKey = args[2]
		privateKey = args[3]

		ttl = config.TTL // seconds
		logFilePath = config.LogFilePath + "logs.txt"
		httpAddress = config.Host
		httpPort = config.Port

		redisConfig := config.Redis
		redisAddress = fmt.Sprintf("%s:%v", redisConfig.Host, redisConfig.Port)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("%s file is not readable", logFilePath)
		return
	}
	logger := util.NewInfoLogger(logFile)
	//endpoint creation
	marvelClient := marvel.NewClient(http.DefaultClient, publicKey, privateKey, logFile)
	redisClient := redis.NewClient(redisAddress)

	//service creation
	marvelService := service.NewService(marvelClient, redisClient, logFile, ttl)

	//router creation
	mux := mux.NewRouter()
	marvelServer := httpserver.NewServer(marvelService, logFile, mux)

	//define routes
	marvelServer.Serve()
	logger.Println("Serving at " + fmt.Sprintf("%s:%v", httpAddress, httpPort))

	http.ListenAndServe(fmt.Sprintf("%s:%v", httpAddress, httpPort), mux)
}
