// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

// This file previously contained convertJobSubmissionToAPI function which has been removed
// due to incompatibility with the anonymous struct types in the v0.0.41 generated client.
//
// The v0.0.41 API version uses SlurmV0041PostJobSubmitJSONBody with an anonymous struct
// for the Job field, making it difficult to work with in conversion helpers.
//
// Future work: Either regenerate the v0.0.41 client.go with named types or implement
// proper conversion directly in job_manager_impl.go.
//
// For now, job submission is marked as unsupported for v0.0.41. Users should use a
// newer API version (v0.0.42 or later) for full job submission functionality.
