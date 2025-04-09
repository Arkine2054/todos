package handler

import (
	"context"
	logging "todos3/todos-api/pkg"
	"todos3/todos-api/server"
	"todos3/todos-api/token"
)

type handler struct {
	ctx        context.Context
	server     *server.Server
	TokenMaker *token.JwtMaker
	logger     logging.Logger
}

func NewHandler(server *server.Server, logger logging.Logger, secretKey string) *handler {
	return &handler{
		ctx:        context.Background(),
		server:     server,
		TokenMaker: token.NewJwtMaker(secretKey),
		logger:     logger,
	}
}
