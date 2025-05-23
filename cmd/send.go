package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func send(m Model) (string, string, string) {
	var console strings.Builder
	console.WriteString("=== REQUEST DETAILS ===\n")
	variablesJSON := loadVariables()
	method := strings.ToUpper(strings.TrimSpace(m.methodField.Value()))
	URL := strings.TrimSpace(m.urlField.Value())
	headersJSON := strings.TrimSpace(m.tabContent[2].Value())
	paramsJSON := strings.TrimSpace(m.tabContent[1].Value())
	
	console.WriteString(fmt.Sprintf("Method: %s\n", method))
	console.WriteString(fmt.Sprintf("URL: %s\n", URL))

	// Parse variables into a map
	var variables map[string]string
	er := json.Unmarshal([]byte(variablesJSON), &variables)
	if er != nil {
		console.WriteString("ERROR: Failed to parse environment variables\n")
		return "\n Error parsing Env Variables", "Incorrect Env Variables", console.String()
	}

	URL = replacePlaceholders(URL, variables)
	headersJSON = replacePlaceholders(headersJSON, variables)
	paramsJSON = replacePlaceholders(paramsJSON, variables)

	// Parse JSON into a map
	var headers map[string]string
	err := json.Unmarshal([]byte(headersJSON), &headers)
	if err != nil {
		console.WriteString(fmt.Sprintf("ERROR: Failed to parse headers: %s\n", err.Error()))
		return " \n Error parsing Headers \n\n Correct the Headers format", " Incorrect Headers ", console.String()
	}
	
	console.WriteString("Headers:\n")
	for k, v := range headers {
		console.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
	}

	if paramsJSON != "" {
		var params map[string]string
		Err := json.Unmarshal([]byte(paramsJSON), &params)
		if Err != nil {
			console.WriteString("ERROR: Failed to parse params\n")
			return " \n Error parsing Params \n\n Correct the Params format", " Incorrect Params ", console.String()
		}

		console.WriteString("Query Params:\n")
		for k, v := range params {
			console.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}

		// Create a URL object
		parsedURL, err := url.Parse(URL)
		if err != nil {
			console.WriteString("ERROR: Failed to parse URL\n")
			return " \n Error parsing Params \n\n Correct the Params format", " Incorrect Params ", console.String()
		}

		// Add query parameters to the URL
		q := parsedURL.Query()
		for key, value := range params {
			q.Set(key, value)
		}
		parsedURL.RawQuery = q.Encode()

		URL = parsedURL.String()
		console.WriteString(fmt.Sprintf("Final URL: %s\n", URL))
	}

	switch method {
	case "GET":
		console.WriteString("\n=== SENDING REQUEST ===\n")
		start := time.Now()
		
		// Make the GET request
		resp, err := http.Get(URL)
		duration := time.Since(start)
		
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Request failed after %v\n%s\n", duration, err.Error()))
			return "Failed to make request\n\n" + err.Error(), "", console.String()
		}

		console.WriteString(fmt.Sprintf("Request completed in %v\n", duration))
		console.WriteString(fmt.Sprintf("Status: %s\n", resp.Status))

		// Note: http.Get() doesn't support custom headers
		if len(headers) > 0 {
			console.WriteString("WARNING: GET requests with http.Get() cannot send custom headers\n")
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Failed to read response body\n%s\n", err.Error()))
			return "Failed to read response body\n\n" + err.Error(), "", console.String()
		}

		console.WriteString(fmt.Sprintf("Response size: %d bytes\n", len(body)))
		return string(body), resp.Status, console.String()

	case "POST":
		content := replacePlaceholders(m.tabContent[0].Value(), variables)
		payload := []byte(content)
		
		console.WriteString(fmt.Sprintf("Body: %s\n", content))
		console.WriteString("\n=== SENDING REQUEST ===\n")
		start := time.Now()

		client := &http.Client{}
		req, err := http.NewRequest("POST", URL, bytes.NewBuffer(payload))
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Failed to create request\n%s\n", err.Error()))
			return "Failed to make request\n\n" + err.Error(), "", console.String()
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		duration := time.Since(start)
		
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Request failed after %v\n%s\n", duration, err.Error()))
			return "Failed to make request\n\n" + err.Error(), "", console.String()
		}
		
		console.WriteString(fmt.Sprintf("Request completed in %v\n", duration))
		console.WriteString(fmt.Sprintf("Status: %s\n", resp.Status))
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Failed to read response body\n%s\n", err.Error()))
			return "Failed to read response body\n\n" + err.Error(), "", console.String()
		}
		
		console.WriteString(fmt.Sprintf("Response size: %d bytes\n", len(body)))
		return string(body), resp.Status, console.String()

	case "PUT":
		content := replacePlaceholders(m.tabContent[0].Value(), variables)
		payload := []byte(content)
		
		console.WriteString(fmt.Sprintf("Body: %s\n", content))
		console.WriteString("\n=== SENDING REQUEST ===\n")
		start := time.Now()

		client := &http.Client{}
		req, err := http.NewRequest("PUT", URL, bytes.NewBuffer(payload))
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Failed to create request\n%s\n", err.Error()))
			return "Failed to make request\n\n" + err.Error(), "", console.String()
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		duration := time.Since(start)
		
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Request failed after %v\n%s\n", duration, err.Error()))
			return "Failed to make request\n\n" + err.Error(), "", console.String()
		}
		
		console.WriteString(fmt.Sprintf("Request completed in %v\n", duration))
		console.WriteString(fmt.Sprintf("Status: %s\n", resp.Status))
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Failed to read response body\n%s\n", err.Error()))
			return "Failed to read response body\n\n" + err.Error(), "", console.String()
		}
		
		console.WriteString(fmt.Sprintf("Response size: %d bytes\n", len(body)))
		return string(body), resp.Status, console.String()

	case "DELETE":
		console.WriteString("\n=== SENDING REQUEST ===\n")
		start := time.Now()
		
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", URL, nil)
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Failed to create request\n%s\n", err.Error()))
			return "Failed to make request\n\n" + err.Error(), "", console.String()
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		duration := time.Since(start)
		
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Request failed after %v\n%s\n", duration, err.Error()))
			return "Failed to make request\n\n" + err.Error(), "", console.String()
		}
		
		console.WriteString(fmt.Sprintf("Request completed in %v\n", duration))
		console.WriteString(fmt.Sprintf("Status: %s\n", resp.Status))
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Failed to read response body\n%s\n", err.Error()))
			return "Failed to read response body\n\n" + err.Error(), "", console.String()
		}
		
		console.WriteString(fmt.Sprintf("Response size: %d bytes\n", len(body)))
		return string(body), resp.Status, console.String()

	case "HEAD":
		console.WriteString("\n=== SENDING REQUEST ===\n")
		start := time.Now()
		
		resp, err := http.Head(URL)
		duration := time.Since(start)
		
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Request failed after %v\n%s\n", duration, err.Error()))
			return "Failed to make request\n\n" + err.Error(), "", console.String()
		}
		
		console.WriteString(fmt.Sprintf("Request completed in %v\n", duration))
		console.WriteString(fmt.Sprintf("Status: %s\n", resp.Status))

		// Note: http.Head() doesn't support custom headers
		if len(headers) > 0 {
			console.WriteString("WARNING: HEAD requests with http.Head() cannot send custom headers\n")
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Failed to read response body\n%s\n", err.Error()))
			return "Failed to read response body\n\n" + err.Error(), "", console.String()
		}
		
		console.WriteString(fmt.Sprintf("Response size: %d bytes\n", len(body)))
		return string(body), resp.Status, console.String()

	case "PATCH":
		content := replacePlaceholders(m.tabContent[0].Value(), variables)
		payload := []byte(content)
		
		console.WriteString(fmt.Sprintf("Body: %s\n", content))
		console.WriteString("\n=== SENDING REQUEST ===\n")
		start := time.Now()

		client := &http.Client{}
		req, err := http.NewRequest("PATCH", URL, bytes.NewBuffer(payload))
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Failed to create request\n%s\n", err.Error()))
			return "Failed to make request\n\n" + err.Error(), "", console.String()
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		duration := time.Since(start)
		
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Request failed after %v\n%s\n", duration, err.Error()))
			return "Failed to make request\n\n" + err.Error(), "", console.String()
		}
		
		console.WriteString(fmt.Sprintf("Request completed in %v\n", duration))
		console.WriteString(fmt.Sprintf("Status: %s\n", resp.Status))
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			console.WriteString(fmt.Sprintf("ERROR: Failed to read response body\n%s\n", err.Error()))
			return "Failed to read response body\n\n" + err.Error(), "", console.String()
		}
		
		console.WriteString(fmt.Sprintf("Response size: %d bytes\n", len(body)))
		return string(body), resp.Status, console.String()

	default:
		console.WriteString("ERROR: Unsupported HTTP method\n")
		return "Request Method or Url is set incorrectly", " Incorrect Request ", console.String()
	}
}
