package easynet

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func Test() {
	fmt.Println("Hello")
}

func executeRequest(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

func Get(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating GET request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Post(url string, data []byte, headers map[string]string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Put(url string, data []byte, headers map[string]string) (string, error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating PUT request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Delete(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating DELETE request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Patch(url string, data []byte, headers map[string]string) (string, error) {
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error creating PATCH request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Options(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating OPTIONS request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}

func Head(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating HEAD request: %w", err)
	}
	addHeaders(req, headers)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing HEAD request: %w", err)
	}
	defer resp.Body.Close()

	return resp.Status, nil
}

func Trace(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("TRACE", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating TRACE request: %w", err)
	}
	addHeaders(req, headers)
	return executeRequest(req)
}
