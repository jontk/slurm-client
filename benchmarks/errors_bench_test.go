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
	for range b.N {
		_ = errors.NewSlurmError(
			errors.ErrorCodeInvalidRequest,
			"Invalid request format",
		)
	}
}

func BenchmarkWrapError(b *testing.B) {
	originalErr := fmt.Errorf("connection timeout")

	b.ResetTimer()
	for range b.N {
		_ = errors.WrapError(originalErr)
	}
}

func BenchmarkGetErrorCode(b *testing.B) {
	err := errors.NewSlurmError(
		errors.ErrorCodeUnauthorized,
		"Authentication failed",
	)

	b.ResetTimer()
	for range b.N {
		_ = errors.GetErrorCode(err)
	}
}

func BenchmarkErrorFormatting(b *testing.B) {
	err := errors.NewSlurmError(
		errors.ErrorCodeResourceNotFound,
		"Job not found",
	)

	b.ResetTimer()
	for range b.N {
		_ = err.Error()
	}
}

func BenchmarkNestedErrorWrapping(b *testing.B) {
	baseErr := fmt.Errorf("connection failed")

	b.ResetTimer()
	for range b.N {
		err1 := errors.WrapError(baseErr)
		err2 := errors.NewSlurmErrorWithCause(errors.ErrorCodeNetworkTimeout, "network timeout", err1)
		_ = errors.NewSlurmErrorWithCause(errors.ErrorCodeServiceUnavailable, "service unavailable", err2)
	}
}
