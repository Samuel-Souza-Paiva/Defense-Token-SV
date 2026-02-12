package main

import (
	"os"
	"sync"

	"token-service/internal/app"
	"token-service/internal/pkg/config"
	"token-service/internal/pkg/defense"
	"token-service/internal/pkg/logger"
	"token-service/internal/pkg/store"

	env "github.com/joho/godotenv"
)

func init() {
	lgr := &logger.Logger{
		Prefix: "[ENV]",
	}
	err := env.Load()
	if err != nil {
		lgr.PrintWarning("Error on load .env file\n" + err.Error())
	}

	if os.Getenv("ENVIRONMENT") != "DEVELOPMENT" {
		os.Setenv("HOST-NAME", "DESKTOP")
	}
}

func main() {
	lgr := &logger.Logger{
		Prefix: "[TS]",
	}

	config := &config.Config{}

	err := config.Load()
	if err != nil {
		lgr.PrintError("Error on load config", err)
		return
	}

	kv := store.NewMemoryStore()
	app.StartHTTPServer(config.HTTP_ADDR, config.HTTP_API_KEY, kv, lgr)

	defense := &defense.Defense{
		Host:               config.DEFENSE_HOST,
		Port:               config.DEFENSE_PORT,
		Scheme:             config.DEFENSE_SCHEME,
		InsecureSkipVerify: config.DEFENSE_INSECURE_SKIP_VERIFY,
	}

	tokenService := &app.TokenService{
		Config:  config,
		Defense: defense,
		Logger:  lgr,
		Store:   kv,
	}
	tokenService.SetInfo()

	go tokenService.Go()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
