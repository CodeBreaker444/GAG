package handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)
func root (w http.ResponseWriter, r *http.Request) {
	route := r.PathValue("route")
	log.Println("Root handler: ",r.Body)
	log.Println("Root handler: ",route)
	forwardRequest(w, r)
	

}

func RootHandler(rootRouter *http.ServeMux) {

	rootRouter.HandleFunc("/test/{route}", testRoute)

	rootRouter.HandleFunc("/{route}", root)
	// rootRouter.HandleFunc("/", root)
}
func testRoute(wclient http.ResponseWriter, rclient *http.Request) {
	log.Println("Microservice:",rclient.URL)

	// set the entire header same as the request header
	for k, v := range rclient.Header {
		wclient.Header().Set(k, v[0])
	}
	// set the body of the response as the body of the request
	log.Println("Request from client body: ",rclient.Body )
	_, err := io.Copy(wclient, rclient.Body)
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

// forward the request to microservice tp http:://localhost:8080
func forwardRequest(w http.ResponseWriter, r *http.Request) {
	// get method, path, body, headers from the request
	method := r.Method
	body := r.Body
	urlPath :=r.PathValue("route")
	//print host
	log.Println("ForwardRequest:",r.Host)

	log.Println("ForwardRequest:",method)

	// create a new request with the method, path, body, headers
	log.Println("ForwardRequest:", method, "Path:", body)


	
	reverseProxy:= httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host: "localhost:8000",
	})
	log.Println("urlPath: ",urlPath)
	r.URL.Path = "/test/"+urlPath

	reverseProxy.ServeHTTP(w,r)

	



}
