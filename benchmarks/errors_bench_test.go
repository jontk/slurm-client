// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	"fmt"
	"testing"

	"github.com/jontk/slurm-client/pkg/errors"
)

func BenchmarkNewSlurmError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errors.NewSlurmError(
			errors.ErrorCodeInvalidRequest,
			"Invalid request format",
			map[string]interface{}{
				"field":  "job_name",
				"reason": "too long",
			},
		)
	}
}

func BenchmarkWrapError(b *testing.B) {
	originalErr := fmt.Errorf("connection timeout")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errors.WrapError(originalErr, errors.ErrorCodeTimeout)
	}
}

func BenchmarkIsErrorCode(b *testing.B) {
	err := errors.NewSlurmError(
		errors.ErrorCodeUnauthorized,
		"Authentication failed",
		nil,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errors.IsErrorCode(err, errors.ErrorCodeUnauthorized)
	}
}

func BenchmarkErrorFormatting(b *testing.B) {
	err := errors.NewSlurmError(
		errors.ErrorCodeResourceNotFound,
		"Job not found",
		map[string]interface{}{
			"job_id": 12345,
			"user":   "testuser",
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkNestedErrorWrapping(b *testing.B) {
	baseErr := fmt.Errorf("database connection failed")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err1 := errors.WrapError(baseErr, errors.ErrorCodeDatabaseError)
		err2 := errors.WrapError(err1, errors.ErrorCodeInternalError)
		_ = errors.WrapError(err2, errors.ErrorCodeServiceUnavailable)
	}
}