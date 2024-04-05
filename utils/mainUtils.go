package utils

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"gopkg.in/yaml.v3"
)

type Config struct {
    AuthenticatedPrefix   string `yaml:"authenticated-prefix"`
    UnauthenticatedPrefix string `yaml:"unauthenticated-prefix"` // 1. CHANGE THESE TO SNAKE CASE
	JwtRSAPublicKey       string `yaml:"jwt-rsa-public-key"`
    JwtRSAPrivateKey      string `yaml:"jwt-rsa-private-key"`
    ServerAddress         string `yaml:"server-address"`
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
