package service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/golang-jwt/jwt"
	"github.com/klovercloud-ci-cd/terminal/config"
	"log"
	"time"
)

var publicKey *rsa.PublicKey

func SetPublicKey() {
	if publicKey == nil {
		block, _ := pem.Decode([]byte(config.PublicKey))
		publicKeyImported, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			log.Print(err.Error())
			panic(err)
		}
		publicKey = publicKeyImported
	}
}

func ValidateToken(tokenString string) (bool, *jwt.Token) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		log.Print("[ERROR]: Token is invalid! ", err.Error())
		return false, nil
	}
	var tm time.Time
	switch iat := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		tm = time.Unix(v, 0)
	}
	if time.Now().UTC().After(tm) {
		return false, nil
	}
	return true, token
}
