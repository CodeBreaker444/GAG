package middleware

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"

	mainutils "github.com/codebreaker444/gag/utils"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	//import type config from main.go
)

func PrimaryMiddleware(next http.Handler, Configdata mainutils.Config ) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// print Authorization header if it exists
		authorization := r.Header.Get("Authorization")
		if authorization != "" {
			log.WithFields(log.Fields{

				"Authorization": authorization,
			}).Debug("Authorization header")
		}
		// load public key from config
		publicKey,err := ioutil.ReadFile(Configdata.JwtRSAPublicKey)
		privateKey,err := ioutil.ReadFile(Configdata.JwtRSAPrivateKey)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error reading public key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// convert public key to string
		fmt.Println(string(publicKey))

		rsaPublicKey, err := mainutils.VerifyPublicKeyFormat(string(publicKey))
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error verifying public key format")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("rsaPublicKey: ",rsaPublicKey)
		block, _ := pem.Decode(privateKey)
		decodedPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error verifying private key format")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// rsaPrivateKey, ok := decodedPrivateKey.(*rsa.PrivateKey)
	
		clientClaims := jwt.MapClaims{
			"name": "John Doe",
			"admin": true,
		}


		generatedToken,err:=mainutils.GenerateJWTToken(clientClaims,decodedPrivateKey)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error generating")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Println("generatedToken: ",generatedToken)
		// split token
		au := authorization[7:]

		_,err=mainutils.VerifyTokenRSA(au,rsaPublicKey )
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error verifying token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// print request URL
		log.WithFields(log.Fields{
			"URL": r.URL,
		}).Debug("Request URL")

		next.ServeHTTP(w, r)
	})
}