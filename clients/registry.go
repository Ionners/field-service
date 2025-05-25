package clients

import (
	"field-service/clients/config"
	clients "field-service/clients/user"
	config2 "field-service/config"
	"fmt"
)

type ClientRegistry struct {
}

type IClientRegistry interface {
	GetUser() clients.IUserClient
}

func NewClientRegistry() IClientRegistry {
	return &ClientRegistry{}
}

func (c *ClientRegistry) GetUser() clients.IUserClient {
	fmt.Println("📦 [CLIENT-REGISTRY-INIT] AuthService BaseURL:", config2.Config.InternalService.User.Host)
	fmt.Println("📦 [CLIENT-REGISTRY-INIT] SignatureKey:", config2.Config.InternalService.User.SignatureKey)

	return clients.NewUserClient(
		config.NewClientConfig(
			config.WithBaseURL(config2.Config.InternalService.User.Host),
			config.WithSignatureKey(config2.Config.InternalService.User.SignatureKey),
		))
}
