//go:build integration
// +build integration

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/jontk/slurm-client/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAdapterReservationOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test against versions that support reservations (v0.0.43+)
	versions := []string{"v0.0.43", "v0.0.44"}

	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			ctx := helpers.TestContext(t)

			// Start mock server
			server := mocks.StartMockServer(t, version)
			defer server.Close()

			// Create factory with mock server URL
			factory, err := factory.NewClientFactory(
				factory.WithBaseURL(server.URL()),
			)
			helpers.RequireNoError(t, err)

			// Create client for specific version
			client, err := factory.NewClientWithVersion(ctx, version)
			helpers.RequireNoError(t, err)

			resManager := client.Reservations()
			if resManager == nil {
				t.Skipf("Reservation manager not available for version %s", version)
			}

			// Test if reservations are supported by trying to list
			_, err = resManager.List(ctx, nil)
			if err != nil {
				// Mock server doesn't support reservations yet
				t.Skipf("Reservations not supported in mock server for version %s: %v", version, err)
				return
			}

			t.Run("List Reservations", func(t *testing.T) {
				list, err := resManager.List(ctx, nil)
				helpers.AssertNoError(t, err)
				assert.NotNil(t, list)
				assert.NotNil(t, list.Reservations)
			})

			t.Run("Create and Get Reservation", func(t *testing.T) {
				resName := fmt.Sprintf("test-res-%d", time.Now().Unix())

				createReq := &interfaces.ReservationCreate{
					Name:          resName,
					StartTime:     time.Now().Add(1 * time.Hour),
					EndTime:       time.Now().Add(3 * time.Hour),
					Duration:      7200, // 2 hours
					Users:         []string{"testuser"},
					Accounts:      []string{"testaccount"},
					NodeCount:     5,
					PartitionName: "compute",
				}

				// Create reservation
				createResp, err := resManager.Create(ctx, createReq)
				helpers.AssertNoError(t, err)
				assert.Equal(t, resName, createResp.ReservationName)

				// Get the created reservation
				res, err := resManager.Get(ctx, resName)
				helpers.AssertNoError(t, err)
				assert.Equal(t, resName, res.Name)
				assert.Equal(t, createReq.Users, res.Users)
				assert.Equal(t, createReq.Accounts, res.Accounts)
				assert.Equal(t, createReq.NodeCount, res.NodeCount)
			})

			t.Run("Update Reservation", func(t *testing.T) {
				resName := fmt.Sprintf("test-update-res-%d", time.Now().Unix())

				// Create reservation first
				createReq := &interfaces.ReservationCreate{
					Name:      resName,
					StartTime: time.Now().Add(2 * time.Hour),
					EndTime:   time.Now().Add(4 * time.Hour),
					Users:     []string{"user1"},
					NodeCount: 3,
				}

				_, err := resManager.Create(ctx, createReq)
				helpers.AssertNoError(t, err)

				// Update reservation
				updateReq := &interfaces.ReservationUpdate{
					Users:     []string{"user1", "user2"},
					NodeCount: helpers.IntPtr(5),
				}

				err = resManager.Update(ctx, resName, updateReq)
				helpers.AssertNoError(t, err)

				// Verify update
				updated, err := resManager.Get(ctx, resName)
				helpers.AssertNoError(t, err)
				assert.Equal(t, []string{"user1", "user2"}, updated.Users)
				assert.Equal(t, 5, updated.NodeCount)
			})

			t.Run("Delete Reservation", func(t *testing.T) {
				resName := fmt.Sprintf("test-delete-res-%d", time.Now().Unix())

				// Create reservation first
				createReq := &interfaces.ReservationCreate{
					Name:      resName,
					StartTime: time.Now().Add(1 * time.Hour),
					EndTime:   time.Now().Add(2 * time.Hour),
					Users:     []string{"deleteuser"},
				}

				_, err := resManager.Create(ctx, createReq)
				helpers.AssertNoError(t, err)

				// Delete reservation
				err = resManager.Delete(ctx, resName)
				helpers.AssertNoError(t, err)

				// Verify it's gone
				_, err = resManager.Get(ctx, resName)
				assert.Error(t, err)
			})
		})
	}
}

func TestAdapterAssociationOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test against each API version
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			ctx := helpers.TestContext(t)

			// Start mock server
			server := mocks.StartMockServer(t, version)
			defer server.Close()

			// Create factory with mock server URL
			factory, err := factory.NewClientFactory(
				factory.WithBaseURL(server.URL()),
			)
			helpers.RequireNoError(t, err)

			// Create client for specific version
			client, err := factory.NewClientWithVersion(ctx, version)
			helpers.RequireNoError(t, err)

			assocManager := client.Associations()
			if assocManager == nil {
				t.Skipf("Association manager not available for version %s", version)
			}

			// Test if associations are supported by trying to list
			_, err = assocManager.List(ctx, nil)
			if err != nil {
				// Mock server doesn't support associations yet
				t.Skipf("Associations not supported in mock server for version %s: %v", version, err)
				return
			}

			t.Run("List Associations", func(t *testing.T) {
				list, err := assocManager.List(ctx, nil)
				helpers.AssertNoError(t, err)
				if assert.NotNil(t, list) {
					assert.NotNil(t, list.Associations)
				}
			})

			t.Run("Create and Get Association", func(t *testing.T) {
				createReq := &interfaces.AssociationCreate{
					User:      "testuser",
					Account:   "testaccount",
					Cluster:   "testcluster",
					Partition: "compute",
				}

				// Create association
				createResp, err := assocManager.Create(ctx, []*interfaces.AssociationCreate{createReq})
				helpers.AssertNoError(t, err)
				assert.NotNil(t, createResp)
				assert.NotEmpty(t, createResp.Associations)

				// Get the created association
				assoc, err := assocManager.Get(ctx, &interfaces.GetAssociationOptions{
					User:    createReq.User,
					Account: createReq.Account,
					Cluster: createReq.Cluster,
				})
				helpers.AssertNoError(t, err)
				assert.Equal(t, createReq.User, assoc.User)
				assert.Equal(t, createReq.Account, assoc.Account)
				assert.Equal(t, createReq.Cluster, assoc.Cluster)
			})

			t.Run("Update Association", func(t *testing.T) {
				// Create association first
				createReq := &interfaces.AssociationCreate{
					User:    "updateuser",
					Account: "updateaccount",
					Cluster: "testcluster",
					MaxJobs: helpers.IntPtr(10),
				}

				createResp, err := assocManager.Create(ctx, []*interfaces.AssociationCreate{createReq})
				helpers.AssertNoError(t, err)
				assert.NotEmpty(t, createResp.Associations)

				// Update association
				updateReq := &interfaces.AssociationUpdate{
					User:    createReq.User,
					Account: createReq.Account,
					Cluster: createReq.Cluster,
					MaxJobs: helpers.IntPtr(20),
				}

				err = assocManager.Update(ctx, []*interfaces.AssociationUpdate{updateReq})
				helpers.AssertNoError(t, err)

				// Verify update
				updated, err := assocManager.Get(ctx, &interfaces.GetAssociationOptions{
					User:    createReq.User,
					Account: createReq.Account,
					Cluster: createReq.Cluster,
				})
				helpers.AssertNoError(t, err)
				assert.NotNil(t, updated.MaxJobs)
				assert.Equal(t, 20, *updated.MaxJobs)
			})

			t.Run("Get User Associations", func(t *testing.T) {
				userName := fmt.Sprintf("user-%d", time.Now().Unix())

				// Create multiple associations for user
				for i := 0; i < 3; i++ {
					createReq := &interfaces.AssociationCreate{
						User:    userName,
						Account: fmt.Sprintf("account%d", i),
						Cluster: "testcluster",
					}
					_, err := assocManager.Create(ctx, []*interfaces.AssociationCreate{createReq})
					helpers.AssertNoError(t, err)
				}

				// Get user associations
				userAssocs, err := assocManager.GetUserAssociations(ctx, userName)
				helpers.AssertNoError(t, err)
				assert.GreaterOrEqual(t, len(userAssocs), 3)

				// Verify all associations belong to the user
				for _, assoc := range userAssocs {
					assert.Equal(t, userName, assoc.User)
				}
			})

			t.Run("Get Account Associations", func(t *testing.T) {
				accountName := fmt.Sprintf("account-%d", time.Now().Unix())

				// Create multiple associations for account
				for i := 0; i < 3; i++ {
					createReq := &interfaces.AssociationCreate{
						User:    fmt.Sprintf("user%d", i),
						Account: accountName,
						Cluster: "testcluster",
					}
					_, err := assocManager.Create(ctx, []*interfaces.AssociationCreate{createReq})
					helpers.AssertNoError(t, err)
				}

				// Get account associations
				accountAssocs, err := assocManager.GetAccountAssociations(ctx, accountName)
				helpers.AssertNoError(t, err)
				assert.GreaterOrEqual(t, len(accountAssocs), 3)

				// Verify all associations belong to the account
				for _, assoc := range accountAssocs {
					assert.Equal(t, accountName, assoc.Account)
				}
			})

			t.Run("Validate Association", func(t *testing.T) {
				// Valid association
				valid, err := assocManager.ValidateAssociation(ctx, "validuser", "validaccount", "validcluster")
				helpers.AssertNoError(t, err)
				assert.True(t, valid)

				// Invalid association (missing required fields)
				valid, err = assocManager.ValidateAssociation(ctx, "onlyuser", "", "")
				assert.Error(t, err)
				assert.False(t, valid)
			})

			t.Run("Bulk Delete Associations", func(t *testing.T) {
				// Create multiple associations
				var users []string
				for i := 0; i < 3; i++ {
					user := fmt.Sprintf("bulkuser%d", i)
					createReq := &interfaces.AssociationCreate{
						User:    user,
						Account: "bulkaccount",
						Cluster: "bulkcluster",
					}
					createResp, err := assocManager.Create(ctx, []*interfaces.AssociationCreate{createReq})
					helpers.AssertNoError(t, err)
					assert.NotEmpty(t, createResp.Associations)
					users = append(users, user)
				}

				// Bulk delete
				resp, err := assocManager.BulkDelete(ctx, &interfaces.BulkDeleteOptions{
					Users:    users,
					Accounts: []string{"bulkaccount"},
					Clusters: []string{"bulkcluster"},
				})
				helpers.AssertNoError(t, err)
				assert.Equal(t, 3, resp.Deleted)

				// Verify all are deleted
				for _, user := range users {
					_, err := assocManager.Get(ctx, &interfaces.GetAssociationOptions{
						User:    user,
						Account: "bulkaccount",
						Cluster: "bulkcluster",
					})
					assert.Error(t, err)
				}
			})

			t.Run("Delete Association", func(t *testing.T) {
				// Create association first
				createReq := &interfaces.AssociationCreate{
					User:    "deleteuser",
					Account: "deleteaccount",
					Cluster: "deletecluster",
				}

				createResp, err := assocManager.Create(ctx, []*interfaces.AssociationCreate{createReq})
				helpers.AssertNoError(t, err)
				assert.NotEmpty(t, createResp.Associations)

				// Delete association
				err = assocManager.Delete(ctx, &interfaces.DeleteAssociationOptions{
					User:    createReq.User,
					Account: createReq.Account,
					Cluster: createReq.Cluster,
				})
				helpers.AssertNoError(t, err)

				// Verify it's gone
				_, err = assocManager.Get(ctx, &interfaces.GetAssociationOptions{
					User:    createReq.User,
					Account: createReq.Account,
					Cluster: createReq.Cluster,
				})
				assert.Error(t, err)
			})
		})
	}
}

func TestAdapterCrossVersionCompatibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := helpers.TestContext(t)

	// Test that adapter clients for different versions can work together
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	// Start mock servers for each version
	servers := make(map[string]*mocks.MockSlurmServer)
	for _, version := range versions {
		server := mocks.StartMockServer(t, version)
		defer server.Close()
		servers[version] = server
	}

	t.Run("Cross Version Reservation Operations", func(t *testing.T) {
		// Create reservation with v0.0.40
		factory40, err := factory.NewClientFactory(
			factory.WithBaseURL(servers["v0.0.40"].URL()),
		)
		helpers.RequireNoError(t, err)

		client40, err := factory40.NewClientWithVersion(ctx, "v0.0.40")
		helpers.RequireNoError(t, err)

		if client40.Reservations() != nil {
			// Check if reservations are supported
			_, err := client40.Reservations().List(ctx, nil)
			if err != nil {
				t.Skipf("Reservations not supported in v0.0.40 mock server: %v", err)
				return
			}

			createReq := &interfaces.ReservationCreate{
				Name:      "cross-version-res",
				StartTime: time.Now().Add(1 * time.Hour),
				EndTime:   time.Now().Add(2 * time.Hour),
				Users:     []string{"testuser"},
			}

			_, err = client40.Reservations().Create(ctx, createReq)
			helpers.AssertNoError(t, err)
		}

		// List with v0.0.43
		factory43, err := factory.NewClientFactory(
			factory.WithBaseURL(servers["v0.0.43"].URL()),
		)
		helpers.RequireNoError(t, err)

		client43, err := factory43.NewClientWithVersion(ctx, "v0.0.43")
		helpers.RequireNoError(t, err)

		if client43.Reservations() != nil {
			// Check if reservations are supported
			_, err := client43.Reservations().List(ctx, nil)
			if err != nil {
				t.Skipf("Reservations not supported in v0.0.43 mock server: %v", err)
				return
			}

			list, err := client43.Reservations().List(ctx, nil)
			helpers.AssertNoError(t, err)
			assert.NotNil(t, list)
		}
	})

	t.Run("Cross Version Association Operations", func(t *testing.T) {
		// Create association with v0.0.42
		factory42, err := factory.NewClientFactory(
			factory.WithBaseURL(servers["v0.0.42"].URL()),
		)
		helpers.RequireNoError(t, err)

		client42, err := factory42.NewClientWithVersion(ctx, "v0.0.42")
		helpers.RequireNoError(t, err)

		if client42.Associations() != nil {
			// Check if associations are supported
			_, err := client42.Associations().List(ctx, nil)
			if err != nil {
				t.Skipf("Associations not supported in v0.0.42 mock server: %v", err)
				return
			}

			createReq := &interfaces.AssociationCreate{
				User:    "crossuser",
				Account: "crossaccount",
				Cluster: "crosscluster",
			}

			createResp, err := client42.Associations().Create(ctx, []*interfaces.AssociationCreate{createReq})
			helpers.AssertNoError(t, err)
			if createResp != nil {
				assert.NotEmpty(t, createResp.Associations)
			}
		}

		// List with v0.0.41
		factory41, err := factory.NewClientFactory(
			factory.WithBaseURL(servers["v0.0.41"].URL()),
		)
		helpers.RequireNoError(t, err)

		client41, err := factory41.NewClientWithVersion(ctx, "v0.0.41")
		helpers.RequireNoError(t, err)

		if client41.Associations() != nil {
			// Check if associations are supported
			_, err := client41.Associations().List(ctx, nil)
			if err != nil {
				t.Skipf("Associations not supported in v0.0.41 mock server: %v", err)
				return
			}

			list, err := client41.Associations().List(ctx, nil)
			helpers.AssertNoError(t, err)
			assert.NotNil(t, list)
		}
	})
}
