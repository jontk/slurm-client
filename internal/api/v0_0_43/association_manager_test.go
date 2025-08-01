// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"net/http"
	"testing"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

func TestAssociationManagerImpl_List(t *testing.T) {
	tests := []struct {
		name    string
		client  *WrapperClient
		opts    *interfaces.ListAssociationsOptions
		wantErr bool
		errType string
	}{
		{
			name:    "nil client",
			client:  nil,
			opts:    nil,
			wantErr: true,
			errType: "client error",
		},
		{
			name: "nil api client",
			client: &WrapperClient{
				apiClient: nil,
			},
			opts:    nil,
			wantErr: true,
			errType: "client error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    &interfaces.ListAssociationsOptions{},
			wantErr: true,
			errType: "client error",
		},
		{
			name: "with all filter options",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts: &interfaces.ListAssociationsOptions{
				Users:           []string{"user1", "user2"},
				Accounts:        []string{"account1", "account2"},
				Clusters:        []string{"cluster1"},
				Partitions:      []string{"partition1"},
				ParentAccounts:  []string{"parent1"},
				QoS:             []string{"normal", "high"},
				WithDeleted:     true,
				WithUsage:       true,
				WithSubAccounts: true,
				OnlyDefaults:    true,
				Offset:          10,
				Limit:           50,
			},
			wantErr: true,
			errType: "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AssociationManagerImpl{
				client: tt.client,
			}

			result, err := a.List(context.Background(), tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("List() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("List() expected client error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("List() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("List() expected nil result on error")
			}
		})
	}
}

func TestAssociationManagerImpl_Get(t *testing.T) {
	tests := []struct {
		name    string
		client  *WrapperClient
		opts    *interfaces.GetAssociationOptions
		wantErr bool
		errType string
	}{
		{
			name:    "nil client",
			client:  nil,
			opts:    &interfaces.GetAssociationOptions{User: "user1", Account: "account1"},
			wantErr: true,
			errType: "client error",
		},
		{
			name: "nil opts",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    nil,
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "empty user",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    &interfaces.GetAssociationOptions{User: "", Account: "account1"},
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "empty account",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    &interfaces.GetAssociationOptions{User: "user1", Account: ""},
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    &interfaces.GetAssociationOptions{User: "user1", Account: "account1"},
			wantErr: true,
			errType: "client error",
		},
		{
			name: "with all options",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts: &interfaces.GetAssociationOptions{
				User:      "user1",
				Account:   "account1",
				Cluster:   "cluster1",
				Partition: "partition1",
				WithUsage: true,
				WithTRES:  true,
			},
			wantErr: true,
			errType: "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AssociationManagerImpl{
				client: tt.client,
			}

			result, err := a.Get(context.Background(), tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Get() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("Get() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("Get() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Get() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("Get() expected nil result on error")
			}
		})
	}
}

func TestAssociationManagerImpl_Create(t *testing.T) {
	tests := []struct {
		name         string
		client       *WrapperClient
		associations []*interfaces.AssociationCreate
		wantErr      bool
		errType      string
	}{
		{
			name:         "nil client",
			client:       nil,
			associations: []*interfaces.AssociationCreate{{User: "user1", Account: "account1"}},
			wantErr:      true,
			errType:      "client error",
		},
		{
			name: "nil associations",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			associations: nil,
			wantErr:      true,
			errType:      "validation error",
		},
		{
			name: "empty associations",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			associations: []*interfaces.AssociationCreate{},
			wantErr:      true,
			errType:      "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			associations: []*interfaces.AssociationCreate{
				{User: "user1", Account: "account1"},
			},
			wantErr: true,
			errType: "client error",
		},
		{
			name: "multiple associations with all fields",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			associations: []*interfaces.AssociationCreate{
				{
					User:            "user1",
					Account:         "account1",
					Cluster:         "cluster1",
					Partition:       "partition1",
					ParentAccount:   "parent1",
					IsDefault:       true,
					Comment:         "Test association",
					SharesRaw:       intPtr(100),
					Priority:        uint32Ptr(10),
					MaxJobs:         intPtr(50),
					MaxJobsAccrue:   intPtr(100),
					MaxSubmitJobs:   intPtr(200),
					MaxWallDuration: intPtr(3600),
					GrpJobs:         intPtr(500),
					MaxTRESPerJob:   map[string]string{"cpu": "100", "mem": "1024M"},
					DefaultQoS:      "normal",
					QoSList:         []string{"normal", "high"},
					Flags:           []string{"flag1", "flag2"},
				},
				{
					User:    "user2",
					Account: "account2",
				},
			},
			wantErr: true,
			errType: "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AssociationManagerImpl{
				client: tt.client,
			}

			result, err := a.Create(context.Background(), tt.associations)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Create() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("Create() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("Create() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Create() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("Create() expected nil result on error")
			}
		})
	}
}

func TestAssociationManagerImpl_Update(t *testing.T) {
	tests := []struct {
		name         string
		client       *WrapperClient
		associations []*interfaces.AssociationUpdate
		wantErr      bool
		errType      string
	}{
		{
			name:         "nil client",
			client:       nil,
			associations: []*interfaces.AssociationUpdate{{User: "user1", Account: "account1"}},
			wantErr:      true,
			errType:      "client error",
		},
		{
			name: "nil associations",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			associations: nil,
			wantErr:      true,
			errType:      "validation error",
		},
		{
			name: "empty associations",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			associations: []*interfaces.AssociationUpdate{},
			wantErr:      true,
			errType:      "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			associations: []*interfaces.AssociationUpdate{
				{
					User:      "user1",
					Account:   "account1",
					IsDefault: boolPtr(true),
					Comment:   stringPtr("Updated comment"),
				},
			},
			wantErr: true,
			errType: "client error",
		},
		{
			name: "update with all fields",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			associations: []*interfaces.AssociationUpdate{
				{
					User:            "user1",
					Account:         "account1",
					Cluster:         "cluster1",
					Partition:       "partition1",
					IsDefault:       boolPtr(false),
					Comment:         stringPtr("Updated association"),
					SharesRaw:       intPtr(200),
					Priority:        uint32Ptr(20),
					MaxJobs:         intPtr(100),
					MaxJobsAccrue:   intPtr(200),
					MaxSubmitJobs:   intPtr(400),
					MaxWallDuration: intPtr(7200),
					GrpJobs:         intPtr(1000),
					MaxTRESPerJob:   map[string]string{"cpu": "200", "mem": "2048M"},
					DefaultQoS:      stringPtr("high"),
					QoSList:         []string{"normal", "high", "critical"},
					Flags:           []string{"updated_flag"},
				},
			},
			wantErr: true,
			errType: "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AssociationManagerImpl{
				client: tt.client,
			}

			err := a.Update(context.Background(), tt.associations)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Update() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("Update() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("Update() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Update() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAssociationManagerImpl_Delete(t *testing.T) {
	tests := []struct {
		name    string
		client  *WrapperClient
		opts    *interfaces.DeleteAssociationOptions
		wantErr bool
		errType string
	}{
		{
			name:    "nil client",
			client:  nil,
			opts:    &interfaces.DeleteAssociationOptions{User: "user1", Account: "account1"},
			wantErr: true,
			errType: "client error",
		},
		{
			name: "nil opts",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    nil,
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "empty user",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    &interfaces.DeleteAssociationOptions{User: "", Account: "account1"},
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "empty account",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    &interfaces.DeleteAssociationOptions{User: "user1", Account: ""},
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    &interfaces.DeleteAssociationOptions{User: "user1", Account: "account1"},
			wantErr: true,
			errType: "client error",
		},
		{
			name: "with all options",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts: &interfaces.DeleteAssociationOptions{
				User:      "user1",
				Account:   "account1",
				Cluster:   "cluster1",
				Partition: "partition1",
				Force:     true,
			},
			wantErr: true,
			errType: "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AssociationManagerImpl{
				client: tt.client,
			}

			err := a.Delete(context.Background(), tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Delete() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("Delete() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("Delete() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Delete() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAssociationManagerImpl_BulkDelete(t *testing.T) {
	tests := []struct {
		name    string
		client  *WrapperClient
		opts    *interfaces.BulkDeleteOptions
		wantErr bool
		errType string
	}{
		{
			name:    "nil client",
			client:  nil,
			opts:    &interfaces.BulkDeleteOptions{Users: []string{"user1"}},
			wantErr: true,
			errType: "client error",
		},
		{
			name: "nil opts",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    nil,
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts:    &interfaces.BulkDeleteOptions{Users: []string{"user1", "user2"}},
			wantErr: true,
			errType: "client error",
		},
		{
			name: "with all options",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			opts: &interfaces.BulkDeleteOptions{
				Users:      []string{"user1", "user2"},
				Accounts:   []string{"account1", "account2"},
				Clusters:   []string{"cluster1"},
				Partitions: []string{"partition1"},
				OnlyIfIdle: true,
				Force:      true,
			},
			wantErr: true,
			errType: "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AssociationManagerImpl{
				client: tt.client,
			}

			result, err := a.BulkDelete(context.Background(), tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("BulkDelete() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("BulkDelete() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("BulkDelete() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("BulkDelete() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("BulkDelete() expected nil result on error")
			}
		})
	}
}

func TestAssociationManagerImpl_GetUserAssociations(t *testing.T) {
	tests := []struct {
		name     string
		client   *WrapperClient
		userName string
		wantErr  bool
		errType  string
	}{
		{
			name:     "nil client",
			client:   nil,
			userName: "user1",
			wantErr:  true,
			errType:  "client error",
		},
		{
			name: "empty user name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "",
			wantErr:  true,
			errType:  "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			userName: "user1",
			wantErr:  true,
			errType:  "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AssociationManagerImpl{
				client: tt.client,
			}

			result, err := a.GetUserAssociations(context.Background(), tt.userName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserAssociations() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetUserAssociations() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetUserAssociations() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetUserAssociations() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetUserAssociations() expected nil result on error")
			}
		})
	}
}

func TestAssociationManagerImpl_GetAccountAssociations(t *testing.T) {
	tests := []struct {
		name        string
		client      *WrapperClient
		accountName string
		wantErr     bool
		errType     string
	}{
		{
			name:        "nil client",
			client:      nil,
			accountName: "account1",
			wantErr:     true,
			errType:     "client error",
		},
		{
			name: "empty account name",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "",
			wantErr:     true,
			errType:     "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			accountName: "account1",
			wantErr:     true,
			errType:     "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AssociationManagerImpl{
				client: tt.client,
			}

			result, err := a.GetAccountAssociations(context.Background(), tt.accountName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetAccountAssociations() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("GetAccountAssociations() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("GetAccountAssociations() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("GetAccountAssociations() unexpected error: %v", err)
				}
			}

			if result != nil && tt.wantErr {
				t.Errorf("GetAccountAssociations() expected nil result on error")
			}
		})
	}
}

func TestAssociationManagerImpl_ValidateAssociation(t *testing.T) {
	tests := []struct {
		name    string
		client  *WrapperClient
		user    string
		account string
		cluster string
		want    bool
		wantErr bool
		errType string
	}{
		{
			name:    "nil client",
			client:  nil,
			user:    "user1",
			account: "account1",
			cluster: "cluster1",
			want:    false,
			wantErr: true,
			errType: "client error",
		},
		{
			name: "empty user",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			user:    "",
			account: "account1",
			cluster: "cluster1",
			want:    false,
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "empty account",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			user:    "user1",
			account: "",
			cluster: "cluster1",
			want:    false,
			wantErr: true,
			errType: "validation error",
		},
		{
			name: "valid input - nil api response",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			user:    "user1",
			account: "account1",
			cluster: "cluster1",
			want:    false,
			wantErr: true,
			errType: "client error",
		},
		{
			name: "without cluster",
			client: &WrapperClient{
				apiClient: &ClientWithResponses{},
			},
			user:    "user1",
			account: "account1",
			cluster: "",
			want:    false,
			wantErr: true,
			errType: "client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AssociationManagerImpl{
				client: tt.client,
			}

			got, err := a.ValidateAssociation(context.Background(), tt.user, tt.account, tt.cluster)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateAssociation() expected error but got none")
					return
				}

				// Check error type
				switch tt.errType {
				case "client error":
					if !errors.IsClientError(err) {
						t.Errorf("ValidateAssociation() expected client error, got %T: %v", err, err)
					}
				case "validation error":
					if !errors.IsValidationError(err) {
						t.Errorf("ValidateAssociation() expected validation error, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateAssociation() unexpected error: %v", err)
				}
			}

			if got != tt.want {
				t.Errorf("ValidateAssociation() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test handleAPIError function
func TestHandleAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   *V0043OpenapiAssocsResp
		wantErr    bool
	}{
		{
			name:       "400 Bad Request",
			statusCode: http.StatusBadRequest,
			response:   nil,
			wantErr:    true,
		},
		{
			name:       "401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   nil,
			wantErr:    true,
		},
		{
			name:       "403 Forbidden",
			statusCode: http.StatusForbidden,
			response:   nil,
			wantErr:    true,
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			response:   nil,
			wantErr:    true,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			response:   nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleAPIError(tt.statusCode, tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleAPIError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper functions for test pointers
// Note: intPtr, stringPtr, and boolPtr are already defined in other test files

func uint32Ptr(i uint32) *uint32 {
	return &i
}
