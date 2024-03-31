package main

import (
    "io/ioutil"
    "net/http"
    "os"
	"fmt"

    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v3"
)

type config struct {
    AuthenticatedPrefix   string `yaml:"authenticated-prefix"`
    UnauthenticatedPrefix string `yaml:"unauthenticated-prefix"`
}

func parseYamlFile(yamlFile string) (config, error) {
    yamlData, err := ioutil.ReadFile(yamlFile)
    if err != nil {
        return config{}, err
	}

    var data config
    err = yaml.Unmarshal(yamlData, &data)
    if err != nil {
        return config{}, err
    }

    if data.AuthenticatedPrefix == "" {
        return config{}, fmt.Errorf("missing key: authenticated-prefix")
    }
    if data.UnauthenticatedPrefix == "" {
        return config{}, fmt.Errorf("missing key: unauthenticated-prefix")
    }

    return data, nil
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: go run main.go <yaml_file>")
    }
    log.SetLevel(log.DebugLevel)

    data, err := parseYamlFile(os.Args[1])
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Error parsing YAML file")
    }

    log.WithFields(log.Fields{
        "config": data,
    }).Info("Parsed config")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })

    err = http.ListenAndServe(":8000", nil)
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Error starting server")
    }
}