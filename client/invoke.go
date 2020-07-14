package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Invoke(function string) error {
	// send HTTP request to server to invoke function AKA spin
	// up new docker container and run the code
	client := &http.Client{}
	r, err := http.NewRequest("POST", "http://localhost:5555/functions/rkneills/"+function, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// get the JSON response
	var result map[string]string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	return nil
}
