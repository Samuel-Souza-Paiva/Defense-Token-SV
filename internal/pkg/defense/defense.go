package defense

import (
	"errors"
)

type Defense struct {
	Host                 string
	Port                 string
	Scheme               string
	InsecureSkipVerify   bool
}

func (defense *Defense) Auth(username, password, rsaPK string) (AuthRes, error) {
	const AUTH_ENDPOINT = "/brms/api/v1.0/accounts/authorize"

	encData := &EncData{}
	auth := &Auth{}

	defenseApi := api{
		Host:               defense.Host,
		Port:               defense.Port,
		Scheme:             defense.Scheme,
		InsecureSkipVerify: defense.InsecureSkipVerify,
	}

	encData.CreatePayload(username)

	_, encDataRes, err := defenseApi.Post(AUTH_ENDPOINT, encData.Payload)
	if err != nil {
		return auth.Res, errors.New("error on get encrypted data with defense\n" + err.Error())
	}

	encData.SetRes(encDataRes)
	if rsaPK != "" {
		encData.Res.Publickey = rsaPK
	}
	auth.CreateSignature(username, password, encData.Res.Realm, encData.Res.RandomKey)
	auth.CreatePayload(auth.Signature, username, encData.Res.RandomKey, encData.Res.Publickey)

	_, authRes, err := defenseApi.Post(AUTH_ENDPOINT, auth.Payload)
	if err != nil {
		return auth.Res, errors.New("error on authenticate with defense\n" + err.Error())
	}

	err = auth.SetRes(authRes)
	if err != nil {
		return auth.Res, errors.New("error on set auth response\n" + err.Error())
	}

	if auth.Res.Token == "" {
		return auth.Res, errors.New("defense token not created")
	}

	return auth.Res, nil
}

func (defense *Defense) Refresh(token string) (bool, error) {
	if token == "" {
		return false, errors.New("invalid defense token")
	}

	const REFRESH_ENDPOINT = "/brms/api/v1.0/accounts/keepalive"

	defenseApi := api{
		Host:               defense.Host,
		Port:               defense.Port,
		Scheme:             defense.Scheme,
		InsecureSkipVerify: defense.InsecureSkipVerify,
		Token:              token,
	}

	code, _, err := defenseApi.Put(REFRESH_ENDPOINT, nil)

	if err != nil {
		return false, errors.New("error on refresh token with defense\n" + err.Error())
	}

	return code == 200, nil
}
