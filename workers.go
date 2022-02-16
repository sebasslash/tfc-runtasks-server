package main

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/google/jsonapi"
)

func CallbackWorker(URL string, accessToken string, body *RunTaskResponse, timeout string) error {
	out := bytes.NewBuffer(nil)
	if err := jsonapi.MarshalPayload(out, body); err != nil {
		return err
	}

	client := http.Client{}

	request, err := http.NewRequest("PATCH", URL, out)
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
