package client

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

type Config struct {
	Url   string `env:"SERVICE_URL,notEmpty"`
	Token string `env:"SERVICE_TOKEN,notEmpty"`
}

type IntegrationClient struct {
	restyClient *resty.Client
}

func Post() {

}

func NewClient() *IntegrationClient {
	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatalf("[ERROR] - Error on configure env on http client: %s", err)
	}

	client := resty.New()
	client.RetryCount = 3
	client.RetryWaitTime = 10 * time.Second
	client.RetryConditions = []resty.RetryConditionFunc{
		func(response *resty.Response, err error) bool {
			return response.StatusCode() == 500
		},
	}
	client.SetHostURL(config.Url)
	client.SetAuthScheme("Basic")
	client.SetAuthToken(config.Token)

	return &IntegrationClient{restyClient: client}
}
