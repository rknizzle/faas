package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func Invoke(function string, input []byte) (string, error) {
	// send HTTP request to server to invoke function AKA spin
	// up new docker container and run the code
	client := &http.Client{}
	r, err := http.NewRequest("POST", "http://localhost:5555/functions/"+function, bytes.NewBuffer(input))
	if err != nil {
		return "", err
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// get the JSON response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var invokeErr struct {
		Message string `json:"message"`
	}
	var invokeRes struct {
		Response string `json:"response"`
	}

	// check status code here and unmarshal into the appropriate struct then return the correct value
	if resp.StatusCode == 200 {
		err = json.Unmarshal(body, &invokeRes)
		if err != nil {
			return "", err
		}
		return invokeRes.Response, nil
	} else {
		err = json.Unmarshal(body, &invokeErr)
		if err != nil {
			return "", err
		}
		return "", errors.New(invokeErr.Message)
	}
}
