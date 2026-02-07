// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// Account represents a SLURM Account.
type Account struct {
	Associations []AssocShort `json:"associations,omitempty"` // Associations involving this account (only populated if requested)
	Coordinators []Coord `json:"coordinators,omitempty"` // List of users that are a coordinator of this account (only populated if...
	Description string `json:"description"` // Arbitrary string describing the account
	Flags []AccountFlagsValue `json:"flags,omitempty"` // Flags associated with this account
	Name string `json:"name"` // Account name
	Organization string `json:"organization"` // Organization to which the account belongs
}


// AccountFlagsValue represents possible values for AccountFlags field.
type AccountFlagsValue string

// AccountFlagsValue constants.
const (
	AccountFlagsDeleted AccountFlagsValue = "DELETED"
	AccountFlagsWithassociations AccountFlagsValue = "WithAssociations"
	AccountFlagsWithcoordinators AccountFlagsValue = "WithCoordinators"
	AccountFlagsNousersarecoords AccountFlagsValue = "NoUsersAreCoords"
	AccountFlagsUsersarecoords AccountFlagsValue = "UsersAreCoords"
)
