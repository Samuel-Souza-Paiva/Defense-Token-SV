package app

import (
	"time"

	"token-service/internal/pkg/config"
	"token-service/internal/pkg/defense"
	"token-service/internal/pkg/logger"
	"token-service/internal/pkg/store"

	env "github.com/joho/godotenv"
)

const DELAY_TO_REFRESH = 20
const DELAY_TO_RETRY = 5

type TokenService struct {
	Config  *config.Config
	Defense *defense.Defense
	Logger  *logger.Logger
	Store   store.KeyValueStore
}

func (ts *TokenService) Go() {
	for {
		code, err := ts.createToken(ts.Config.DEFENSE_USER, ts.Config.DEFENSE_PASS)
		if err != nil && code != 2001 {
			ts.Logger.PrintError("Error on create defense token", err)

			time.Sleep(time.Minute)
			continue

		} else if code == 2001 {

			err = ts.Store.SetKey("defenseConnection", "false")
			if err != nil {
				ts.Logger.PrintError("error on set defense token on store", err)
			}

			oldPass := ts.Config.DEFENSE_PASS
			ts.Logger.PrintError("Changed credentials", nil)

			for {

				time.Sleep(DELAY_TO_RETRY * time.Second)
				loadErr := env.Load()
				if loadErr != nil {
					ts.Logger.PrintWarning("Error on reload .env file\n" + loadErr.Error())
				}
				ts.Config.Load()

				if ts.Config.DEFENSE_PASS != oldPass {
					break
				}
				ts.Logger.PrintError("Password haven't Changed ", loadErr)

			}
			continue
		} else if code == 2002 {
			//TODO: ESPERAR O TEMPO QUE O USUARIO ESTÁ LOCKADO
			time.Sleep(3 * time.Minute)
			continue

		}
		ts.SetInfo()

		ts.Logger.PrintSuccess("Defense token and credential created successfully")

		for {
			time.Sleep(DELAY_TO_REFRESH * time.Second)

			token, err := ts.Store.GetKey("defenseToken")
			if err != nil {

				ts.Logger.PrintError("Error on get defense token from store", err)
				break
			}

			refreshed, err := ts.refreshToken(token)
			if err != nil {
				ts.Logger.PrintError("Error on refresh defense token", err)

				break
			}

			if !refreshed {
				ts.Logger.PrintWarning("Defense token not refreshed")
				break
			}

			updatedAt := time.Now().UTC().Format(time.RFC3339)
			if err := ts.Store.SetKey("defenseTokenUpdatedAt", updatedAt); err != nil {
				ts.Logger.PrintError("Error on set defense token updatedAt on store", err)
			}
			ts.Store.SetKey("defenseConnection", "true")
			ts.Logger.PrintLog("Defense token refreshed")

		}
	}
}
