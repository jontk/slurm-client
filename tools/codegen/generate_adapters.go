//go:build ignore

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// generate_adapters.go - Adapter generator for Slurm SDK
// Generates adapter method bodies that call generated converters
// Usage: go run generate_adapters.go -version=v0_0_44

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Config structures
type AdapterConfig struct {
	Generation GenerationConfig            `yaml:"generation"`
	Patterns   map[string]PatternConfig    `yaml:"patterns"`
	Entities   map[string]EntityDef        `yaml:"entities"`
	Versions   map[string]VersionDef       `yaml:"versions"`
}

type GenerationConfig struct {
	OutputSuffix       string `yaml:"output_suffix"`
	GenerateConverters bool   `yaml:"generate_converters"`
	GenerateAdapters   bool   `yaml:"generate_adapters"`
}

type PatternConfig struct {
	Description string   `yaml:"description"`
	Steps       []string `yaml:"steps"`
}

type EntityDef struct {
	Identifier         string                       `yaml:"identifier"`
	IdentifierType     string                       `yaml:"identifier_type"`
	IdentifierDisplay  string                       `yaml:"identifier_display"`
	ListOptions        string                       `yaml:"list_options"`
	ListResult         string                       `yaml:"list_result"`
	ListField          string                       `yaml:"list_field"` // Field name in list result
	CreateInput        string                       `yaml:"create_input"`
	CreateResult       string                       `yaml:"create_result"`
	UpdateInput        string                       `yaml:"update_input"`
	WatchOptions       string                       `yaml:"watch_options"`
	AssociationRequest string                       `yaml:"association_request"`
	AssociationResult  string                       `yaml:"association_result"`
	Methods            []string                     `yaml:"methods"`
	Validation         map[string]ValidationConfig  `yaml:"validation"`
}

type ValidationConfig struct {
	NilError        string              `yaml:"nil_error"`
	RequiredFields  []RequiredField     `yaml:"required_fields"`
	AtLeastOneOf    []string            `yaml:"at_least_one_of"`
	AtLeastOneError string              `yaml:"at_least_one_error"`
}

type RequiredField struct {
	Field     string `yaml:"field"`
	Error     string `yaml:"error"`
	CheckType string `yaml:"check_type"` // "string", "time", "slice"
}

type VersionDef struct {
	APIPrefix   string                       `yaml:"api_prefix"`
	APIPackage  string                       `yaml:"api_package"`
	APIMethods  map[string]EntityAPIMethods  `yaml:"api_methods"`
}

type EntityAPIMethods struct {
	List              string `yaml:"list"`
	Get               string `yaml:"get"`
	Create            string `yaml:"create"`
	Update            string `yaml:"update"`
	Delete            string `yaml:"delete"`
	Submit            string `yaml:"submit"`
	Allocate          string `yaml:"allocate"`
	ListParams        string `yaml:"list_params"`
	GetParams         string `yaml:"get_params"`
	ResponseList      string `yaml:"response_list"`
	APIType           string `yaml:"api_type"`           // API type for list/get responses
	InputAPIType      string `yaml:"input_api_type"`     // API type for create/update input (if different from APIType)
	RequestBodyType   string `yaml:"request_body_type"`
	RequestListField  string `yaml:"request_list_field"`
	SubmitRequestBody   string `yaml:"submit_request_body"`
	AllocateRequestBody string `yaml:"allocate_request_body"`
	CreateNeedsParams   bool   `yaml:"create_needs_params"` // API requires params argument for create
	DeleteNeedsParams   bool   `yaml:"delete_needs_params"` // API requires params argument for delete
}

var (
	versionFlag = flag.String("version", "", "API version to generate (e.g., v0_0_44)")
	allFlag     = flag.Bool("all", false, "Generate for all versions")
	configFlag  = flag.String("config", "tools/codegen/adapter_config.yaml", "Path to config file")
	dryRunFlag  = flag.Bool("dry-run", false, "Print generated code without writing files")
)

func main() {
	flag.Parse()

	config, err := loadConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	versions := []string{}
	if *allFlag {
		for v := range config.Versions {
			versions = append(versions, v)
		}
	} else if *versionFlag != "" {
		versions = []string{*versionFlag}
	} else {
		fmt.Fprintln(os.Stderr, "Error: specify -version or -all")
		os.Exit(1)
	}

	for _, version := range versions {
		fmt.Printf("Generating adapters for %s...\n", version)
		if err := generateAdaptersForVersion(version, config); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating %s: %v\n", version, err)
			continue
		}
		fmt.Printf("  ✓ Done: %s\n", version)
	}
}

func loadConfig(path string) (*AdapterConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config AdapterConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func generateAdaptersForVersion(version string, config *AdapterConfig) error {
	versionDef, ok := config.Versions[version]
	if !ok {
		return fmt.Errorf("version %s not found in config", version)
	}

	for entityName, entityDef := range config.Entities {
		apiMethods, ok := versionDef.APIMethods[entityName]
		if !ok {
			// Skip entities not configured for this version
			continue
		}

		// Generate adapter file
		code, err := generateEntityAdapter(version, entityName, entityDef, versionDef, apiMethods)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: Could not generate %s adapter: %v\n", entityName, err)
			continue
		}

		outputPath := filepath.Join("internal/adapters", version,
			strings.ToLower(entityName)+"_adapter.gen.go")

		if *dryRunFlag {
			fmt.Printf("Would write to %s:\n%s\n", outputPath, code)
		} else {
			if err := os.WriteFile(outputPath, []byte(code), 0644); err != nil {
				return fmt.Errorf("writing %s: %w", outputPath, err)
			}
			fmt.Printf("  Generated: %s\n", outputPath)
		}

		// Generate validation file if entity has validation config
		if len(entityDef.Validation) > 0 {
			validationCode, err := generateEntityValidation(version, entityName, entityDef)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: Could not generate %s validation: %v\n", entityName, err)
				continue
			}

			validationPath := filepath.Join("internal/adapters", version,
				strings.ToLower(entityName)+"_validation.gen.go")

			if *dryRunFlag {
				fmt.Printf("Would write to %s:\n%s\n", validationPath, validationCode)
			} else {
				if err := os.WriteFile(validationPath, []byte(validationCode), 0644); err != nil {
					return fmt.Errorf("writing %s: %w", validationPath, err)
				}
				fmt.Printf("  Generated: %s\n", validationPath)
			}
		}

		// Generate helper file for special methods
		// Skip if there's a manual _extra.go file that may contain the implementations
		extraPath := filepath.Join("internal/adapters", version,
			strings.ToLower(entityName)+"_helpers_extra.go")
		if _, err := os.Stat(extraPath); err == nil {
			// Manual helper file exists, skip generation to avoid conflicts
			fmt.Printf("  Skipping helpers (manual file exists): %s\n", extraPath)
		} else {
			helperCode := generateEntityHelpers(version, entityName, entityDef, versionDef, apiMethods)
			if helperCode != "" {
				helperPath := filepath.Join("internal/adapters", version,
					strings.ToLower(entityName)+"_helpers.gen.go")

				if *dryRunFlag {
					fmt.Printf("Would write to %s:\n%s\n", helperPath, helperCode)
				} else {
					if err := os.WriteFile(helperPath, []byte(helperCode), 0644); err != nil {
						return fmt.Errorf("writing %s: %w", helperPath, err)
					}
					fmt.Printf("  Generated: %s\n", helperPath)
				}
			}
		}

		// Generate test file
		testCode, err := generateEntityTests(version, entityName, entityDef, versionDef, apiMethods)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: Could not generate %s tests: %v\n", entityName, err)
			continue
		}

		testPath := filepath.Join("internal/adapters", version,
			strings.ToLower(entityName)+"_adapter_gen_test.go")

		if *dryRunFlag {
			fmt.Printf("Would write to %s:\n%s\n", testPath, testCode)
		} else {
			if err := os.WriteFile(testPath, []byte(testCode), 0644); err != nil {
				return fmt.Errorf("writing %s: %w", testPath, err)
			}
			fmt.Printf("  Generated: %s\n", testPath)
		}
	}

	return nil
}

// generateEntityValidation generates the validation file for an entity
func generateEntityValidation(version, entityName string, entityDef EntityDef) (string, error) {
	var buf bytes.Buffer

	// Write header
	buf.WriteString(fmt.Sprintf(`// Code generated by generate_adapters.go. DO NOT EDIT.
// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package %s

import (
	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/pkg/errors"
)

`, version))

	// Generate validation methods for each operation (sorted for deterministic output)
	validationOps := make([]string, 0, len(entityDef.Validation))
	for op := range entityDef.Validation {
		validationOps = append(validationOps, op)
	}
	sort.Strings(validationOps)
	for _, operation := range validationOps {
		validationConfig := entityDef.Validation[operation]
		methodCode := generateValidationMethod(entityName, operation, validationConfig, entityDef)
		buf.WriteString(methodCode)
		buf.WriteString("\n")
	}

	// Format the code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), nil // Return unformatted on error
	}

	return string(formatted), nil
}

// generateEntityHelpers generates helper stub file for special methods
func generateEntityHelpers(version, entityName string, entityDef EntityDef, versionDef VersionDef, apiMethods EntityAPIMethods) string {
	var buf bytes.Buffer
	hasHelpers := false

	// Check which helpers are needed based on methods
	needsAssociationHelper := false
	needsJobHelpers := false
	needsNodeHelpers := false

	for _, method := range entityDef.Methods {
		switch method {
		case "create_association":
			needsAssociationHelper = true
		case "hold", "signal", "notify", "requeue", "allocate":
			if entityName == "Job" {
				needsJobHelpers = true
			}
		case "drain", "resume":
			if entityName == "Node" {
				needsNodeHelpers = true
			}
		case "watch":
			if entityName == "Job" {
				needsJobHelpers = true
			}
			if entityName == "Node" {
				needsNodeHelpers = true
			}
		}
	}

	if !needsAssociationHelper && !needsJobHelpers && !needsNodeHelpers {
		return ""
	}

	// Determine if api import is needed (Job and Node helpers use it)
	needsAPIImport := needsJobHelpers || needsNodeHelpers

	// Write header with conditional imports
	buf.WriteString(fmt.Sprintf(`// Code generated by generate_adapters.go. DO NOT EDIT.
// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package %s

import (
	"context"
`, version))

	if needsAPIImport {
		buf.WriteString(fmt.Sprintf(`
	api "github.com/jontk/slurm-client/internal/openapi/%s"`, version))
	}

	buf.WriteString(`
	types "github.com/jontk/slurm-client/api"`)

	// Only include errors package if needed (job helpers and association helpers use it)
	needsErrors := needsJobHelpers || needsAssociationHelper
	if needsErrors {
		buf.WriteString(`
	"github.com/jontk/slurm-client/pkg/errors"`)
	}

	// Add additional imports for job helpers (need strconv and common)
	if needsJobHelpers {
		buf.WriteString(`
	"strconv"
	"github.com/jontk/slurm-client/internal/common"`)
	}

	// Add common import for node helpers
	if needsNodeHelpers && !needsJobHelpers {
		buf.WriteString(`
	"github.com/jontk/slurm-client/internal/common"`)
	}

	buf.WriteString(`
)

`)

	// Generate association helper for Account and User
	if needsAssociationHelper && (entityName == "Account" || entityName == "User") {
		hasHelpers = true
		requestType := entityDef.AssociationRequest
		if requestType == "" {
			requestType = entityName + "AssociationRequest"
		}
		buf.WriteString(fmt.Sprintf(`// createAssociationImpl implements the CreateAssociation method
func (a *%sAdapter) createAssociationImpl(ctx context.Context, req *types.%s) (*types.AssociationCreateResponse, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "association request is required", "request", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// TODO: Implement association creation
	return nil, errors.NewClientError(errors.ErrorCodeUnsupportedOperation, "create association not yet implemented", "operation not supported")
}

`, entityName, requestType))
	}

	// Generate job helpers
	if needsJobHelpers {
		hasHelpers = true
		apiPrefix := versionDef.APIPrefix
		version := strings.Replace(versionDef.APIPackage, "_", ".", -1)

		// Version-specific constant name for FEDERATION_REQUEUE flag
		// v0_0_42 and v0_0_43 use short names (FEDERATIONREQUEUE)
		// v0_0_44+ uses full names (SlurmV0044DeleteJobParamsFlagsFEDERATIONREQUEUE)
		requeueFlagConst := "FEDERATIONREQUEUE"
		if versionDef.APIPackage == "v0_0_44" {
			requeueFlagConst = fmt.Sprintf("Slurm%sDeleteJobParamsFlagsFEDERATIONREQUEUE", apiPrefix)
		}

		buf.WriteString(fmt.Sprintf(`// holdJobImpl implements the Hold method
func (a *JobAdapter) holdJobImpl(ctx context.Context, req *types.JobHoldRequest) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate job ID
	if err := a.ValidateResourceID(req.JobId, "jobID"); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Build job update request with hold flag
	updateReq := api.%sJobDescMsg{
		Hold: &req.Hold,
	}

	// If priority is specified, include it
	if req.Priority != 0 {
		priority := api.%sUint32NoValStruct{
			Set:    ptrBool(true),
			Number: ptrInt32(req.Priority),
		}
		updateReq.Priority = &priority
	}

	// Call the API to update the job
	resp, err := a.client.Slurm%sPostJobWithResponse(ctx, strconv.Itoa(int(req.JobId)), updateReq)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "%s"); err != nil {
		return err
	}

	return nil
}

// signalJobImpl implements the Signal method
func (a *JobAdapter) signalJobImpl(ctx context.Context, req *types.JobSignalRequest) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate job ID
	if err := a.ValidateResourceID(req.JobId, "jobID"); err != nil {
		return err
	}

	// Validate signal
	if req.Signal == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "signal is required", "signal", nil, nil)
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// The SLURM REST API uses DELETE with a signal parameter to send signals to jobs
	params := &api.Slurm%sDeleteJobParams{
		Signal: &req.Signal,
	}

	// Call the API to signal the job
	resp, err := a.client.Slurm%sDeleteJobWithResponse(ctx, strconv.Itoa(int(req.JobId)), params)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "%s"); err != nil {
		return err
	}

	return nil
}

// notifyJobImpl implements the Notify method
func (a *JobAdapter) notifyJobImpl(ctx context.Context, req *types.JobNotifyRequest) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate job ID
	if err := a.ValidateResourceID(req.JobId, "jobID"); err != nil {
		return err
	}

	// Validate message
	if req.Message == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "message is required", "message", nil, nil)
	}

	// Job notification (scontrol notify) is not supported in SLURM REST API
	return errors.NewClientError(
		errors.ErrorCodeUnsupportedOperation,
		"Job notification is not supported via REST API %s. Use 'scontrol notify' command instead",
		"operation not supported in REST API",
	)
}

// requeueJobImpl implements the Requeue method
func (a *JobAdapter) requeueJobImpl(ctx context.Context, jobID int32) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate job ID
	if err := a.ValidateResourceID(jobID, "jobID"); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Job requeue uses DELETE with FEDERATION_REQUEUE flag
	requeueFlag := api.%s
	params := &api.Slurm%sDeleteJobParams{
		Flags: &requeueFlag,
	}

	// Call the API to requeue the job
	resp, err := a.client.Slurm%sDeleteJobWithResponse(ctx, strconv.Itoa(int(jobID)), params)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "%s"); err != nil {
		return err
	}

	return nil
}

// NOTE: watchJobsImpl is intentionally NOT generated here.
// The real implementation is in job_watch_extra.go which uses polling.

// allocateJobImpl implements the Allocate method
func (a *JobAdapter) allocateJobImpl(ctx context.Context, req *types.JobAllocateRequest) (*types.JobAllocateResponse, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	if req == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "allocate request is required", "request", nil, nil)
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Build job allocation request
	jobDesc := &api.%sJobDescMsg{}

	// Basic metadata
	if req.Name != "" {
		jobDesc.Name = &req.Name
	}
	if req.Account != "" {
		jobDesc.Account = &req.Account
	}
	if req.Partition != "" {
		jobDesc.Partition = &req.Partition
	}
	if req.QoS != "" {
		jobDesc.Qos = &req.QoS
	}

	// Resource requirements
	if req.Nodes != "" {
		// Nodes can be a number (e.g., "2") or range (e.g., "1-4")
		jobDesc.Nodes = &req.Nodes
	}
	if req.Cpus > 0 {
		jobDesc.MinimumCpus = &req.Cpus
	}
	if req.Memory != "" {
		// Parse memory specification (e.g., "4096" for MB, "4G" for GB)
		memVal, err := common.ParseMemory(req.Memory)
		if err == nil && memVal > 0 {
			memPerNode := &api.%sUint64NoValStruct{}
			memPerNode.Number = &memVal
			memPerNode.Set = ptrBool(true)
			jobDesc.MemoryPerNode = memPerNode
		}
	}
	if req.TimeLimit > 0 {
		timeLimit := &api.%sUint32NoValStruct{}
		val := int32(req.TimeLimit)
		timeLimit.Number = &val
		timeLimit.Set = ptrBool(true)
		jobDesc.TimeLimit = timeLimit
	}

	// Environment and execution
	if req.WorkingDir != "" {
		jobDesc.CurrentWorkingDirectory = &req.WorkingDir
	}
	if len(req.Environment) > 0 {
		envArray := &api.%sStringArray{}
		for key, value := range req.Environment {
			envStr := key + "=" + value
			*envArray = append(*envArray, envStr)
		}
		jobDesc.Environment = envArray
	}
	if len(req.Command) > 0 {
		cmdArray := &api.%sStringArray{}
		*cmdArray = append(*cmdArray, req.Command...)
		jobDesc.Argv = cmdArray
	}

	allocReq := api.%sJobAllocReq{
		Job: jobDesc,
	}

	// Call the API
	resp, err := a.client.Slurm%sPostJobAllocateWithResponse(ctx, allocReq)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "%s"); err != nil {
		return nil, err
	}

	// Convert response
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeInvalidRequest, "empty response from server", "no response body")
	}

	return a.convertAPIJobAllocateResponseToCommon(resp.JSON200), nil
}

// convertCommonJobUpdateToAPIRequestBody converts JobUpdate to the API request body type
func (a *JobAdapter) convertCommonJobUpdateToAPIRequestBody(update *types.JobUpdate) api.Slurm%sPostJobJSONRequestBody {
	if update == nil {
		return api.Slurm%sPostJobJSONRequestBody{}
	}
	// JobUpdate is an alias for JobCreate, so we can use the same converter
	// This delegates to the goverter-generated ConvertCommonJobCreateToAPI
	result := jobWriteConverter.ConvertCommonJobCreateToAPI(update)
	if result == nil {
		return api.Slurm%sPostJobJSONRequestBody{}
	}
	return *result
}

// convertAPIJobSubmitResponseToCommon converts API job submit response to common type
func (a *JobAdapter) convertAPIJobSubmitResponseToCommon(resp *api.%sOpenapiJobSubmitResponse) *types.JobSubmitResponse {
	if resp == nil {
		return nil
	}

	result := &types.JobSubmitResponse{}

	// Extract job ID
	if resp.JobId != nil {
		result.JobId = *resp.JobId
	}

	// Extract step ID
	if resp.StepId != nil {
		result.StepId = *resp.StepId
	}

	// Extract user message
	if resp.JobSubmitUserMsg != nil {
		result.JobSubmitUserMsg = *resp.JobSubmitUserMsg
	}

	// Extract errors
	if resp.Errors != nil {
		for _, e := range *resp.Errors {
			if e.Error != nil {
				result.Error = append(result.Error, *e.Error)
			}
		}
	}

	// Extract warnings
	if resp.Warnings != nil {
		for _, w := range *resp.Warnings {
			if w.Description != nil {
				result.Warning = append(result.Warning, *w.Description)
			}
		}
	}

	return result
}

// convertAPIJobAllocateResponseToCommon converts API job allocate response to common type
func (a *JobAdapter) convertAPIJobAllocateResponseToCommon(resp *api.%sOpenapiJobAllocResp) *types.JobAllocateResponse {
	if resp == nil {
		return nil
	}

	result := &types.JobAllocateResponse{
		Status: "allocated",
	}

	// Extract job ID
	if resp.JobId != nil {
		result.JobId = *resp.JobId
	}

	// Extract message
	if resp.JobSubmitUserMsg != nil {
		result.Message = *resp.JobSubmitUserMsg
	}

	return result
}

// Helper functions for pointer creation
func ptrBool(b bool) *bool {
	return &b
}

func ptrInt32(i int32) *int32 {
	return &i
}

`, apiPrefix, apiPrefix, apiPrefix, apiPrefix, version, apiPrefix, apiPrefix, apiPrefix, version, version, requeueFlagConst, apiPrefix, apiPrefix, apiPrefix, version, apiPrefix, apiPrefix, apiPrefix, apiPrefix, apiPrefix, apiPrefix, apiPrefix, apiPrefix, apiPrefix, apiPrefix, apiPrefix, apiPrefix, apiPrefix, apiPrefix))
	}

	// Generate node helpers
	if needsNodeHelpers {
		hasHelpers = true
		apiPrefix := versionDef.APIPrefix
		version := strings.Replace(versionDef.APIPackage, "_", ".", -1)

		// v0_0_42 uses string slices for node state, later versions use enums
		var drainStateCode string
		var resumeStateCode string
		if versionDef.APIPackage == "v0_0_42" {
			// v0_0_42 uses []string for state
			drainStateCode = fmt.Sprintf(`updateReq := api.%sUpdateNodeMsg{
		State: &[]string{"DRAIN"},
	}`, apiPrefix)
			resumeStateCode = fmt.Sprintf(`updateReq := api.%sUpdateNodeMsg{
		State: &[]string{"RESUME"},
	}`, apiPrefix)
		} else {
			// v0_0_43+ use enum constants
			drainStateCode = fmt.Sprintf(`drainState := api.%sUpdateNodeMsgStateDRAIN
	updateReq := api.%sUpdateNodeMsg{
		State: &[]api.%sUpdateNodeMsgState{drainState},
	}`, apiPrefix, apiPrefix, apiPrefix)
			resumeStateCode = fmt.Sprintf(`resumeState := api.%sUpdateNodeMsgStateRESUME
	updateReq := api.%sUpdateNodeMsg{
		State: &[]api.%sUpdateNodeMsgState{resumeState},
	}`, apiPrefix, apiPrefix, apiPrefix)
		}

		buf.WriteString(fmt.Sprintf(`// drainNodeImpl implements the Drain method
func (a *NodeAdapter) drainNodeImpl(ctx context.Context, nodeName string, reason string) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate node name
	if err := a.ValidateResourceName(nodeName, "nodeName"); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Build node update request with DRAIN state
	%s

	if reason != "" {
		updateReq.Reason = &reason
	}

	// Call the API to update the node
	resp, err := a.client.Slurm%sPostNodeWithResponse(ctx, nodeName, updateReq)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "%s"); err != nil {
		return err
	}

	return nil
}

// resumeNodeImpl implements the Resume method
func (a *NodeAdapter) resumeNodeImpl(ctx context.Context, nodeName string) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate node name
	if err := a.ValidateResourceName(nodeName, "nodeName"); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Build node update request with RESUME state
	%s

	// Call the API to update the node
	resp, err := a.client.Slurm%sPostNodeWithResponse(ctx, nodeName, updateReq)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "%s"); err != nil {
		return err
	}

	return nil
}

// NOTE: watchNodesImpl is intentionally NOT generated here.
// The real implementation is in node_watch_extra.go which uses polling.

// convertCommonNodeUpdateToAPIRequestBody converts NodeUpdate to the API request body type
func (a *NodeAdapter) convertCommonNodeUpdateToAPIRequestBody(update *types.NodeUpdate) api.Slurm%sPostNodeJSONRequestBody {
	if update == nil {
		return api.Slurm%sPostNodeJSONRequestBody{}
	}
	// Delegate to the goverter-generated ConvertCommonNodeUpdateToAPI
	result := nodeWriteConverter.ConvertCommonNodeUpdateToAPI(update)
	if result == nil {
		return api.Slurm%sPostNodeJSONRequestBody{}
	}
	return *result
}

`, drainStateCode, apiPrefix, apiPrefix, version, resumeStateCode, apiPrefix, apiPrefix, version, apiPrefix, apiPrefix, apiPrefix))
	}

	if !hasHelpers {
		return ""
	}

	// Format the code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String() // Return unformatted on error
	}

	return string(formatted)
}

// NOTE: generateEntityConverters has been removed - converters are now generated
// by generate_converters_v2.go which creates actual field mappings instead of stubs

// generateValidationMethod generates a single validation method
func generateValidationMethod(entityName, operation string, config ValidationConfig, entityDef EntityDef) string {
	var buf bytes.Buffer

	// Determine method name and input type
	var methodName, inputType string
	switch operation {
	case "create":
		methodName = fmt.Sprintf("Validate%sCreate", entityName)
		inputType = entityDef.CreateInput
	case "update":
		methodName = fmt.Sprintf("validate%sUpdate", entityName)
		inputType = entityDef.UpdateInput
	default:
		return fmt.Sprintf("// Unknown operation: %s\n", operation)
	}

	if inputType == "" {
		return fmt.Sprintf("// No input type for %s %s\n", entityName, operation)
	}

	// Determine parameter name (lowercase entity + operation suffix)
	paramName := strings.ToLower(operation)

	buf.WriteString(fmt.Sprintf(`// %s validates %s %s data
func (a *%sAdapter) %s(%s *types.%s) error {
`, methodName, strings.ToLower(entityName), operation, entityName, methodName, paramName, inputType))

	// Nil check
	if config.NilError != "" {
		buf.WriteString(fmt.Sprintf(`	if %s == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "%s", "%s", nil, nil)
	}
`, paramName, config.NilError, paramName))
	}

	// Required fields checks
	for _, req := range config.RequiredFields {
		switch req.CheckType {
		case "time":
			buf.WriteString(fmt.Sprintf(`	if %s.%s.IsZero() {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "%s", "%s", nil, nil)
	}
`, paramName, req.Field, req.Error, strings.ToLower(req.Field)))
		case "slice":
			buf.WriteString(fmt.Sprintf(`	if len(%s.%s) == 0 {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "%s", "%s", nil, nil)
	}
`, paramName, req.Field, req.Error, strings.ToLower(req.Field)))
		case "pointer":
			// For pointer fields, check nil or empty dereferenced value
			buf.WriteString(fmt.Sprintf(`	if %s.%s == nil || *%s.%s == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "%s", "%s", nil, nil)
	}
`, paramName, req.Field, paramName, req.Field, req.Error, strings.ToLower(req.Field)))
		default:
			// Default to empty string check
			buf.WriteString(fmt.Sprintf(`	if %s.%s == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "%s", "%s", nil, nil)
	}
`, paramName, req.Field, req.Error, strings.ToLower(req.Field)))
		}
	}

	// At least one field check
	if len(config.AtLeastOneOf) > 0 {
		var conditions []string
		for _, field := range config.AtLeastOneOf {
			// Handle pointer fields (check for nil) vs value fields
			conditions = append(conditions, fmt.Sprintf("%s.%s == nil", paramName, field))
		}

		// For slice fields, check length instead
		for i, field := range config.AtLeastOneOf {
			switch field {
			case "Accounts", "Users", "QoSList", "Groups", "Flags", "Features":
				// These are slice fields - check len()
				conditions[i] = fmt.Sprintf("len(%s.%s) == 0", paramName, field)
			}
		}

		errorMsg := config.AtLeastOneError
		if errorMsg == "" {
			errorMsg = "at least one field must be provided for update"
		}

		buf.WriteString(fmt.Sprintf(`	if %s {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "%s", "%s", nil, nil)
	}
`, strings.Join(conditions, " && "), errorMsg, paramName))
	}

	buf.WriteString(`	return nil
}
`)

	return buf.String()
}

func generateEntityAdapter(version, entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	var buf bytes.Buffer

	// Determine if we need strconv import (for int32 identifiers or Association update)
	needsStrconv := entityDef.IdentifierType == "int32" || entityName == "Association"

	// Write header with conditional imports
	buf.WriteString(fmt.Sprintf(`// Code generated by generate_adapters.go. DO NOT EDIT.
// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package %s

import (
	"context"
	"fmt"
`, version))

	if needsStrconv {
		buf.WriteString(`	"strconv"
`)
	}

	buf.WriteString(fmt.Sprintf(`
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/%s"
	"github.com/jontk/slurm-client/internal/common"
	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/pkg/errors"
)

`, version))

	// Generate adapter struct
	buf.WriteString(fmt.Sprintf(`// %sAdapter implements the %s adapter interface for %s
type %sAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// New%sAdapter creates a new %s adapter
func New%sAdapter(client *api.ClientWithResponses) *%sAdapter {
	return &%sAdapter{
		BaseManager: adapterbase.NewBaseManager("%s", "%s"),
		client:      client,
	}
}

`, entityName, entityName, version, entityName, entityName, entityName,
		entityName, entityName, entityName,
		strings.Replace(version, "_", ".", -1), entityName))

	// Generate methods based on entity definition
	for _, method := range entityDef.Methods {
		methodCode, err := generateMethod(method, entityName, entityDef, versionDef, apiMethods)
		if err != nil {
			return "", fmt.Errorf("generating method %s: %w", method, err)
		}
		buf.WriteString(methodCode)
		buf.WriteString("\n")
	}

	// Format the code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), nil // Return unformatted on error
	}

	return string(formatted), nil
}

func generateMethod(method, entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	// Special handling for Association entity (uses query params instead of path params)
	if entityName == "Association" {
		switch method {
		case "get":
			return generateGetMethodQueryParams(entityName, entityDef, versionDef, apiMethods)
		case "delete":
			return generateDeleteMethodQueryParams(entityName, entityDef, versionDef, apiMethods)
		}
	}

	// Special handling for Partition entity (Create/Update/Delete not supported by API)
	if entityName == "Partition" {
		switch method {
		case "create":
			return generateUnsupportedMethod("create", entityName,
				"Create(ctx context.Context, partition *types.PartitionCreate) (*types.PartitionCreateResponse, error)",
				"nil, "), nil
		case "update":
			return generateUnsupportedMethod("update", entityName,
				"Update(ctx context.Context, partitionName string, update *types.PartitionUpdate) error",
				""), nil
		case "delete":
			return generateUnsupportedMethod("delete", entityName,
				"Delete(ctx context.Context, partitionName string) error",
				""), nil
		}
	}

	// Special handling for Reservation entity (uses different request body structure)
	if entityName == "Reservation" {
		switch method {
		case "create":
			return generateReservationCreateMethod(entityDef, versionDef, apiMethods)
		case "update":
			return generateReservationUpdateMethod(entityDef, versionDef, apiMethods)
		}
	}

	switch method {
	case "list":
		return generateListMethod(entityName, entityDef, versionDef, apiMethods)
	case "get":
		return generateGetMethod(entityName, entityDef, versionDef, apiMethods)
	case "create":
		return generateCreateMethod(entityName, entityDef, versionDef, apiMethods)
	case "update":
		return generateUpdateMethod(entityName, entityDef, versionDef, apiMethods)
	case "delete":
		return generateDeleteMethod(entityName, entityDef, versionDef, apiMethods)
	case "submit":
		return generateSubmitMethod(entityName, entityDef, versionDef, apiMethods)
	case "cancel":
		return generateCancelMethod(entityName, entityDef, versionDef, apiMethods)
	case "hold":
		return generateHoldMethod(entityName, entityDef, versionDef, apiMethods)
	case "release":
		return generateReleaseMethod(entityName, entityDef, versionDef, apiMethods)
	case "signal":
		return generateSignalMethod(entityName, entityDef, versionDef, apiMethods)
	case "notify":
		return generateNotifyMethod(entityName, entityDef, versionDef, apiMethods)
	case "requeue":
		return generateRequeueMethod(entityName, entityDef, versionDef, apiMethods)
	case "watch":
		return generateWatchMethod(entityName, entityDef, versionDef, apiMethods)
	case "allocate":
		return generateAllocateMethod(entityName, entityDef, versionDef, apiMethods)
	case "drain":
		return generateDrainMethod(entityName, entityDef, versionDef, apiMethods)
	case "resume":
		return generateResumeMethod(entityName, entityDef, versionDef, apiMethods)
	case "create_association":
		return generateCreateAssociationMethod(entityName, entityDef, versionDef, apiMethods)
	default:
		// Skip methods not yet implemented
		return fmt.Sprintf("// %s method not yet generated\n", method), nil
	}
}

var listMethodTemplate = template.Must(template.New("list").Parse(`
// List retrieves a list of {{.EntityLower}}s with optional filtering
func (a *{{.Entity}}Adapter) List(ctx context.Context, opts *types.{{.ListOptions}}) (*types.{{.ListResult}}, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &api.{{.ListParams}}{}

	// Call the API
	resp, err := a.client.{{.APIMethod}}(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.{{.APIPrefix}}OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "{{.Version}}"); err != nil {
		return nil, err
	}

	// Check for nil response
	if err := a.CheckNilResponse(resp.JSON200, "List {{.Entity}}s"); err != nil {
		return nil, err
	}

	// Convert response to common types
	items := make([]types.{{.Entity}}, 0, len(resp.JSON200.{{.ResponseList}}))
	for _, apiItem := range resp.JSON200.{{.ResponseList}} {
		item := a.convertAPI{{.Entity}}ToCommon(apiItem)
		items = append(items, *item)
	}
{{if eq .Entity "Job"}}
	// Apply filtering before pagination
	{{.EntityLower}}BaseManager := adapterbase.New{{.Entity}}BaseManager("{{.Version}}")
	items = {{.EntityLower}}BaseManager.Filter{{.Entity}}List(items, opts)
{{end}}
	// Apply pagination
	listOpts := adapterbase.ListOptions{}
	if opts != nil {
		listOpts.Limit = opts.Limit
		listOpts.Offset = opts.Offset
	}

	start := listOpts.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(items) {
		return &types.{{.ListResult}}{
			{{.ListField}}: []types.{{.Entity}}{},
			Total: len(items),
		}, nil
	}

	end := len(items)
	if listOpts.Limit > 0 {
		end = start + listOpts.Limit
		if end > len(items) {
			end = len(items)
		}
	}

	return &types.{{.ListResult}}{
		{{.ListField}}: items[start:end],
		Total: len(items),
	}, nil
}
`))

func generateListMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	// Determine list field name (default is Entity+"s")
	listField := entityDef.ListField
	if listField == "" {
		listField = entityName + "s"
	}

	data := map[string]string{
		"Entity":       entityName,
		"EntityLower":  strings.ToLower(entityName),
		"ListOptions":  entityDef.ListOptions,
		"ListResult":   entityDef.ListResult,
		"ListField":    listField,
		"ListParams":   apiMethods.ListParams,
		"APIMethod":    apiMethods.List,
		"APIPrefix":    versionDef.APIPrefix,
		"Version":      strings.Replace(versionDef.APIPackage, "_", ".", -1),
		"ResponseList": apiMethods.ResponseList,
	}

	var buf bytes.Buffer
	if err := listMethodTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var getMethodTemplate = template.Must(template.New("get").Parse(`
// Get retrieves a specific {{.EntityLower}} by {{.Identifier}}
func (a *{{.Entity}}Adapter) Get(ctx context.Context, {{.Identifier}} {{.IdentifierType}}) (*types.{{.Entity}}, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
{{if eq .IdentifierType "int32"}}	if err := a.ValidateResourceID({{.Identifier}}, "{{.Identifier}}"); err != nil {
		return nil, err
	}
{{else}}	if err := a.ValidateResourceName({{.Identifier}}, "{{.Identifier}}"); err != nil {
		return nil, err
	}
{{end}}	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
{{if .GetParams}}
	// Prepare parameters
	params := &api.{{.GetParams}}{}

	// Call the API
{{if eq .IdentifierType "int32"}}	resp, err := a.client.{{.APIMethod}}(ctx, strconv.Itoa(int({{.Identifier}})), params)
{{else}}	resp, err := a.client.{{.APIMethod}}(ctx, {{.Identifier}}, params)
{{end}}{{else}}
	// Call the API (no params)
{{if eq .IdentifierType "int32"}}	resp, err := a.client.{{.APIMethod}}(ctx, strconv.Itoa(int({{.Identifier}})))
{{else}}	resp, err := a.client.{{.APIMethod}}(ctx, {{.Identifier}})
{{end}}{{end}}	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.{{.APIPrefix}}OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "{{.Version}}"); err != nil {
		return nil, err
	}

	// Check for nil response
	if err := a.CheckNilResponse(resp.JSON200, "Get {{.Entity}}"); err != nil {
		return nil, err
	}

	// Check if entity exists
	if len(resp.JSON200.{{.ResponseList}}) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("{{.Entity}} {{.FormatVerb}} not found", {{.Identifier}}))
	}

	// Convert and return
	return a.convertAPI{{.Entity}}ToCommon(resp.JSON200.{{.ResponseList}}[0]), nil
}
`))

func generateGetMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	// Use %d for int types, %s for strings
	formatVerb := "%s"
	if entityDef.IdentifierType == "int32" || entityDef.IdentifierType == "int64" || entityDef.IdentifierType == "int" {
		formatVerb = "%d"
	}

	data := map[string]string{
		"Entity":         entityName,
		"EntityLower":    strings.ToLower(entityName),
		"Identifier":     entityDef.Identifier,
		"IdentifierType": entityDef.IdentifierType,
		"GetParams":      apiMethods.GetParams,
		"APIMethod":      apiMethods.Get,
		"APIPrefix":      versionDef.APIPrefix,
		"Version":        strings.Replace(versionDef.APIPackage, "_", ".", -1),
		"ResponseList":   apiMethods.ResponseList,
		"FormatVerb":     formatVerb,
	}

	var buf bytes.Buffer
	if err := getMethodTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var deleteMethodTemplateNoParams = template.Must(template.New("delete_no_params").Parse(`
// Delete removes a {{.EntityLower}} by {{.Identifier}}
func (a *{{.Entity}}Adapter) Delete(ctx context.Context, {{.Identifier}} {{.IdentifierType}}) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
{{if eq .IdentifierType "int32"}}	if err := a.ValidateResourceID({{.Identifier}}, "{{.Identifier}}"); err != nil {
		return err
	}
{{else}}	if err := a.ValidateResourceName({{.Identifier}}, "{{.Identifier}}"); err != nil {
		return err
	}
{{end}}	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the API
{{if eq .IdentifierType "int32"}}	resp, err := a.client.{{.APIMethod}}(ctx, strconv.Itoa(int({{.Identifier}})))
{{else}}	resp, err := a.client.{{.APIMethod}}(ctx, {{.Identifier}})
{{end}}	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.{{.APIPrefix}}OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "{{.Version}}")
}
`))

var deleteMethodTemplateWithParams = template.Must(template.New("delete_with_params").Parse(`
// Delete removes a {{.EntityLower}} by {{.Identifier}}
func (a *{{.Entity}}Adapter) Delete(ctx context.Context, {{.Identifier}} {{.IdentifierType}}) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
{{if eq .IdentifierType "int32"}}	if err := a.ValidateResourceID({{.Identifier}}, "{{.Identifier}}"); err != nil {
		return err
	}
{{else}}	if err := a.ValidateResourceName({{.Identifier}}, "{{.Identifier}}"); err != nil {
		return err
	}
{{end}}	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the API (params is nil for default behavior)
{{if eq .IdentifierType "int32"}}	resp, err := a.client.{{.APIMethod}}(ctx, strconv.Itoa(int({{.Identifier}})), nil)
{{else}}	resp, err := a.client.{{.APIMethod}}(ctx, {{.Identifier}}, nil)
{{end}}	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.{{.APIPrefix}}OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "{{.Version}}")
}
`))

func generateDeleteMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	data := map[string]string{
		"Entity":         entityName,
		"EntityLower":    strings.ToLower(entityName),
		"Identifier":     entityDef.Identifier,
		"IdentifierType": entityDef.IdentifierType,
		"APIMethod":      apiMethods.Delete,
		"APIPrefix":      versionDef.APIPrefix,
		"Version":        strings.Replace(versionDef.APIPackage, "_", ".", -1),
	}

	var buf bytes.Buffer
	// Select template based on whether API needs params
	tmpl := deleteMethodTemplateNoParams
	if apiMethods.DeleteNeedsParams {
		tmpl = deleteMethodTemplateWithParams
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var createMethodTemplateNoParams = template.Must(template.New("create_no_params").Parse(`
// Create creates a new {{.EntityLower}}
func (a *{{.Entity}}Adapter) Create(ctx context.Context, input *types.{{.CreateInput}}) (*types.{{.CreateResult}}, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if input == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "input is required", "input", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API type
	apiInput := a.convertCommon{{.Entity}}CreateToAPI(input)

	// Create request body with the entity wrapped in a slice
	reqBody := api.{{.RequestBodyType}}{
		{{.RequestListField}}: []api.{{.APIType}}{*apiInput},
	}

	// Call the API
	resp, err := a.client.{{.APIMethod}}(ctx, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.{{.APIPrefix}}OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "{{.Version}}"); err != nil {
		return nil, err
	}

	return &types.{{.CreateResult}}{}, nil
}
`))

var createMethodTemplateWithParams = template.Must(template.New("create_with_params").Parse(`
// Create creates a new {{.EntityLower}}
func (a *{{.Entity}}Adapter) Create(ctx context.Context, input *types.{{.CreateInput}}) (*types.{{.CreateResult}}, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if input == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "input is required", "input", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API type
	apiInput := a.convertCommon{{.Entity}}CreateToAPI(input)

	// Create request body with the entity wrapped in a slice
	reqBody := api.{{.RequestBodyType}}{
		{{.RequestListField}}: []api.{{.APIType}}{*apiInput},
	}

	// Call the API (params is nil for default behavior)
	resp, err := a.client.{{.APIMethod}}(ctx, nil, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.{{.APIPrefix}}OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "{{.Version}}"); err != nil {
		return nil, err
	}

	return &types.{{.CreateResult}}{}, nil
}
`))

func generateCreateMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityDef.CreateInput == "" || apiMethods.Create == "" {
		return fmt.Sprintf("// Create method not configured for %s\n", entityName), nil
	}

	// Check for request body type - required for proper API calls
	if apiMethods.RequestBodyType == "" || apiMethods.RequestListField == "" {
		return fmt.Sprintf("// Create method: request_body_type or request_list_field not configured for %s\n", entityName), nil
	}

	createResult := entityDef.CreateResult
	if createResult == "" {
		createResult = entityName + "CreateResponse"
	}

	data := map[string]string{
		"Entity":           entityName,
		"EntityLower":      strings.ToLower(entityName),
		"CreateInput":      entityDef.CreateInput,
		"CreateResult":     createResult,
		"APIMethod":        apiMethods.Create,
		"APIPrefix":        versionDef.APIPrefix,
		"Version":          strings.Replace(versionDef.APIPackage, "_", ".", -1),
		"APIType":          apiMethods.APIType,
		"RequestBodyType":  apiMethods.RequestBodyType,
		"RequestListField": apiMethods.RequestListField,
	}

	var buf bytes.Buffer
	// Select template based on whether API needs params
	tmpl := createMethodTemplateNoParams
	if apiMethods.CreateNeedsParams {
		tmpl = createMethodTemplateWithParams
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var updateMethodTemplate = template.Must(template.New("update").Parse(`
// Update updates an existing {{.EntityLower}}
func (a *{{.Entity}}Adapter) Update(ctx context.Context, {{.Identifier}} {{.IdentifierType}}, update *types.{{.UpdateInput}}) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName({{.Identifier}}, "{{.Identifier}}"); err != nil {
		return err
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update is required", "update", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert to API type (SLURM API handles partial updates)
	apiInput := a.convertCommon{{.Entity}}UpdateToAPI(update)
	// Set the identifier for the update
	apiInput.{{.IdentifierField}} = {{.Identifier}}

	// Create request body with the entity wrapped in a slice
	reqBody := api.{{.RequestBodyType}}{
		{{.RequestListField}}: []api.{{.APIType}}{*apiInput},
	}

	// Call the API (POST is used for updates in SLURM)
	resp, err := a.client.{{.APIMethod}}(ctx, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.{{.APIPrefix}}OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "{{.Version}}")
}
`))

func generateUpdateMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityDef.UpdateInput == "" {
		return fmt.Sprintf("// Update method not configured for %s\n", entityName), nil
	}

	// Special handling for entities with different update patterns
	if entityName == "Job" {
		return generateJobUpdateMethod(entityDef, versionDef, apiMethods)
	}
	if entityName == "Node" {
		return generateNodeUpdateMethod(entityDef, versionDef, apiMethods)
	}
	if entityName == "Association" {
		return generateAssociationUpdateMethod(entityDef, versionDef, apiMethods)
	}
	if entityName == "QoS" {
		return generateQoSUpdateMethod(entityDef, versionDef, apiMethods)
	}

	// Update typically uses the same Create endpoint in Slurm
	apiMethod := apiMethods.Create
	if apiMethod == "" {
		return fmt.Sprintf("// Update method: no API method configured for %s\n", entityName), nil
	}

	// Check for request body type - required for proper API calls
	if apiMethods.RequestBodyType == "" || apiMethods.RequestListField == "" {
		return fmt.Sprintf("// Update method: request_body_type or request_list_field not configured for %s\n", entityName), nil
	}

	// Map identifier to API field name
	identifierField := identifierToAPIField(entityDef.Identifier)

	data := map[string]string{
		"Entity":           entityName,
		"EntityLower":      strings.ToLower(entityName),
		"Identifier":       entityDef.Identifier,
		"IdentifierType":   entityDef.IdentifierType,
		"IdentifierField":  identifierField,
		"UpdateInput":      entityDef.UpdateInput,
		"APIMethod":        apiMethod,
		"APIPrefix":        versionDef.APIPrefix,
		"Version":          strings.Replace(versionDef.APIPackage, "_", ".", -1),
		"APIType":          apiMethods.APIType,
		"RequestBodyType":  apiMethods.RequestBodyType,
		"RequestListField": apiMethods.RequestListField,
	}

	var buf bytes.Buffer
	if err := updateMethodTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// generateJobUpdateMethod generates the Job-specific Update method
func generateJobUpdateMethod(entityDef EntityDef, versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {
	apiMethod := apiMethods.Update
	if apiMethod == "" {
		apiMethod = "SlurmV" + versionDef.APIPrefix[1:] + "PostJobWithResponse"
	}

	return fmt.Sprintf(`
// Update updates an existing job
func (a *JobAdapter) Update(ctx context.Context, jobID int32, update *types.JobUpdate) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceID(jobID, "jobID"); err != nil {
		return err
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update is required", "update", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert to API type
	reqBody := a.convertCommonJobUpdateToAPIRequestBody(update)

	// Call the API with job ID in path
	resp, err := a.client.%s(ctx, strconv.Itoa(int(jobID)), reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "%s")
}
`, apiMethod, versionDef.APIPrefix, strings.Replace(versionDef.APIPackage, "_", ".", -1)), nil
}

// generateNodeUpdateMethod generates the Node-specific Update method
func generateNodeUpdateMethod(entityDef EntityDef, versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {
	apiMethod := apiMethods.Update
	if apiMethod == "" {
		apiMethod = "SlurmV" + versionDef.APIPrefix[1:] + "PostNodeWithResponse"
	}

	return fmt.Sprintf(`
// Update updates an existing node
func (a *NodeAdapter) Update(ctx context.Context, nodeName string, update *types.NodeUpdate) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(nodeName, "nodeName"); err != nil {
		return err
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update is required", "update", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert to API type
	reqBody := a.convertCommonNodeUpdateToAPIRequestBody(update)

	// Call the API with node name in path
	resp, err := a.client.%s(ctx, nodeName, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "%s")
}
`, apiMethod, versionDef.APIPrefix, strings.Replace(versionDef.APIPackage, "_", ".", -1)), nil
}

// generateAssociationUpdateMethod generates the Association-specific Update method
// Association has a string identifier but the API expects *int32 for the Id field
func generateAssociationUpdateMethod(entityDef EntityDef, versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {
	apiMethod := apiMethods.Create
	if apiMethod == "" {
		apiMethod = "SlurmdbV" + versionDef.APIPrefix[1:] + "PostAssociationsWithResponse"
	}

	return fmt.Sprintf(`
// Update updates an existing association
func (a *AssociationAdapter) Update(ctx context.Context, associationID string, update *types.AssociationUpdate) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(associationID, "associationID"); err != nil {
		return err
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update is required", "update", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert association ID from string to int32
	id, err := strconv.ParseInt(associationID, 10, 32)
	if err != nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "invalid association ID format", "associationID", nil, err)
	}
	idInt32 := int32(id)

	// Convert to API type
	apiInput := a.convertCommonAssociationUpdateToAPI(update)
	// Set the ID for the update (API uses *int32)
	apiInput.Id = &idInt32

	// Create request body with the association wrapped in a slice
	reqBody := api.SlurmdbV%sPostAssociationsJSONRequestBody{
		Associations: []api.V%sAssoc{*apiInput},
	}

	// Call the API
	resp, err := a.client.%s(ctx, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "%s")
}
`, versionDef.APIPrefix[1:], versionDef.APIPrefix[1:], apiMethod, versionDef.APIPrefix, strings.Replace(versionDef.APIPackage, "_", ".", -1)), nil
}

// generateQoSUpdateMethod generates the QoS-specific Update method
// QoS uses *string for Name field, so needs pointer assignment
func generateQoSUpdateMethod(entityDef EntityDef, versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {
	apiMethod := apiMethods.Create
	if apiMethod == "" {
		apiMethod = "SlurmdbV" + versionDef.APIPrefix[1:] + "PostQosWithResponse"
	}

	return fmt.Sprintf(`
// Update updates an existing QoS
func (a *QoSAdapter) Update(ctx context.Context, qosName string, update *types.QoSUpdate) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(qosName, "qosName"); err != nil {
		return err
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update is required", "update", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert to API type
	apiInput := a.convertCommonQoSUpdateToAPI(update)
	// Set the name for the update (API uses *string)
	apiInput.Name = &qosName

	// Create request body with the QoS wrapped in a slice
	reqBody := api.SlurmdbV%sPostQosJSONRequestBody{
		Qos: []api.V%sQos{*apiInput},
	}

	// Call the API (params is nil for default behavior)
	resp, err := a.client.%s(ctx, nil, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "%s")
}
`, versionDef.APIPrefix[1:], versionDef.APIPrefix[1:], apiMethod, versionDef.APIPrefix, strings.Replace(versionDef.APIPackage, "_", ".", -1)), nil
}

// identifierToAPIField maps common identifier names to their API field names
func identifierToAPIField(identifier string) string {
	mapping := map[string]string{
		"accountName":     "Name",
		"userName":        "Name",
		"qosName":         "Name",
		"nodeName":        "Name",
		"partitionName":   "Name",
		"reservationName": "Name",
		"clusterName":     "Name",
		"associationID":   "Id",
		"jobID":           "JobId",
	}
	if field, ok := mapping[identifier]; ok {
		return field
	}
	return "Name" // Default fallback
}

// Reservation-specific methods (uses V0044ReservationDescMsg for input, not V0044ReservationInfo)

// generateReservationCreateMethod generates the Reservation-specific Create method
func generateReservationCreateMethod(entityDef EntityDef, versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {
	// Check if Create API is configured
	if apiMethods.Create == "" {
		// Generate stub implementation that returns unsupported error
		return fmt.Sprintf(`
// Create creates a new reservation
func (a *ReservationAdapter) Create(ctx context.Context, input *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
	return nil, errors.NewClientError(errors.ErrorCodeUnsupportedOperation, "Create is not supported for Reservation in %s", "operation not supported")
}
`, strings.Replace(versionDef.APIPackage, "_", ".", -1)), nil
	}
	apiMethod := apiMethods.Create

	return fmt.Sprintf(`
// Create creates a new reservation
func (a *ReservationAdapter) Create(ctx context.Context, input *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if input == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "input is required", "input", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API type (Reservation uses V%sReservationDescMsg for input)
	apiInput := a.convertCommonReservationCreateToAPI(input)

	// Create request body - Reservation uses a different structure
	list := api.V%sReservationDescMsgList{*apiInput}
	reqBody := api.SlurmV%sPostReservationsJSONRequestBody{
		Reservations: &list,
	}

	// Call the API
	resp, err := a.client.%s(ctx, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "%s"); err != nil {
		return nil, err
	}

	return &types.ReservationCreateResponse{}, nil
}
`, versionDef.APIPrefix[1:], versionDef.APIPrefix[1:], versionDef.APIPrefix[1:], apiMethod, versionDef.APIPrefix, strings.Replace(versionDef.APIPackage, "_", ".", -1)), nil
}

// generateReservationUpdateMethod generates the Reservation-specific Update method
func generateReservationUpdateMethod(entityDef EntityDef, versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {
	// Check if Update API is configured (Reservation uses Create API method for updates)
	if apiMethods.Create == "" && apiMethods.Update == "" {
		// Generate stub implementation that returns unsupported error
		return fmt.Sprintf(`
// Update updates an existing reservation
func (a *ReservationAdapter) Update(ctx context.Context, reservationName string, update *types.ReservationUpdate) error {
	return errors.NewClientError(errors.ErrorCodeUnsupportedOperation, "Update is not supported for Reservation in %s", "operation not supported")
}
`, strings.Replace(versionDef.APIPackage, "_", ".", -1)), nil
	}
	apiMethod := apiMethods.Update
	if apiMethod == "" {
		apiMethod = apiMethods.Create // Fallback to Create if Update isn't defined
	}

	return fmt.Sprintf(`
// Update updates an existing reservation
func (a *ReservationAdapter) Update(ctx context.Context, reservationName string, update *types.ReservationUpdate) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(reservationName, "reservationName"); err != nil {
		return err
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update is required", "update", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert to API type (Reservation uses V%sReservationDescMsg for input)
	apiInput := a.convertCommonReservationUpdateToAPI(update)
	// Set the name for the update
	apiInput.Name = &reservationName

	// Create request body - Reservation uses a different structure
	list := api.V%sReservationDescMsgList{*apiInput}
	reqBody := api.SlurmV%sPostReservationsJSONRequestBody{
		Reservations: &list,
	}

	// Call the API
	resp, err := a.client.%s(ctx, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "%s")
}
`, versionDef.APIPrefix[1:], versionDef.APIPrefix[1:], versionDef.APIPrefix[1:], apiMethod, versionDef.APIPrefix, strings.Replace(versionDef.APIPackage, "_", ".", -1)), nil
}

// Job-specific method templates

var submitMethodTemplate = template.Must(template.New("submit").Parse(`
// Submit submits a new job
func (a *JobAdapter) Submit(ctx context.Context, job *types.JobCreate) (*types.JobSubmitResponse, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if job == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "job is required", "job", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API type
	apiJob := a.convertCommonJobCreateToAPI(job)

	// Call the API
	resp, err := a.client.{{.APIMethod}}(ctx, apiJob)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Handle response errors — check JSONDefault on non-200 responses so that
	// the actual SLURM rejection reason is surfaced instead of a generic HTTP error.
	var apiErrors *api.{{.APIPrefix}}OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "{{.Version}}"); err != nil {
		return nil, err
	}

	// Extract job ID from response
	return a.convertAPIJobSubmitResponseToCommon(resp.JSON200), nil
}
`))

func generateSubmitMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Job" || apiMethods.List == "" {
		return "// Submit method only available for Job\n", nil
	}

	// Get submit method name from config, or use default
	submitMethod := apiMethods.Submit
	if submitMethod == "" {
		submitMethod = "SlurmV" + versionDef.APIPrefix[1:] + "PostJobSubmitWithResponse"
	}

	data := map[string]string{
		"APIMethod": submitMethod,
		"APIPrefix": versionDef.APIPrefix,
		"Version":   strings.Replace(versionDef.APIPackage, "_", ".", -1),
	}

	var buf bytes.Buffer
	if err := submitMethodTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var cancelMethodTemplate = template.Must(template.New("cancel").Parse(`
// Cancel cancels a job
func (a *JobAdapter) Cancel(ctx context.Context, jobID int32, opts *types.JobCancelRequest) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceID(jobID, "jobID"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the delete/cancel API
	resp, err := a.client.{{.APIMethod}}(ctx, strconv.Itoa(int(jobID)), nil)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response
	if resp.StatusCode() == 204 || resp.StatusCode() == 200 {
		return nil
	}

	var apiErrors *api.{{.APIPrefix}}OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "{{.Version}}")
}
`))

func generateCancelMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Job" {
		return "// Cancel method only available for Job\n", nil
	}

	cancelMethod := "SlurmV" + versionDef.APIPrefix[1:] + "DeleteJobWithResponse"

	data := map[string]string{
		"APIMethod": cancelMethod,
		"APIPrefix": versionDef.APIPrefix,
		"Version":   strings.Replace(versionDef.APIPackage, "_", ".", -1),
	}

	var buf bytes.Buffer
	if err := cancelMethodTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func generateHoldMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Job" {
		return "// Hold method only available for Job\n", nil
	}

	return `
// Hold places a job on hold
func (a *JobAdapter) Hold(ctx context.Context, req *types.JobHoldRequest) error {
	return a.holdJobImpl(ctx, req)
}
`, nil
}

func generateReleaseMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Job" {
		return "// Release method only available for Job\n", nil
	}

	// Release is not in the interface - Hold handles both hold and release
	return "// Release is handled via Hold method with release flag\n", nil
}

func generateSignalMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Job" {
		return "// Signal method only available for Job\n", nil
	}

	return `
// Signal sends a signal to a job
func (a *JobAdapter) Signal(ctx context.Context, req *types.JobSignalRequest) error {
	return a.signalJobImpl(ctx, req)
}
`, nil
}

func generateNotifyMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Job" {
		return "// Notify method only available for Job\n", nil
	}

	return `
// Notify sends a notification to a job
func (a *JobAdapter) Notify(ctx context.Context, req *types.JobNotifyRequest) error {
	return a.notifyJobImpl(ctx, req)
}
`, nil
}

func generateRequeueMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Job" {
		return "// Requeue method only available for Job\n", nil
	}

	return `
// Requeue requeues a job
func (a *JobAdapter) Requeue(ctx context.Context, jobID int32) error {
	return a.requeueJobImpl(ctx, jobID)
}
`, nil
}

func generateWatchMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Job" && entityName != "Node" {
		return "// Watch method only available for Job and Node\n", nil
	}

	watchOpts := entityDef.WatchOptions
	if watchOpts == "" {
		watchOpts = entityName + "WatchOptions"
	}

	eventType := entityName + "WatchEvent"

	return fmt.Sprintf(`
// Watch watches for %s events
func (a *%sAdapter) Watch(ctx context.Context, opts *types.%s) (<-chan types.%s, error) {
	return a.watch%ssImpl(ctx, opts)
}
`, strings.ToLower(entityName), entityName, watchOpts, eventType, entityName), nil
}

func generateAllocateMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Job" {
		return "// Allocate method only available for Job\n", nil
	}

	return `
// Allocate allocates resources for a job
func (a *JobAdapter) Allocate(ctx context.Context, req *types.JobAllocateRequest) (*types.JobAllocateResponse, error) {
	return a.allocateJobImpl(ctx, req)
}
`, nil
}

func generateDrainMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Node" {
		return "// Drain method only available for Node\n", nil
	}

	return `
// Drain drains a node
func (a *NodeAdapter) Drain(ctx context.Context, nodeName string, reason string) error {
	return a.drainNodeImpl(ctx, nodeName, reason)
}
`, nil
}

func generateResumeMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Node" {
		return "// Resume method only available for Node\n", nil
	}

	return `
// Resume resumes a drained node
func (a *NodeAdapter) Resume(ctx context.Context, nodeName string) error {
	return a.resumeNodeImpl(ctx, nodeName)
}
`, nil
}

func generateCreateAssociationMethod(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	if entityName != "Account" && entityName != "User" {
		return "// CreateAssociation method only available for Account and User\n", nil
	}

	requestType := entityDef.AssociationRequest
	if requestType == "" {
		requestType = entityName + "AssociationRequest"
	}

	resultType := entityDef.AssociationResult
	if resultType == "" {
		resultType = "AssociationCreateResponse"
	}

	return fmt.Sprintf(`
// CreateAssociation creates an association for this %s
func (a *%sAdapter) CreateAssociation(ctx context.Context, req *types.%s) (*types.%s, error) {
	return a.createAssociationImpl(ctx, req)
}
`, strings.ToLower(entityName), entityName, requestType, resultType), nil
}

// Override Get method for entities that use query params instead of path params (like Association)
func generateGetMethodQueryParams(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	data := map[string]string{
		"Entity":         entityName,
		"EntityLower":    strings.ToLower(entityName),
		"Identifier":     entityDef.Identifier,
		"IdentifierType": entityDef.IdentifierType,
		"GetParams":      apiMethods.GetParams,
		"APIMethod":      apiMethods.Get,
		"APIPrefix":      versionDef.APIPrefix,
		"Version":        strings.Replace(versionDef.APIPackage, "_", ".", -1),
		"ResponseList":   apiMethods.ResponseList,
	}

	return fmt.Sprintf(`
// Get retrieves a specific %s by %s
// Note: This API uses query parameters instead of path parameters
func (a *%sAdapter) Get(ctx context.Context, %s %s) (*types.%s, error) {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(%s, "%s"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters - Association uses query params with Id field
	params := &api.%s{
		Id: &%s,
	}

	// Call the API (identifier passed via query params, not path)
	resp, err := a.client.%s(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "%s"); err != nil {
		return nil, err
	}

	// Check for nil response
	if err := a.CheckNilResponse(resp.JSON200, "Get %s"); err != nil {
		return nil, err
	}

	// Check if entity exists
	if len(resp.JSON200.%s) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("%s %%s not found", %s))
	}

	// Convert and return
	return a.convertAPI%sToCommon(resp.JSON200.%s[0]), nil
}
`, data["EntityLower"], data["Identifier"],
		data["Entity"], data["Identifier"], data["IdentifierType"], data["Entity"],
		data["Identifier"], data["Identifier"],
		data["GetParams"], data["Identifier"],
		data["APIMethod"],
		data["APIPrefix"],
		data["Version"],
		data["Entity"],
		data["ResponseList"],
		data["Entity"], data["Identifier"],
		data["Entity"], data["ResponseList"]), nil
}

// Override Delete method for entities that use query params instead of path params
func generateDeleteMethodQueryParams(entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	data := map[string]string{
		"Entity":         entityName,
		"EntityLower":    strings.ToLower(entityName),
		"Identifier":     entityDef.Identifier,
		"IdentifierType": entityDef.IdentifierType,
		"APIMethod":      apiMethods.Delete,
		"APIPrefix":      versionDef.APIPrefix,
		"Version":        strings.Replace(versionDef.APIPackage, "_", ".", -1),
	}

	return fmt.Sprintf(`
// Delete removes a %s by %s
// Note: This API uses query parameters instead of path parameters
func (a *%sAdapter) Delete(ctx context.Context, %s %s) error {
	// Validate context
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(%s, "%s"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the API (uses query params, not path params)
	// Pass the ID as a query parameter
	resp, err := a.client.%s(ctx, &api.%s{Id: &%s})
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Handle response errors
	var apiErrors *api.%sOpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "%s")
}
`, data["EntityLower"], data["Identifier"],
		data["Entity"], data["Identifier"], data["IdentifierType"],
		data["Identifier"], data["Identifier"],
		data["APIMethod"], strings.Replace(data["APIMethod"], "WithResponse", "Params", 1), data["Identifier"],
		data["APIPrefix"],
		data["Version"]), nil
}

// Generate unsupported method stub
func generateUnsupportedMethod(method, entityName, signature, returnType string) string {
	return fmt.Sprintf(`
// %s is not supported in this API version
func (a *%sAdapter) %s {
	return %serrors.NewClientError(
		errors.ErrorCodeUnsupportedOperation,
		"%s %s not supported in this version",
		"Method not allowed (405)")
}
`, strings.Title(method), entityName, signature, returnType,
		strings.ToLower(entityName), method)
}

// generateEntityTests generates the test file for an entity
func generateEntityTests(version, entityName string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) (string, error) {

	var buf bytes.Buffer

	// Check if we need time import for validation tests
	needsTimeImport := false
	for _, validationConfig := range entityDef.Validation {
		for _, field := range validationConfig.RequiredFields {
			if field.CheckType == "time" {
				needsTimeImport = true
				break
			}
		}
		if needsTimeImport {
			break
		}
	}

	// Write header
	if needsTimeImport {
		buf.WriteString(fmt.Sprintf(`// Code generated by generate_adapters.go. DO NOT EDIT.
// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package %s

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/jontk/slurm-client/internal/openapi/%s"
	types "github.com/jontk/slurm-client/api"
)

`, version, version))
	} else {
		buf.WriteString(fmt.Sprintf(`// Code generated by generate_adapters.go. DO NOT EDIT.
// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package %s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/jontk/slurm-client/internal/openapi/%s"
	types "github.com/jontk/slurm-client/api"
)

`, version, version))
	}

	// Generate constructor test
	buf.WriteString(fmt.Sprintf(`func TestNew%sAdapter(t *testing.T) {
	client := &api.ClientWithResponses{}
	adapter := New%sAdapter(client)
	require.NotNil(t, adapter)
}

`, entityName, entityName))

	// Generate context validation test
	buf.WriteString(fmt.Sprintf(`func Test%sAdapter_ValidateContext(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	// Test nil context
	//lint:ignore SA1012 intentionally testing nil context validation
	err := adapter.ValidateContext(nil)
	assert.Error(t, err)

	// Test valid context
	err = adapter.ValidateContext(context.Background())
	assert.NoError(t, err)
}

`, entityName, entityName))

	// Generate method-specific tests
	for _, method := range entityDef.Methods {
		testCode := generateMethodTest(entityName, method, entityDef, versionDef, apiMethods)
		buf.WriteString(testCode)
	}

	// Generate converter round-trip tests
	converterTests := generateConverterRoundTripTests(entityName, entityDef, apiMethods)
	buf.WriteString(converterTests)

	// Generate validation tests if entity has validation config
	if len(entityDef.Validation) > 0 {
		validationTests := generateValidationTests(entityName, entityDef)
		buf.WriteString(validationTests)
	}

	// Format the code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), nil // Return unformatted on error
	}

	return string(formatted), nil
}

// generateMethodTest generates tests for a specific method
func generateMethodTest(entityName, method string, entityDef EntityDef,
	versionDef VersionDef, apiMethods EntityAPIMethods) string {

	var buf bytes.Buffer

	switch method {
	case "list":
		buf.WriteString(fmt.Sprintf(`func Test%sAdapter_List_NilContext(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	//lint:ignore SA1012 intentionally testing nil context
	result, err := adapter.List(nil, nil)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func Test%sAdapter_List_ClientNotInitialized(t *testing.T) {
	adapter := New%sAdapter(nil)

	result, err := adapter.List(context.Background(), nil)
	assert.Nil(t, result)
	assert.Error(t, err)
}

`, entityName, entityName, entityName, entityName))

	case "get":
		identifierType := entityDef.IdentifierType
		var emptyValue, testValue string
		if identifierType == "int32" || identifierType == "int64" {
			emptyValue = "0"
			testValue = "123"
		} else {
			emptyValue = `""`
			testValue = `"test-` + strings.ToLower(entityName) + `"`
		}

		buf.WriteString(fmt.Sprintf(`func Test%sAdapter_Get_NilContext(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	//lint:ignore SA1012 intentionally testing nil context
	result, err := adapter.Get(nil, %s)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func Test%sAdapter_Get_EmptyIdentifier(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	result, err := adapter.Get(context.Background(), %s)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func Test%sAdapter_Get_ClientNotInitialized(t *testing.T) {
	adapter := New%sAdapter(nil)

	result, err := adapter.Get(context.Background(), %s)
	assert.Nil(t, result)
	assert.Error(t, err)
}

`, entityName, entityName, testValue,
			entityName, entityName, emptyValue,
			entityName, entityName, testValue))

	case "create":
		createInput := entityDef.CreateInput
		if createInput == "" {
			createInput = entityName + "Create"
		}

		buf.WriteString(fmt.Sprintf(`func Test%sAdapter_Create_NilContext(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	//lint:ignore SA1012 intentionally testing nil context
	result, err := adapter.Create(nil, &types.%s{})
	assert.Nil(t, result)
	assert.Error(t, err)
}

func Test%sAdapter_Create_NilInput(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	result, err := adapter.Create(context.Background(), nil)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func Test%sAdapter_Create_ClientNotInitialized(t *testing.T) {
	adapter := New%sAdapter(nil)

	result, err := adapter.Create(context.Background(), &types.%s{})
	assert.Nil(t, result)
	assert.Error(t, err)
}

`, entityName, entityName, createInput,
			entityName, entityName,
			entityName, entityName, createInput))

	case "update":
		updateInput := entityDef.UpdateInput
		if updateInput == "" {
			updateInput = entityName + "Update"
		}
		identifierType := entityDef.IdentifierType
		var testValue string
		if identifierType == "int32" || identifierType == "int64" {
			testValue = "123"
		} else {
			testValue = `"test-` + strings.ToLower(entityName) + `"`
		}

		buf.WriteString(fmt.Sprintf(`func Test%sAdapter_Update_NilContext(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	//lint:ignore SA1012 intentionally testing nil context
	err := adapter.Update(nil, %s, &types.%s{})
	assert.Error(t, err)
}

func Test%sAdapter_Update_NilInput(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	err := adapter.Update(context.Background(), %s, nil)
	assert.Error(t, err)
}

func Test%sAdapter_Update_ClientNotInitialized(t *testing.T) {
	adapter := New%sAdapter(nil)

	err := adapter.Update(context.Background(), %s, &types.%s{})
	assert.Error(t, err)
}

`, entityName, entityName, testValue, updateInput,
			entityName, entityName, testValue,
			entityName, entityName, testValue, updateInput))

	case "delete":
		identifierType := entityDef.IdentifierType
		var emptyValue, testValue string
		if identifierType == "int32" || identifierType == "int64" {
			emptyValue = "0"
			testValue = "123"
		} else {
			emptyValue = `""`
			testValue = `"test-` + strings.ToLower(entityName) + `"`
		}

		buf.WriteString(fmt.Sprintf(`func Test%sAdapter_Delete_NilContext(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	//lint:ignore SA1012 intentionally testing nil context
	err := adapter.Delete(nil, %s)
	assert.Error(t, err)
}

func Test%sAdapter_Delete_EmptyIdentifier(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	err := adapter.Delete(context.Background(), %s)
	assert.Error(t, err)
}

func Test%sAdapter_Delete_ClientNotInitialized(t *testing.T) {
	adapter := New%sAdapter(nil)

	err := adapter.Delete(context.Background(), %s)
	assert.Error(t, err)
}

`, entityName, entityName, testValue,
			entityName, entityName, emptyValue,
			entityName, entityName, testValue))

	case "submit":
		if entityName == "Job" {
			buf.WriteString(`func TestJobAdapter_Submit_NilContext(t *testing.T) {
	adapter := NewJobAdapter(&api.ClientWithResponses{})

	//lint:ignore SA1012 intentionally testing nil context
	result, err := adapter.Submit(nil, &types.JobCreate{})
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestJobAdapter_Submit_NilInput(t *testing.T) {
	adapter := NewJobAdapter(&api.ClientWithResponses{})

	result, err := adapter.Submit(context.Background(), nil)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestJobAdapter_Submit_ClientNotInitialized(t *testing.T) {
	adapter := NewJobAdapter(nil)

	result, err := adapter.Submit(context.Background(), &types.JobCreate{})
	assert.Nil(t, result)
	assert.Error(t, err)
}

`)
		}

	case "cancel":
		if entityName == "Job" {
			buf.WriteString(`func TestJobAdapter_Cancel_NilContext(t *testing.T) {
	adapter := NewJobAdapter(&api.ClientWithResponses{})

	//lint:ignore SA1012 intentionally testing nil context
	err := adapter.Cancel(nil, 123, nil)
	assert.Error(t, err)
}

func TestJobAdapter_Cancel_InvalidJobID(t *testing.T) {
	adapter := NewJobAdapter(&api.ClientWithResponses{})

	err := adapter.Cancel(context.Background(), 0, nil)
	assert.Error(t, err)
}

func TestJobAdapter_Cancel_ClientNotInitialized(t *testing.T) {
	adapter := NewJobAdapter(nil)

	err := adapter.Cancel(context.Background(), 123, nil)
	assert.Error(t, err)
}

`)
		}
	}

	return buf.String()
}

// generateConverterRoundTripTests generates tests for converter round-trip conversions
func generateConverterRoundTripTests(entityName string, entityDef EntityDef, apiMethods EntityAPIMethods) string {
	var buf bytes.Buffer

	// Generate read converter test (API → Common)
	if apiMethods.APIType != "" {
		buf.WriteString(generateReadConverterTest(entityName, apiMethods))
	}

	// Generate write converter tests (Common → API)
	if entityDef.CreateInput != "" && apiMethods.Create != "" {
		buf.WriteString(generateWriteConverterTest(entityName, entityDef, apiMethods, "Create"))
	}

	// Only generate Update converter test if the entity has a standard Update converter
	// Job and Node use special convertCommon*UpdateToAPIRequestBody converters instead
	if entityDef.UpdateInput != "" && (apiMethods.Update != "" || apiMethods.Create != "") {
		// Skip Update converter tests for entities with custom request body converters
		if entityName != "Job" && entityName != "Node" {
			buf.WriteString(generateWriteConverterTest(entityName, entityDef, apiMethods, "Update"))
		}
	}

	return buf.String()
}

// generateReadConverterTest generates a test for API → Common converter
func generateReadConverterTest(entityName string, apiMethods EntityAPIMethods) string {
	var buf bytes.Buffer

	// Generate entity-specific test based on what fields we can test
	testCode, assertions := generateReadConverterTestCodeWithAssertions(entityName, apiMethods.APIType)

	buf.WriteString(fmt.Sprintf(`func Test%sAdapter_ReadConverter(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

%s

	// Convert API → Common
	commonObj := adapter.convertAPI%sToCommon(apiObj)
	require.NotNil(t, commonObj)

	// Verify field values were correctly converted
%s}

`, entityName, entityName, testCode, entityName, assertions))

	return buf.String()
}

// generateReadConverterTestCode generates the test object creation and assertions (legacy)
func generateReadConverterTestCode(entityName, apiType string) string {
	code, _ := generateReadConverterTestCodeWithAssertions(entityName, apiType)
	return code
}

// generateReadConverterTestCodeWithAssertions generates test object creation and field assertions
func generateReadConverterTestCodeWithAssertions(entityName, apiType string) (testCode, assertions string) {
	// For most entities, create a minimal test object with multiple fields
	// Each entity type needs specific handling based on its structure
	switch entityName {
	case "Account":
		testCode = fmt.Sprintf(`	// Create test API object with known values
	testName := "test-account"
	testDesc := "test description"
	testOrg := "test-org"
	apiObj := api.%s{
		Name:         testName,
		Description:  testDesc,
		Organization: testOrg,
	}`, apiType)
		assertions = `	assert.Equal(t, testName, commonObj.Name)
	assert.Equal(t, testDesc, commonObj.Description)
	assert.Equal(t, testOrg, commonObj.Organization)
`

	case "User":
		testCode = fmt.Sprintf(`	// Create test API object with known values
	testName := "test-user"
	testOldName := "old-user"
	apiObj := api.%s{
		Name:    testName,
		OldName: &testOldName,
	}`, apiType)
		assertions = `	assert.Equal(t, testName, commonObj.Name)
	if commonObj.OldName != nil {
		assert.Equal(t, testOldName, *commonObj.OldName)
	}
`

	case "Cluster":
		testCode = fmt.Sprintf(`	// Create test API object with known values
	testName := "test-cluster"
	apiObj := api.%s{
		Name: &testName,
	}`, apiType)
		assertions = `	if commonObj.Name != nil {
		assert.Equal(t, testName, *commonObj.Name)
	}
`

	case "QoS":
		testCode = fmt.Sprintf(`	// Create test API object with known values
	testName := "test-qos"
	testDesc := "test description"
	testID := int32(42)
	apiObj := api.%s{
		Name:        &testName,
		Description: &testDesc,
		Id:          &testID,
	}`, apiType)
		assertions = `	if commonObj.Name != nil {
		assert.Equal(t, testName, *commonObj.Name)
	}
	if commonObj.Description != nil {
		assert.Equal(t, testDesc, *commonObj.Description)
	}
	if commonObj.ID != nil {
		assert.Equal(t, testID, *commonObj.ID)
	}
`

	case "Node":
		testCode = fmt.Sprintf(`	// Create test API object with known values
	testName := "test-node"
	testArch := "x86_64"
	apiObj := api.%s{
		Name:         &testName,
		Architecture: &testArch,
	}`, apiType)
		assertions = `	if commonObj.Name != nil {
		assert.Equal(t, testName, *commonObj.Name)
	}
	if commonObj.Architecture != nil {
		assert.Equal(t, testArch, *commonObj.Architecture)
	}
`

	case "Partition":
		testCode = fmt.Sprintf(`	// Create test API object with known values
	testName := "test-partition"
	testCluster := "test-cluster"
	apiObj := api.%s{
		Name:    &testName,
		Cluster: &testCluster,
	}`, apiType)
		assertions = `	if commonObj.Name != nil {
		assert.Equal(t, testName, *commonObj.Name)
	}
	if commonObj.Cluster != nil {
		assert.Equal(t, testCluster, *commonObj.Cluster)
	}
`

	case "Reservation":
		testCode = fmt.Sprintf(`	// Create test API object with known values
	testName := "test-reservation"
	apiObj := api.%s{
		Name: &testName,
	}`, apiType)
		assertions = `	if commonObj.Name != nil {
		assert.Equal(t, testName, *commonObj.Name)
	}
`

	case "Job":
		testCode = fmt.Sprintf(`	// Create test API object with known values
	testID := int32(12345)
	testName := "test-job"
	testUserName := "testuser"
	apiObj := api.%s{
		JobId:    &testID,
		Name:     &testName,
		UserName: &testUserName,
	}`, apiType)
		assertions = `	if commonObj.JobID != nil {
		assert.Equal(t, testID, *commonObj.JobID)
	}
	if commonObj.Name != nil {
		assert.Equal(t, testName, *commonObj.Name)
	}
	if commonObj.UserName != nil {
		assert.Equal(t, testUserName, *commonObj.UserName)
	}
`

	case "Association":
		testCode = fmt.Sprintf(`	// Create test API object with known values
	testID := int32(123)
	testAccount := "test-account"
	testCluster := "test-cluster"
	testUser := "test-user"
	apiObj := api.%s{
		Id:      &testID,
		Account: &testAccount,
		Cluster: &testCluster,
		User:    testUser,
	}`, apiType)
		assertions = `	if commonObj.ID != nil {
		assert.Equal(t, testID, *commonObj.ID)
	}
	if commonObj.Account != nil {
		assert.Equal(t, testAccount, *commonObj.Account)
	}
	if commonObj.Cluster != nil {
		assert.Equal(t, testCluster, *commonObj.Cluster)
	}
	assert.Equal(t, testUser, commonObj.User)
`

	case "Wckey":
		testCode = fmt.Sprintf(`	// Create test API object with known values
	testName := "test-wckey"
	testUser := "test-user"
	testCluster := "test-cluster"
	apiObj := api.%s{
		Name:    testName,
		User:    testUser,
		Cluster: testCluster,
	}`, apiType)
		assertions = `	assert.Equal(t, testName, commonObj.Name)
	assert.Equal(t, testUser, commonObj.User)
	assert.Equal(t, testCluster, commonObj.Cluster)
`

	default:
		// Generic fallback - just create empty object
		testCode = fmt.Sprintf(`	// Create test API object
	apiObj := api.%s{}`, apiType)
		assertions = ""
	}
	return testCode, assertions
}

// generateWriteConverterTest generates a test for Common → API converter
func generateWriteConverterTest(entityName string, entityDef EntityDef, apiMethods EntityAPIMethods, operation string) string {
	var buf bytes.Buffer

	var inputType string
	var converterFunc string

	if operation == "Create" {
		inputType = entityDef.CreateInput
		converterFunc = fmt.Sprintf("convertCommon%sCreateToAPI", entityName)
	} else {
		inputType = entityDef.UpdateInput
		converterFunc = fmt.Sprintf("convertCommon%sUpdateToAPI", entityName)
	}

	if inputType == "" {
		return ""
	}

	// Determine test values based on input type structure
	testValueCode := generateTestInputForType(entityName, inputType, operation)

	buf.WriteString(fmt.Sprintf(`func Test%sAdapter_%sConverter(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

	// Create test input with known values
	%s

	// Convert Common → API
	apiObj := adapter.%s(input)
	require.NotNil(t, apiObj)
}

`, entityName, operation, entityName, testValueCode, converterFunc))

	return buf.String()
}

// generateTestInputForType creates appropriate test input based on entity and operation
func generateTestInputForType(entityName, inputType, operation string) string {
	// Entity-specific test input generation
	switch entityName {
	case "Account":
		if operation == "Create" {
			return `testName := "test-account"
	input := &types.AccountCreate{
		Name: testName,
	}`
		} else {
			// AccountUpdate doesn't have a Name field
			return `testDesc := "test description"
	input := &types.AccountUpdate{
		Description: &testDesc,
	}`
		}

	case "User":
		if operation == "Create" {
			return `testName := "test-user"
	input := &types.UserCreate{
		Name: testName,
	}`
		} else {
			// UserUpdate doesn't have a Name field
			return `testAccount := "test-account"
	input := &types.UserUpdate{
		DefaultAccount: &testAccount,
	}`
		}

	case "Cluster":
		if operation == "Create" {
			return `testName := "test-cluster"
	input := &types.ClusterCreate{
		Name: testName,
	}`
		} else {
			// ClusterUpdate structure
			return `input := &types.ClusterUpdate{}`
		}

	case "QoS":
		if operation == "Create" {
			return `testName := "test-qos"
	input := &types.QoSCreate{
		Name: testName,
	}`
		} else {
			return `testDesc := "test description"
	input := &types.QoSUpdate{
		Description: &testDesc,
	}`
		}

	case "Association":
		if operation == "Create" {
			return `testAccount := "test-account"
	input := &types.AssociationCreate{
		Account: testAccount,
	}`
		} else {
			return `testComment := "test comment"
	input := &types.AssociationUpdate{
		Comment: &testComment,
	}`
		}

	case "Wckey":
		if operation == "Create" {
			return `testName := "test-wckey"
	input := &types.WckeyCreate{
		Name: testName,
	}`
		} else {
			return `input := &types.WckeyUpdate{}`
		}

	case "Node":
		// Node doesn't have a standard Update converter, it uses NodeUpdate request body
		return `input := &types.NodeUpdate{}`

	case "Reservation":
		if operation == "Create" {
			return `testName := "test-reservation"
	input := &types.ReservationCreate{
		Name: &testName,
	}`
		} else {
			// ReservationUpdate fields
			return `testDuration := int32(3600)
	input := &types.ReservationUpdate{
		Duration: &testDuration,
	}`
		}

	default:
		// Generic fallback
		return fmt.Sprintf(`input := &types.%s{}`, inputType)
	}
}

// generateValidationTests generates tests for validation methods
func generateValidationTests(entityName string, entityDef EntityDef) string {
	var buf bytes.Buffer

	// Generate tests for each validation operation (sorted for deterministic output)
	testOps := make([]string, 0, len(entityDef.Validation))
	for op := range entityDef.Validation {
		testOps = append(testOps, op)
	}
	sort.Strings(testOps)
	for _, operation := range testOps {
		config := entityDef.Validation[operation]
		var inputType string
		var methodName string
		switch operation {
		case "create":
			inputType = entityDef.CreateInput
			if inputType == "" {
				inputType = entityName + "Create"
			}
			methodName = fmt.Sprintf("Validate%sCreate", entityName) // Exported
		case "update":
			inputType = entityDef.UpdateInput
			if inputType == "" {
				inputType = entityName + "Update"
			}
			methodName = fmt.Sprintf("validate%sUpdate", entityName) // Unexported
		default:
			continue
		}

		// Generate validation test function
		titleOp := strings.Title(operation)
		buf.WriteString(fmt.Sprintf(`func Test%sAdapter_Validate%s(t *testing.T) {
	adapter := New%sAdapter(&api.ClientWithResponses{})

`, entityName, titleOp, entityName))

		// Test nil input
		if config.NilError != "" {
			buf.WriteString(fmt.Sprintf(`	t.Run("nil input returns error", func(t *testing.T) {
		err := adapter.%s(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "%s")
	})

`, methodName, config.NilError))
		}

		// Test required fields
		for _, req := range config.RequiredFields {
			fieldCheck := generateValidationTestEmptyField(inputType, req, config.RequiredFields)
			buf.WriteString(fmt.Sprintf(`	t.Run("missing %s returns error", func(t *testing.T) {
		input := %s
		err := adapter.%s(input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "%s")
	})

`, req.Field, fieldCheck, methodName, req.Error))
		}

		// Test at_least_one_of validation
		if len(config.AtLeastOneOf) > 0 && config.AtLeastOneError != "" {
			buf.WriteString(fmt.Sprintf(`	t.Run("empty fields returns error", func(t *testing.T) {
		input := &types.%s{}
		err := adapter.%s(input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "%s")
	})

`, inputType, methodName, config.AtLeastOneError))
		}

		// Test valid input passes
		validInput := generateValidationTestValidInput(entityName, inputType, operation, config)
		buf.WriteString(fmt.Sprintf(`	t.Run("valid input passes", func(t *testing.T) {
		input := %s
		err := adapter.%s(input)
		require.NoError(t, err)
	})
}

`, validInput, methodName))
	}

	return buf.String()
}

// generateValidationTestEmptyField generates code for an input struct with the specified field empty
func generateValidationTestEmptyField(inputType string, emptyField RequiredField, allFields []RequiredField) string {
	// For simple types with one required field, just return empty struct
	if len(allFields) == 1 {
		return fmt.Sprintf("&types.%s{}", inputType)
	}

	// For types with multiple required fields, fill all except the target
	var fields []string
	for _, f := range allFields {
		if f.Field == emptyField.Field {
			continue // Skip the field we want empty
		}
		switch f.CheckType {
		case "time":
			fields = append(fields, fmt.Sprintf("%s: time.Now()", f.Field))
		case "slice":
			fields = append(fields, fmt.Sprintf("%s: []string{\"test\"}", f.Field))
		case "pointer":
			fields = append(fields, fmt.Sprintf("%s: ptrString(\"test\")", f.Field))
		default:
			fields = append(fields, fmt.Sprintf("%s: \"test\"", f.Field))
		}
	}

	if len(fields) == 0 {
		return fmt.Sprintf("&types.%s{}", inputType)
	}

	return fmt.Sprintf("&types.%s{%s}", inputType, strings.Join(fields, ", "))
}

// generateValidationTestValidInput generates code for a valid input struct
func generateValidationTestValidInput(entityName, inputType, operation string, config ValidationConfig) string {
	// If there's at_least_one_of, just set one of those fields
	if len(config.AtLeastOneOf) > 0 {
		field := config.AtLeastOneOf[0]
		if field == "Accounts" || field == "Users" || field == "QoSList" {
			return fmt.Sprintf("&types.%s{%s: []string{\"test\"}}", inputType, field)
		}
		// Handle special types
		fieldValue := getFieldValueForEntity(entityName, field)
		return fmt.Sprintf("&types.%s{%s: %s}", inputType, field, fieldValue)
	}

	// Generate based on required fields
	var fields []string
	for _, f := range config.RequiredFields {
		switch f.CheckType {
		case "time":
			fields = append(fields, fmt.Sprintf("%s: time.Now()", f.Field))
		case "slice":
			fields = append(fields, fmt.Sprintf("%s: []string{\"test\"}", f.Field))
		case "pointer":
			fields = append(fields, fmt.Sprintf("%s: ptrString(\"test\")", f.Field))
		default:
			fields = append(fields, fmt.Sprintf("%s: \"test\"", f.Field))
		}
	}

	if len(fields) == 0 {
		return fmt.Sprintf("&types.%s{}", inputType)
	}

	return fmt.Sprintf("&types.%s{%s}", inputType, strings.Join(fields, ", "))
}

// getFieldValueForEntity returns the appropriate test value for a field based on entity
func getFieldValueForEntity(entityName, fieldName string) string {
	// Handle special types per entity
	switch entityName {
	case "Node":
		switch fieldName {
		case "State", "NextStateAfterReboot":
			// NodeUpdate.State is []NodeState, not *NodeState
			return "[]types.NodeState{types.NodeStateIdle}"
		}
	}
	// Default: pointer to string
	return "ptrString(\"test\")"
}
