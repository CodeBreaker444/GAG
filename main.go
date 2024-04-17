package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/codebreaker444/gag/api/handlers"
	"github.com/codebreaker444/gag/api/middleware"
	utils "github.com/codebreaker444/gag/utils"
	log "github.com/sirupsen/logrus"
)
type slashFix struct {
    mux http.Handler
}

func (h *slashFix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    r.URL.Path = strings.Replace(r.URL.Path, "//", "/", -1)
    h.mux.ServeHTTP(w, r)
}
func main() {
    var rsaKeys *utils.RSAkeys = &utils.RSAkeys{}
    configFile := flag.String("config", "", "Path to the YAML configuration file")
    flag.Parse()

    if *configFile == "" {
        log.Fatal("Usage: gag --config=<yaml_file>")
    }

    log.SetLevel(log.DebugLevel)

    configData, err := utils.ParseYamlFile(*configFile)
    if err != nil {
        log.WithError(err).Fatal("Error parsing YAML file")
    }
    log.WithField("config", configData).Info("Parsed config")
    if configData.Mode == "GAG" {
    _, rsaKeys = utils.VerifyAllKeys(configData)
    log.Info("Processed Public and Private keys")
    }
    
    rootRouter := http.NewServeMux()
    handler := &handlers.Handler{
        Config: configData,
    }
    handler.RootHandler(rootRouter)
    stackMiddleware := utils.MiddlewareStack(
        
        middleware.MiddlewareSwitch(configData, *rsaKeys),
    )

    server := http.Server{
        Addr:    configData.ServerAddress,
        Handler: &slashFix{stackMiddleware(rootRouter)},
    }

    if err := server.ListenAndServe(); err != nil {
        log.WithError(err).Fatal("Error starting server")
    }

    fmt.Printf("Server started at %s\n", configData.ServerAddress)
}