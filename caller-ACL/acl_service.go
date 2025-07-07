package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http" // Added for os.Getenv example
	"time"
)

// OPAInput represents the structure of the input required by your Rego policy.
type ACLInput struct {
	Input struct {
		Method string `json:"method"`
		User   struct {
			Role string `json:"role"`
		} `json:"user"`
	} `json:"input"`
}

// OPAResponse represents the structure of the response from OPA.
type ACLResponse struct {
	Result bool `json:"result"`
}

// OPAService represents the OPA client and its configuration.
// This is the struct that your "injector service" would ideally return.
type ACLService struct {
	serverURL  string
	httpClient *http.Client
}

// NewOPAService creates and returns a new OPAService instance.
// This function mimics what your "injector service" would do.
// It could take other parameters for configuration (e.g., custom http.Client, timeouts).
func NewACLService(url string) *ACLService {
	return &ACLService{
		serverURL:  url,
		httpClient: &http.Client{Timeout: 5 * time.Second}, // Configure client once
	}
}

// EvaluatePolicy sends an authorization request to the OPA server
// associated with this OPAService instance and returns the boolean decision.
func (s *ACLService) Authorize(method, userRole string) (bool, error) {
	// Construct the full URL to your policy endpoint
	// For your `authz.rego` with `package authz`, the endpoint is /v1/data/authz/allow
	endpoint := fmt.Sprintf("%s/v1/data/authz/allow", s.serverURL)

	// Prepare the input data according to your Rego policy's expectations
	input := ACLInput{}
	input.Input.Method = method
	input.Input.User.Role = userRole

	// Marshal the input struct to JSON
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return false, fmt.Errorf("failed to marshal OPA input: %w", err)
	}

	// Create the HTTP POST request using the service's HTTP client
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonInput))
	if err != nil {
		return false, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send request to OPA: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read OPA response body: %w", err)
	}

	// Check for non-2xx status codes
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("OPA server returned non-OK status: %s (body: %s)", resp.Status, string(body))
	}

	// Unmarshal the JSON response
	var opaResponse ACLResponse
	err = json.Unmarshal(body, &opaResponse)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal OPA response: %w (body: %s)", err, string(body))
	}

	return opaResponse.Result, nil
}
