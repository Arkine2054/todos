package main

import (
	"github.com/ianschenck/envflag"
	"todos3/db"
	"todos3/todos-api/handler"
	logging "todos3/todos-api/pkg"
	"todos3/todos-api/server"
	"todos3/todos-api/storer"
)

const minSecretKeySize = 32

func main() {
	var secretKey = envflag.String(
		"SECRET_KEY", "01234567890123456789012345678901", "secret key for JWT signing")

	logger := logging.GetLogger()

	if len(*secretKey) < minSecretKeySize {
		logger.Fatalf("SECRET_KEY must be at least %d characters", minSecretKeySize)
	}
	db, err := db.NewDatabase()
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Info("Successfully connected to database")

	//do something with the database
	st := storer.NewPostgresStorer(db.GetDb())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv, logger, *secretKey)
	handler.RegisterRouter(hdl)
	handler.Start(":8080")
}
