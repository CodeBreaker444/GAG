package tests

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/codebreaker444/gag/api/handlers"
    "github.com/codebreaker444/gag/utils"
)


func TestForwardRequest(t *testing.T) {
    // Create a request to pass to our handler
    
    req, err := http.NewRequest("GET", "/test", bytes.NewBuffer([]byte("test body")))
    if err != nil {
        t.Fatal(err)
    }

    // Set some headers
    // req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer test")

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()
    handler := &handlers.Handler{
        // Initialize your handler here
        Config: utils.Config{
            Mode: "CORS",
            CorsApiKey: "defaultCorsApiKey",
        },



    }
    // method is GET, schema is https, corsurl is api.sampleapis.com/coffee/hot
	
    handler.ForwardRequest(rr, req, "https", "reqres.in/api/users?page=2")

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Check the response body
    // expected := `{"message":"test"}`
    // if rr.Body.String() != expected {
    //     t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    // }

    // Check the response headers
    if contentType := rr.Header().Get("Access-Control-Allow-Origin"); contentType != "*" {
        t.Errorf("handler returned wrong Access-Control-Allow-Origin header: got %v want %v", contentType, "*")
    }
    
    if methods := rr.Header().Get("Access-Control-Allow-Methods"); methods != "GET, POST, PUT, DELETE, OPTIONS" {
        t.Errorf("handler returned wrong Access-Control-Allow-Methods header: got %v want %v", methods, "POST, GET, OPTIONS, PUT, DELETE")
    }
    
  

    if authorization := rr.Header().Get("Authorization"); authorization == "Bearer test" {
        t.Errorf("handler returned wrong authorization header: got %v want %v", authorization, "Bearer test")
    }
}