package slurm

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// makeRequest makes an HTTP request with retry logic
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}
	
	// Build full URL
	fullURL, err := c.buildURL(endpoint)
	if err != nil {
		return fmt.Errorf("failed to build URL: %w", err)
	}
	
	var lastErr error
	var lastResp *http.Response
	
	for attempt := 0; attempt <= c.retry.MaxRetries(); attempt++ {
		// Create request
		req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		
		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.config.UserAgent)
		
		// Add authentication
		if c.auth != nil {
			if err := c.auth.Authenticate(ctx, req); err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
		}
		
		// Make the request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if c.retry.ShouldRetry(ctx, nil, err, attempt) {
				waitTime := c.retry.WaitTime(attempt)
				if c.config.Debug {
					fmt.Printf("Request failed (attempt %d), retrying in %v: %v\n", attempt+1, waitTime, err)
				}
				time.Sleep(waitTime)
				continue
			}
			return fmt.Errorf("request failed: %w", err)
		}
		
		lastResp = resp
		
		// Check if we should retry based on status code
		if c.retry.ShouldRetry(ctx, resp, nil, attempt) {
			waitTime := c.retry.WaitTime(attempt)
			if c.config.Debug {
				fmt.Printf("Request returned %d (attempt %d), retrying in %v\n", resp.StatusCode, attempt+1, waitTime)
			}
			resp.Body.Close()
			time.Sleep(waitTime)
			continue
		}
		
		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		
		// Check for HTTP errors
		if resp.StatusCode >= 400 {
			var slurmErr SlurmError
			if err := json.Unmarshal(respBody, &slurmErr); err != nil {
				// If we can't parse the error, create a generic one
				slurmErr = SlurmError{
					Code:    resp.StatusCode,
					Message: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(respBody)),
					Source:  "slurm-client",
				}
			}
			return &slurmErr
		}
		
		// Parse response if result is provided
		if result != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("failed to unmarshal response: %w", err)
			}
		}
		
		return nil
	}
	
	if lastErr != nil {
		return lastErr
	}
	
	if lastResp != nil {
		return fmt.Errorf("max retries exceeded, last status: %d", lastResp.StatusCode)
	}
	
	return fmt.Errorf("max retries exceeded")
}

// buildURL builds the full URL for an endpoint
func (c *Client) buildURL(endpoint string) (string, error) {
	baseURL := c.baseURL
	if baseURL == "" {
		baseURL = c.config.BaseURL
	}
	
	// Parse base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	
	// Parse endpoint
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint: %w", err)
	}
	
	// Resolve relative URL
	resolved := base.ResolveReference(endpointURL)
	
	return resolved.String(), nil
}

// configureHTTPClient configures the HTTP client with TLS and timeout settings
func (c *Client) configureHTTPClient() {
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}
	
	// Set timeout
	c.httpClient.Timeout = c.config.Timeout
	
	// Configure TLS
	if c.config.InsecureSkipVerify {
		if c.httpClient.Transport == nil {
			c.httpClient.Transport = &http.Transport{}
		}
		
		if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
			if transport.TLSClientConfig == nil {
				transport.TLSClientConfig = &tls.Config{}
			}
			transport.TLSClientConfig.InsecureSkipVerify = true
		}
	}
}