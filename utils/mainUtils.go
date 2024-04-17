package utils

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

type Config struct {
    AuthenticatedPrefix   string `yaml:"GAG_AUTHENTICATED_PREFIX"`
    UnauthenticatedPrefix string `yaml:"GAG_UNATHETICATED_PREFIX"` // 1. CHANGE THESE TO SNAKE CASE
	JwtRSAPublicKey       string `yaml:"GAG_JWT_RSA_PUBLIC_KEY"`
    JwtRSAPrivateKey      string `yaml:"GAG_JWT_RSA_PRIVATE_KEY"`
    ServerAddress         string `yaml:"GAG_SERVER_ADDRESS"`
	DestinationURL        string `yaml:"GAG_DESTINATION_URL"`
	CorsApiKey 		  	  string `yaml:"CORS_API_KEY"`
	Mode 				  string `yaml:"MODE"`
}

type RSAkeys struct {
    PublicKey  *rsa.PublicKey
    PrivateKey *rsa.PrivateKey
}


type Middleware func (http.Handler) http.Handler
func MiddlewareStack(middlewares ...Middleware) Middleware {
    return func(final http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            final = middlewares[i](final)
        }
        return final
    }
}
func checkAllFieldsPresent(data Config) error {
	log.Info("Checking if all fields are in valid format")
	v := reflect.ValueOf(data)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("yaml")
		value := v.Field(i).Interface().(string)
		if value == "" {
			return fmt.Errorf("missing field '%s' in YAML file", key)
		}
	}
	formatErr := checkAllValuesFormat(data)
	if formatErr != nil {
		return formatErr
	}
	
	return nil
}

func checkAllValuesFormat(data Config) error {
	if (data.Mode != "GAG" && data.Mode != "CORS"){
		return fmt.Errorf("mode should be either GAG or CORS")
	}
	if  data.Mode == "GAG" {
		if (strings.EqualFold(data.AuthenticatedPrefix, data.UnauthenticatedPrefix)){
			return fmt.Errorf("prefixes cannot be the same")
		}
		// convert prefix to lowercase



		if (data.AuthenticatedPrefix == "" || data.UnauthenticatedPrefix == ""){
			return fmt.Errorf("prefixes cannot be empty")
		}
		if (data.AuthenticatedPrefix[0] != '/' || data.UnauthenticatedPrefix[0] != '/'){
			return fmt.Errorf("prefixes should start with '/'")
		}
		

	}

	
	// check if the values are in the correct format
	return nil
}

func ParseYamlFile(yamlFile string) (Config, error) { // 2. SHIFT IT TO utils/mainUtils.go
    yamlData, err := ioutil.ReadFile(yamlFile)
    if err != nil {
        return Config{}, err
	}

    var data Config
    err = yaml.Unmarshal(yamlData, &data)
    if err != nil {
        return Config{}, err
    }
	if err := checkAllFieldsPresent(data); err != nil {
		return Config{}, err
	}
    return data, nil
}

