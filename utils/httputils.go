package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func GetJSON(url string, v interface{}) (bool, error) {
	resp, err := http.Get(url)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return false, errors.New("response status not 200")
	}

	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return false, err
	}
	return true, nil
}

func PostJSON(url string, v interface{}) (bool, error) {
	jsonReq, err := json.Marshal(v)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		//bodyBytes, _ := ioutil.ReadAll(resp.Body)
		//bodyString := string(bodyBytes)
		return false, errors.New("POST response status not 200")
	}
	return true, nil
}
