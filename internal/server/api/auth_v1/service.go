package auth

import (
	desc "github.com/igortoigildin/goph-keeper/pkg/auth_v1"

	service "github.com/igortoigildin/goph-keeper/internal/server/service"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
