// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// CronEntry represents a SLURM CronEntry.
type CronEntry struct {
	Command *string `json:"command,omitempty"` // Command to run
	DayOfMonth *string `json:"day_of_month,omitempty"` // Ranged string specifying eligible day of month values (e.g. 0-10,29)
	DayOfWeek *string `json:"day_of_week,omitempty"` // Ranged string specifying eligible day of week values (e.g.0-3,7)
	Flags []FlagsValue `json:"flags,omitempty"` // Flags
	Hour *string `json:"hour,omitempty"` // Ranged string specifying eligible hour values (e.g. 0-5,23)
	Line *CronEntryLine `json:"line,omitempty"`
	Minute *string `json:"minute,omitempty"` // Ranged string specifying eligible minute values (e.g. 0-10,50)
	Month *string `json:"month,omitempty"` // Ranged string specifying eligible month values (e.g. 0-5,12)
	Specification *string `json:"specification,omitempty"` // Complete time specification (* means valid for all allowed values) - minute...
}


// CronEntryLine is a nested type within its parent.
type CronEntryLine struct {
	End *int32 `json:"end,omitempty"` // End of this entry in file
	Start *int32 `json:"start,omitempty"` // Start of this entry in file
}
