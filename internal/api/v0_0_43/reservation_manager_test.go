package v0_0_43

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReservationManagerImpl_List(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/slurm/v0.0.43/reservations", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			// Return mock response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"meta": {
					"slurm": {
						"version": {
							"major": 25,
							"micro": 5,
							"minor": 0
						},
						"release": "25.05.0"
					}
				},
				"reservations": [
					{
						"name": "test-reservation",
						"state": ["ACTIVE"],
						"node": "node001,node002",
						"node_count": 2,
						"core_count": 64,
						"start_time": {
							"set": true,
							"infinite": false,
							"number": 1735257600
						},
						"end_time": {
							"set": true,
							"infinite": false,
							"number": 1735344000
						},
						"duration": {
							"set": true,
							"infinite": false,
							"number": 86400
						},
						"users": "user1,user2",
						"accounts": "account1",
						"features": "gpu,highmem",
						"partition": "partition1",
						"flags": ["MAINT", "IGNORE_JOBS"],
						"licenses": "software:10",
						"groups": "group1"
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test List
		mgr := NewReservationManagerImpl(client)
		result, err := mgr.List(context.Background(), nil)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.Total)
		assert.Len(t, result.Reservations, 1)

		// Verify reservation data
		res := result.Reservations[0]
		assert.Equal(t, "test-reservation", res.Name)
		assert.Equal(t, "ACTIVE", res.State)
		assert.Equal(t, []string{"node001", "node002"}, res.Nodes)
		assert.Equal(t, 2, res.NodeCount)
		assert.Equal(t, 64, res.CoreCount)
		assert.Equal(t, []string{"user1", "user2"}, res.Users)
		assert.Equal(t, []string{"account1"}, res.Accounts)
		assert.Equal(t, []string{"gpu", "highmem"}, res.Features)
		assert.Equal(t, []string{"partition1"}, res.Partitions)
		assert.Equal(t, []string{"MAINT", "IGNORE_JOBS"}, res.Flags)
		assert.Equal(t, "software:10", res.Licenses)
		assert.Equal(t, []string{"group1"}, res.Groups)
		assert.Equal(t, time.Unix(1735257600, 0), res.StartTime)
		assert.Equal(t, time.Unix(1735344000, 0), res.EndTime)
		assert.Equal(t, 86400, res.Duration)
	})

	t.Run("list with filtering", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"reservations": [
					{
						"name": "res1",
						"state": ["ACTIVE"],
						"users": "user1"
					},
					{
						"name": "res2",
						"state": ["INACTIVE"],
						"users": "user2"
					},
					{
						"name": "res3",
						"state": ["ACTIVE"],
						"users": "user1,user3"
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test with filtering
		mgr := NewReservationManagerImpl(client)
		opts := &interfaces.ListReservationsOptions{
			Users: []string{"user1"},
			State: "ACTIVE",
		}
		result, err := mgr.List(context.Background(), opts)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, result.Total) // Only res1 and res3 match
		assert.Equal(t, "res1", result.Reservations[0].Name)
		assert.Equal(t, "res3", result.Reservations[1].Name)
	})

	t.Run("error handling", func(t *testing.T) {
		// Create test server that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{
				"errors": [
					{
						"error": "Internal Server Error",
						"error_number": 500,
						"source": "reservation list",
						"description": "Failed to retrieve reservations"
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test error
		mgr := NewReservationManagerImpl(client)
		result, err := mgr.List(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestReservationManagerImpl_Get(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/slurm/v0.0.43/reservation/test-reservation", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			// Return mock response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"reservations": [
					{
						"name": "test-reservation",
						"state": ["ACTIVE"],
						"node": "node001,node002",
						"node_count": 2
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test Get
		mgr := NewReservationManagerImpl(client)
		result, err := mgr.Get(context.Background(), "test-reservation")
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "test-reservation", result.Name)
		assert.Equal(t, "ACTIVE", result.State)
	})

	t.Run("not found", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"reservations": []
			}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test not found
		mgr := NewReservationManagerImpl(client)
		result, err := mgr.Get(context.Background(), "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestReservationManagerImpl_Create(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/slurm/v0.0.43/reservation", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			// Return success response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"reservations": [
					{
						"name": "new-reservation"
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test Create
		mgr := NewReservationManagerImpl(client)
		create := &interfaces.ReservationCreate{
			Name:      "new-reservation",
			StartTime: time.Now().Add(1 * time.Hour),
			Duration:  3600, // 1 hour
			Nodes:     []string{"node001"},
			Users:     []string{"user1"},
		}
		result, err := mgr.Create(context.Background(), create)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Created)
		assert.Equal(t, "new-reservation", result.Name)
	})

	t.Run("validation errors", func(t *testing.T) {
		client, err := createTestClient("http://localhost")
		require.NoError(t, err)
		mgr := NewReservationManagerImpl(client)

		// Test missing name
		create := &interfaces.ReservationCreate{
			StartTime: time.Now(),
			Duration:  3600,
		}
		result, err := mgr.Create(context.Background(), create)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name is required")

		// Test missing start time
		create = &interfaces.ReservationCreate{
			Name: "test",
		}
		result, err = mgr.Create(context.Background(), create)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "start time is required")

		// Test missing duration/end time
		create = &interfaces.ReservationCreate{
			Name:      "test",
			StartTime: time.Now(),
		}
		result, err = mgr.Create(context.Background(), create)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "end time or duration is required")
	})
}

func TestReservationManagerImpl_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/slurm/v0.0.43/reservation", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			// Return success response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test Update
		mgr := NewReservationManagerImpl(client)
		newDuration := 7200
		update := &interfaces.ReservationUpdate{
			Duration: &newDuration,
		}
		err = mgr.Update(context.Background(), "test-reservation", update)
		require.NoError(t, err)
	})
}

func TestReservationManagerImpl_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/slurm/v0.0.43/reservation/test-reservation", r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)

			// Return success response
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		// Create client
		client, err := createTestClient(server.URL)
		require.NoError(t, err)

		// Test Delete
		mgr := NewReservationManagerImpl(client)
		err = mgr.Delete(context.Background(), "test-reservation")
		require.NoError(t, err)
	})
}