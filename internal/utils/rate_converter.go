package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiResponse struct {
	StatusCode int `json:"status_code"`
	Data       struct {
		Base      string  `json:"base"`
		Target    string  `json:"target"`
		Mid       float64 `json:"mid"`
		Unit      int     `json:"unit"`
		Timestamp string  `json:"timestamp"`
	} `json:"data"`
}

func GetLatestCurrencyRate(currencyCode string) (*ApiResponse, error) {
	if currencyCode == "IDR" {
		return nil, nil
	}

	url := fmt.Sprintf("https://hexarate.paikama.co/api/rates/latest/%s?target=IDR", currencyCode)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
