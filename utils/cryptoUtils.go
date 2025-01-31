package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	log "github.com/sirupsen/logrus"

)

func identifyEncryptionAlgorithm(tokenString string) (jwt.SigningMethod, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	return token.Method, nil
}

func VerifyPublicKeyFormat(publicKey string) (*rsa.PublicKey, error) {
	// decode public key in pub format
	
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, fmt.Errorf("public key is not in PEM format")
	}
	// parse public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA")
	}
	return rsaPublicKey, nil

}
func VerifyPrivateKeyFormat(privateKey string) (*rsa.PrivateKey, error) {
	
	block, _ := pem.Decode([]byte(privateKey))
	decodedPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error verifying private key format")
		return nil, err
	}
	return decodedPrivateKey, nil
}

func VerifyTokenRSA(tokenString string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	// check if the signing method is RSA
	signingMethod, err := identifyEncryptionAlgorithm(tokenString)
	if err != nil {
		return nil, err
	}
	if signingMethod != jwt.SigningMethodRS256 {
		return nil, fmt.Errorf("unexpected signing method: %v", signingMethod.Alg())
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}


func GenerateJWTToken(claims jwt.Claims, privateKey *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func VerifyAllKeys(Configdata Config)(error,*RSAkeys){
	publicKey,err := ioutil.ReadFile(Configdata.JwtRSAPublicKey)
	privateKey,err := ioutil.ReadFile(Configdata.JwtRSAPrivateKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error reading public key")
		return err,nil
	}
	// convert public key to string
	fmt.Println(string(publicKey))

	rsaPublicKey, err := VerifyPublicKeyFormat(string(publicKey))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error verifying public key format")
		return err,nil
	}
	// fmt.Println("rsaPublicKey: ",rsaPublicKey)

	clientClaims := jwt.MapClaims{
		"name": "John Doe",
		"admin": true,
	}
	decodedPrivateKey, err := VerifyPrivateKeyFormat(string(privateKey))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error verifying private key format")
		return err,nil
	}


	generatedToken,err:=GenerateJWTToken(clientClaims,decodedPrivateKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error generating")
		return err,nil
	}
	fmt.Println("generatedToken: ",generatedToken)
	return nil,&RSAkeys{PublicKey:rsaPublicKey,PrivateKey:decodedPrivateKey}

}

