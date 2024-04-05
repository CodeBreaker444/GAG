package main

import (
    "fmt"
    "net/http"
    "os"

    "github.com/codebreaker444/gag/api/handlers"
    "github.com/codebreaker444/gag/api/middleware"
    utils "github.com/codebreaker444/gag/utils"
    log "github.com/sirupsen/logrus"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: go run main.go <yaml_file>")
    }

    log.SetLevel(log.DebugLevel)

    configData, err := utils.ParseYamlFile(os.Args[1])
    if err != nil {
        log.WithError(err).Fatal("Error parsing YAML file")
    }

    log.WithField("config", configData).Info("Parsed config")

    _, rsaKeys := utils.VerifyAllKeys(configData)
    log.Info("Processed Public and Private keys")

    rootRouter := http.NewServeMux()
    handlers.RootHandler(rootRouter)

    stackMiddleware := utils.MiddlewareStack(
        middleware.PrimaryMiddleware(configData, *rsaKeys),
        middleware.CorsMiddleware,
    )

    server := http.Server{
        Addr:    configData.ServerAddress,
        Handler: stackMiddleware(rootRouter),
    }

    if err := server.ListenAndServe(); err != nil {
        log.WithError(err).Fatal("Error starting server")
    }

    fmt.Printf("Server started at %s\n", configData.ServerAddress)
}