package handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	utils "github.com/codebreaker444/gag/utils"
)
type Handler struct {
	Config utils.Config
}

func (h *Handler) RootHandler(rootRouter *http.ServeMux) {
	log.Println("RootHandler",h.Config.AuthenticatedPrefix)
	rootRouter.HandleFunc("/test/{routetest}", testRoute)
	if h.Config.Mode == "CORS" {
		rootRouter.HandleFunc("/", h.corsRoot)
	}else{
		rootRouter.HandleFunc("/{route}", h.root)
	}
	// rootRouter.HandleFunc("/", root)

}

func (h *Handler) root (w http.ResponseWriter, r *http.Request) {
	route := r.PathValue("route")
	log.Println("Root handler: ",route)
	log.Println("Root handler: ",route)
	h.forwardRequest(w, r, "http", route)
	
}
func (h *Handler) corsRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("CORS Root handler: ",r.URL)
    url, err := url.Parse(r.URL.String())
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    segments := strings.Split(url.Path, "/")

    if len(segments) < 3 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
	schema := strings.Split(segments[1], ":")

    secondURL := strings.Join(segments[2:], "/")
	log.Println("CORS Root handler: ",segments)
	log.Println("CORS Root handler: ",secondURL)
	log.Println("CORS Root handler Schema: ",schema[0])
	if schema[0] != "http" && schema[0] != "https" {
		http.Error(w, "Invalid URL Schema", http.StatusBadRequest)
		return
	}
	h.forwardRequest(w, r, schema[0], secondURL)
}
func (h *Handler) forwardRequest(w http.ResponseWriter, r *http.Request, schema string, corsurl string) {
	// get method, path, body, headers from the request
	method := r.Method
	body := r.Body
	var urlPath string
	log.Println("ForwardRequest:", method, "Body:", body, "Headers:", r.Header)
	var destinationUrl string
	if h.Config.Mode == "CORS" {
		u, err := url.Parse(schema+"://"+corsurl)
		if err != nil {
			// handle error
			log.Println("Error in parsing URL: ", err)
		}
	
		destinationUrl = u.Host
		urlPath = u.Path
		schema = u.Scheme
		log.Println("Destination URL: ",destinationUrl,"URL Path: ",urlPath, "Schema: ",schema)
	}else{
		destinationUrl = h.Config.DestinationURL
		urlPath =r.PathValue("route")

		
	}
	r.Header.Del("X-Gag-Api-Key")
	r.URL = &url.URL{
		Scheme: schema,
		Host: destinationUrl,
		Path: urlPath,
	}
	log.Println(r)
		
	reverseProxy:= httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: schema,
		Host: destinationUrl,
	})
	log.Println("urlPath: ",reverseProxy)
	reverseProxy.ModifyResponse = func(response *http.Response) error {
		_ = h.reverseProxyResponseModifier(response)
		return nil
	}	
	reverseProxy.ServeHTTP(w,r)
}

func (h *Handler) reverseProxyResponseModifier(response *http.Response) error {
	// log.Println("Response from RPMODIFIER: ", response)
	// set cors headers to the response\
	if h.Config.Mode != "CORS" {
		return nil
	}
	response.Header.Set("Access-Control-Allow-Origin", "*")
	response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	response.Header.Set("X-API-GATEWAY", "github.com/codebreaker444/gag")
	return nil
}

func testRoute(wclient http.ResponseWriter, rclient *http.Request) {
	log.Println("Microservice:",rclient.URL)

	// set the entire header same as the request header
	wclient.Header().Set("X-API-GATEWAY-TEST", "github.com/codebreaker444/gag")

	for k, v := range rclient.Header {
		wclient.Header().Set(k, v[0])
	}
	// set the body of the response as the body of the request
	log.Println("Request from client body: ",rclient.Body )
	_, err := io.Copy(wclient, rclient.Body)

	// set header
	if err != nil {
		wclient.WriteHeader(http.StatusInternalServerError)
		return
	}
	// print pointer address of the rcient
	log.Println("Response from microservice: ", &rclient)

	wclient.WriteHeader(http.StatusOK)
	// print body of the request
	// write response body same as the request body
	
	
}