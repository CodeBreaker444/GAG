package middleware

import (
	"net/http"
	"context"
	utils "github.com/codebreaker444/gag/utils"
	log "github.com/sirupsen/logrus"
	//import type config from main.go
)

func PrimaryMiddleware(Configdata utils.Config, rsaKeys utils.RSAkeys ) utils.Middleware {
return func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// print Authorization header if it exists
		ctx := r.Context()
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			// set request context to unauthenticated
			log.WithFields(log.Fields{
				"Authorization": authorization,
			}).Debug("No Authorization header")
			ctx = context.WithValue(ctx, "authenticated", false)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			

		}
		// load public key from config
		
		// rsaPrivateKey, ok := decodedPrivateKey.(*rsa.PrivateKey)
	
		
		// split token
		au := authorization[7:]

		_,err:=utils.VerifyTokenRSA(au, rsaKeys.PublicKey)
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
		w.Header().Set("X-API-GATEWAY", "github.com/codebreaker444/gag")

		next.ServeHTTP(w, r)

	
})
}}

func CorsMiddleware(Configdata utils.Config) utils.Middleware {
	return func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Headers: ",r.Header)

		// check access key in headers
		accessKey := r.Header.Get("x-gag-api-key")
		if accessKey != Configdata.CorsApiKey {
			log.WithFields(log.Fields{
				"accessKey": accessKey,
			}).Error("Invalid access key")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		
		w.Header().Set("X-API-GATEWAY", "github.com/codebreaker444/gag")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
}



func MiddlewareSwitch(configData utils.Config, rsaKeys utils.RSAkeys) utils.Middleware {
    switch configData.Mode {
    case "GAG":
        log.Info("GAG mode selected")
        return PrimaryMiddleware(configData, rsaKeys)
    default:
        log.Info("CORS mode selected")
        return CorsMiddleware(configData)
    }
}