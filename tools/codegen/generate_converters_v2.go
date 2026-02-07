//go:build ignore

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// generate_converters_v2.go - Enhanced converter generator with pattern detection
// Usage: go run generate_converters_v2.go converter_patterns.go converter_helpers.go -version=v0_0_44

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Config structures for YAML loading
type GeneratorConfig struct {
	Generation      GenerationSettings           `yaml:"generation"`
	FieldTransforms map[string]string            `yaml:"field_transforms"`
	Versions        map[string]VersionConfig     `yaml:"versions"`
	Patterns        map[string]PatternDefinition `yaml:"patterns"`
}

type GenerationSettings struct {
	AutoDetectPatterns bool     `yaml:"auto_detect_patterns"`
	GenerateTests      bool     `yaml:"generate_tests"`
	OutputSuffix       string   `yaml:"output_suffix"`
	SkipFields         []string `yaml:"skip_fields"`
}

type VersionConfig struct {
	Skip       bool                      `yaml:"skip"`        // Skip generation for this version
	APIPackage string                    `yaml:"api_package"`
	APIPrefix  string                    `yaml:"api_prefix"`
	Entities   map[string]EntityConfig   `yaml:"entities"`
}

type EntityConfig struct {
	APIType        string                  `yaml:"api_type"`
	CommonType     string                  `yaml:"common_type"`
	CreateType     string                  `yaml:"create_type"`     // e.g., "AccountCreate"
	UpdateType     string                  `yaml:"update_type"`     // e.g., "AccountUpdate"
	CreateAPIType  string                  `yaml:"create_api_type"` // Override API type for create (defaults to APIType)
	UpdateAPIType  string                  `yaml:"update_api_type"` // Override API type for update (defaults to APIType)
	CustomFields   map[string]FieldConfig  `yaml:"custom_fields"`
	CreateFields   map[string]FieldConfig  `yaml:"create_fields"` // Custom mappings for Create
	UpdateFields   map[string]FieldConfig  `yaml:"update_fields"` // Custom mappings for Update
}

type FieldConfig struct {
	APIField    string `yaml:"api_field"`
	CommonField string `yaml:"common_field"`
	Pattern     string `yaml:"pattern"`
	APIType     string `yaml:"api_type"`
	CommonType  string `yaml:"common_type"`
	Converter   string `yaml:"converter"`
	Helper      string `yaml:"helper"`
}

type PatternDefinition struct {
	Description string `yaml:"description"`
	Template    string `yaml:"template"`
}

// TypeInfo holds parsed type information
type TypeInfo struct {
	Name   string
	Fields map[string]FieldInfo
}

// TypeRegistry holds both struct types and type aliases for resolution
type TypeRegistry struct {
	Structs map[string]*TypeInfo // Struct type definitions
	Aliases map[string]ast.Expr  // Type alias definitions (name -> underlying type expression)
}

// FieldInfo holds parsed field information
type FieldInfo struct {
	Name              string
	Type              string
	IsPtr             bool
	IsSlice           bool
	ElemType          string
	IsAnonymousStruct bool // True if field type is an inline anonymous struct
	IsTypeAlias       bool // True if type is an alias (e.g., V0042Flags = []string)
	UnderlyingType    string // The underlying type if this is an alias
}

var (
	versionFlag = flag.String("version", "", "API version to generate (e.g., v0_0_44)")
	allFlag     = flag.Bool("all", false, "Generate for all versions")
	dryRunFlag  = flag.Bool("dry-run", false, "Print generated code without writing files")
	configFlag  = flag.String("config", "tools/codegen/converter_config_enhanced.yaml", "Path to config file")
	debugFlag   = flag.Bool("debug", false, "Enable debug output")
)

func main() {
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	versions := []string{"v0_0_40", "v0_0_41", "v0_0_42", "v0_0_43", "v0_0_44"}

	if *versionFlag != "" {
		versions = []string{*versionFlag}
	} else if !*allFlag {
		fmt.Println("Usage: go run generate_converters_v2.go converter_patterns.go converter_helpers.go -version=v0_0_44")
		fmt.Println("       go run generate_converters_v2.go converter_patterns.go converter_helpers.go -all")
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, version := range versions {
		fmt.Printf("Generating converters for %s...\n", version)
		if err := generateConvertersForVersion(version, config); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating %s: %v\n", version, err)
			continue
		}
		fmt.Printf("  ✓ Done: %s\n", version)
	}
}

func loadConfig(path string) (*GeneratorConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var config GeneratorConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &config, nil
}

func generateConvertersForVersion(version string, config *GeneratorConfig) error {
	versionConfig, ok := config.Versions[version]
	if !ok {
		return fmt.Errorf("no configuration for version %s", version)
	}

	// Skip generation if version is marked as skip
	if versionConfig.Skip {
		fmt.Printf("  Skipping %s (marked skip: true in config)\n", version)
		return nil
	}

	// Parse API types from generated client
	apiDir := filepath.Join("internal", "openapi", version)
	apiTypes, err := parseAPITypes(apiDir, versionConfig.APIPrefix)
	if err != nil {
		return fmt.Errorf("parsing API types: %w", err)
	}

	// Parse common types (defined in api/ package)
	commonDir := "api"
	commonTypes, err := parseCommonTypes(commonDir)
	if err != nil {
		return fmt.Errorf("parsing common types: %w", err)
	}

	// Generate shared helpers file first
	if err := generateHelpersFile(version, versionConfig.APIPrefix, config); err != nil {
		return fmt.Errorf("generating helpers: %w", err)
	}

	// Generate converters for each entity
	// Note: Order matters for some dependencies, but most entities are independent
	entities := []string{"Account", "Association", "Job", "Node", "Partition", "QoS", "Reservation", "User", "WCKey", "Cluster"}

	for _, entity := range entities {
		entityConfig, ok := versionConfig.Entities[entity]
		if !ok {
			// Skip entities not in config
			continue
		}

		code, err := generateEntityConverter(version, entity, entityConfig, apiTypes, commonTypes, versionConfig.APIPrefix, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: Could not generate %s: %v\n", entity, err)
			continue
		}

		if code == "" {
			continue
		}

		outputPath := filepath.Join("internal", "adapters", version, strings.ToLower(entity)+"_converters"+config.Generation.OutputSuffix)

		if *dryRunFlag {
			fmt.Printf("--- %s ---\n%s\n", outputPath, code)
		} else {
			if err := os.WriteFile(outputPath, []byte(code), 0644); err != nil {
				return fmt.Errorf("writing %s: %w", outputPath, err)
			}
			fmt.Printf("  Generated: %s\n", outputPath)
		}

		// Generate Common→API converters for Create/Update if configured
		if entityConfig.CreateType != "" || entityConfig.UpdateType != "" {
			reverseCode, err := generateEntityReverseConverters(version, entity, entityConfig, apiTypes, commonTypes, versionConfig.APIPrefix, config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: Could not generate reverse converters for %s: %v\n", entity, err)
			} else if reverseCode != "" {
				reverseOutputPath := filepath.Join("internal", "adapters", version, strings.ToLower(entity)+"_converters_write"+config.Generation.OutputSuffix)
				if *dryRunFlag {
					fmt.Printf("--- %s ---\n%s\n", reverseOutputPath, reverseCode)
				} else {
					if err := os.WriteFile(reverseOutputPath, []byte(reverseCode), 0644); err != nil {
						return fmt.Errorf("writing %s: %w", reverseOutputPath, err)
					}
					fmt.Printf("  Generated: %s\n", reverseOutputPath)
				}
			}
		}
	}

	return nil
}

func generateHelpersFile(version, apiPrefix string, config *GeneratorConfig) error {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(`// Code generated by generate_converters_v2.go. DO NOT EDIT.
// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package %s

import (
	api "github.com/jontk/slurm-client/internal/openapi/%s"
	types "github.com/jontk/slurm-client/api"
)

`, version, version))

	// Generate helper functions
	helpers := GenerateHelpers(version, apiPrefix)
	buf.WriteString(helpers)

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not format helpers file: %v\n", err)
		formatted = buf.Bytes()
	}

	outputPath := filepath.Join("internal", "adapters", version, "converter_helpers"+config.Generation.OutputSuffix)

	if *dryRunFlag {
		fmt.Printf("--- %s ---\n%s\n", outputPath, string(formatted))
	} else {
		if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
			return fmt.Errorf("writing helpers: %w", err)
		}
		fmt.Printf("  Generated: %s\n", outputPath)
	}

	return nil
}

func generateEntityConverter(
	version, entity string,
	entityConfig EntityConfig,
	apiTypes, commonTypes map[string]*TypeInfo,
	apiPrefix string,
	config *GeneratorConfig,
) (string, error) {
	apiType, ok := apiTypes[entityConfig.APIType]
	if !ok {
		return "", fmt.Errorf("API type %s not found", entityConfig.APIType)
	}

	commonType, ok := commonTypes[entityConfig.CommonType]
	if !ok {
		return "", fmt.Errorf("common type %s not found", entityConfig.CommonType)
	}

	// Build field mappings
	mappings := buildFieldMappings(apiType, commonType, entityConfig, config, apiPrefix)

	// Check if any field needs time import (only for non-skipped patterns)
	needsTime := false
	for _, mapping := range mappings {
		if mapping.Pattern == "skip" {
			continue
		}
		if mapping.Pattern == "time_novalnumber" || strings.Contains(mapping.CommonType, "time.Time") {
			needsTime = true
			break
		}
	}

	// Generate file header
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(`// Code generated by generate_converters_v2.go. DO NOT EDIT.
// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package %s

import (
`, version))

	if needsTime {
		buf.WriteString("\t\"time\"\n\n")
	}

	buf.WriteString(fmt.Sprintf(`	api "github.com/jontk/slurm-client/internal/openapi/%s"
	types "github.com/jontk/slurm-client/api"
)

`, version))

	// Generate main converter function
	objName := "apiObj"
	buf.WriteString(fmt.Sprintf(`// convertAPI%sToCommon converts API %s to common %s type
func (a *%sAdapter) convertAPI%sToCommon(%s api.%s) *types.%s {
	result := &types.%s{}

`, entity, entityConfig.APIType, entity, entity, entity, objName, entityConfig.APIType, entity, entity))

	// Generate field conversions
	for _, mapping := range mappings {
		conversionCode, err := generateFieldConversion(mapping, config)
		if err != nil {
			return "", fmt.Errorf("generating field %s: %w", mapping.APIField, err)
		}
		buf.WriteString(conversionCode)
	}

	buf.WriteString("\treturn result\n}\n")

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// Return unformatted code with error for debugging
		fmt.Fprintf(os.Stderr, "Warning: Could not format generated code: %v\n", err)
		return buf.String(), nil
	}

	return string(formatted), nil
}

// generateEntityReverseConverters generates Common→API converters for Create and Update operations
func generateEntityReverseConverters(
	version, entity string,
	entityConfig EntityConfig,
	apiTypes, commonTypes map[string]*TypeInfo,
	apiPrefix string,
	config *GeneratorConfig,
) (string, error) {
	apiType, ok := apiTypes[entityConfig.APIType]
	if !ok {
		return "", fmt.Errorf("API type %s not found", entityConfig.APIType)
	}

	// First, generate the converters to a temp buffer to detect what imports are needed
	var converterCode bytes.Buffer
	needsStrings := false

	// Generate Create converter if configured
	if entityConfig.CreateType != "" {
		createType, ok := commonTypes[entityConfig.CreateType]
		if ok {
			// Use CreateAPIType if specified, otherwise use default APIType
			targetAPIType := apiType
			if entityConfig.CreateAPIType != "" {
				if customAPIType, ok := apiTypes[entityConfig.CreateAPIType]; ok {
					targetAPIType = customAPIType
				} else {
					return "", fmt.Errorf("create_api_type %s not found", entityConfig.CreateAPIType)
				}
			}
			code := generateReverseConverter(entity, entityConfig.CreateType, "Create", createType, targetAPIType, entityConfig, apiPrefix, config)
			converterCode.WriteString(code)
			if strings.Contains(code, "strings.Join") {
				needsStrings = true
			}
		}
	}

	// Generate Update converter if configured
	if entityConfig.UpdateType != "" {
		updateType, ok := commonTypes[entityConfig.UpdateType]
		if ok {
			// Use UpdateAPIType if specified, otherwise use default APIType
			targetAPIType := apiType
			if entityConfig.UpdateAPIType != "" {
				if customAPIType, ok := apiTypes[entityConfig.UpdateAPIType]; ok {
					targetAPIType = customAPIType
				} else {
					return "", fmt.Errorf("update_api_type %s not found", entityConfig.UpdateAPIType)
				}
			}
			code := generateReverseConverter(entity, entityConfig.UpdateType, "Update", updateType, targetAPIType, entityConfig, apiPrefix, config)
			converterCode.WriteString(code)
			if strings.Contains(code, "strings.Join") {
				needsStrings = true
			}
		}
	}

	if converterCode.Len() == 0 {
		return "", nil
	}

	var buf bytes.Buffer

	// File header with conditional imports
	buf.WriteString(fmt.Sprintf(`// Code generated by generate_converters_v2.go. DO NOT EDIT.
// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package %s

import (
`, version))

	if needsStrings {
		buf.WriteString("\t\"strings\"\n\n")
	}

	buf.WriteString(fmt.Sprintf(`	api "github.com/jontk/slurm-client/internal/openapi/%s"
	types "github.com/jontk/slurm-client/api"
)

`, version))

	// Add the converter code
	buf.Write(converterCode.Bytes())

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not format reverse converter code: %v\n", err)
		return buf.String(), nil
	}

	return string(formatted), nil
}

// generateReverseConverter generates a single Common→API converter function
func generateReverseConverter(
	entity, commonTypeName, operation string,
	commonType, apiType *TypeInfo,
	entityConfig EntityConfig,
	apiPrefix string,
	config *GeneratorConfig,
) string {
	var buf bytes.Buffer

	funcName := fmt.Sprintf("convertCommon%s%sToAPI", entity, operation)
	inputType := fmt.Sprintf("types.%s", commonTypeName)
	outputType := fmt.Sprintf("api.%s", apiType.Name)

	buf.WriteString(fmt.Sprintf(`// %s converts common %s to API %s type
func (a *%sAdapter) %s(input *%s) *%s {
	if input == nil {
		return nil
	}
	result := &%s{}

`, funcName, commonTypeName, apiType.Name, entity, funcName, inputType, outputType, outputType))

	// Build reverse field mappings
	mappings := buildReverseFieldMappings(commonType, apiType, entityConfig, config, apiPrefix, operation)

	// Generate field conversions
	for _, mapping := range mappings {
		conversionCode := generateReverseFieldConversion(mapping, apiPrefix)
		buf.WriteString(conversionCode)
	}

	buf.WriteString("\treturn result\n}\n\n")

	return buf.String()
}

// ReverseFieldMapping holds mapping info for Common→API conversion
type ReverseFieldMapping struct {
	CommonField    string
	APIField       string
	CommonType     string
	APIType        string
	CommonIsPtr    bool
	APIIsPtr       bool
	CommonIsSlice  bool
	APIIsSlice     bool
	APIIsTypeAlias bool   // True if API type is a type alias (e.g., V0042Flags = []string)
	APIElemType    string // Element type if API is slice (e.g., "string" for []string)
	Pattern        string // reverse pattern
}

// buildReverseFieldMappings builds field mappings for Common→API conversion
func buildReverseFieldMappings(
	commonType, apiType *TypeInfo,
	entityConfig EntityConfig,
	config *GeneratorConfig,
	apiPrefix string,
	operation string,
) []ReverseFieldMapping {
	var mappings []ReverseFieldMapping

	// Get custom fields for this operation
	var customFields map[string]FieldConfig
	if operation == "Create" && entityConfig.CreateFields != nil {
		customFields = entityConfig.CreateFields
	} else if operation == "Update" && entityConfig.UpdateFields != nil {
		customFields = entityConfig.UpdateFields
	}

	// Sort common fields for consistent output
	var commonFieldNames []string
	for name := range commonType.Fields {
		commonFieldNames = append(commonFieldNames, name)
	}
	sort.Strings(commonFieldNames)

	for _, commonFieldName := range commonFieldNames {
		commonField := commonType.Fields[commonFieldName]

		// Skip fields starting with underscore
		if strings.HasPrefix(commonFieldName, "_") {
			continue
		}

		// Check for custom field configuration
		if customFields != nil {
			if cfg, ok := customFields[commonFieldName]; ok {
				if cfg.Pattern == "skip" {
					continue
				}
			}
		}

		// Apply field transforms to find API field
		apiFieldName := commonFieldName
		// Reverse transform lookup
		for apiName, commonName := range config.FieldTransforms {
			if commonName == commonFieldName {
				apiFieldName = apiName
				break
			}
		}

		// Find corresponding API field
		apiField, exists := apiType.Fields[apiFieldName]
		if !exists {
			// Try case-insensitive match
			for aName, aField := range apiType.Fields {
				if strings.EqualFold(aName, apiFieldName) {
					apiField = aField
					apiFieldName = aName
					exists = true
					break
				}
			}
		}

		if !exists {
			// Field doesn't exist in API type, skip
			continue
		}

		mapping := ReverseFieldMapping{
			CommonField:    commonFieldName,
			APIField:       apiFieldName,
			CommonType:     commonField.Type,
			APIType:        apiField.Type,
			CommonIsPtr:    commonField.IsPtr,
			APIIsPtr:       apiField.IsPtr,
			CommonIsSlice:  commonField.IsSlice,
			APIIsSlice:     apiField.IsSlice,
			APIIsTypeAlias: apiField.IsTypeAlias,
			APIElemType:    apiField.ElemType,
		}

		// Determine reverse pattern
		mapping.Pattern = detectReversePattern(mapping, apiPrefix)

		if *debugFlag {
			fmt.Printf("    [DEBUG REVERSE] %s.%s -> %s: common{Type:%s, IsPtr:%v, IsSlice:%v} api{Type:%s, IsPtr:%v, IsSlice:%v} => pattern: %s\n",
				operation, commonFieldName, apiFieldName,
				commonField.Type, commonField.IsPtr, commonField.IsSlice,
				apiField.Type, apiField.IsPtr, apiField.IsSlice,
				mapping.Pattern)
		}

		mappings = append(mappings, mapping)
	}

	return mappings
}

// detectReversePattern determines the conversion pattern for Common→API
func detectReversePattern(mapping ReverseFieldMapping, apiPrefix string) string {
	// Check if API type is a NoValStruct
	if strings.Contains(mapping.APIType, "NoValStruct") {
		return "to_novalstruct"
	}

	// Check for CsvString type (comma-separated string) - []string -> *CsvString
	if strings.Contains(mapping.APIType, "CsvString") {
		if mapping.CommonIsSlice && (mapping.CommonType == "string" || mapping.CommonType == "[]string") {
			return "to_csv_string"
		}
	}

	// Check for CoordList type ([]string -> []Coord)
	if strings.Contains(mapping.APIType, "CoordList") {
		if mapping.CommonIsSlice && (mapping.CommonType == "string" || mapping.CommonType == "[]string") {
			return "to_coord_list"
		}
	}

	// Check for StringList type
	if strings.Contains(mapping.APIType, "StringList") {
		if mapping.CommonIsSlice {
			return "reverse_slice_to_ptr"
		}
	}

	// Check for map types - these need manual handling
	if strings.HasPrefix(mapping.CommonType, "map[") {
		return "reverse_skip_complex"
	}

	// Check for enum state conversions (e.g., *NodeState -> *[]V0044NodeState)
	if strings.Contains(mapping.CommonType, "State") && strings.Contains(mapping.APIType, "State") {
		return "reverse_skip_complex" // State enums need manual handling due to enum mapping
	}

	// Check for Limits type - common QoSLimits is different from API inline struct
	if mapping.CommonType == "QoSLimits" || strings.HasSuffix(mapping.CommonType, ".QoSLimits") {
		return "reverse_skip_complex" // QoSLimits needs manual conversion due to different structure
	}

	// Check for enum slice conversions (e.g., []string -> *[]V0044QosFlags)
	// Detect when API type contains "Flags" and is a pointer to slice
	// Handle both []string (Create) and *[]string (Update) cases
	if mapping.CommonIsSlice && (mapping.CommonType == "string" || mapping.CommonType == "[]string") {
		// Check if API type looks like an enum type (e.g., V0044QosFlags, V0044AccountFlags)
		if strings.Contains(mapping.APIType, "Flags") && mapping.APIIsSlice {
			// If API type is a type alias for []string (like V0042ClusterRecFlags = []string),
			// use direct assignment instead of element-wise casting
			if mapping.APIIsTypeAlias && mapping.APIElemType == "string" {
				return "to_slice_type_alias"
			}
			return "to_enum_flags_slice"
		}
	}

	// Slice handling - check BEFORE pointer handling
	if mapping.CommonIsSlice {
		if mapping.APIIsSlice {
			return "reverse_slice_copy"
		}
		if mapping.APIIsPtr {
			// Common is []T, API is *[]T - need to take address of slice
			return "reverse_slice_to_ptr"
		}
		return "reverse_slice_copy"
	}

	// Simple cases for non-slice types
	if !mapping.CommonIsPtr && !mapping.APIIsPtr {
		return "reverse_simple_copy"
	}
	if mapping.CommonIsPtr && mapping.APIIsPtr {
		return "reverse_ptr_copy"
	}
	if !mapping.CommonIsPtr && mapping.APIIsPtr {
		return "reverse_take_addr"
	}
	if mapping.CommonIsPtr && !mapping.APIIsPtr {
		return "reverse_deref"
	}

	return "reverse_simple_copy"
}

// generateReverseFieldConversion generates code for a single Common→API field
func generateReverseFieldConversion(mapping ReverseFieldMapping, apiPrefix string) string {
	switch mapping.Pattern {
	case "to_novalstruct":
		// Handle NoValStruct: need to create struct with Set=true and Number=value
		return generateNoValStructAssignment(mapping, apiPrefix)
	case "to_csv_string":
		// Convert []string to CsvString (comma-separated)
		return generateCsvStringConversion(mapping, apiPrefix)
	case "to_coord_list":
		// Convert []string to []Coord
		return generateCoordListConversion(mapping, apiPrefix)
	case "to_enum_flags_slice":
		// Convert []string to *[]EnumFlagsType
		return generateEnumFlagsSliceConversion(mapping, apiPrefix)
	case "to_slice_type_alias":
		// Convert []string to *TypeAlias where TypeAlias = []string (direct assignment)
		return generateSliceTypeAliasConversion(mapping, apiPrefix)
	case "reverse_simple_copy":
		return fmt.Sprintf("\tresult.%s = input.%s\n", mapping.APIField, mapping.CommonField)
	case "reverse_ptr_copy":
		return fmt.Sprintf("\tif input.%s != nil {\n\t\tresult.%s = input.%s\n\t}\n", mapping.CommonField, mapping.APIField, mapping.CommonField)
	case "reverse_take_addr":
		// Non-zero check for value types
		if strings.HasPrefix(mapping.CommonType, "int") || strings.HasPrefix(mapping.CommonType, "uint") || strings.HasPrefix(mapping.CommonType, "float") {
			return fmt.Sprintf("\tif input.%s != 0 {\n\t\tval := input.%s\n\t\tresult.%s = &val\n\t}\n", mapping.CommonField, mapping.CommonField, mapping.APIField)
		}
		if mapping.CommonType == "string" {
			return fmt.Sprintf("\tif input.%s != \"\" {\n\t\tval := input.%s\n\t\tresult.%s = &val\n\t}\n", mapping.CommonField, mapping.CommonField, mapping.APIField)
		}
		return fmt.Sprintf("\tval%s := input.%s\n\tresult.%s = &val%s\n", mapping.CommonField, mapping.CommonField, mapping.APIField, mapping.CommonField)
	case "reverse_deref":
		return fmt.Sprintf("\tif input.%s != nil {\n\t\tresult.%s = *input.%s\n\t}\n", mapping.CommonField, mapping.APIField, mapping.CommonField)
	case "reverse_slice_copy":
		return fmt.Sprintf("\tif len(input.%s) > 0 {\n\t\tresult.%s = input.%s\n\t}\n", mapping.CommonField, mapping.APIField, mapping.CommonField)
	case "reverse_slice_to_ptr":
		// Common is []T, API is *[]T
		return fmt.Sprintf("\tif len(input.%s) > 0 {\n\t\tslice := input.%s\n\t\tresult.%s = &slice\n\t}\n", mapping.CommonField, mapping.CommonField, mapping.APIField)
	case "reverse_skip_complex":
		// Skip complex types that need manual handling
		return fmt.Sprintf("\t// %s - SKIPPED: Complex type requiring manual conversion\n", mapping.CommonField)
	default:
		return fmt.Sprintf("\t// %s -> %s (pattern: %s - not implemented)\n", mapping.CommonField, mapping.APIField, mapping.Pattern)
	}
}

// generateEnumFlagsSliceConversion generates code to convert []string or *[]string to *[]EnumFlagsType
func generateEnumFlagsSliceConversion(mapping ReverseFieldMapping, apiPrefix string) string {
	// Extract the enum type from API type (e.g., "*[]V0044QosFlags" -> "api.V0044QosFlags")
	apiTypeClean := strings.TrimPrefix(mapping.APIType, "*")
	apiTypeClean = strings.TrimPrefix(apiTypeClean, "[]")
	if !strings.HasPrefix(apiTypeClean, "api.") {
		apiTypeClean = "api." + apiTypeClean
	}

	// Handle pointer-to-slice case (e.g., *[]string in Update types)
	if mapping.CommonIsPtr {
		return fmt.Sprintf(`	if input.%s != nil && len(*input.%s) > 0 {
		flags := make([]%s, len(*input.%s))
		for i, f := range *input.%s {
			flags[i] = %s(f)
		}
		result.%s = &flags
	}
`, mapping.CommonField, mapping.CommonField, apiTypeClean, mapping.CommonField, mapping.CommonField, apiTypeClean, mapping.APIField)
	}

	// Non-pointer slice case (e.g., []string in Create types)
	return fmt.Sprintf(`	if len(input.%s) > 0 {
		flags := make([]%s, len(input.%s))
		for i, f := range input.%s {
			flags[i] = %s(f)
		}
		result.%s = &flags
	}
`, mapping.CommonField, apiTypeClean, mapping.CommonField, mapping.CommonField, apiTypeClean, mapping.APIField)
}

// generateSliceTypeAliasConversion generates code to convert []string to *TypeAlias where TypeAlias = []string
// This handles cases like V0042ClusterRecFlags = []string where we can directly assign without element casting
func generateSliceTypeAliasConversion(mapping ReverseFieldMapping, apiPrefix string) string {
	// The API type is a type alias for []string, so we can directly assign after converting to the alias type
	apiTypeClean := strings.TrimPrefix(mapping.APIType, "*")
	if !strings.HasPrefix(apiTypeClean, "api.") {
		apiTypeClean = "api." + apiTypeClean
	}

	// Handle pointer-to-slice case (e.g., *[]string in Update types)
	if mapping.CommonIsPtr {
		return fmt.Sprintf(`	if input.%s != nil && len(*input.%s) > 0 {
		slice := %s(*input.%s)
		result.%s = &slice
	}
`, mapping.CommonField, mapping.CommonField, apiTypeClean, mapping.CommonField, mapping.APIField)
	}

	// Non-pointer slice case (e.g., []string in Create types)
	return fmt.Sprintf(`	if len(input.%s) > 0 {
		slice := %s(input.%s)
		result.%s = &slice
	}
`, mapping.CommonField, apiTypeClean, mapping.CommonField, mapping.APIField)
}

// generateCsvStringConversion generates code to convert []string to CsvString
func generateCsvStringConversion(mapping ReverseFieldMapping, apiPrefix string) string {
	// Extract the CsvString type (e.g., "*api.V0044CsvString" -> "api.V0044CsvString")
	apiTypeClean := strings.TrimPrefix(mapping.APIType, "*")

	return fmt.Sprintf(`	if len(input.%s) > 0 {
		csv := strings.Join(input.%s, ",")
		csvStr := %s(csv)
		result.%s = &csvStr
	}
`, mapping.CommonField, mapping.CommonField, apiTypeClean, mapping.APIField)
}

// generateCoordListConversion generates code to convert []string to []Coord
func generateCoordListConversion(mapping ReverseFieldMapping, apiPrefix string) string {
	// Determine the Coord type name based on API prefix
	coordType := fmt.Sprintf("api.%sCoord", apiPrefix)

	return fmt.Sprintf(`	if len(input.%s) > 0 {
		coords := make([]%s, len(input.%s))
		for i, name := range input.%s {
			coords[i] = %s{Name: name}
		}
		result.%s = &coords
	}
`, mapping.CommonField, coordType, mapping.CommonField, mapping.CommonField, coordType, mapping.APIField)
}

// generateNoValStructAssignment generates code to create a NoValStruct
func generateNoValStructAssignment(mapping ReverseFieldMapping, apiPrefix string) string {
	// Extract the NoValStruct type name and add api. prefix
	// mapping.APIType is like "*V0044Uint32NoValStruct" or "V0044Uint32NoValStruct"
	apiTypeClean := strings.TrimPrefix(mapping.APIType, "*")
	// Add api. prefix if not already present
	if !strings.HasPrefix(apiTypeClean, "api.") {
		apiTypeClean = "api." + apiTypeClean
	}

	// Determine the value type based on the API NoValStruct type
	// The Number field type depends on the NoValStruct variant:
	// - Uint32NoValStruct.Number is *int32
	// - Uint64NoValStruct.Number is *int64
	// - Float64NoValStruct.Number is *float64
	valueType := "int64" // default for Uint64NoValStruct
	if strings.Contains(mapping.APIType, "Uint32NoValStruct") {
		valueType = "int32"
	} else if strings.Contains(mapping.APIType, "Float64NoValStruct") {
		valueType = "float64"
	}

	if mapping.CommonIsPtr {
		return fmt.Sprintf(`	if input.%s != nil {
		setTrue := true
		num := %s(*input.%s)
		result.%s = &%s{
			Set:    &setTrue,
			Number: &num,
		}
	}
`, mapping.CommonField, valueType, mapping.CommonField, mapping.APIField, apiTypeClean)
	}

	// Non-pointer common field - check for non-zero
	zeroCheck := "input." + mapping.CommonField + " != 0"
	if mapping.CommonType == "string" {
		zeroCheck = "input." + mapping.CommonField + " != \"\""
	}

	return fmt.Sprintf(`	if %s {
		setTrue := true
		num := %s(input.%s)
		result.%s = &%s{
			Set:    &setTrue,
			Number: &num,
		}
	}
`, zeroCheck, valueType, mapping.CommonField, mapping.APIField, apiTypeClean)
}

type FieldMappingInfo struct {
	APIField       string
	CommonField    string
	APIType        string
	CommonType     string
	Pattern        string
	Converter      string
	Helper         string
	CommonElemType string // For slice conversions
}

func buildFieldMappings(
	apiType, commonType *TypeInfo,
	entityConfig EntityConfig,
	config *GeneratorConfig,
	apiPrefix string,
) []FieldMappingInfo {
	var mappings []FieldMappingInfo
	detector := NewPatternDetector(apiPrefix)

	// Sort API fields for consistent output
	var apiFieldNames []string
	for name := range apiType.Fields {
		apiFieldNames = append(apiFieldNames, name)
	}
	sort.Strings(apiFieldNames)

	for _, apiFieldName := range apiFieldNames {
		apiField := apiType.Fields[apiFieldName]

		// Skip fields starting with underscore
		if strings.HasPrefix(apiFieldName, "_") {
			continue
		}

		// Determine common field name
		commonFieldName := apiFieldName
		if transformed, ok := config.FieldTransforms[apiFieldName]; ok {
			commonFieldName = transformed
		}

		// Check for custom field configuration
		var fieldConfig FieldConfig
		var hasCustomConfig bool
		if entityConfig.CustomFields != nil {
			if cfg, ok := entityConfig.CustomFields[commonFieldName]; ok {
				fieldConfig = cfg
				hasCustomConfig = true
				// Override field names if specified
				if fieldConfig.APIField != "" {
					apiFieldName = fieldConfig.APIField
				}
				if fieldConfig.CommonField != "" {
					commonFieldName = fieldConfig.CommonField
				}
			}
		}

		// Find corresponding common field
		commonField, exists := commonType.Fields[commonFieldName]
		if !exists {
			// Try case-insensitive match
			for cName, cField := range commonType.Fields {
				if strings.EqualFold(cName, commonFieldName) {
					commonField = cField
					commonFieldName = cName
					exists = true
					break
				}
			}
		}

		if !exists {
			// Field doesn't exist in common type, skip
			continue
		}

		mapping := FieldMappingInfo{
			APIField:    apiFieldName,
			CommonField: commonFieldName,
			APIType:     apiField.Type,
			CommonType:  commonField.Type,
		}

		// Determine pattern
		if hasCustomConfig && fieldConfig.Pattern != "" {
			mapping.Pattern = fieldConfig.Pattern
			mapping.Converter = fieldConfig.Converter
			mapping.Helper = fieldConfig.Helper
			if *debugFlag {
				fmt.Printf("    [DEBUG] Field %s -> %s: using custom pattern '%s'\n", apiFieldName, commonFieldName, mapping.Pattern)
			}
		} else if config.Generation.AutoDetectPatterns {
			// Auto-detect pattern
			pattern := detector.DetectPattern(apiField, commonField)
			mapping.Pattern = patternToString(pattern)
			if *debugFlag {
				fmt.Printf("    [DEBUG] Field %s -> %s: api{Type:%s, IsPtr:%v, IsSlice:%v, ElemType:%s} common{Type:%s, IsPtr:%v, IsSlice:%v, ElemType:%s} => pattern: %s\n",
					apiFieldName, commonFieldName,
					apiField.Type, apiField.IsPtr, apiField.IsSlice, apiField.ElemType,
					commonField.Type, commonField.IsPtr, commonField.IsSlice, commonField.ElemType,
					mapping.Pattern)
			}
		} else {
			mapping.Pattern = "simple_pointer_copy"
		}

		// Extract element type for slice conversions
		if commonField.IsSlice {
			mapping.CommonElemType = strings.TrimPrefix(commonField.ElemType, "types.")
		}

		mappings = append(mappings, mapping)
	}

	return mappings
}

func patternToString(pattern FieldPattern) string {
	switch pattern {
	case PatternSimpleCopy:
		return "simple_copy"
	case PatternSimplePointerCopy:
		return "simple_pointer_copy"
	case PatternPointerDereference:
		return "ptr_deref"
	case PatternPointerReference:
		return "ptr_ref"
	case PatternPtrSliceDereference:
		return "ptr_slice_deref"
	case PatternSliceCast:
		return "enum_slice"
	case PatternNestedStruct:
		return "nested_struct"
	case PatternTimeConversion:
		return "time_novalnumber"
	case PatternNoValStructUnwrap:
		return "novalnumber_uint32"
	case PatternEnumConversion:
		return "enum_conversion"
	default:
		return "custom"
	}
}

func generateFieldConversion(mapping FieldMappingInfo, config *GeneratorConfig) (string, error) {
	// Check for custom helper
	if mapping.Helper != "" {
		return fmt.Sprintf("\t// %s -> %s (custom helper: %s)\n\tresult.%s = %s(apiObj.%s)\n\n",
			mapping.APIField, mapping.CommonField, mapping.Helper,
			mapping.CommonField, mapping.Helper, mapping.APIField), nil
	}

	// Get pattern template
	patternDef, ok := config.Patterns[mapping.Pattern]
	if !ok {
		// Fallback to simple copy
		return fmt.Sprintf("\t// %s -> %s (pattern '%s' not found, using simple copy)\n\tif api%s.%s != nil {\n\t\tresult.%s = api%s.%s\n\t}\n\n",
			mapping.APIField, mapping.CommonField, mapping.Pattern,
			mapping.CommonField, mapping.APIField, mapping.CommonField, mapping.CommonField, mapping.APIField), nil
	}

	// Execute template
	tmpl, err := template.New("field").Parse(patternDef.Template)
	if err != nil {
		return "", fmt.Errorf("parsing template for %s: %w", mapping.Pattern, err)
	}

	// Strip pointer prefix for base type (e.g., "*uint32" -> "uint32")
	commonBaseType := strings.TrimPrefix(mapping.CommonType, "*")

	data := map[string]string{
		"APIField":       mapping.APIField,
		"CommonField":    mapping.CommonField,
		"APIType":        mapping.APIType,
		"CommonType":     mapping.CommonType,
		"CommonBaseType": commonBaseType,
		"Converter":      mapping.Converter,
		"Helper":         mapping.Helper,
		"CommonElemType": mapping.CommonElemType,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// Type parsing (reuse from original generate_converters.go)

func parseAPITypes(dir, prefix string) (map[string]*TypeInfo, error) {
	registry := &TypeRegistry{
		Structs: make(map[string]*TypeInfo),
		Aliases: make(map[string]ast.Expr),
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

	// First pass: collect all type aliases
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					name := typeSpec.Name.Name
					if prefix != "" && !strings.HasPrefix(name, prefix) {
						continue
					}

					// Check if it's NOT a struct (i.e., it's a type alias)
					if _, ok := typeSpec.Type.(*ast.StructType); !ok {
						registry.Aliases[name] = typeSpec.Type
						if *debugFlag {
							fmt.Printf("  [DEBUG] Found type alias: %s\n", name)
						}
					}
				}
			}
		}
	}
	if *debugFlag {
		fmt.Printf("  [DEBUG] Total type aliases collected: %d\n", len(registry.Aliases))
	}

	// Second pass: parse struct types with alias resolution
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					name := typeSpec.Name.Name
					if prefix != "" && !strings.HasPrefix(name, prefix) {
						continue
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					typeInfo := &TypeInfo{
						Name:   name,
						Fields: make(map[string]FieldInfo),
					}

					for _, field := range structType.Fields.List {
						if len(field.Names) == 0 {
							continue
						}
						fieldName := field.Names[0].Name
						fieldInfo := parseFieldTypeWithRegistry(field.Type, registry)
						fieldInfo.Name = fieldName
						typeInfo.Fields[fieldName] = fieldInfo
					}

					registry.Structs[name] = typeInfo
				}
			}
		}
	}

	return registry.Structs, nil
}

func parseCommonTypes(dir string) (map[string]*TypeInfo, error) {
	registry := &TypeRegistry{
		Structs: make(map[string]*TypeInfo),
		Aliases: make(map[string]ast.Expr),
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

	// First pass: collect type aliases
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					// Check if it's NOT a struct (i.e., it's a type alias)
					if _, ok := typeSpec.Type.(*ast.StructType); !ok {
						registry.Aliases[typeSpec.Name.Name] = typeSpec.Type
					}
				}
			}
		}
	}

	// Second pass: parse struct types
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					name := typeSpec.Name.Name
					typeInfo := &TypeInfo{
						Name:   name,
						Fields: make(map[string]FieldInfo),
					}

					for _, field := range structType.Fields.List {
						if len(field.Names) == 0 {
							continue
						}
						fieldName := field.Names[0].Name
						fieldInfo := parseFieldTypeWithRegistry(field.Type, registry)
						fieldInfo.Name = fieldName
						typeInfo.Fields[fieldName] = fieldInfo
					}

					registry.Structs[name] = typeInfo
				}
			}
		}
	}

	return registry.Structs, nil
}

func parseFieldType(expr ast.Expr) FieldInfo {
	info := FieldInfo{}

	switch t := expr.(type) {
	case *ast.StarExpr:
		info.IsPtr = true
		inner := parseFieldType(t.X)
		info.Type = "*" + inner.Type
		info.ElemType = inner.Type
	case *ast.ArrayType:
		info.IsSlice = true
		inner := parseFieldType(t.Elt)
		info.Type = "[]" + inner.Type
		info.ElemType = inner.Type
	case *ast.Ident:
		info.Type = t.Name
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			info.Type = x.Name + "." + t.Sel.Name
		}
	default:
		info.Type = "unknown"
	}

	return info
}

// parseFieldTypeWithRegistry parses a field type and resolves type aliases
func parseFieldTypeWithRegistry(expr ast.Expr, registry *TypeRegistry) FieldInfo {
	info := FieldInfo{}

	switch t := expr.(type) {
	case *ast.StarExpr:
		info.IsPtr = true
		inner := parseFieldTypeWithRegistry(t.X, registry)
		info.Type = inner.Type // Keep resolved type without adding *
		info.ElemType = inner.ElemType
		info.IsSlice = inner.IsSlice
		info.IsAnonymousStruct = inner.IsAnonymousStruct
		info.IsTypeAlias = inner.IsTypeAlias
		info.UnderlyingType = inner.UnderlyingType

	case *ast.ArrayType:
		info.IsSlice = true
		inner := parseFieldTypeWithRegistry(t.Elt, registry)
		info.Type = inner.Type // Keep element type
		info.ElemType = inner.Type

	case *ast.Ident:
		// Check if this identifier is a type alias
		if aliasExpr, ok := registry.Aliases[t.Name]; ok {
			// Recursively resolve the alias
			resolved := parseFieldTypeWithRegistry(aliasExpr, registry)
			// Preserve the slice/ptr nature but use the alias name for display
			info.IsSlice = resolved.IsSlice
			info.IsPtr = resolved.IsPtr
			info.ElemType = resolved.ElemType
			info.Type = t.Name // Keep the alias name
			info.IsTypeAlias = true
			// Build the underlying type string
			if resolved.IsSlice {
				info.UnderlyingType = "[]" + resolved.ElemType
			} else if resolved.IsPtr {
				info.UnderlyingType = "*" + resolved.Type
			} else {
				info.UnderlyingType = resolved.Type
			}
		} else {
			info.Type = t.Name
		}

	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			info.Type = x.Name + "." + t.Sel.Name
		}

	case *ast.StructType:
		// Anonymous struct - mark as such
		info.Type = "struct{...}"
		info.IsAnonymousStruct = true

	default:
		info.Type = "unknown"
	}

	return info
}
