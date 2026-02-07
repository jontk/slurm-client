// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package types provides common type definitions for SLURM entities.
// Core entity types (Job, Node, User, etc.) are generated in *.gen.go files.
// This file contains operation types (Create, Update, List, etc.).
package api

// QoSCreate represents the data needed to create a new QoS
type QoSCreate struct {
	Name              string
	Description       string
	Priority          int
	Flags             []string
	PreemptMode       []string
	PreemptList       []string
	PreemptExemptTime *int
	GraceTime         int // Changed to non-pointer for validation
	UsageFactor       float64
	UsageThreshold    float64
	ParentQoS         string
	MaxTRESPerUser    string
	MaxTRESPerAccount string
	MaxTRESPerJob     string
	Limits            *QoSLimits
}

// QoSUpdate represents fields that can be updated on a QoS
type QoSUpdate struct {
	Description       *string
	Priority          *int
	Flags             *[]string
	PreemptMode       *[]string
	PreemptList       []string
	PreemptExemptTime *int
	GraceTime         *int
	UsageFactor       *float64
	UsageThreshold    *float64
	ParentQoS         *string
	MaxTRESPerUser    *string
	MaxTRESPerAccount *string
	MaxTRESPerJob     *string
	Limits            *QoSLimits
}

// QoSListOptions represents options for listing QoS entries
type QoSListOptions struct {
	Names    []string
	Accounts []string
	Users    []string

	// Limit specifies the maximum number of QoS entries to return.
	// WARNING: Due to SLURM REST API limitations, this is CLIENT-SIDE pagination.
	// The full QoS list is fetched from the server, then sliced. Consider using
	// filtering options (Names, Accounts, Users) to reduce the dataset before pagination.
	Limit int

	// Offset specifies the number of QoS entries to skip before returning results.
	// WARNING: This is CLIENT-SIDE pagination - see Limit field documentation.
	Offset int
}

// QoSList represents a list of QoS entries
type QoSList struct {
	QoS   []QoS
	Total int
}

// QoSCreateResponse represents the response from creating a QoS
type QoSCreateResponse struct {
	QoSName string
}
