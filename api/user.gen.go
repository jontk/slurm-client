// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// User represents a SLURM User.
type User struct {
	AdministratorLevel []AdministratorLevelValue `json:"administrator_level,omitempty"` // AdminLevel granted to the user
	Associations []AssocShort `json:"associations,omitempty"` // Associations created for this user
	Coordinators []Coord `json:"coordinators,omitempty"` // Accounts this user is a coordinator for
	Default *UserDefault `json:"default,omitempty"`
	Flags []UserDefaultFlagsValue `json:"flags,omitempty"` // Flags associated with this user
	Name string `json:"name"` // User name
	OldName *string `json:"old_name,omitempty"` // Previous user name
	Wckeys []WCKey `json:"wckeys,omitempty"` // List of available WCKeys
}


// UserDefault is a nested type within its parent.
type UserDefault struct {
	Account *string `json:"account,omitempty"` // Default account
	QoS *int32 `json:"qos,omitempty"` // Default QOS
	Wckey *string `json:"wckey,omitempty"` // Default WCKey
}


// AdministratorLevelValue represents possible values for AdministratorLevel field.
type AdministratorLevelValue string

// AdministratorLevelValue constants.
const (
	AdministratorLevelNotSet AdministratorLevelValue = "Not Set"
	AdministratorLevelNone AdministratorLevelValue = "None"
	AdministratorLevelOperator AdministratorLevelValue = "Operator"
	AdministratorLevelAdministrator AdministratorLevelValue = "Administrator"
)


// UserDefaultFlagsValue represents possible values for UserDefaultFlags field.
type UserDefaultFlagsValue string

// UserDefaultFlagsValue constants.
const (
	UserDefaultFlagsNone UserDefaultFlagsValue = "NONE"
	UserDefaultFlagsDeleted UserDefaultFlagsValue = "DELETED"
)
