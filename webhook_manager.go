package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrInvalidTimeout     = errors.New("invalid timeout, must be an int greator than zero")
	ErrInvalidRequestBody = errors.New("failed to parse request body")
)

type TfcWebhookManager struct {
	jobs chan *CallbackJob
}

func (m *TfcWebhookManager) registerWorkers() {
	for {
		select {
		case job := <-m.jobs:
			log.Printf("Received job with message: %s \n", job.Response.Message)
			go func() {
				log.Println("Register worker")
				err := CallbackWorker(job.CallbackUrl, job.AccessToken, job.Response, job.Timeout)
				if err != nil {
					log.Printf("callback error: %s", err)
				}
			}()
		default:
		}
	}
}

func (m *TfcWebhookManager) SuccessfulRunTask(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get("timeout")
	if err := m.validTimeout(timeout); err == ErrInvalidTimeout {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	body, err := m.parseBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	runTaskResp := &RunTaskResponse{
		Status:  "passed",
		Message: fmt.Sprintf("Successful Run Task Integration initiated by %s", body.RunID),
		Url:     fmt.Sprintf("https://%s/success", Host),
	}

	m.jobs <- &CallbackJob{
		AccessToken: body.AccessToken,
		CallbackUrl: body.TaskResultCallbackUrl,
		Response:    runTaskResp,
		Timeout:     timeout,
	}

	w.WriteHeader(http.StatusOK)
}

func (m *TfcWebhookManager) FailedRunTask(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get("timeout")
	if err := m.validTimeout(timeout); err == ErrInvalidTimeout {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	body, err := m.parseBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	runTaskResp := &RunTaskResponse{
		Status:  "failed",
		Message: fmt.Sprintf("Failed Run Task Integration initiated by %s", body.RunID),
		Url:     fmt.Sprintf("https://%s/failed", Host),
	}

	m.jobs <- &CallbackJob{
		AccessToken: body.AccessToken,
		CallbackUrl: body.TaskResultCallbackUrl,
		Response:    runTaskResp,
		Timeout:     timeout,
	}

	w.WriteHeader(http.StatusOK)
}

func (m *TfcWebhookManager) validTimeout(timeout string) error {
	if timeout != "" {
		i, err := strconv.Atoi(timeout)
		if err != nil || i <= 0 {
			return ErrInvalidTimeout
		}
	}
	return nil
}

func (m *TfcWebhookManager) parseBody(r *http.Request) (*RunTaskRequest, error) {
	var runTaskReq RunTaskRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, ErrInvalidRequestBody
	}

	err = json.Unmarshal(reqBody, &runTaskReq)
	if err != nil {
		return nil, ErrInvalidRequestBody
	}

	return &runTaskReq, nil
}
