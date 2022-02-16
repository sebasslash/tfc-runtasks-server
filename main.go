package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	HmacKey string
	Host    string
)

func hmacValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sha := r.Header.Get("x-tfc-task-signature"); sha != "" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
			}

			h := hmac.New(sha512.New, []byte(HmacKey))
			h.Write(body)

			expectedSha := hex.EncodeToString(h.Sum(nil))
			if !hmac.Equal([]byte(expectedSha), []byte(sha)) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			// since the body was read, we need to restore it
			r.Body = ioutil.NopCloser(bytes.NewReader(body))
		}

		next.ServeHTTP(w, r)
	})
}

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TFC Run Tasks Server root")
}

func handleRequests() {
	if Host = os.Getenv("RUN_TASK_HOST"); Host == "" {
		Host = "localhost"
	}

	manager := &TfcWebhookManager{
		jobs: make(chan *CallbackJob, 1000),
	}

	go manager.registerWorkers()

	r := mux.NewRouter()
	r.Use(hmacValidationMiddleware)
	r.HandleFunc("/", Root).Methods("GET")
	r.HandleFunc("/success", manager.SuccessfulRunTask).Methods("POST")
	r.HandleFunc("/failed", manager.FailedRunTask).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	if HmacKey = os.Getenv("TFC_TASK_HMAC_KEY"); HmacKey == "" {
		HmacKey = "hashicorp"
	}
	handleRequests()
}
