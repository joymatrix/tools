package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	client = http.Client{}
)

func InitHttp() {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
	}
	//client.Timeout = time.Duration(appConfig.HttpClient.TimeoutSeconds)
	client.Transport = transport
}

func HttpRequest(method, endpoint string, data []byte) ([]byte, error) {
	var req *http.Request
	var resp *http.Response
	var err error
	body := make([]byte, 0)

	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(method, endpoint, nil)
	case http.MethodPost:
		req, err = http.NewRequest(method, endpoint, bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer 6d5bd52c1ba448319b8e301a34576c8d")
	default:
		req, err = http.NewRequest(method, endpoint, bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
	}

	if err != nil {
		fmt.Printf("new request err:%+v", err.Error())
		return body, err
	}

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("do request err:%+v", err.Error())
		return body, err
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read err:%+v", err.Error())
		return body, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 0 {
		err := errors.New("resp.StatusCode not 200")
		return body, err
	}

	return body, nil
}
