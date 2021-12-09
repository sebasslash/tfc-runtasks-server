package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/jsonapi"
)

func TfcCallback(callbackUrl string, accessToken string, body *RunTaskResponse) error {
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

	_, err = client.Do(request)
	return err
}

func SuccessfulRunTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var runTaskReq RunTaskRequest
	json.Unmarshal(reqBody, &runTaskReq)

	runTaskResp := &RunTaskResponse{
		Status:  "passed",
		Message: fmt.Sprintf("Successful Run Task Integration initiated by %s", runTaskReq.RunID),
		Url:     "http://localhost:10000/success",
	}

	if err := TfcCallback(runTaskReq.TaskResultCallbackUrl, runTaskReq.AccessToken, runTaskResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func FailedRunTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var runTaskReq RunTaskRequest
	json.Unmarshal(reqBody, &runTaskReq)

	runTaskResp := &RunTaskResponse{
		Status:  "failed",
		Message: fmt.Sprintf("Failed Run Task Integration initiated by %s", runTaskReq.RunID),
		Url:     "http://localhost:10000/failed",
	}

	if err := TfcCallback(runTaskReq.TaskResultCallbackUrl, runTaskReq.AccessToken, runTaskResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "TFC Run Tasks Server root")
}

func handleRequests() {
	http.HandleFunc("/", Root)
	http.HandleFunc("/success", SuccessfulRunTask)
	http.HandleFunc("/failed", FailedRunTask)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
	log.Println("Started server on port 10000")
}
