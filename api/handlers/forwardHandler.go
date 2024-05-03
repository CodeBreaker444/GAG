package handlers
import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

)

func (h *Handler) ForwardRequest(w http.ResponseWriter, r *http.Request, schema string, corsurl string) {
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
	log.Println("Response from RPMODIFIER: ", response)
	if h.Config.Mode != "CORS" {
		return nil
	}
	response.Header.Set("Access-Control-Allow-Origin", "*")
	response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	response.Header.Set("X-API-GATEWAY", "github.com/codebreaker444/gag")
	return nil
}
