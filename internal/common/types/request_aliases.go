// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package types

// Request type aliases for backward compatibility and adapter usage

// AccountCreateRequest is an alias for AccountCreate
type AccountCreateRequest = AccountCreate

// AccountUpdateRequest is an alias for AccountUpdate  
type AccountUpdateRequest = AccountUpdate

// AssociationCreateRequest is an alias for AssociationCreate
type AssociationCreateRequest = AssociationCreate

// AssociationUpdateRequest is an alias for AssociationUpdate
type AssociationUpdateRequest = AssociationUpdate

// JobSubmitRequest is an alias for JobCreate
type JobSubmitRequest = JobCreate

// JobUpdateRequest is an alias for JobUpdate
type JobUpdateRequest = JobUpdate

// PartitionCreateRequest is an alias for PartitionCreate
type PartitionCreateRequest = PartitionCreate

// NodeUpdateRequest is an alias for NodeUpdate
type NodeUpdateRequest = NodeUpdate

// UserCreateRequest is an alias for UserCreate
type UserCreateRequest = UserCreate

// UserUpdateRequest is an alias for UserUpdate
type UserUpdateRequest = UserUpdate

// Additional request types
type QoSCreateRequest = QoSCreate
type QoSUpdateRequest = QoSUpdate
type NodeCreateRequest = NodeUpdate  // Nodes are usually not created programmatically
type ReservationCreateRequest = ReservationCreate
type ReservationUpdateRequest = ReservationUpdate
type PartitionUpdateRequest = PartitionUpdate
