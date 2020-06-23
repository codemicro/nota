package authentication

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

var (
	PubKey  *rsa.PublicKey
	PrivKey *rsa.PrivateKey
)

func LoadKeys() error {
	pubKeyCont, err := ioutil.ReadFile("./public.key")
	privKeyCont, err := ioutil.ReadFile("./private.key")
	if err != nil {
		return err
	}

	PubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyCont)
	PrivKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyCont)
	if err != nil {
		return err
	}

	return nil
}
