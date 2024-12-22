package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func send(m Model) (string, string) {

	method := strings.ToUpper(strings.TrimSpace(m.methodField.Value()))
	url := strings.TrimSpace(m.urlField.Value())
	headersJSON := strings.TrimSpace(m.tabContent[3].Value())

	// Parse JSON into a map
	var headers map[string]string
	err := json.Unmarshal([]byte(headersJSON), &headers)
	if err != nil {
		return "Error parsing Headers", "Incorrect Headers"
	}

	switch method {
	case "GET":
		// Make the GET request
		resp, err := http.Get(url)
		if err != nil {
			return "Failed to make request", " Incorrect Request "
		}

		// Set headers
		for key, value := range headers {
			resp.Header.Set(key, value)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body", ""
		}

		return string(body), resp.Status

	case "POST":

		payload := []byte(m.tabContent[0].Value())

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return "Failed to make request", " Incorrect Request "
		}
		defer resp.Body.Close()

		// Set headers
		for key, value := range headers {
			resp.Header.Set(key, value)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body", ""
		}
		return string(body), resp.Status

	case "PUT":
		payload := []byte(m.tabContent[0].Value())

		client := &http.Client{}
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
		if err != nil {
			return "Failed to create request", " Incorrect Request "
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			return "Failed to make request", " Incorrect Request "
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body", ""
		}
		return string(body), resp.Status

	case "DELETE":
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return "Failed to create request", " Incorrect Request "
		}

		resp, err := client.Do(req)
		if err != nil {
			return "Failed to make request", " Incorrect Request "
		}
		defer resp.Body.Close()

		// Set headers
		for key, value := range headers {
			resp.Header.Set(key, value)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body", ""
		}
		return string(body), resp.Status

	case "HEAD":
		resp, err := http.Head(url)
		if err != nil {
			return "Failed to make request", " Incorrect Request "
		}
		defer resp.Body.Close()

		// Set headers
		for key, value := range headers {
			resp.Header.Set(key, value)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body", ""
		}
		return string(body), resp.Status

	case "PATCH":
		payload := []byte(m.tabContent[0].Value())

		client := &http.Client{}
		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payload))
		if err != nil {
			return "Failed to create request", " Incorrect Request "
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			return "Failed to make request", " Incorrect Request "
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body", ""
		}

		return string(body), resp.Status

	default:
		return "Request Method or Url is set incorrectly", " Incorrect Request "
	}
}
