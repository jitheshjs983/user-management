package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type PanResponse struct {
	Result struct {
		Name string `json:"name"`
	} `json:"result"`
	RequestID  string `json:"request_id"`
	StatusCode string `json:"status-code"` // Note the dash in the key
}

func GetNameFromPan(pan string) (*PanResponse, error) {
	url := os.Getenv("KARZA_BASE_URL") + "pan"
	key := os.Getenv("KARZA_KEY")
	if key == "" {
		return nil, errors.New("KARZA_KEY environment variable is not set")
	}

	payload := map[string]interface{}{
		"pan":     pan,
		"consent": "Y",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-karza-key", key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("resp Data %+v\n", resp)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode JSON response into struct
	var response PanResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
