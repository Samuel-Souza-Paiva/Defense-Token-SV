package app

import (
	"errors"
	"time"
)

func (ts *TokenService) createToken(user, pass string) (int, error) {
	res, err := ts.Defense.Auth(user, pass, "")
	if err != nil {
		return res.Code, errors.New("error on create defense token\n" + err.Error())
	}

	err = ts.Store.SetKey("defenseToken", res.Token)

	if err != nil {
		return res.Code, errors.New("error on set defense token on store\n" + err.Error())
	}

	err = ts.Store.SetKey("defenseCredential", res.Credential)
	if err != nil {
		return res.Code, errors.New("error on set defense credential on store\n" + err.Error())
	}

	updatedAt := time.Now().UTC().Format(time.RFC3339)
	err = ts.Store.SetKey("defenseTokenUpdatedAt", updatedAt)
	if err != nil {
		return res.Code, errors.New("error on set defense token updatedAt on store\n" + err.Error())
	}

	ts.Logger.PrintLog("Defense token: " + res.Token)

	return res.Code, nil
}
