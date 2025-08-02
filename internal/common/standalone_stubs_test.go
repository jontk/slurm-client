// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandaloneOperationsStub(t *testing.T) {
	stub := &StandaloneOperationsStub{Version: "v0.0.40"}
	ctx := context.Background()
	
	// Test GetLicenses
	licenses, err := stub.GetLicenses(ctx)
	require.Error(t, err)
	assert.Nil(t, licenses)
	assert.Contains(t, err.Error(), "not implemented")
	
	// Test GetShares
	shares, err := stub.GetShares(ctx, nil)
	require.Error(t, err)
	assert.Nil(t, shares)
	assert.Contains(t, err.Error(), "not implemented")
	
	// Test GetConfig
	config, err := stub.GetConfig(ctx)
	require.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "not implemented")
	
	// Test GetDiagnostics
	diagnostics, err := stub.GetDiagnostics(ctx)
	require.Error(t, err)
	assert.Nil(t, diagnostics)
	assert.Contains(t, err.Error(), "not implemented")
	
	// Test GetDBDiagnostics
	dbDiagnostics, err := stub.GetDBDiagnostics(ctx)
	require.Error(t, err)
	assert.Nil(t, dbDiagnostics)
	assert.Contains(t, err.Error(), "not implemented")
	
	// Test GetInstance
	instance, err := stub.GetInstance(ctx, nil)
	require.Error(t, err)
	assert.Nil(t, instance)
	assert.Contains(t, err.Error(), "not implemented")
	
	// Test GetInstances
	instances, err := stub.GetInstances(ctx, nil)
	require.Error(t, err)
	assert.Nil(t, instances)
	assert.Contains(t, err.Error(), "not implemented")
	
	// Test GetTRES
	tres, err := stub.GetTRES(ctx)
	require.Error(t, err)
	assert.Nil(t, tres)
	assert.Contains(t, err.Error(), "not implemented")
	
	// Test CreateTRES
	createdTres, err := stub.CreateTRES(ctx, nil)
	require.Error(t, err)
	assert.Nil(t, createdTres)
	assert.Contains(t, err.Error(), "not implemented")
	
	// Test Reconfigure
	reconfigResp, err := stub.Reconfigure(ctx)
	require.Error(t, err)
	assert.Nil(t, reconfigResp)
	assert.Contains(t, err.Error(), "not implemented")
}