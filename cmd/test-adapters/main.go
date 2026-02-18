// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Test v0.0.42 adapter methods for Shares, TRES, WCKeys, Accounts, Users, Associations
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/config"
)

// userTokenAuth implements authentication with both username and token headers
type userTokenAuth struct {
	username string
	token    string
}

func (u *userTokenAuth) Authenticate(ctx context.Context, req *http.Request) error {
	req.Header.Set("X-SLURM-USER-NAME", u.username)
	req.Header.Set("X-SLURM-USER-TOKEN", u.token)
	return nil
}

func (u *userTokenAuth) Type() string {
	return "user-token"
}

func main() {
	ctx := context.Background()

	// Get JWT token from environment
	jwtToken := os.Getenv("SLURM_JWT")
	if jwtToken == "" {
		log.Fatal("SLURM_JWT environment variable is required")
	}

	username := os.Getenv("SLURM_USER")
	if username == "" {
		username = "root" // Default username for testing
	}

	// Create client
	c, err := createClient(ctx, jwtToken, username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create client: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Testing v0.0.42 Adapter Methods ===")

	// Run all tests
	testGetShares(ctx, c)
	testGetTRES(ctx, c)
	testWCKeys(ctx, c)
	testAccounts(ctx, c)
	testUsers(ctx, c)
	testAssociations(ctx, c)

	fmt.Println("=== Test Complete ===")
}

func createClient(ctx context.Context, jwtToken string, username string) (types.SlurmClient, error) {
	cfg := config.NewDefault()
	cfg.BaseURL = "http://localhost:6820"
	cfg.Debug = false

	authProvider := &userTokenAuth{
		username: username,
		token:    jwtToken,
	}

	clientFactory, err := factory.NewClientFactory(
		factory.WithConfig(cfg),
		factory.WithAuth(authProvider),
		factory.WithBaseURL(cfg.BaseURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create factory: %w", err)
	}

	return clientFactory.NewClientWithVersion(ctx, "v0.0.42")
}

func printResult(label string, count int, sample interface{}, err error) {
	if err != nil {
		fmt.Printf("   ❌ ERROR: %v\n", err)
	} else {
		fmt.Printf("   ✅ SUCCESS: Retrieved %d %s\n", count, label)
		if sample != nil {
			if jsonData, jsonErr := json.MarshalIndent(sample, "      ", "  "); jsonErr == nil {
				fmt.Printf("      Sample: %s\n", string(jsonData))
			}
		}
	}
	fmt.Println()
}

func testGetShares(ctx context.Context, c types.SlurmClient) {
	fmt.Println("1. Testing GetShares()...")
	shares, err := c.GetShares(ctx, &types.GetSharesOptions{})
	var sample interface{}
	count := 0
	if shares != nil {
		count = len(shares.Shares)
		if count > 0 {
			sample = shares.Shares[0]
		}
	}
	printResult("shares", count, sample, err)
}

func testGetTRES(ctx context.Context, c types.SlurmClient) {
	fmt.Println("2. Testing GetTRES()...")
	tres, err := c.GetTRES(ctx)
	var sample interface{}
	count := 0
	if tres != nil {
		count = len(tres.TRES)
		if count > 0 {
			sample = tres.TRES[0]
		}
	}
	printResult("TRES", count, sample, err)
}

func testWCKeys(ctx context.Context, c types.SlurmClient) {
	fmt.Println("3. Testing WCKeys.List()...")
	wckeys, err := c.WCKeys().List(ctx, nil)
	var sample interface{}
	count := 0
	if wckeys != nil {
		count = len(wckeys.WCKeys)
		if count > 0 {
			sample = wckeys.WCKeys[0]
		}
	}
	printResult("WCKeys", count, sample, err)
}

func testAccounts(ctx context.Context, c types.SlurmClient) {
	fmt.Println("4. Testing Accounts.List()...")
	accounts, err := c.Accounts().List(ctx, nil)
	var sample interface{}
	count := 0
	if accounts != nil {
		count = len(accounts.Accounts)
		if count > 0 {
			sample = accounts.Accounts[0]
		}
	}
	printResult("accounts", count, sample, err)
}

func testUsers(ctx context.Context, c types.SlurmClient) {
	fmt.Println("5. Testing Users.List()...")
	users, err := c.Users().List(ctx, nil)
	var sample interface{}
	count := 0
	if users != nil {
		count = len(users.Users)
		if count > 0 {
			sample = users.Users[0]
		}
	}
	printResult("users", count, sample, err)
}

func testAssociations(ctx context.Context, c types.SlurmClient) {
	fmt.Println("6. Testing Associations.List()...")
	assocs, err := c.Associations().List(ctx, nil)
	var sample interface{}
	count := 0
	if assocs != nil {
		count = len(assocs.Associations)
		if count > 0 {
			sample = assocs.Associations[0]
		}
	}
	printResult("associations", count, sample, err)
}
