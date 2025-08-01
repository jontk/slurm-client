// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	"net/http"
	"testing"

	"github.com/jontk/slurm-client/pkg/auth"
)

func BenchmarkTokenAuth(b *testing.B) {
	ta := auth.NewTokenAuth("test-jwt-token-with-reasonable-length")
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ta.Apply(req)
	}
}

func BenchmarkAPIKeyAuth(b *testing.B) {
	aka := auth.NewAPIKeyAuth("X-API-Key", "test-api-key-value")
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aka.Apply(req)
	}
}

func BenchmarkBasicAuth(b *testing.B) {
	ba := auth.NewBasicAuth("username", "password")
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ba.Apply(req)
	}
}

func BenchmarkAuthChain(b *testing.B) {
	chain := auth.Chain(
		auth.NewTokenAuth("token"),
		auth.NewAPIKeyAuth("X-API-Key", "key"),
		auth.NewBasicAuth("user", "pass"),
	)
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain.Apply(req)
	}
}