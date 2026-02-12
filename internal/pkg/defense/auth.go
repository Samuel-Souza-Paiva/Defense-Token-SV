package defense

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type Auth struct {
	Payload   AuthPayload
	Res       AuthRes
	Signature string
}

func (auth *Auth) CreatePayload(signature, userName, randomKey, publicKey string) {
	auth.Payload = AuthPayload{
		Signature:   signature,
		UserName:    userName,
		RandomKey:   randomKey,
		PublicKey:   publicKey,
		EncryptType: "MD5",
		IpAddress:   "",
		ClientType:  "WINPC_V2",
		UserType:    "0",
	}
}

func (auth *Auth) CreateSignature(username, password, realm, randomKey string) {
	hash := auth.hash(password)
	hash = auth.hash(username + hash)
	hash = auth.hash(hash)
	hash = auth.hash(username + ":" + realm + ":" + hash)
	hash = auth.hash(hash + ":" + randomKey)

	auth.Signature = hash
}

func (auth *Auth) SetRes(data []byte) error {
	err := json.Unmarshal(data, &auth.Res)
	if err != nil {
		return errors.New("error on unmarshall auth response\n" + err.Error())
	}

	return nil
}

func (auth *Auth) hash(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
