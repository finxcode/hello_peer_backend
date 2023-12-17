package helper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	url = "https://www.hellopeer.net/moment/task/complete"
)

type CompleteRequest struct {
	TaskId int `json:"taskId"`
	UserId int `json:"userId"`
}

type CompleteResponse struct {
	Status    int    `json:"status"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Data      bool   `json:"data"`
}

func CompleteTask(cReq CompleteRequest, token string) (*CompleteResponse, error) {
	// Marshal the request data to JSON
	jsonReq, err := json.Marshal(cReq)
	if err != nil {
		return nil, err
	}

	// Make the POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	// Set the content type and the authorization headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+token)

	// Send the request
	client := &http.Client{}
	//dumpRequest, _ := httputil.DumpRequest(req, true)
	//fmt.Println(string(dumpRequest))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data
	var compResp CompleteResponse
	err = json.Unmarshal(body, &compResp)
	if err != nil {
		return nil, err
	}

	return &compResp, nil
}
