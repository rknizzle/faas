package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func Invoke(function string) (map[string]string, error) {
	// send HTTP request to server to invoke function AKA spin
	// up new docker container and run the code
	client := &http.Client{}
	r, err := http.NewRequest("POST", "http://localhost:5555/functions/"+function, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// get the JSON response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var invokeErr struct {
		message string
	}
	var invokeRes struct {
		response map[string]string
	}

	// check status code here and unmarshal into the appropriate struct then return the correct value
	if resp.StatusCode == 200 {
		err = json.Unmarshal(body, &invokeRes)
		if err != nil {
			return nil, err
		}
		return invokeRes.response, nil
	} else {
		err = json.Unmarshal(body, &invokeErr)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(invokeErr.message)
	}
}
