//go:build ignore

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// generate_converters.go generates type converter functions for each API version
// Usage: go run generate_converters.go [version]
// Example: go run generate_converters.go v0_0_44

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// Config holds the generator configuration
type Config struct {
	Version    string
	APIPackage string
	APIPrefix  string
	OutputDir  string
}

// FieldMapping defines how a field maps between API and common types
type FieldMapping struct {
	APIField    string
	CommonField string
	APIType     string
	CommonType  string
	Converter   string // Name of converter template to use
	Direct      bool   // Direct assignment (same pointer types)
}

// EntityMapping defines type conversion for an entity
type EntityMapping struct {
	APIType    string
	CommonType string
	Fields     []FieldMapping
}

// TypeInfo holds parsed type information
type TypeInfo struct {
	Name   string
	Fields map[string]FieldInfo
}

// FieldInfo holds parsed field information
type FieldInfo struct {
	Name     string
	Type     string
	IsPtr    bool
	IsSlice  bool
	ElemType string
}

var (
	versionFlag = flag.String("version", "", "API version to generate (e.g., v0_0_44)")
	allFlag     = flag.Bool("all", false, "Generate for all versions")
	dryRunFlag  = flag.Bool("dry-run", false, "Print generated code without writing files")
)

func main() {
	flag.Parse()

	versions := []string{"v0_0_40", "v0_0_41", "v0_0_42", "v0_0_43", "v0_0_44"}

	if *versionFlag != "" {
		versions = []string{*versionFlag}
	} else if !*allFlag {
		fmt.Println("Usage: go run generate_converters.go -version=v0_0_44")
		fmt.Println("       go run generate_converters.go -all")
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, version := range versions {
		fmt.Printf("Generating converters for %s...\n", version)
		if err := generateConverters(version); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating %s: %v\n", version, err)
			continue
		}
		fmt.Printf("  Done: %s\n", version)
	}
}

func generateConverters(version string) error {
	// Parse API types from the generated client
	apiDir := filepath.Join("internal", "api", version)
	apiTypes, err := parseAPITypes(apiDir, getAPIPrefix(version))
	if err != nil {
		return fmt.Errorf("parsing API types: %w", err)
	}

	// Parse common types
	commonDir := filepath.Join("internal", "common", "types")
	commonTypes, err := parseCommonTypes(commonDir)
	if err != nil {
		return fmt.Errorf("parsing common types: %w", err)
	}

	// Generate converters for each entity
	entities := []string{"Account", "Association", "Job", "Node", "Partition", "Reservation", "QoS", "User", "Cluster"}

	for _, entity := range entities {
		code, err := generateEntityConverters(version, entity, apiTypes, commonTypes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: Could not generate %s converters: %v\n", entity, err)
			continue
		}

		if code == "" {
			continue
		}

		outputPath := filepath.Join("internal", "adapters", version, strings.ToLower(entity)+"_converters.gen.go")

		if *dryRunFlag {
			fmt.Printf("--- %s ---\n%s\n", outputPath, code)
		} else {
			if err := os.WriteFile(outputPath, []byte(code), 0644); err != nil {
				return fmt.Errorf("writing %s: %w", outputPath, err)
			}
		}
	}

	return nil
}

func getAPIPrefix(version string) string {
	// v0_0_44 -> V0044
	parts := strings.Split(version, "_")
	if len(parts) != 3 {
		return "V0000"
	}
	return fmt.Sprintf("V%s%s%s", parts[0][1:], parts[1], parts[2])
}

func parseAPITypes(dir, prefix string) (map[string]*TypeInfo, error) {
	types := make(map[string]*TypeInfo)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

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
					if !strings.HasPrefix(name, prefix) {
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
						fieldInfo := parseFieldType(field.Type)
						fieldInfo.Name = fieldName
						typeInfo.Fields[fieldName] = fieldInfo
					}

					types[name] = typeInfo
				}
			}
		}
	}

	return types, nil
}

func parseCommonTypes(dir string) (map[string]*TypeInfo, error) {
	types := make(map[string]*TypeInfo)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

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
						fieldInfo := parseFieldType(field.Type)
						fieldInfo.Name = fieldName
						typeInfo.Fields[fieldName] = fieldInfo
					}

					types[name] = typeInfo
				}
			}
		}
	}

	return types, nil
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

func generateEntityConverters(version, entity string, apiTypes, commonTypes map[string]*TypeInfo) (string, error) {
	prefix := getAPIPrefix(version)
	apiTypeName := prefix + getAPITypeName(entity)
	commonTypeName := entity

	apiType, ok := apiTypes[apiTypeName]
	if !ok {
		return "", fmt.Errorf("API type %s not found", apiTypeName)
	}

	commonType, ok := commonTypes[commonTypeName]
	if !ok {
		return "", fmt.Errorf("common type %s not found", commonTypeName)
	}

	// Generate field mappings
	mappings := generateFieldMappings(apiType, commonType, prefix)

	// Generate code using template
	data := struct {
		Version      string
		Package      string
		APIPrefix    string
		Entity       string
		APIType      string
		CommonType   string
		Mappings     []FieldMapping
		HasMappings  bool
	}{
		Version:     version,
		Package:     version,
		APIPrefix:   prefix,
		Entity:      entity,
		APIType:     apiTypeName,
		CommonType:  commonTypeName,
		Mappings:    mappings,
		HasMappings: len(mappings) > 0,
	}

	var buf bytes.Buffer
	if err := converterTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

func getAPITypeName(entity string) string {
	// Map entity names to API type suffixes
	// These vary by entity and sometimes by version
	switch entity {
	case "Job":
		return "JobInfo"
	case "Association":
		return "Assoc"
	case "Cluster":
		return "ClusterRec"
	case "Reservation":
		return "ReservationInfo"
	case "QoS":
		return "Qos"
	case "Partition":
		return "PartitionInfo"
	default:
		return entity
	}
}

func generateFieldMappings(apiType, commonType *TypeInfo, prefix string) []FieldMapping {
	var mappings []FieldMapping

	// Known field name transformations
	transforms := map[string]string{
		"Id":         "ID",
		"JobId":      "JobID",
		"UserId":     "UserID",
		"GroupId":    "GroupID",
		"Cpus":       "CPUs",
		"JobState":   "State",
	}

	// Sort field names for consistent output
	var apiFields []string
	for name := range apiType.Fields {
		apiFields = append(apiFields, name)
	}
	sort.Strings(apiFields)

	for _, apiFieldName := range apiFields {
		apiField := apiType.Fields[apiFieldName]

		// Skip internal/computed fields
		if strings.HasPrefix(apiFieldName, "_") {
			continue
		}

		// Find corresponding common field
		commonFieldName := apiFieldName
		if transformed, ok := transforms[apiFieldName]; ok {
			commonFieldName = transformed
		}

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
			continue
		}

		mapping := FieldMapping{
			APIField:    apiFieldName,
			CommonField: commonFieldName,
			APIType:     apiField.Type,
			CommonType:  commonField.Type,
		}

		// Determine conversion strategy
		mapping.Direct = canDirectAssign(apiField, commonField)
		if !mapping.Direct {
			mapping.Converter = determineConverter(apiField, commonField, prefix)
		}

		mappings = append(mappings, mapping)
	}

	return mappings
}

func canDirectAssign(api, common FieldInfo) bool {
	// Same type
	if api.Type == common.Type {
		return true
	}
	// Both are same pointer type
	if api.IsPtr && common.IsPtr && api.ElemType == common.ElemType {
		return true
	}
	return false
}

func determineConverter(api, common FieldInfo, prefix string) string {
	apiType := api.Type
	_ = prefix // May be used in future for type-specific converters

	// Pointer to value
	if api.IsPtr && !common.IsPtr {
		return "ptrToValue"
	}

	// Value to pointer
	if !api.IsPtr && common.IsPtr {
		return "valueToPtr"
	}

	// Slice type conversions
	if api.IsSlice && common.IsSlice {
		// State slices
		if strings.Contains(apiType, "State") {
			return "stateSlice"
		}
		// Coord slices
		if strings.Contains(apiType, "Coord") {
			return "coordSlice"
		}
		// Flag slices
		if strings.Contains(apiType, "Flags") {
			return "flagsSlice"
		}
	}

	return "custom"
}

var converterTemplate = template.Must(template.New("converter").Parse(`// Code generated by generate_converters.go. DO NOT EDIT.
// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package {{.Package}}

import (
	api "github.com/jontk/slurm-client/internal/openapi/{{.Version}}"
	types "github.com/jontk/slurm-client/api"
)

{{if .HasMappings}}
// convertAPI{{.Entity}}ToCommon converts API {{.APIType}} to common {{.CommonType}}
func convertAPI{{.Entity}}ToCommon(apiObj api.{{.APIType}}) *types.{{.CommonType}} {
	result := &types.{{.CommonType}}{}
	{{range .Mappings}}
	{{if .Direct}}// {{.APIField}} -> {{.CommonField}} (direct)
	result.{{.CommonField}} = apiObj.{{.APIField}}
	{{else if eq .Converter "ptrToValue"}}// {{.APIField}} -> {{.CommonField}} (ptr to value)
	if apiObj.{{.APIField}} != nil {
		result.{{.CommonField}} = *apiObj.{{.APIField}}
	}
	{{else if eq .Converter "valueToPtr"}}// {{.APIField}} -> {{.CommonField}} (value to ptr)
	v{{.APIField}} := apiObj.{{.APIField}}
	result.{{.CommonField}} = &v{{.APIField}}
	{{else if eq .Converter "stateSlice"}}// {{.APIField}} -> {{.CommonField}} (state slice)
	if len(apiObj.{{.APIField}}) > 0 {
		_ = apiObj.{{.APIField}} // state slice conversion - implement based on type
		// result.{{.CommonField}} = converted states
	}
	{{else if eq .Converter "coordSlice"}}// {{.APIField}} -> {{.CommonField}} (coord slice)
	if apiObj.{{.APIField}} != nil && len(*apiObj.{{.APIField}}) > 0 {
		coords := make([]types.Coord, len(*apiObj.{{.APIField}}))
		for i, c := range *apiObj.{{.APIField}} {
			coords[i] = types.Coord{Name: c.Name}
		}
		result.{{.CommonField}} = coords
	}
	{{else}}// {{.APIField}} -> {{.CommonField}} (needs custom conversion: {{.Converter}})
	// TODO: Implement conversion from {{.APIType}} to {{.CommonType}}
	{{end}}
	{{end}}
	return result
}
{{end}}
`))
