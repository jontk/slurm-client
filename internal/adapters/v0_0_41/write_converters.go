// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"encoding/json"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// =============================================================================
// Account Write Converters
// =============================================================================

// convertAccountCreateToAPI converts AccountCreate to the v0.0.41 API request body.
// Uses JSON marshaling to work around anonymous struct types.
func (a *AccountAdapter) convertAccountCreateToAPI(input *types.AccountCreate) (api.SlurmdbV0041PostAccountsJSONRequestBody, error) {
	if input == nil {
		return api.SlurmdbV0041PostAccountsJSONRequestBody{}, nil
	}

	// Build account structure
	accountMap := make(map[string]interface{})
	accountMap["name"] = input.Name

	if input.Description != "" {
		accountMap["description"] = input.Description
	}
	if input.Organization != "" {
		accountMap["organization"] = input.Organization
	}
	if len(input.Coordinators) > 0 {
		coords := make([]map[string]interface{}, len(input.Coordinators))
		for i, name := range input.Coordinators {
			coords[i] = map[string]interface{}{"name": name}
		}
		accountMap["coordinators"] = coords
	}

	// Build request body
	bodyMap := map[string]interface{}{
		"accounts": []interface{}{accountMap},
	}

	// Marshal and unmarshal to API type
	jsonBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return api.SlurmdbV0041PostAccountsJSONRequestBody{}, err
	}

	var body api.SlurmdbV0041PostAccountsJSONRequestBody
	if err := json.Unmarshal(jsonBytes, &body); err != nil {
		return api.SlurmdbV0041PostAccountsJSONRequestBody{}, err
	}

	return body, nil
}

// convertAccountUpdateToAPI converts Account to the v0.0.41 API request body for updates.
func (a *AccountAdapter) convertAccountUpdateToAPI(account *types.Account) (api.SlurmdbV0041PostAccountsJSONRequestBody, error) {
	if account == nil {
		return api.SlurmdbV0041PostAccountsJSONRequestBody{}, nil
	}

	// Build account structure
	accountMap := make(map[string]interface{})
	accountMap["name"] = account.Name
	accountMap["description"] = account.Description
	accountMap["organization"] = account.Organization

	if len(account.Coordinators) > 0 {
		coords := make([]map[string]interface{}, len(account.Coordinators))
		for i, coord := range account.Coordinators {
			coords[i] = map[string]interface{}{"name": coord.Name}
		}
		accountMap["coordinators"] = coords
	}

	// Build request body
	bodyMap := map[string]interface{}{
		"accounts": []interface{}{accountMap},
	}

	// Marshal and unmarshal to API type
	jsonBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return api.SlurmdbV0041PostAccountsJSONRequestBody{}, err
	}

	var body api.SlurmdbV0041PostAccountsJSONRequestBody
	if err := json.Unmarshal(jsonBytes, &body); err != nil {
		return api.SlurmdbV0041PostAccountsJSONRequestBody{}, err
	}

	return body, nil
}

// =============================================================================
// QoS Write Converters
// =============================================================================

// convertQoSCreateToAPI converts QoSCreate to the v0.0.41 API request body.
func (a *QoSAdapter) convertQoSCreateToAPI(input *types.QoSCreate) (api.SlurmdbV0041PostQosJSONRequestBody, error) {
	if input == nil {
		return api.SlurmdbV0041PostQosJSONRequestBody{}, nil
	}

	// Build QoS structure
	qosMap := make(map[string]interface{})
	qosMap["name"] = input.Name

	if input.Description != "" {
		qosMap["description"] = input.Description
	}
	if input.Priority > 0 {
		qosMap["priority"] = map[string]interface{}{
			"set":    true,
			"number": int32(input.Priority),
		}
	}
	if input.GraceTime > 0 {
		qosMap["grace_time"] = int32(input.GraceTime)
	}
	if len(input.PreemptMode) > 0 {
		qosMap["preempt"] = map[string]interface{}{
			"mode": input.PreemptMode,
		}
	}

	// Build request body
	bodyMap := map[string]interface{}{
		"qos": []interface{}{qosMap},
	}

	// Marshal and unmarshal to API type
	jsonBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return api.SlurmdbV0041PostQosJSONRequestBody{}, err
	}

	var body api.SlurmdbV0041PostQosJSONRequestBody
	if err := json.Unmarshal(jsonBytes, &body); err != nil {
		return api.SlurmdbV0041PostQosJSONRequestBody{}, err
	}

	return body, nil
}

// convertQoSUpdateToAPI converts QoS to the v0.0.41 API request body for updates.
func (a *QoSAdapter) convertQoSUpdateToAPI(qos *types.QoS) (api.SlurmdbV0041PostQosJSONRequestBody, error) {
	if qos == nil {
		return api.SlurmdbV0041PostQosJSONRequestBody{}, nil
	}

	// Build QoS structure
	qosMap := make(map[string]interface{})
	if qos.Name != nil {
		qosMap["name"] = *qos.Name
	}
	if qos.Description != nil {
		qosMap["description"] = *qos.Description
	}
	if qos.Priority != nil {
		qosMap["priority"] = map[string]interface{}{
			"set":    true,
			"number": int32(*qos.Priority),
		}
	}
	if qos.Limits != nil && qos.Limits.GraceTime != nil {
		qosMap["grace_time"] = *qos.Limits.GraceTime
	}
	if qos.Preempt != nil && len(qos.Preempt.Mode) > 0 {
		modes := make([]string, len(qos.Preempt.Mode))
		for i, m := range qos.Preempt.Mode {
			modes[i] = string(m)
		}
		qosMap["preempt"] = map[string]interface{}{
			"mode": modes,
		}
	}

	// Build request body
	bodyMap := map[string]interface{}{
		"qos": []interface{}{qosMap},
	}

	// Marshal and unmarshal to API type
	jsonBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return api.SlurmdbV0041PostQosJSONRequestBody{}, err
	}

	var body api.SlurmdbV0041PostQosJSONRequestBody
	if err := json.Unmarshal(jsonBytes, &body); err != nil {
		return api.SlurmdbV0041PostQosJSONRequestBody{}, err
	}

	return body, nil
}

// =============================================================================
// Association Write Converters
// =============================================================================

// convertAssociationCreateToAPI converts AssociationCreate to the v0.0.41 API request body.
func (a *AssociationAdapter) convertAssociationCreateToAPI(input *types.AssociationCreate) (api.SlurmdbV0041PostAssociationsJSONRequestBody, error) {
	if input == nil {
		return api.SlurmdbV0041PostAssociationsJSONRequestBody{}, nil
	}

	// Build association structure
	assocMap := make(map[string]interface{})

	if input.Account != "" {
		assocMap["account"] = input.Account
	}
	if input.User != "" {
		assocMap["user"] = input.User
	}
	if input.Cluster != "" {
		assocMap["cluster"] = input.Cluster
	}
	if input.Partition != "" {
		assocMap["partition"] = input.Partition
	}
	if input.DefaultQoS != "" {
		assocMap["default"] = map[string]interface{}{
			"qos": input.DefaultQoS,
		}
	}
	if input.SharesRaw != 0 {
		assocMap["shares_raw"] = input.SharesRaw
	}

	// Build request body
	bodyMap := map[string]interface{}{
		"associations": []interface{}{assocMap},
	}

	// Marshal and unmarshal to API type
	jsonBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return api.SlurmdbV0041PostAssociationsJSONRequestBody{}, err
	}

	var body api.SlurmdbV0041PostAssociationsJSONRequestBody
	if err := json.Unmarshal(jsonBytes, &body); err != nil {
		return api.SlurmdbV0041PostAssociationsJSONRequestBody{}, err
	}

	return body, nil
}

// convertAssociationUpdateToAPI converts AssociationUpdate to the v0.0.41 API request body.
func (a *AssociationAdapter) convertAssociationUpdateToAPI(id string, update *types.AssociationUpdate) (api.SlurmdbV0041PostAssociationsJSONRequestBody, error) {
	if update == nil {
		return api.SlurmdbV0041PostAssociationsJSONRequestBody{}, nil
	}

	// Build association structure with ID info
	assocMap := make(map[string]interface{})

	// Parse ID to get account/user/cluster/partition
	// ID format: "account:user:cluster[:partition]"
	// We need these to identify the association to update

	if update.DefaultQoS != nil {
		assocMap["default"] = map[string]interface{}{
			"qos": *update.DefaultQoS,
		}
	}
	if update.SharesRaw != nil {
		assocMap["shares_raw"] = *update.SharesRaw
	}

	// Build request body
	bodyMap := map[string]interface{}{
		"associations": []interface{}{assocMap},
	}

	// Marshal and unmarshal to API type
	jsonBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return api.SlurmdbV0041PostAssociationsJSONRequestBody{}, err
	}

	var body api.SlurmdbV0041PostAssociationsJSONRequestBody
	if err := json.Unmarshal(jsonBytes, &body); err != nil {
		return api.SlurmdbV0041PostAssociationsJSONRequestBody{}, err
	}

	return body, nil
}

// =============================================================================
// Job Write Converters
// =============================================================================

// convertJobUpdateToAPI converts JobUpdate to the v0.0.41 API request body.
func (a *JobAdapter) convertJobUpdateToAPI(update *types.JobUpdate) (api.SlurmV0041PostJobJSONRequestBody, error) {
	if update == nil {
		return api.SlurmV0041PostJobJSONRequestBody{}, nil
	}

	// Build job update structure
	jobMap := make(map[string]interface{})

	// String fields
	if update.Account != nil {
		jobMap["account"] = *update.Account
	}
	if update.Name != nil {
		jobMap["name"] = *update.Name
	}
	if update.Partition != nil {
		jobMap["partition"] = *update.Partition
	}
	if update.Comment != nil {
		jobMap["comment"] = *update.Comment
	}
	if update.QoS != nil {
		jobMap["qos"] = *update.QoS
	}

	// Numeric fields with no_val struct wrapper
	if update.Priority != nil {
		jobMap["priority"] = map[string]interface{}{
			"set":    true,
			"number": int32(*update.Priority),
		}
	}
	if update.TimeLimit != nil {
		jobMap["time_limit"] = map[string]interface{}{
			"set":    true,
			"number": int32(*update.TimeLimit),
		}
	}

	// Boolean fields
	if update.Hold != nil {
		jobMap["hold"] = *update.Hold
	}
	if update.Requeue != nil {
		jobMap["requeue"] = *update.Requeue
	}

	// Marshal and unmarshal to API type
	jsonBytes, err := json.Marshal(jobMap)
	if err != nil {
		return api.SlurmV0041PostJobJSONRequestBody{}, err
	}

	var body api.SlurmV0041PostJobJSONRequestBody
	if err := json.Unmarshal(jsonBytes, &body); err != nil {
		return api.SlurmV0041PostJobJSONRequestBody{}, err
	}

	return body, nil
}
