// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSimpleJobSubmit(t *testing.T) {
	if os.Getenv("SLURM_REAL_SERVER_TEST") != "true" {
		t.Skip("Skipping real server tests")
	}

	// Fetch JWT token
	token, err := fetchJWTTokenViaSSH()
	require.NoError(t, err, "Failed to fetch JWT token")

	// Try different job submission payloads
	payloads := []struct {
		name string
		body map[string]interface{}
	}{
		{
			name: "minimal",
			body: map[string]interface{}{
				"job": map[string]interface{}{
					"script": "#!/bin/bash\necho test\n",
				},
			},
		},
		{
			name: "with_name",
			body: map[string]interface{}{
				"job": map[string]interface{}{
					"name":   "test-job",
					"script": "#!/bin/bash\necho test\n",
				},
			},
		},
		{
			name: "with_partition",
			body: map[string]interface{}{
				"job": map[string]interface{}{
					"name":      "test-job",
					"script":    "#!/bin/bash\necho test\n",
					"partition": "debug",
				},
			},
		},
		{
			name: "with_resources",
			body: map[string]interface{}{
				"job": map[string]interface{}{
					"name":         "test-job",
					"script":       "#!/bin/bash\necho test\n",
					"partition":    "debug",
					"minimum_cpus": 1,
				},
			},
		},
		{
			name: "full_payload",
			body: map[string]interface{}{
				"job": map[string]interface{}{
					"name":                       "test-job",
					"script":                     "#!/bin/bash\necho test\n",
					"partition":                  "debug",
					"minimum_cpus":               1,
					"current_working_directory":  "/tmp",
					"memory_per_node": map[string]interface{}{
						"set":    true,
						"number": 1024,
					},
					"time_limit": map[string]interface{}{
						"set":    true,
						"number": 5,
					},
				},
			},
		},
	}

	for _, payload := range payloads {
		t.Run(payload.name, func(t *testing.T) {
			jsonData, err := json.Marshal(payload.body)
			require.NoError(t, err)

			url := "http://rocky9:6820/slurm/v0.0.43/job/submit"
			req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(jsonData))
			require.NoError(t, err)

			req.Header.Set("X-SLURM-USER-TOKEN", token)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{
				Timeout: 30 * time.Second,
			}

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			t.Logf("Payload: %s", payload.name)
			t.Logf("Status Code: %d", resp.StatusCode)
			t.Logf("Response Body: %s", string(body))

			// Try to parse the response
			var result map[string]interface{}
			if err := json.Unmarshal(body, &result); err != nil {
				t.Logf("Failed to parse response as JSON: %v", err)
			} else {
				prettyJSON, _ := json.MarshalIndent(result, "", "  ")
				t.Logf("Response JSON:\n%s", string(prettyJSON))
			}

			// Don't fail the test, just log results
			if resp.StatusCode != 200 && resp.StatusCode != 201 {
				t.Logf("Request failed with status %d", resp.StatusCode)
			}
			
			fmt.Println("---")
		})
	}
}
