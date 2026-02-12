package app

import "errors"

func (ts *TokenService) refreshToken(token string) (bool, error) {
  refreshed, err := ts.Defense.Refresh(token)

  if err != nil {
    return false, errors.New("error on refresh defense token\n" + err.Error())
  }

  return refreshed, nil
}
