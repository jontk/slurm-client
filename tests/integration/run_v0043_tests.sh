#!/bin/bash

# SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
# SPDX-License-Identifier: Apache-2.0

set -e

echo "=== V0.0.43 New Features Integration Test Runner ==="
echo "Server: localhost:6820"
echo "API Version: v0.0.43"
echo

# Set environment variables for our tests
export SLURM_V0043_NEW_FEATURES_TEST=true
export SLURM_REAL_SERVER_TEST=true

# Change to project root directory
cd "$(dirname "$0")/../.."

echo "Running V0.0.43 new features integration tests..."
echo "Press Ctrl+C to stop"
echo

# Run the new features test suite
go test -v -timeout 10m ./tests/integration -run TestV0043NewFeaturesSuite

echo
echo "=== Test Run Complete ==="