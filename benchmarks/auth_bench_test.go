// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	"context"
	"net/http"
	"testing"

	"github.com/jontk/slurm-client/pkg/auth"
)

func BenchmarkTokenAuth(b *testing.B) {
	ta := auth.NewTokenAuth("test-jwt-token-with-reasonable-length")
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ta.Authenticate(ctx, req)
	}
}

func BenchmarkBasicAuth(b *testing.B) {
	ba := auth.NewBasicAuth("username", "password")
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ba.Authenticate(ctx, req)
	}
}

func BenchmarkNoAuth(b *testing.B) {
	na := auth.NewNoAuth()
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		na.Authenticate(ctx, req)
	}
}
