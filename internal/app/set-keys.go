package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

func (ts *TokenService) SetInfo() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return errors.New("error on generate rsa key pair\n" + err.Error())
	}
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	publicKeyPEM := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyDER,
	}
	publicKeyPEMString := string(pem.EncodeToMemory(&publicKeyPEM))
	publicKeyLines := strings.Split(publicKeyPEMString, "\n")
	var publicKeyString strings.Builder
	for _, line := range publicKeyLines {
		if !strings.Contains(line, "PUBLIC KEY") && line != "" {
			publicKeyString.WriteString(line)
		}
	}

	res, err := ts.Defense.Auth(ts.Config.DEFENSE_USER, ts.Config.DEFENSE_PASS, publicKeyString.String())
	if err != nil {
		return errors.New("error on create defense token\n" + err.Error())
	}
	sv, err := decryptRSA(privateKey, res.SecretVector)
	if err != nil {
		return errors.New("error on decrypt secret vector\n" + err.Error())
	}

	sk, err := decryptRSA(privateKey, res.SecretKey)
	if err != nil {
		return errors.New("error on decrypt secret key\n" + err.Error())
	}

	err = ts.Store.SetKey("secretKey", string(sk))
	if err != nil {
		return errors.New("error on set secret key on store\n" + err.Error())
	}
	err = ts.Store.SetKey("secretVector", string(sv))
	if err != nil {
		return errors.New("error on set secret vector on store\n" + err.Error())
	}
	ts.Logger.PrintSuccess("defense info set successfully")

	return nil
}

func decryptRSA(privateKey *rsa.PrivateKey, cipherText string) (string, error) {
	cipherText = strings.TrimSpace(cipherText)
	if cipherText == "" {
		return "", errors.New("empty ciphertext")
	}

	data := []byte(cipherText)
	decoded, err := base64.StdEncoding.DecodeString(cipherText)
	if err == nil {
		data = decoded
	}

	plain, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err == nil {
		return string(plain), nil
	}

	plain, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
