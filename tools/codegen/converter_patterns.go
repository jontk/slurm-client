//go:build ignore

// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// converter_patterns.go defines converter patterns and their detection logic

package main

import (
	"fmt"
	"strings"
	"text/template"
)

// FieldPattern represents different conversion patterns
type FieldPattern int

const (
	PatternUnknown FieldPattern = iota
	PatternSimpleCopy           // T -> T (non-pointer, direct assignment)
	PatternSimplePointerCopy    // *T -> *T (pointer, with nil check)
	PatternPointerDereference   // *T -> T (dereference pointer)
	PatternPointerReference     // T -> *T (take address)
	PatternSliceCast            // []APIEnum -> []CommonEnum
	PatternPtrSliceDereference  // *[]T -> []T
	PatternNestedStruct         // Nested struct conversion
	PatternTimeConversion       // Unix timestamp to time.Time
	PatternNoValStructUnwrap    // NoValStruct{Set, Number/Value} -> *T
	PatternEnumConversion       // API enum value to common enum value
	PatternCustomLogic          // Requires custom implementation
	PatternSliceTypeAlias       // TypeAlias = []T -> []T (type alias for slice, direct assignment)
)

func (p FieldPattern) String() string {
	switch p {
	case PatternSimpleCopy:
		return "SimpleCopy"
	case PatternSimplePointerCopy:
		return "SimplePointerCopy"
	case PatternPointerDereference:
		return "PointerDereference"
	case PatternPointerReference:
		return "PointerReference"
	case PatternSliceCast:
		return "SliceCast"
	case PatternPtrSliceDereference:
		return "PtrSliceDereference"
	case PatternNestedStruct:
		return "NestedStruct"
	case PatternTimeConversion:
		return "TimeConversion"
	case PatternNoValStructUnwrap:
		return "NoValStructUnwrap"
	case PatternEnumConversion:
		return "EnumConversion"
	case PatternCustomLogic:
		return "CustomLogic"
	case PatternSliceTypeAlias:
		return "SliceTypeAlias"
	default:
		return "Unknown"
	}
}

// FieldInfo is declared in generate_converters_v2.go (shared between files)

// PatternDetector detects conversion patterns between field types
type PatternDetector struct {
	apiPrefix string
}

func NewPatternDetector(apiPrefix string) *PatternDetector {
	return &PatternDetector{apiPrefix: apiPrefix}
}

// DetectPattern determines the conversion pattern for a field mapping
func (pd *PatternDetector) DetectPattern(apiField, commonField FieldInfo) FieldPattern {
	// Check for NoValStruct pattern (e.g., V0044Uint64NoValStruct)
	// IMPORTANT: This must come before slice check since NoValStruct fields shouldn't be treated as slices
	if pd.isNoValStruct(apiField.Type) {
		// If the common type is time.Time, use time conversion instead of raw unwrap
		if commonField.Type == "time.Time" {
			return PatternTimeConversion
		}
		return PatternNoValStructUnwrap
	}

	// Check for time conversion (Unix timestamp struct to time.Time) - for non-NoValStruct time types
	if pd.isTimeStruct(apiField.Type) && commonField.Type == "time.Time" {
		return PatternTimeConversion
	}

	// IMPORTANT: Check slices BEFORE nested struct, because type aliases like
	// "V0042AssocShortList" have IsSlice=true but their Type name doesn't start with "[]"
	// Slice conversions
	if apiField.IsSlice && commonField.IsSlice {
		// *[]T -> []T (pointer to slice, common is just slice)
		if apiField.IsPtr && !commonField.IsPtr {
			return PatternPtrSliceDereference
		}
		// Check for type alias for slice with same element type (e.g., V0042ClusterRecFlags = []string)
		// When API type is a type alias for a slice and element types match, use direct assignment
		if apiField.IsTypeAlias && apiField.ElemType == commonField.ElemType {
			return PatternSliceTypeAlias
		}
		// []APIType -> []CommonType (enum or struct cast)
		if apiField.ElemType != commonField.ElemType {
			return PatternSliceCast
		}
		// []T -> []T (direct copy, both slices)
		if apiField.IsPtr == commonField.IsPtr {
			return PatternSimpleCopy
		}
		return PatternSimplePointerCopy
	}

	// Check for enum conversion (API enum to common enum)
	if pd.isAPIEnum(apiField.Type) && pd.isCommonEnum(commonField.Type) {
		return PatternEnumConversion
	}

	// Check for nested struct conversion (after slice check)
	if pd.isStructType(apiField.Type) && pd.isStructType(commonField.Type) {
		return PatternNestedStruct
	}

	// Pointer conversions
	if apiField.IsPtr && !commonField.IsPtr {
		return PatternPointerDereference
	}
	if !apiField.IsPtr && commonField.IsPtr {
		return PatternPointerReference
	}

	// Same type - check if pointer or not
	if apiField.Type == commonField.Type {
		// Both are pointers (with nil check)
		if apiField.IsPtr && commonField.IsPtr {
			return PatternSimplePointerCopy
		}
		// Neither are pointers (direct assignment, no nil check)
		if !apiField.IsPtr && !commonField.IsPtr {
			return PatternSimpleCopy
		}
	}

	// Default to custom logic
	return PatternCustomLogic
}

// isNoValStruct checks if type is a NoValStruct (e.g., V0044Uint64NoValStruct)
func (pd *PatternDetector) isNoValStruct(typeName string) bool {
	return strings.Contains(typeName, "NoValStruct")
}

// isTimeStruct checks if type is a time-related struct (e.g., V0044Uint64)
func (pd *PatternDetector) isTimeStruct(typeName string) bool {
	// Time fields in SLURM API often use Uint64NoValStruct with Number field
	return strings.Contains(typeName, "NoValStruct")
}

// isAPIEnum checks if type is an API enum (starts with version prefix)
func (pd *PatternDetector) isAPIEnum(typeName string) bool {
	return strings.HasPrefix(typeName, pd.apiPrefix) &&
		(strings.Contains(typeName, "State") ||
			strings.Contains(typeName, "Flags") ||
			strings.Contains(typeName, "Type"))
}

// isCommonEnum checks if type is a common enum
func (pd *PatternDetector) isCommonEnum(typeName string) bool {
	return strings.Contains(typeName, "State") ||
		strings.Contains(typeName, "Flags") ||
		strings.Contains(typeName, "Type")
}

// isStructType checks if type is a struct (not primitive, enum, or slice)
func (pd *PatternDetector) isStructType(typeName string) bool {
	// Remove pointer marker
	t := strings.TrimPrefix(typeName, "*")

	// Primitives
	primitives := []string{"string", "int", "int32", "int64", "uint32", "uint64", "float32", "float64", "bool"}
	for _, p := range primitives {
		if t == p {
			return false
		}
	}

	// Standard library types
	if strings.HasPrefix(t, "time.") {
		return false
	}

	// Slices
	if strings.HasPrefix(t, "[]") {
		return false
	}

	return true
}

// Template data for code generation
type TemplateData struct {
	APIField     string
	CommonField  string
	APIType      string
	CommonType   string
	Pattern      FieldPattern
	APIPrefix    string
	ElemType     string // For slice conversions
	Converter    string // Helper function name (if custom)
}

// GenerateConversionCode generates the conversion code for a field
func GenerateConversionCode(data TemplateData) (string, error) {
	tmpl, ok := conversionTemplates[data.Pattern]
	if !ok {
		return "", fmt.Errorf("no template for pattern %v", data.Pattern)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// Conversion templates for each pattern
var conversionTemplates = map[FieldPattern]*template.Template{
	PatternSimpleCopy: template.Must(template.New("simpleCopy").Parse(`	// {{.APIField}} -> {{.CommonField}} (direct copy)
	result.{{.CommonField}} = apiObj.{{.APIField}}

`)),

	PatternSimplePointerCopy: template.Must(template.New("simplePointer").Parse(`	// {{.APIField}} -> {{.CommonField}} (pointer copy with nil check)
	if apiObj.{{.APIField}} != nil {
		result.{{.CommonField}} = apiObj.{{.APIField}}
	}

`)),

	PatternPointerDereference: template.Must(template.New("ptrDeref").Parse(`	// {{.APIField}} -> {{.CommonField}} (dereference)
	if apiObj.{{.APIField}} != nil {
		result.{{.CommonField}} = *apiObj.{{.APIField}}
	}
`)),

	PatternPointerReference: template.Must(template.New("ptrRef").Parse(`	// {{.APIField}} -> {{.CommonField}} (take address)
	if apiObj.{{.APIField}} != "" {
		val := apiObj.{{.APIField}}
		result.{{.CommonField}} = &val
	}
`)),

	PatternSliceCast: template.Must(template.New("sliceCast").Parse(`	// {{.APIField}} -> {{.CommonField}} (slice cast)
	if apiObj.{{.APIField}} != nil && len(*apiObj.{{.APIField}}) > 0 {
		items := *apiObj.{{.APIField}}
		result := make([]{{.CommonType}}, len(items))
		for i, item := range items {
			result[i] = {{.CommonType}}(item)
		}
		result.{{.CommonField}} = result
	}
`)),

	PatternPtrSliceDereference: template.Must(template.New("ptrSliceDeref").Parse(`	// {{.APIField}} -> {{.CommonField}} (dereference slice pointer)
	if apiObj.{{.APIField}} != nil {
		result.{{.CommonField}} = *apiObj.{{.APIField}}
	}
`)),

	PatternTimeConversion: template.Must(template.New("timeConv").Parse(`	// {{.APIField}} -> {{.CommonField}} (time conversion)
	if apiObj.{{.APIField}} != nil && apiObj.{{.APIField}}.Number != nil && *apiObj.{{.APIField}}.Number > 0 {
		timestamp := time.Unix(*apiObj.{{.APIField}}.Number, 0)
		result.{{.CommonField}} = timestamp
	}
`)),

	PatternNoValStructUnwrap: template.Must(template.New("noValUnwrap").Parse(`	// {{.APIField}} -> {{.CommonField}} (unwrap NoValStruct)
	if apiObj.{{.APIField}} != nil && apiObj.{{.APIField}}.Set != nil && *apiObj.{{.APIField}}.Set && apiObj.{{.APIField}}.Number != nil {
		val := {{.CommonBaseType}}(*apiObj.{{.APIField}}.Number)
		result.{{.CommonField}} = &val
	}
`)),

	PatternEnumConversion: template.Must(template.New("enumConv").Parse(`	// {{.APIField}} -> {{.CommonField}} (enum conversion)
	if apiObj.{{.APIField}} != "" {
		result.{{.CommonField}} = {{.CommonType}}(apiObj.{{.APIField}})
	}
`)),

	PatternNestedStruct: template.Must(template.New("nestedStruct").Parse(`	// {{.APIField}} -> {{.CommonField}} (nested struct - requires custom converter)
	if apiObj.{{.APIField}} != nil {
		result.{{.CommonField}} = convert{{.APIType}}To{{.CommonType}}(apiObj.{{.APIField}})
	}
`)),

	PatternCustomLogic: template.Must(template.New("custom").Parse(`	// {{.APIField}} -> {{.CommonField}} (custom logic required)
	// TODO: Implement custom conversion from {{.APIType}} to {{.CommonType}}
`)),

	PatternSliceTypeAlias: template.Must(template.New("sliceTypeAlias").Parse(`	// {{.APIField}} -> {{.CommonField}} (slice type alias - direct assignment)
	if apiObj.{{.APIField}} != nil {
		result.{{.CommonField}} = ([]{{.ElemType}})(*apiObj.{{.APIField}})
	}
`)),
}
