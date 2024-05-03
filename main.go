package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
    "os"
	"github.com/codebreaker444/gag/api/handlers"
	"github.com/codebreaker444/gag/api/middleware"
	utils "github.com/codebreaker444/gag/utils"
    "github.com/fatih/color"
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
    var configData utils.Config
    flag.Usage = func() {
        helpText := []string{
            "Usage of GAG:\n",
            "  -config string",
            "        Path to the YAML configuration file (only when not using flags)",
            "  -authenticatedPrefix string",
            "        Authenticated prefix (default \"/defaultAuthenticatedPrefix\")",
            "  -unauthenticatedPrefix string",
            "        Unauthenticated prefix (default \"/defaultUnauthenticatedPrefix\")",
            "  -jwtRSAPublicKey string",
            "        Path to the JWT RSA public key (default \"/path/to/default/public/key\")",
            "  -jwtRSAPrivateKey string",
            "        Path to the JWT RSA private key (default \"/path/to/default/private/key\")",
            "  -serverAddress string",
            "        Server address (default \"localhost:8080\")",
            "  -destinationURL string",
            "        Destination URL (default \"http://defaultDestination.com\")",
            "  -corsApiKey string",
            "        CORS API key (default \"defaultCorsApiKey\")",
            "  -mode string",
            "        Mode (default \"GAG\")",
        }

        color.Cyan(helpText[0])
        for i := 1; i < len(helpText); i++ {
            if i%2 == 0 {
                color.Green(helpText[i])
            } else {
                color.Yellow(helpText[i])
            }
        }
    }
    configFile := flag.String("config", "", "Path to the YAML configuration file")
        // Define flags for each configuration parameter
    authenticatedPrefix := flag.String("authenticatedPrefix", "/defaultAuthenticatedPrefix", "Authenticated prefix")
    unauthenticatedPrefix := flag.String("unauthenticatedPrefix", "/defaultUnauthenticatedPrefix", "Unauthenticated prefix")
    jwtRSAPublicKey := flag.String("jwtRSAPublicKey", "/path/to/default/public/key", "Path to the JWT RSA public key")
    jwtRSAPrivateKey := flag.String("jwtRSAPrivateKey", "/path/to/default/private/key", "Path to the JWT RSA private key")
    serverAddress := flag.String("serverAddress", "localhost:8080", "Server address")
    destinationURL := flag.String("destinationURL", "http://defaultDestination.com", "Destination URL")
    corsApiKey := flag.String("corsApiKey", "defaultCorsApiKey", "CORS API key")
    mode := flag.String("mode", "GAG", "Mode")
    flag.Parse()
    if len(os.Args) == 1 {
        flag.Usage()
        os.Exit(1)
    }

    if *configFile == "" {
        configData = utils.Config{
            AuthenticatedPrefix:   *authenticatedPrefix,
            UnauthenticatedPrefix: *unauthenticatedPrefix,
            JwtRSAPublicKey:       *jwtRSAPublicKey,
            JwtRSAPrivateKey:      *jwtRSAPrivateKey,
            ServerAddress:         *serverAddress,
            DestinationURL:        *destinationURL,
            CorsApiKey:            *corsApiKey,
            Mode:                  *mode,
        }
    }else{
        var err error
        configData, err = utils.ParseYamlFile(*configFile)
        if err != nil {
            log.WithError(err).Fatal("Error parsing YAML file")
        }
    }

    log.SetLevel(log.DebugLevel)

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