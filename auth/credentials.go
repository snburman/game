package auth

import "github.com/snburman/game/config"

type ClientCredentialsDTO struct {
	ClientID     string      `json:"client_id"`
	ClientSecret string      `json:"client_secret"`
	Data         interface{} `json:"data"`
}

func NewClientCredentialsDTO[T any](data T) ClientCredentialsDTO {
	return ClientCredentialsDTO{
		ClientID:     config.Env().CLIENT_ID,
		ClientSecret: config.Env().CLIENT_SECRET,
		Data:         data,
	}
}
