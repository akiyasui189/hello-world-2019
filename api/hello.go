package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "encoding/json"
)

var err error

type HelloWorldRequest struct {
    Message string `json:message,string`
}

type HelloWorldResponse struct {
    MessageId    string `json:"messageId,string"`
    RegisteredAt string `json:"registeredAt,string"`
}

// health check
func health (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "OK")
}

// API
func handler (w http.ResponseWriter, r *http.Request) {
    /* validation */
    // method
    if r.Method != http.MethodPost {
        // 405
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    // request header
    if r.Header.Get("Accept") != "application/json" {
        // 406
        w.WriteHeader(http.StatusNotAcceptable)
        return
    }
    if r.Header.Get("Content-Type") != "application/json" {
        // 400
        log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // request body
    contentLength := r.ContentLength
    log.Printf("Content-Length: %d", contentLength)
    if contentLength < 1 {
        // 400
        log.Printf("Content-Length: %d", contentLength)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // read body
    body := make([]byte, contentLength)
    r.Body.Read(body)
    /*
    _, err := r.Body.Read(body)
    if err != nil {
        // 500
        fmt.Printf("error: ", err)
        log.Println("Request Body Read Error :", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    */
    log.Printf("request body: %s", string(body))
    // convert from json to struct
    var request HelloWorldRequest
    json.Unmarshal(body, &request)
    /*
    if err != nil {
        // 500
        log.Println("Request Body JSON Parse Error :", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    */
    log.Printf("request message : %s", request.Message)
    // set response
    var response HelloWorldResponse
    response.MessageId = "1000000001"
    response.RegisteredAt = time.Now().Format("2006-01-02T15:04:05.000")
    responseBody, err := json.Marshal(response);
    if err != nil {
        // 500
        log.Println("Response Body JSON Encode Error :", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    log.Printf("response body : %s", string(responseBody))
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    w.Write(responseBody)
}

func main() {
    log.Printf("API app listening on port 8080")
    http.HandleFunc("/helloworld", handler)
    http.HandleFunc("/health", health)
    http.ListenAndServe(":8080", nil)
}

