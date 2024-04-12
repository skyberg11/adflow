package main

import (
	"adflow/internal/adapters"
	"adflow/internal/adapters/adrepo"
	"adflow/internal/adapters/aduser"
	"adflow/internal/app"
	"adflow/internal/app/auth"
	"os"

	// "adflow/internal/ports/grpc"
	// grpcApp "adflow/internal/ports/grpc/app"
	"adflow/internal/ports/httpgin"
)

func main() {
	db_uri_users, ok := os.LookupEnv("DB_URI_USERS")
	if !ok {
		db_uri_users = "test_users.db"
	}

	db_uri_ads, ok := os.LookupEnv("DB_URI_ADS")
	if !ok {
		db_uri_ads = "test_ads.db"
	}

	signKey, ok := os.LookupEnv("AUTH_SIGNING_KEY")
	if ok {
		auth.JwtKey = []byte(signKey)
	}

	db_users, err := adapters.NewSQLite(db_uri_users)
	if err != nil {
		panic(err)
	}

	db_ads, err := adapters.NewSQLite(db_uri_ads)
	if err != nil {
		panic(err)
	}

	repo, users := adrepo.NewSQLiteAds(db_ads), aduser.NewSQLiteUsers(db_users)
	httpServer := httpgin.NewHTTPServer(":8080", app.NewApp(repo, users))
	// _ = grpc.NewGRPCServer(grpcApp.NewAdService(repo, users))
	err = httpServer.Listen()
	if err != nil {
		panic(err)
	}
}
