package main

import (
    "io/ioutil"
    "net/http"
    "os"
	"fmt"
	"reflect"
	"github.com/codebreaker444/gag/api/handlers"
	"github.com/codebreaker444/gag/api/middleware"
	mainutils"github.com/codebreaker444/gag/utils"
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v3"
)

func checkAllFieldsPresent(data mainutils.Config) error {
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


func parseYamlFile(yamlFile string) (mainutils.Config, error) { // 2. SHIFT IT TO utils/mainUtils.go
    yamlData, err := ioutil.ReadFile(yamlFile)
    if err != nil {
        return mainutils.Config{}, err
	}

    var data mainutils.Config
    err = yaml.Unmarshal(yamlData, &data)
    if err != nil {
        return mainutils.Config{}, err
    }
	if err := checkAllFieldsPresent(data); err != nil {
		return mainutils.Config{}, err
	}



    return data, nil
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: go run main.go <yaml_file>")
    }
    log.SetLevel(log.DebugLevel)
    Configdata, err := parseYamlFile(os.Args[1])
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Error parsing YAML file")
    }
    log.WithFields(log.Fields{
        "config": Configdata,
    }).Info("Parsed config")
	
	publicKey,err := ioutil.ReadFile(Configdata.JwtRSAPublicKey)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error reading public key")
			return
		}
	fmt.Println(string(publicKey))
	_, err = mainutils.VerifyPublicKeyFormat(string(publicKey))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error verifying public key format")
		return
	}
	

	rootRouter := http.NewServeMux()
	handlers.RootHandler(rootRouter)
	server := http.Server{
		Addr:    "localhost:8000",
		Handler: middleware.PrimaryMiddleware(rootRouter,Configdata),
	}
	err = server.ListenAndServe()

    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Error starting server")
    }
	fmt.Println("Server started at localhost:8000")
}