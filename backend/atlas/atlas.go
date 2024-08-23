package atlas

import (
	"github.com/go-resty/resty/v2"
)

type Atlas struct {
	ApiKey string
	Client *resty.Client
}

func InitialiseAtlas(apiKey string) *Atlas {
	client := resty.New()
	return &Atlas{
		ApiKey: apiKey,
		Client: client,
	}
}
