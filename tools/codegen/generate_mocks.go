//go:build ignore

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package main generates mock builders for SLURM API types from OpenAPI specifications.
//
// Architecture:
//
// The generator uses a type detector pattern to handle version-specific differences:
//
//   - TypeDetector interface: Defines how to detect and classify field types
//   - BaseTypeDetector: Implements common detection logic shared across versions
//   - V0040TypeDetector: Handles v0.0.40-specific patterns (_no_val suffix)
//   - V0042TypeDetector: Handles v0.0.42+ patterns (_no_val_struct suffix)
//
// This architecture provides:
//   - Clear separation of version-specific logic
//   - Easier to add new version support
//   - Better maintainability through smaller, focused functions
//   - Extensibility for future API versions
//
// Usage:
//
//	go run tools/codegen/generate_mocks.go v0.0.42
//
// The generator will:
//  1. Parse the OpenAPI spec for the specified version
//  2. Use the appropriate TypeDetector to analyze field types
//  3. Generate builder files with With* methods for each field
//  4. Handle NoVal, array types, and primitive types automatically
//
// Generated builders provide type-safe, fluent interfaces for creating mock objects in tests.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// OpenAPI spec structures
type OpenAPISpec struct {
	Components Components `json:"components"`
}

type Components struct {
	Schemas map[string]Schema `json:"schemas"`
}

type Schema struct {
	Type        string              `json:"type"`
	Description string              `json:"description"`
	Properties  map[string]Property `json:"properties"`
	Items       *Property           `json:"items"`
	Enum        []string            `json:"enum"`
}

type Property struct {
	Type        string    `json:"type"`
	Format      string    `json:"format"`
	Description string    `json:"description"`
	Ref         string    `json:"$ref"`
	Items       *Property `json:"items"`
	Enum        []string  `json:"enum"`
}

// Model represents a type to generate factories for
type Model struct {
	Version        string
	PackageVersion string // e.g., "v0_0_40"
	Name           string // e.g., "JobInfo"
	TypeName       string // e.g., "V0040JobInfo"
	SchemaName     string // e.g., "v0.0.40_job_info"
	Fields         []Field
}

// Field represents a single field in a model
type Field struct {
	Name        string // Go field name (e.g., "Cpus")
	JSONName    string // JSON field name (e.g., "cpus")
	MethodName  string // Method name for With* (e.g., "Cpus")
	GoType      string // Go type for the method parameter
	Description string
	IsPointer   bool
	IsNoVal     bool
	NoValType   string // e.g., "V0040Uint32NoVal"
	IsArray     bool
	ElementType string // For arrays, the element type (e.g., "string")
	ArrayType   string // For enum arrays, the generated type name (e.g., "V0043JobInfoFlags")
	IsSimple    bool   // Simple types like string, int32
}

// DetectionContext provides all information needed for type detection
type DetectionContext struct {
	JSONName       string
	Property       Property
	PkgVersion     string
	ParentTypeName string
	Spec           *OpenAPISpec
}

// TypeDetector knows how to detect and classify field types from OpenAPI properties
type TypeDetector interface {
	// DetectFieldType analyzes a property and returns field metadata
	DetectFieldType(ctx DetectionContext) (*Field, error)
}

// Common errors for type detection
var (
	ErrInlineObjectSkipped = fmt.Errorf("skipping inline object type")
	ErrUnsupportedType     = fmt.Errorf("unsupported property type")
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run generate_mocks.go <version>")
	}

	version := os.Args[1]
	specFile := filepath.Join("openapi-specs", fmt.Sprintf("slurm-%s.json", version))

	log.Printf("Generating mock builders for %s from %s", version, specFile)

	// Parse OpenAPI spec
	spec, err := parseSpec(specFile)
	if err != nil {
		log.Fatalf("Failed to parse spec: %v", err)
	}

	// Normalize version for package names (e.g., "v0.0.40" -> "v0_0_40")
	pkgVersion := normalizeVersion(version)

	generated := 0
	for schemaName, schema := range spec.Components.Schemas {
		// Check if this is a model we want to generate for
		if !shouldGenerateModel(schemaName, version) {
			continue
		}

		model, err := buildModel(version, pkgVersion, schemaName, schema, spec)
		if err != nil {
			log.Printf("Warning: Failed to build model for %s: %v", schemaName, err)
			continue
		}

		if err := generateFactory(model); err != nil {
			log.Printf("Warning: Failed to generate factory for %s: %v", schemaName, err)
			continue
		}

		log.Printf("Generated factory for %s", model.TypeName)
		generated++
	}

	if generated == 0 {
		log.Fatal("No factories generated")
	}

	log.Printf("Successfully generated %d mock factories for %s", generated, version)
}

func shouldGenerateModel(schemaName, version string) bool {
	// Generate for main entity types
	// Schema names use dots like "v0.0.40_job_info"
	exactMatches := []string{
		fmt.Sprintf("%s_job_info", version),
		fmt.Sprintf("%s_node", version),
		fmt.Sprintf("%s_partition_info", version),
	}

	for _, match := range exactMatches {
		if schemaName == match {
			return true
		}
	}

	return false
}

func parseSpec(filename string) (*OpenAPISpec, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file: %w", err)
	}

	var spec OpenAPISpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &spec, nil
}

// BaseTypeDetector implements common type detection logic shared across versions
type BaseTypeDetector struct {
	NoValPatterns []string // Patterns to match for NoVal types (e.g., "_no_val", "_no_val_struct")
}

// DetectFieldType analyzes a property and returns field metadata
func (d *BaseTypeDetector) DetectFieldType(ctx DetectionContext) (*Field, error) {
	field := d.initializeField(ctx)

	// Handle $ref types
	if ctx.Property.Ref != "" {
		return d.detectRefType(ctx, field)
	}

	// Skip inline object types (complex nested structs)
	if ctx.Property.Type == "object" {
		return nil, ErrInlineObjectSkipped
	}

	// Handle inline array types
	if ctx.Property.Type == "array" {
		return d.detectInlineArrayType(ctx, field)
	}

	// Handle primitive types
	if ctx.Property.Type != "" {
		return d.detectPrimitiveType(ctx, field)
	}

	return nil, ErrUnsupportedType
}

// initializeField creates a field with common initial values
func (d *BaseTypeDetector) initializeField(ctx DetectionContext) *Field {
	return &Field{
		Name:        toGoName(ctx.JSONName),
		JSONName:    ctx.JSONName,
		MethodName:  toGoName(ctx.JSONName),
		Description: ctx.Property.Description,
		IsPointer:   true, // Most OpenAPI fields are pointers
	}
}

// detectRefType handles $ref types (references to other schemas)
func (d *BaseTypeDetector) detectRefType(ctx DetectionContext, field *Field) (*Field, error) {
	refSchema := resolveRef(ctx.Property.Ref, ctx.Spec)
	if refSchema == nil {
		return nil, fmt.Errorf("failed to resolve ref: %s", ctx.Property.Ref)
	}

	refName := extractRefName(ctx.Property.Ref)

	// Check if it's a NoVal type
	if d.isNoValType(refName) {
		return d.buildNoValField(ctx, field, refName, refSchema)
	}

	// Check if it's an array type (like job_state)
	if refSchema.Type == "array" {
		return d.buildRefArrayField(ctx, field, refSchema)
	}

	// Other ref types - skip for now as they're complex nested types
	return nil, fmt.Errorf("unsupported ref type: %s", refName)
}

// isNoValType checks if a ref name matches any NoVal pattern
func (d *BaseTypeDetector) isNoValType(refName string) bool {
	for _, pattern := range d.NoValPatterns {
		if strings.HasSuffix(refName, pattern) {
			return true
		}
	}
	return false
}

// buildNoValField constructs a field for NoVal wrapper types
func (d *BaseTypeDetector) buildNoValField(ctx DetectionContext, field *Field, refName string, refSchema *Schema) (*Field, error) {
	field.IsNoVal = true
	field.NoValType = toGoTypeName(refName, ctx.PkgVersion)
	// Determine the number type from the NoVal schema
	field.GoType = getNoValNumberType(refSchema)
	return field, nil
}

// buildRefArrayField constructs a field for reference array types
func (d *BaseTypeDetector) buildRefArrayField(ctx DetectionContext, field *Field, refSchema *Schema) (*Field, error) {
	field.IsArray = true
	if refSchema.Items != nil {
		// Check if items are a ref to another complex type
		if refSchema.Items.Ref != "" {
			// This is a complex array type (array of structs), skip it
			return nil, fmt.Errorf("array with complex items (ref): %s", refSchema.Items.Ref)
		}
		// Check if items have enum - these are simple string enums
		if len(refSchema.Items.Enum) > 0 {
			field.ElementType = "string"
		} else {
			field.ElementType = mapPrimitiveType(refSchema.Items.Type, refSchema.Items.Format)
		}
	} else {
		field.ElementType = "string"
	}
	field.GoType = field.ElementType
	return field, nil
}

// detectInlineArrayType handles inline array types (not references)
func (d *BaseTypeDetector) detectInlineArrayType(ctx DetectionContext, field *Field) (*Field, error) {
	field.IsArray = true
	if ctx.Property.Items != nil {
		// Check if items are a ref to another complex type
		if ctx.Property.Items.Ref != "" {
			// This is a complex array type (array of structs), skip it
			return nil, fmt.Errorf("inline array with complex items (ref): %s", ctx.Property.Items.Ref)
		}
		// Check if items have enum - these are special enum array types
		// oapi-codegen generates type aliases like V0043JobInfoFlags for these
		if len(ctx.Property.Items.Enum) > 0 {
			field.ElementType = "string"
			// Generate the type name for this enum array
			// Pattern: {ParentTypeName}{FieldName} (e.g., V0043JobInfoFlags)
			field.ArrayType = ctx.ParentTypeName + field.Name
		} else {
			field.ElementType = mapPrimitiveType(ctx.Property.Items.Type, ctx.Property.Items.Format)
		}
	} else {
		field.ElementType = "string"
	}
	field.GoType = field.ElementType
	return field, nil
}

// detectPrimitiveType handles primitive types (string, int, bool, etc.)
func (d *BaseTypeDetector) detectPrimitiveType(ctx DetectionContext, field *Field) (*Field, error) {
	field.IsSimple = true
	field.GoType = mapPrimitiveType(ctx.Property.Type, ctx.Property.Format)
	return field, nil
}

// V0040TypeDetector handles v0.0.40-specific type detection
type V0040TypeDetector struct {
	*BaseTypeDetector
}

// NewV0040TypeDetector creates a detector for v0.0.40 API
func NewV0040TypeDetector() *V0040TypeDetector {
	return &V0040TypeDetector{
		BaseTypeDetector: &BaseTypeDetector{
			NoValPatterns: []string{"_no_val"}, // v0.0.40 uses _no_val suffix
		},
	}
}

// V0042TypeDetector handles v0.0.42+ type detection
type V0042TypeDetector struct {
	*BaseTypeDetector
}

// NewV0042TypeDetector creates a detector for v0.0.42+ APIs
func NewV0042TypeDetector() *V0042TypeDetector {
	return &V0042TypeDetector{
		BaseTypeDetector: &BaseTypeDetector{
			// v0.0.42+ uses both _no_val and _no_val_struct suffixes
			NoValPatterns: []string{"_no_val", "_no_val_struct"},
		},
	}
}

// getDetectorForVersion returns the appropriate type detector for a given version
func getDetectorForVersion(pkgVersion string) TypeDetector {
	// v0.0.42, v0.0.43, v0.0.44 use _no_val_struct pattern
	switch pkgVersion {
	case "v0_0_42", "v0_0_43", "v0_0_44":
		return NewV0042TypeDetector()
	default:
		// v0.0.40 and others use _no_val pattern
		return NewV0040TypeDetector()
	}
}

func buildModel(version, pkgVersion, schemaName string, schema Schema, spec *OpenAPISpec) (*Model, error) {
	if schema.Type != "object" {
		return nil, fmt.Errorf("schema is not an object type")
	}

	// Extract model name from schema name
	// Schema uses dots: "v0.0.40_job_info"
	// Package uses underscores: "v0_0_40"
	// Need to trim version prefix with dots
	namePart := strings.TrimPrefix(schemaName, version+"_")
	modelName := toGoName(namePart)

	// Build type name (e.g., "V0040JobInfo")
	typeName := strings.ToUpper(pkgVersion[0:1]) + strings.ReplaceAll(pkgVersion[1:], "_", "") + modelName

	model := &Model{
		Version:        version,
		PackageVersion: pkgVersion,
		Name:           modelName,
		TypeName:       typeName,
		SchemaName:     schemaName,
		Fields:         []Field{},
	}

	// Process all properties
	for jsonName, prop := range schema.Properties {
		field, err := buildField(jsonName, prop, pkgVersion, typeName, spec)
		if err != nil {
			log.Printf("Warning: Skipping field %s.%s: %v", schemaName, jsonName, err)
			continue
		}
		model.Fields = append(model.Fields, field)
	}

	// Sort fields alphabetically for consistent output
	sort.Slice(model.Fields, func(i, j int) bool {
		return model.Fields[i].Name < model.Fields[j].Name
	})

	return model, nil
}

func buildField(jsonName string, prop Property, pkgVersion, parentTypeName string, spec *OpenAPISpec) (Field, error) {
	// Get the appropriate detector for this version
	detector := getDetectorForVersion(pkgVersion)

	// Create detection context
	ctx := DetectionContext{
		JSONName:       jsonName,
		Property:       prop,
		PkgVersion:     pkgVersion,
		ParentTypeName: parentTypeName,
		Spec:           spec,
	}

	// Delegate to detector
	field, err := detector.DetectFieldType(ctx)
	if err != nil {
		return Field{}, err
	}

	return *field, nil
}

func getNoValNumberType(schema *Schema) string {
	// Check the number field's format in the NoVal schema
	if schema.Properties != nil {
		if numberProp, ok := schema.Properties["number"]; ok {
			if numberProp.Type == "number" {
				if numberProp.Format == "float" {
					return "float32"
				}
				if numberProp.Format == "double" {
					return "float64"
				}
				return "float64"
			}
			if numberProp.Type == "integer" {
				// Check format to determine correct integer size
				switch numberProp.Format {
				case "int32":
					return "int32"
				case "int16":
					return "int16"
				case "int64":
					return "int64"
				default:
					return "int64" // default to int64 if format not specified
				}
			}
		}
	}
	return "int64" // default
}

func resolveRef(ref string, spec *OpenAPISpec) *Schema {
	// Extract schema name from ref like "#/components/schemas/v0_0_40_job_state"
	parts := strings.Split(ref, "/")
	if len(parts) != 4 {
		return nil
	}
	schemaName := parts[3]
	schema, ok := spec.Components.Schemas[schemaName]
	if !ok {
		return nil
	}
	return &schema
}

func extractRefName(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func toGoName(jsonName string) string {
	// Convert snake_case to PascalCase
	parts := strings.Split(jsonName, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][0:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

func toGoTypeName(schemaName, pkgVersion string) string {
	// Convert schema name to Go type name
	// Schema uses dots: "v0.0.40_uint32_no_val"
	// Package uses underscores: "v0_0_40"
	// Need to extract just the type part after the version

	// Convert version to match schema format (v0_0_40 -> v0.0.40)
	versionWithDots := strings.ReplaceAll(pkgVersion, "_", ".")
	namePart := strings.TrimPrefix(schemaName, versionWithDots+"_")
	name := toGoName(namePart)
	prefix := strings.ToUpper(pkgVersion[0:1]) + strings.ReplaceAll(pkgVersion[1:], "_", "")
	return prefix + name
}

func mapPrimitiveType(typeStr, format string) string {
	if typeStr == "integer" {
		if format == "int32" {
			return "int32"
		}
		if format == "int64" {
			return "int64"
		}
		return "int"
	}
	if typeStr == "string" {
		return "string"
	}
	if typeStr == "boolean" {
		return "bool"
	}
	if typeStr == "number" {
		if format == "float" {
			return "float32"
		}
		return "float64"
	}
	return "interface{}"
}

func normalizeVersion(version string) string {
	return strings.ReplaceAll(version, ".", "_")
}

func generateFactory(model *Model) error {
	// Create output directory
	outputDir := filepath.Join("tests", "mocks", "generated", model.PackageVersion)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate factory file
	outputFile := filepath.Join(outputDir, strings.ToLower(model.Name)+"_factory.go")
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	tmpl := template.Must(template.New("factory").Funcs(templateFuncs).Parse(factoryTemplate))
	if err := tmpl.Execute(f, model); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

var templateFuncs = template.FuncMap{
	"lower": strings.ToLower,
}

const factoryTemplate = `// Code generated by tools/codegen/generate_mocks.go - DO NOT EDIT
// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package {{.PackageVersion}}

import (
	{{.PackageVersion}} "github.com/jontk/slurm-client/internal/openapi/{{.PackageVersion}}"
)

// {{.Name}}Builder provides a fluent interface for building {{.TypeName}} instances for testing.
type {{.Name}}Builder struct {
	obj *{{.PackageVersion}}.{{.TypeName}}
}

// New{{.Name}} creates a new {{.Name}}Builder for building {{.TypeName}} instances.
// Use the With* methods to set fields, then call Build() to get the final object.
func New{{.Name}}() *{{.Name}}Builder {
	return &{{.Name}}Builder{
		obj: &{{.PackageVersion}}.{{.TypeName}}{},
	}
}

// Build returns the constructed {{.TypeName}} instance.
func (b *{{.Name}}Builder) Build() *{{.PackageVersion}}.{{.TypeName}} {
	return b.obj
}

{{range .Fields}}
{{if .IsNoVal}}
// With{{.MethodName}} sets the {{.JSONName}} field with proper NoVal wrapping.
// The value will be marked as "set" with the provided number.
func (b *{{$.Name}}Builder) With{{.MethodName}}(value {{.GoType}}) *{{$.Name}}Builder {
	setTrue := true
	b.obj.{{.Name}} = &{{$.PackageVersion}}.{{.NoValType}}{
		Set:    &setTrue,
		Number: &value,
	}
	return b
}
{{else if .IsArray}}
// With{{.MethodName}} sets the {{.JSONName}} field.
// The value will be wrapped in an array as required by the API.
func (b *{{$.Name}}Builder) With{{.MethodName}}(value {{.ElementType}}) *{{$.Name}}Builder {
	{{if .ArrayType -}}
	b.obj.{{.Name}} = &[]{{$.PackageVersion}}.{{.ArrayType}}{ {{- $.PackageVersion}}.{{.ArrayType}}(value)}
	{{- else -}}
	b.obj.{{.Name}} = &[]{{.ElementType}}{value}
	{{- end}}
	return b
}
{{else if .IsSimple}}
// With{{.MethodName}} sets the {{.JSONName}} field.
func (b *{{$.Name}}Builder) With{{.MethodName}}(value {{.GoType}}) *{{$.Name}}Builder {
	{{if .IsPointer}}
	b.obj.{{.Name}} = &value
	{{else}}
	b.obj.{{.Name}} = value
	{{end}}
	return b
}
{{end}}
{{end}}
`
