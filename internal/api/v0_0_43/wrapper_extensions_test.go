package v0_0_43

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWrapperClient_GetLicenses tests the GetLicenses operation
func TestWrapperClient_GetLicenses(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockClient)
		expectedError  string
		expectedCount  int
		validateResult func(*testing.T, *interfaces.LicenseList)
	}{
		{
			name: "successful licenses retrieval",
			setupMock: func(m *MockClient) {
				licenseName := "test-license"
				total := int32(100)
				used := int32(50)
				available := int32(50)
				reserved := int32(10)
				remote := false
				server := "test-server"

				m.EXPECT().SlurmV0043GetLicensesWithResponse(
					context.Background(),
				).Return(&SlurmV0043GetLicensesResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiLicensesResp{
						Licenses: &[]V0043License{
							{
								LicenseName: &licenseName,
								Total:       &total,
								Used:        &used,
								Available:   &available,
								Reserved:    &reserved,
								Remote:      &remote,
								Server:      &server,
							},
						},
					},
				}, nil)
			},
			expectedCount: 1,
			validateResult: func(t *testing.T, result *interfaces.LicenseList) {
				require.Len(t, result.Licenses, 1)
				license := result.Licenses[0]
				assert.Equal(t, "test-license", license.Name)
				assert.Equal(t, 100, license.Total)
				assert.Equal(t, 50, license.Used)
				assert.Equal(t, 50, license.Available)
				assert.Equal(t, 10, license.Reserved)
				assert.False(t, license.Remote)
				assert.Equal(t, "test-server", license.Server)
				assert.Equal(t, 50.0, license.Percent)
			},
		},
		{
			name: "API error response",
			setupMock: func(m *MockClient) {
				errorNumber := int32(1001)
				errorMsg := "API Error"
				description := "Test API error"
				source := "test"

				m.EXPECT().SlurmV0043GetLicensesWithResponse(
					context.Background(),
				).Return(&SlurmV0043GetLicensesResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiLicensesResp{
						Errors: &[]V0043Error{
							{
								ErrorNumber: &errorNumber,
								Error:       &errorMsg,
								Description: &description,
								Source:      &source,
							},
						},
					},
				}, nil)
			},
			expectedError: "API returned errors",
		},
		{
			name: "HTTP error",
			setupMock: func(m *MockClient) {
				m.EXPECT().SlurmV0043GetLicensesWithResponse(
					context.Background(),
				).Return(&SlurmV0043GetLicensesResponse{
					HTTPResponse: &http.Response{StatusCode: 500},
				}, nil)
			},
			expectedError: "Get licenses failed with status 500",
		},
		{
			name: "uninitialized client",
			setupMock: func(m *MockClient) {
				// No setup needed - client will be nil
			},
			expectedError: "API client not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var client *WrapperClient
			if !strings.Contains(tt.expectedError, "not initialized") {
				mockClient := NewMockClient(t)
				tt.setupMock(mockClient)
				client = &WrapperClient{
					apiClient: &ClientWithResponses{
						ClientInterface: mockClient,
					},
				}
			} else {
				client = &WrapperClient{} // Uninitialized client
			}

			result, err := client.GetLicenses(context.Background())

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
				if tt.expectedCount > 0 {
					assert.Len(t, result.Licenses, tt.expectedCount)
				}
			}
		})
	}
}

// TestWrapperClient_GetShares tests the GetShares operation
func TestWrapperClient_GetShares(t *testing.T) {
	tests := []struct {
		name           string
		options        *interfaces.GetSharesOptions
		setupMock      func(*MockClient, *interfaces.GetSharesOptions)
		expectedError  string
		expectedCount  int
		validateResult func(*testing.T, *interfaces.SharesList)
	}{
		{
			name: "successful shares retrieval",
			options: &interfaces.GetSharesOptions{
				Users:    []string{"user1", "user2"},
				Accounts: []string{"account1"},
				Clusters: []string{"cluster1"},
			},
			setupMock: func(m *MockClient, opts *interfaces.GetSharesOptions) {
				users := strings.Join(opts.Users, ",")
				accounts := strings.Join(opts.Accounts, ",")
				clusters := strings.Join(opts.Clusters, ",")

				name := "test-share"
				user := "user1"
				account := "account1"
				cluster := "cluster1"
				shares := int32(100)
				rawShares := int32(100)
				normShares := 1.0
				rawUsage := int32(50)
				normUsage := 0.5
				effectUsage := 0.5
				fairShare := 0.8
				levelFS := 0.9
				priority := 1000.0
				level := int32(1)

				m.EXPECT().SlurmV0043GetSharesWithResponse(
					context.Background(),
					&SlurmV0043GetSharesParams{
						Users:    &users,
						Accounts: &accounts,
						Clusters: &clusters,
					},
				).Return(&SlurmV0043GetSharesResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiSharesResp{
						Shares: &[]V0043SharesResponseSharesInner{
							{
								Name:        &name,
								User:        &user,
								Account:     &account,
								Cluster:     &cluster,
								Shares:      &shares,
								RawShares:   &rawShares,
								NormShares:  &normShares,
								RawUsage:    &rawUsage,
								NormUsage:   &normUsage,
								EffectUsage: &effectUsage,
								FairShare:   &fairShare,
								LevelFS:     &levelFS,
								Priority:    &priority,
								Level:       &level,
							},
						},
					},
				}, nil)
			},
			expectedCount: 1,
			validateResult: func(t *testing.T, result *interfaces.SharesList) {
				require.Len(t, result.Shares, 1)
				share := result.Shares[0]
				assert.Equal(t, "test-share", share.Name)
				assert.Equal(t, "user1", share.User)
				assert.Equal(t, "account1", share.Account)
				assert.Equal(t, "cluster1", share.Cluster)
				assert.Equal(t, 100, share.Shares)
				assert.Equal(t, 0.8, share.FairShare)
			},
		},
		{
			name:    "nil options",
			options: nil,
			setupMock: func(m *MockClient, opts *interfaces.GetSharesOptions) {
				m.EXPECT().SlurmV0043GetSharesWithResponse(
					context.Background(),
					&SlurmV0043GetSharesParams{},
				).Return(&SlurmV0043GetSharesResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiSharesResp{
						Shares: &[]V0043SharesResponseSharesInner{},
					},
				}, nil)
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(t)
			tt.setupMock(mockClient, tt.options)
			client := &WrapperClient{
				apiClient: &ClientWithResponses{
					ClientInterface: mockClient,
				},
			}

			result, err := client.GetShares(context.Background(), tt.options)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
				if tt.expectedCount >= 0 {
					assert.Len(t, result.Shares, tt.expectedCount)
				}
			}
		})
	}
}

// TestWrapperClient_GetConfig tests the GetConfig operation
func TestWrapperClient_GetConfig(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockClient)
		expectedError string
	}{
		{
			name: "successful config retrieval",
			setupMock: func(m *MockClient) {
				config := map[string]interface{}{
					"cluster_name":      "test-cluster",
					"control_machine":   "controller1",
					"default_partition": "normal",
				}

				m.EXPECT().SlurmdbV0043GetConfigWithResponse(
					context.Background(),
				).Return(&SlurmdbV0043GetConfigResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiSlurmdbdConfigResp{
						Config: &config,
					},
				}, nil)
			},
		},
		{
			name: "HTTP error",
			setupMock: func(m *MockClient) {
				m.EXPECT().SlurmdbV0043GetConfigWithResponse(
					context.Background(),
				).Return(&SlurmdbV0043GetConfigResponse{
					HTTPResponse: &http.Response{StatusCode: 500},
				}, nil)
			},
			expectedError: "Get config failed with status 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(t)
			tt.setupMock(mockClient)
			client := &WrapperClient{
				apiClient: &ClientWithResponses{
					ClientInterface: mockClient,
				},
			}

			result, err := client.GetConfig(context.Background())

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotNil(t, result.Parameters)
			}
		})
	}
}

// TestWrapperClient_GetDiagnostics tests the GetDiagnostics operation
func TestWrapperClient_GetDiagnostics(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockClient)
		expectedError string
	}{
		{
			name: "successful diagnostics retrieval",
			setupMock: func(m *MockClient) {
				reqTime := int64(1000)
				reqTimeStart := int64(500)
				serverThreadCount := int32(10)
				agentCount := int32(5)
				agentThreadCount := int32(50)
				jobsSubmitted := int32(100)
				jobsStarted := int32(95)
				jobsCompleted := int32(90)
				jobsCanceled := int32(3)
				jobsFailed := int32(2)

				m.EXPECT().SlurmV0043GetDiagWithResponse(
					context.Background(),
				).Return(&SlurmV0043GetDiagResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiDiagResp{
						Statistics: &V0043DiagStatistics{
							ReqTime:           &reqTime,
							ReqTimeStart:      &reqTimeStart,
							ServerThreadCount: &serverThreadCount,
							AgentCount:        &agentCount,
							AgentThreadCount:  &agentThreadCount,
							JobsSubmitted:     &jobsSubmitted,
							JobsStarted:       &jobsStarted,
							JobsCompleted:     &jobsCompleted,
							JobsCanceled:      &jobsCanceled,
							JobsFailed:        &jobsFailed,
						},
					},
				}, nil)
			},
		},
		{
			name: "HTTP error",
			setupMock: func(m *MockClient) {
				m.EXPECT().SlurmV0043GetDiagWithResponse(
					context.Background(),
				).Return(&SlurmV0043GetDiagResponse{
					HTTPResponse: &http.Response{StatusCode: 500},
				}, nil)
			},
			expectedError: "Get diagnostics failed with status 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(t)
			tt.setupMock(mockClient)
			client := &WrapperClient{
				apiClient: &ClientWithResponses{
					ClientInterface: mockClient,
				},
			}

			result, err := client.GetDiagnostics(context.Background())

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotNil(t, result.Statistics)
			}
		})
	}
}

// TestWrapperClient_GetDBDiagnostics tests the GetDBDiagnostics operation
func TestWrapperClient_GetDBDiagnostics(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockClient)
		expectedError string
	}{
		{
			name: "successful DB diagnostics retrieval",
			setupMock: func(m *MockClient) {
				statistics := V0043DiagStatistics{}

				m.EXPECT().SlurmdbV0043GetDiagWithResponse(
					context.Background(),
				).Return(&SlurmdbV0043GetDiagResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiSlurmdbdDiagResp{
						Statistics: &statistics,
					},
				}, nil)
			},
		},
		{
			name: "HTTP error",
			setupMock: func(m *MockClient) {
				m.EXPECT().SlurmdbV0043GetDiagWithResponse(
					context.Background(),
				).Return(&SlurmdbV0043GetDiagResponse{
					HTTPResponse: &http.Response{StatusCode: 500},
				}, nil)
			},
			expectedError: "Get DB diagnostics failed with status 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(t)
			tt.setupMock(mockClient)
			client := &WrapperClient{
				apiClient: &ClientWithResponses{
					ClientInterface: mockClient,
				},
			}

			result, err := client.GetDBDiagnostics(context.Background())

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotNil(t, result.Statistics)
				assert.Equal(t, "database", result.Statistics["source"])
			}
		})
	}
}

// TestWrapperClient_GetInstance tests the GetInstance operation
func TestWrapperClient_GetInstance(t *testing.T) {
	tests := []struct {
		name           string
		options        *interfaces.GetInstanceOptions
		setupMock      func(*MockClient, *interfaces.GetInstanceOptions)
		expectedError  string
		validateResult func(*testing.T, *interfaces.Instance)
	}{
		{
			name: "successful instance retrieval",
			options: &interfaces.GetInstanceOptions{
				Cluster:  "cluster1",
				Instance: "instance1",
				NodeList: []string{"node1", "node2"},
			},
			setupMock: func(m *MockClient, opts *interfaces.GetInstanceOptions) {
				cluster := opts.Cluster
				instance := opts.Instance
				nodeList := strings.Join(opts.NodeList, ",")

				clusterName := "cluster1"
				extraInfo := "extra"
				instanceName := "instance1"
				nodeName := "node1"
				timeEnd := int64(1234567890)
				timeStart := int64(1234567800)
				tres := "cpu=4,mem=8G"

				m.EXPECT().SlurmdbV0043GetInstanceWithResponse(
					context.Background(),
					&SlurmdbV0043GetInstanceParams{
						Cluster:  &cluster,
						Instance: &instance,
						NodeList: &nodeList,
					},
				).Return(&SlurmdbV0043GetInstanceResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiInstancesResp{
						Instances: &[]V0043Instance{
							{
								Cluster:   &clusterName,
								Extra:     &extraInfo,
								Instance:  &instanceName,
								NodeName:  &nodeName,
								TimeEnd:   &timeEnd,
								TimeStart: &timeStart,
								Tres:      &tres,
							},
						},
					},
				}, nil)
			},
			validateResult: func(t *testing.T, result *interfaces.Instance) {
				assert.Equal(t, "cluster1", result.Cluster)
				assert.Equal(t, "extra", result.ExtraInfo)
				assert.Equal(t, "instance1", result.Instance)
				assert.Equal(t, "node1", result.NodeName)
				assert.Equal(t, int64(1234567890), result.TimeEnd)
				assert.Equal(t, int64(1234567800), result.TimeStart)
				assert.Equal(t, "cpu=4,mem=8G", result.TRES)
			},
		},
		{
			name:    "nil options",
			options: nil,
			setupMock: func(m *MockClient, opts *interfaces.GetInstanceOptions) {
				m.EXPECT().SlurmdbV0043GetInstanceWithResponse(
					context.Background(),
					&SlurmdbV0043GetInstanceParams{},
				).Return(&SlurmdbV0043GetInstanceResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiInstancesResp{
						Instances: &[]V0043Instance{},
					},
				}, nil)
			},
			expectedError: "Instance not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(t)
			tt.setupMock(mockClient, tt.options)
			client := &WrapperClient{
				apiClient: &ClientWithResponses{
					ClientInterface: mockClient,
				},
			}

			result, err := client.GetInstance(context.Background(), tt.options)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}
		})
	}
}

// TestWrapperClient_GetInstances tests the GetInstances operation
func TestWrapperClient_GetInstances(t *testing.T) {
	tests := []struct {
		name           string
		options        *interfaces.GetInstancesOptions
		setupMock      func(*MockClient, *interfaces.GetInstancesOptions)
		expectedError  string
		expectedCount  int
		validateResult func(*testing.T, *interfaces.InstanceList)
	}{
		{
			name: "successful instances retrieval",
			options: &interfaces.GetInstancesOptions{
				Clusters:  []string{"cluster1", "cluster2"},
				Instances: []string{"instance1"},
			},
			setupMock: func(m *MockClient, opts *interfaces.GetInstancesOptions) {
				clusters := strings.Join(opts.Clusters, ",")
				instances := strings.Join(opts.Instances, ",")

				clusterName := "cluster1"
				instanceName := "instance1"
				nodeName := "node1"
				timeEnd := int64(1234567890)
				timeStart := int64(1234567800)

				m.EXPECT().SlurmdbV0043GetInstancesWithResponse(
					context.Background(),
					&SlurmdbV0043GetInstancesParams{
						Cluster:  &clusters,
						Instance: &instances,
					},
				).Return(&SlurmdbV0043GetInstancesResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiInstancesResp{
						Instances: &[]V0043Instance{
							{
								Cluster:   &clusterName,
								Instance:  &instanceName,
								NodeName:  &nodeName,
								TimeEnd:   &timeEnd,
								TimeStart: &timeStart,
							},
						},
					},
				}, nil)
			},
			expectedCount: 1,
			validateResult: func(t *testing.T, result *interfaces.InstanceList) {
				require.Len(t, result.Instances, 1)
				instance := result.Instances[0]
				assert.Equal(t, "cluster1", instance.Cluster)
				assert.Equal(t, "instance1", instance.Instance)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(t)
			tt.setupMock(mockClient, tt.options)
			client := &WrapperClient{
				apiClient: &ClientWithResponses{
					ClientInterface: mockClient,
				},
			}

			result, err := client.GetInstances(context.Background(), tt.options)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
				if tt.expectedCount >= 0 {
					assert.Len(t, result.Instances, tt.expectedCount)
				}
			}
		})
	}
}

// TestWrapperClient_GetTRES tests the GetTRES operation
func TestWrapperClient_GetTRES(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockClient)
		expectedError  string
		expectedCount  int
		validateResult func(*testing.T, *interfaces.TRESList)
	}{
		{
			name: "successful TRES retrieval",
			setupMock: func(m *MockClient) {
				id := int64(1)
				typeName := "cpu"
				name := "cpu"
				count := int64(100)
				allocSecs := int64(3600)
				description := "CPU resources"

				m.EXPECT().SlurmdbV0043GetTresWithResponse(
					context.Background(),
				).Return(&SlurmdbV0043GetTresResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200: &V0043OpenapiTresResp{
						Tres: &[]V0043Tres{
							{
								Id:          &id,
								Type:        &typeName,
								Name:        &name,
								Count:       &count,
								AllocSecs:   &allocSecs,
								Description: &description,
							},
						},
					},
				}, nil)
			},
			expectedCount: 1,
			validateResult: func(t *testing.T, result *interfaces.TRESList) {
				require.Len(t, result.TRES, 1)
				tres := result.TRES[0]
				assert.Equal(t, uint64(1), tres.ID)
				assert.Equal(t, "cpu", tres.Type)
				assert.Equal(t, "cpu", tres.Name)
				assert.Equal(t, int64(100), tres.Count)
				assert.Equal(t, int64(3600), tres.AllocSecs)
				assert.Equal(t, "CPU resources", tres.Description)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(t)
			tt.setupMock(mockClient)
			client := &WrapperClient{
				apiClient: &ClientWithResponses{
					ClientInterface: mockClient,
				},
			}

			result, err := client.GetTRES(context.Background())

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
				if tt.expectedCount >= 0 {
					assert.Len(t, result.TRES, tt.expectedCount)
				}
			}
		})
	}
}

// TestWrapperClient_CreateTRES tests the CreateTRES operation
func TestWrapperClient_CreateTRES(t *testing.T) {
	tests := []struct {
		name           string
		request        *interfaces.CreateTRESRequest
		setupMock      func(*MockClient, *interfaces.CreateTRESRequest)
		expectedError  string
		validateResult func(*testing.T, *interfaces.TRES)
	}{
		{
			name: "successful TRES creation",
			request: &interfaces.CreateTRESRequest{
				Type:        "gpu",
				Name:        "gpu",
				Description: "GPU resources",
			},
			setupMock: func(m *MockClient, req *interfaces.CreateTRESRequest) {
				id := int64(2)
				typeName := req.Type
				name := req.Name
				description := req.Description

				m.EXPECT().SlurmdbV0043PostTresWithResponse(
					context.Background(),
					SlurmdbV0043PostTresJSONRequestBody{
						Tres: &[]V0043Tres{
							{
								Type:        &typeName,
								Name:        &name,
								Description: &description,
							},
						},
					},
				).Return(&SlurmdbV0043PostTresResponse{
					HTTPResponse: &http.Response{StatusCode: 201},
					JSON200: &V0043OpenapiTresResp{
						Tres: &[]V0043Tres{
							{
								Id:          &id,
								Type:        &typeName,
								Name:        &name,
								Description: &description,
							},
						},
					},
				}, nil)
			},
			validateResult: func(t *testing.T, result *interfaces.TRES) {
				assert.Equal(t, uint64(2), result.ID)
				assert.Equal(t, "gpu", result.Type)
				assert.Equal(t, "gpu", result.Name)
				assert.Equal(t, "GPU resources", result.Description)
			},
		},
		{
			name:          "nil request",
			request:       nil,
			setupMock:     func(m *MockClient, req *interfaces.CreateTRESRequest) {},
			expectedError: "Create TRES request cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(t)
			tt.setupMock(mockClient, tt.request)
			client := &WrapperClient{
				apiClient: &ClientWithResponses{
					ClientInterface: mockClient,
				},
			}

			result, err := client.CreateTRES(context.Background(), tt.request)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}
		})
	}
}

// TestWrapperClient_Reconfigure tests the Reconfigure operation
func TestWrapperClient_Reconfigure(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockClient)
		expectedError  string
		validateResult func(*testing.T, *interfaces.ReconfigureResponse)
	}{
		{
			name: "successful reconfigure",
			setupMock: func(m *MockClient) {
				m.EXPECT().SlurmV0043GetReconfigureWithResponse(
					context.Background(),
				).Return(&SlurmV0043GetReconfigureResponse{
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200:      &V0043OpenapiResp{},
				}, nil)
			},
			validateResult: func(t *testing.T, result *interfaces.ReconfigureResponse) {
				assert.Equal(t, "success", result.Status)
				assert.Equal(t, "Reconfiguration completed successfully", result.Message)
				assert.NotNil(t, result.Meta)
			},
		},
		{
			name: "HTTP error",
			setupMock: func(m *MockClient) {
				m.EXPECT().SlurmV0043GetReconfigureWithResponse(
					context.Background(),
				).Return(&SlurmV0043GetReconfigureResponse{
					HTTPResponse: &http.Response{StatusCode: 500},
				}, nil)
			},
			expectedError: "Reconfigure failed with status 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(t)
			tt.setupMock(mockClient)
			client := &WrapperClient{
				apiClient: &ClientWithResponses{
					ClientInterface: mockClient,
				},
			}

			result, err := client.Reconfigure(context.Background())

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}
		})
	}
}

// TestHelperFunctions tests the helper functions
func TestHelperFunctions(t *testing.T) {
	t.Run("getStringValue", func(t *testing.T) {
		assert.Equal(t, "", getStringValue(nil))
		str := "test"
		assert.Equal(t, "test", getStringValue(&str))
	})

	t.Run("getIntValue", func(t *testing.T) {
		assert.Equal(t, 0, getIntValue(nil))
		val := int32(42)
		assert.Equal(t, 42, getIntValue(&val))
	})

	t.Run("getInt64Value", func(t *testing.T) {
		assert.Equal(t, int64(0), getInt64Value(nil))
		val := int64(42)
		assert.Equal(t, int64(42), getInt64Value(&val))
	})

	t.Run("getUint64Value", func(t *testing.T) {
		assert.Equal(t, uint64(0), getUint64Value(nil))
		val := int64(42)
		assert.Equal(t, uint64(42), getUint64Value(&val))
	})

	t.Run("getFloatValue", func(t *testing.T) {
		assert.Equal(t, 0.0, getFloatValue(nil))
		val := 42.5
		assert.Equal(t, 42.5, getFloatValue(&val))
	})

	t.Run("getBoolValue", func(t *testing.T) {
		assert.False(t, getBoolValue(nil))
		val := true
		assert.True(t, getBoolValue(&val))
	})

	t.Run("parseInt", func(t *testing.T) {
		assert.Equal(t, 0, parseInt(""))
		assert.Equal(t, 42, parseInt("42"))
		assert.Equal(t, 0, parseInt("invalid"))
	})

	t.Run("parseFloat", func(t *testing.T) {
		assert.Equal(t, 0.0, parseFloat(""))
		assert.Equal(t, 42.5, parseFloat("42.5"))
		assert.Equal(t, 0.0, parseFloat("invalid"))
	})
}

// TestStandaloneOperationsIntegration tests integration scenarios
func TestStandaloneOperationsIntegration(t *testing.T) {
	t.Run("all operations with uninitialized client", func(t *testing.T) {
		client := &WrapperClient{} // Uninitialized client

		ctx := context.Background()

		// Test all operations fail with uninitialized client
		_, err := client.GetLicenses(ctx)
		assert.Contains(t, err.Error(), "API client not initialized")

		_, err = client.GetShares(ctx, nil)
		assert.Contains(t, err.Error(), "API client not initialized")

		_, err = client.GetConfig(ctx)
		assert.Contains(t, err.Error(), "API client not initialized")

		_, err = client.GetDiagnostics(ctx)
		assert.Contains(t, err.Error(), "API client not initialized")

		_, err = client.GetDBDiagnostics(ctx)
		assert.Contains(t, err.Error(), "API client not initialized")

		_, err = client.GetInstance(ctx, nil)
		assert.Contains(t, err.Error(), "API client not initialized")

		_, err = client.GetInstances(ctx, nil)
		assert.Contains(t, err.Error(), "API client not initialized")

		_, err = client.GetTRES(ctx)
		assert.Contains(t, err.Error(), "API client not initialized")

		_, err = client.CreateTRES(ctx, nil)
		assert.Contains(t, err.Error(), "API client not initialized")

		_, err = client.Reconfigure(ctx)
		assert.Contains(t, err.Error(), "API client not initialized")
	})

	t.Run("context cancellation", func(t *testing.T) {
		mockClient := NewMockClient(t)
		client := &WrapperClient{
			apiClient: &ClientWithResponses{
				ClientInterface: mockClient,
			},
		}

		// Create a cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Setup mock to return context error
		mockClient.EXPECT().SlurmV0043GetLicensesWithResponse(ctx).
			Return(nil, context.Canceled)

		_, err := client.GetLicenses(ctx)
		assert.Error(t, err)
		var clientErr *errors.ClientError
		assert.True(t, errors.As(err, &clientErr))
	})
}

// BenchmarkStandaloneOperations benchmarks the standalone operations
func BenchmarkStandaloneOperations(b *testing.B) {
	mockClient := NewMockClient(b)
	client := &WrapperClient{
		apiClient: &ClientWithResponses{
			ClientInterface: mockClient,
		},
	}

	// Setup basic mocks for benchmarking
	licenseName := "test-license"
	total := int32(100)
	used := int32(50)
	available := int32(50)

	mockClient.EXPECT().SlurmV0043GetLicensesWithResponse(
		context.Background(),
	).Return(&SlurmV0043GetLicensesResponse{
		HTTPResponse: &http.Response{StatusCode: 200},
		JSON200: &V0043OpenapiLicensesResp{
			Licenses: &[]V0043License{
				{
					LicenseName: &licenseName,
					Total:       &total,
					Used:        &used,
					Available:   &available,
				},
			},
		},
	}, nil).Maybe()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetLicenses(context.Background())
	}
}