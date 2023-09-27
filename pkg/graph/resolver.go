package graph

import "github.com/Zavr22/EMTestTask/pkg/rest/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userService *service.UserService
}

func NewResolver(userService *service.UserService) *Resolver {
	return &Resolver{userService: userService}
}
