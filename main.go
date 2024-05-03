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
            "  -MODE string",
            "        Mode (GAG or CORS) (No default value, must be provided)",
            "  -SERVER_ADDRESS string",
            "        Server address (default \"localhost:8080\")",
            "  -GAG_AUTHENTICATED_PREFIX string",
            "        Authenticated prefix (default \"/defaultAuthenticatedPrefix\")",
            "  -GAG_UNATHETICATED_PREFIX string",
            "        Unauthenticated prefix (default \"/defaultUnauthenticatedPrefix\")",
            "  -GAG_JWT_RSA_PUBLIC_KEY string",
            "        Path to the JWT RSA public key (default \"/path/to/default/public/key\")",
            "  -GAG_JWT_RSA_PRIVATE_KEY string",
            "        Path to the JWT RSA private key (default \"/path/to/default/private/key\")",

            "  -GAG_DESTINATION_URL string",
            "        Destination URL (default \"http://defaultDestination.com\")",
            "  -CORS_API_KEY string",
            "        CORS API key (default \"defaultCorsApiKey\")",
            "----------------------------------------","----------------------------------------",
            "  -config string",
            "        Path to the YAML configuration file (only when not using flags)",
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
    authenticatedPrefix := flag.String("GAG_AUTHENTICATED_PREFIX", "/defaultAuthenticatedPrefix", "Authenticated prefix")
    unauthenticatedPrefix := flag.String("GAG_UNATHETICATED_PREFIX", "/defaultUnauthenticatedPrefix", "Unauthenticated prefix")
    jwtRSAPublicKey := flag.String("GAG_JWT_RSA_PUBLIC_KEY", "/path/to/default/public/key", "Path to the JWT RSA public key")
    jwtRSAPrivateKey := flag.String("GAG_JWT_RSA_PRIVATE_KEY", "/path/to/default/private/key", "Path to the JWT RSA private key")
    serverAddress := flag.String("SERVER_ADDRESS", "localhost:8080", "Server address")
    destinationURL := flag.String("GAG_DESTINATION_URL", "http://defaultDestination.com", "Destination URL")
    corsApiKey := flag.String("CORS_API_KEY", "defaultCorsApiKey", "CORS API key")
    mode := flag.String("MODE", "", "Mode")
    flag.Parse()
    if len(os.Args) == 1 {
        flag.Usage()
        os.Exit(1)
    }
    // check mode is provided
    if *mode == "" {
        log.Fatal("Mode must be provided")
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