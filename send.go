package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func send(m Model) (string, string) {

	method := strings.ToUpper(strings.TrimSpace(m.methodField.Value()))
	URL := strings.TrimSpace(m.urlField.Value())
	headersJSON := strings.TrimSpace(m.tabContent[2].Value())
	paramsJSON := strings.TrimSpace(m.tabContent[1].Value())

	// Parse JSON into a map
	var headers map[string]string
	err := json.Unmarshal([]byte(headersJSON), &headers)
	if err != nil {
		return " \n Error parsing Headers \n\n Correct the Headers format", " Incorrect Headers "
	}

	if paramsJSON != "" {
		var params map[string]string
		Err := json.Unmarshal([]byte(paramsJSON), &params)
		if Err != nil {
			return " \n Error parsing Params \n\n Correct the Params format", " Incorrect Params "
		}

		// Create a URL object
		parsedURL, err := url.Parse(URL)
		if err != nil {
			return " \n Error parsing Params \n\n Correct the Params format", " Incorrect Params "
		}

		// Add query parameters to the URL
		q := parsedURL.Query()
		for key, value := range params {
			q.Set(key, value)
		}
		parsedURL.RawQuery = q.Encode()

		URL = parsedURL.String()
	}

	switch method {
	case "GET":
		// Make the GET request
		resp, err := http.Get(URL)
		if err != nil {
			return "Failed to make request\n\n" + err.Error(), ""
		}

		// Set headers
		for key, value := range headers {
			resp.Header.Set(key, value)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body\n\n" + err.Error(), ""
		}

		return string(body), resp.Status

	case "POST":

		payload := []byte(m.tabContent[0].Value())

		client := &http.Client{}
		req, err := http.NewRequest("POST", URL, bytes.NewBuffer(payload))
		if err != nil {
			return "Failed to make request\n\n" + err.Error(), ""
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			return "Failed to make request\n\n" + err.Error(), ""
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body\n\n" + err.Error(), ""
		}
		return string(body), resp.Status

	case "PUT":
		payload := []byte(m.tabContent[0].Value())

		client := &http.Client{}
		req, err := http.NewRequest("PUT", URL, bytes.NewBuffer(payload))
		if err != nil {
			return "Failed to make request\n\n" + err.Error(), ""
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			return "Failed to make request\n\n" + err.Error(), ""
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body\n\n" + err.Error(), ""
		}
		return string(body), resp.Status

	case "DELETE":
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", URL, nil)
		if err != nil {
			return "Failed to make request\n\n" + err.Error(), ""
		}

		resp, err := client.Do(req)
		if err != nil {
			return "Failed to make request\n\n" + err.Error(), ""
		}
		defer resp.Body.Close()

		// Set headers
		for key, value := range headers {
			resp.Header.Set(key, value)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body\n\n" + err.Error(), ""
		}
		return string(body), resp.Status

	case "HEAD":
		resp, err := http.Head(URL)
		if err != nil {
			return "Failed to make request\n\n" + err.Error(), ""
		}
		defer resp.Body.Close()

		// Set headers
		for key, value := range headers {
			resp.Header.Set(key, value)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body\n\n" + err.Error(), ""
		}
		return string(body), resp.Status

	case "PATCH":
		payload := []byte(m.tabContent[0].Value())

		client := &http.Client{}
		req, err := http.NewRequest("PATCH", URL, bytes.NewBuffer(payload))
		if err != nil {
			return "Failed to make request\n\n" + err.Error(), ""
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			return "Failed to make request\n\n" + err.Error(), ""
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Failed to read response body\n\n" + err.Error(), ""
		}

		return string(body), resp.Status

	default:
		return "Request Method or Url is set incorrectly", " Incorrect Request "
	}
}
