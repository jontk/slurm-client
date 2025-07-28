package base

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/pkg/errors"
)

// BaseManager provides common functionality for all managers
type BaseManager struct {
	version      string
	resourceType string
}

// NewBaseManager creates a new base manager instance
func NewBaseManager(version, resourceType string) *BaseManager {
	return &BaseManager{
		version:      version,
		resourceType: resourceType,
	}
}

// ValidateContext ensures the context is valid
func (b *BaseManager) ValidateContext(ctx context.Context) error {
	if ctx == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"context is required",
			"ctx", nil, nil,
		)
	}
	return nil
}

// ValidateResourceName validates a resource name is not empty
func (b *BaseManager) ValidateResourceName(name string, fieldName string) error {
	if name == "" {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			fmt.Sprintf("%s name is required", b.resourceType),
			fieldName, name, nil,
		)
	}
	return nil
}

// ValidateResourceID validates a resource ID
func (b *BaseManager) ValidateResourceID(id interface{}, fieldName string) error {
	if id == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			fmt.Sprintf("%s ID is required", b.resourceType),
			fieldName, id, nil,
		)
	}

	// Check for zero values using reflection
	v := reflect.ValueOf(id)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("%s ID is required", b.resourceType),
				fieldName, id, nil,
			)
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("%s ID must be greater than 0", b.resourceType),
				fieldName, id, nil,
			)
		}
	case reflect.String:
		if v.String() == "" {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("%s ID cannot be empty", b.resourceType),
				fieldName, id, nil,
			)
		}
	}

	return nil
}

// ValidateNonNegative validates that a numeric value is non-negative
func (b *BaseManager) ValidateNonNegative(value int, fieldName string) error {
	if value < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			fmt.Sprintf("%s must be non-negative", fieldName),
			fieldName, value, nil,
		)
	}
	return nil
}

// HandleAPIError wraps an API error with version information
func (b *BaseManager) HandleAPIError(err error) error {
	return common.WrapAndEnhanceError(err, b.version)
}

// HandleConversionError creates a standardized conversion error
func (b *BaseManager) HandleConversionError(err error, resourceID interface{}) error {
	return common.HandleConversionError(err, b.resourceType, resourceID)
}

// CheckClientInitialized verifies the client is initialized
func (b *BaseManager) CheckClientInitialized(client interface{}) error {
	return common.CheckClientInitialized(client)
}

// CheckNilResponse checks if a response is nil
func (b *BaseManager) CheckNilResponse(response interface{}, operation string) error {
	return common.CheckNilResponse(response, operation)
}

// GetVersion returns the API version
func (b *BaseManager) GetVersion() string {
	return b.version
}

// GetResourceType returns the resource type
func (b *BaseManager) GetResourceType() string {
	return b.resourceType
}