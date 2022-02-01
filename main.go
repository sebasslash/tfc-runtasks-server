package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
)

func (m *TfcWebhookManager) sendWebhookResponse() {
	log.Println("Register worker")
	for {
		select {
		case job := <-m.jobs:
			log.Printf("Received job with message: %s \n", job.Response.Message)
			go func() {
				err := TfcCallback(job.CallbackUrl, job.AccessToken, job.Response, job.Timeout)
				if err != nil {
					log.Printf("Error: %s \n", err.Error())
				}
			}()
		default:
		}
	}

}

func TfcCallback(callbackUrl string, accessToken string, body *RunTaskResponse, timeout string) error {
	out := bytes.NewBuffer(nil)
	if err := jsonapi.MarshalPayload(out, body); err != nil {
		return err
	}

	client := http.Client{}

	request, err := http.NewRequest("PATCH", callbackUrl, out)
	request.Header.Set("Content-Type", jsonapi.MediaType)
	request.Header.Set("Authorization", "Bearer "+accessToken)

	if err != nil {
		return err
	}

	if timeout != "" {
		i, _ := strconv.Atoi(timeout)
		time.Sleep(time.Duration(i) * time.Second)
		_, err = client.Do(request)
	} else {
		_, err = client.Do(request)
	}

	return err
}

func (m *TfcWebhookManager) SuccessfulRunTask(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get("timeout")
	if timeout != "" {
		i, err := strconv.Atoi(timeout)
		if err != nil {
			http.Error(w, "timeout query param must be an integer", http.StatusInternalServerError)
		}
		if i <= 0 {
			http.Error(w, "timeout query param must be greater than 0 seconds", http.StatusInternalServerError)
		}
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var runTaskReq RunTaskRequest
	json.Unmarshal(reqBody, &runTaskReq)

	runTaskResp := &RunTaskResponse{
		Status:  "passed",
		Message: fmt.Sprintf("Successful Run Task Integration initiated by %s [reset]", runTaskReq.RunID),
		Url:     "http://localhost:10000/success",
	}

	m.jobs <- &TfcWebhookJob{
		AccessToken: runTaskReq.AccessToken,
		CallbackUrl: runTaskReq.TaskResultCallbackUrl,
		Response:    runTaskResp,
		Timeout:     timeout,
	}

	w.WriteHeader(http.StatusOK)
}

func (m *TfcWebhookManager) FailedRunTask(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get("timeout")
	if timeout != "" {
		i, err := strconv.Atoi(timeout)
		if err != nil {
			http.Error(w, "timeout query param must be a string", http.StatusInternalServerError)
		}
		if i <= 0 {
			http.Error(w, "timeout query param must be greater than 0 seconds", http.StatusInternalServerError)
		}
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var runTaskReq RunTaskRequest
	json.Unmarshal(reqBody, &runTaskReq)

	runTaskResp := &RunTaskResponse{
		Status:  "failed",
		Message: fmt.Sprintf("Failed Run Task Integration initiated by %s [reset]", runTaskReq.RunID),
		Url:     "http://localhost:10000/failed",
	}

	m.jobs <- &TfcWebhookJob{
		AccessToken: runTaskReq.AccessToken,
		CallbackUrl: runTaskReq.TaskResultCallbackUrl,
		Response:    runTaskResp,
		Timeout:     timeout,
	}

	w.WriteHeader(http.StatusOK)
}

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TFC Run Tasks Server root")
}

func handleRequests() {
	manager := &TfcWebhookManager{
		jobs: make(chan *TfcWebhookJob, 1000),
	}

	go manager.sendWebhookResponse()

	r := mux.NewRouter()
	r.HandleFunc("/", Root).Methods("GET")
	r.HandleFunc("/success", manager.SuccessfulRunTask).Methods("POST")
	r.HandleFunc("/failed", manager.FailedRunTask).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
