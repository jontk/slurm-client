package base

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jontk/slurm-client/pkg/errors"
)

// ListOptions represents common list operation options
type ListOptions struct {
	Limit  int
	Offset int
}

// CRUDManager provides generic CRUD operations for resource managers
type CRUDManager struct {
	*BaseManager
}

// NewCRUDManager creates a new CRUD manager instance
func NewCRUDManager(version, resourceType string) *CRUDManager {
	return &CRUDManager{
		BaseManager: NewBaseManager(version, resourceType),
	}
}

// ProcessListResponse processes a list response with common pagination logic
func (c *CRUDManager) ProcessListResponse(
	items interface{},
	opts ListOptions,
	converter func(interface{}) (interface{}, error),
) ([]interface{}, int, error) {
	// Use reflection to handle slice of any type
	slice, err := convertToSlice(items)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid items type: %w", err)
	}

	// Convert items
	result := make([]interface{}, 0, len(slice))
	for i, item := range slice {
		converted, err := converter(item)
		if err != nil {
			return nil, 0, c.HandleConversionError(err, i)
		}
		result = append(result, converted)
	}

	total := len(result)

	// Apply offset
	if opts.Offset > 0 {
		if opts.Offset >= len(result) {
			return []interface{}{}, total, nil
		}
		result = result[opts.Offset:]
	}

	// Apply limit
	if opts.Limit > 0 && len(result) > opts.Limit {
		result = result[:opts.Limit]
	}

	return result, total, nil
}

// ValidatePaginationOptions validates pagination parameters
func (c *CRUDManager) ValidatePaginationOptions(opts ListOptions) error {
	if opts.Limit < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"limit must be non-negative",
			"limit", opts.Limit, nil,
		)
	}
	if opts.Offset < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"offset must be non-negative",
			"offset", opts.Offset, nil,
		)
	}
	return nil
}

// BatchOperation performs a batch operation with error collection
func (c *CRUDManager) BatchOperation(
	ctx context.Context,
	items []interface{},
	operation func(context.Context, interface{}) error,
	continueOnError bool,
) error {
	if err := c.ValidateContext(ctx); err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	var errs []error
	for i, item := range items {
		if err := operation(ctx, item); err != nil {
			wrappedErr := fmt.Errorf("%s %d: %w", c.resourceType, i, err)
			if !continueOnError {
				return wrappedErr
			}
			errs = append(errs, wrappedErr)
		}
	}

	if len(errs) > 0 {
		// Create a composite error message
		errMsg := fmt.Sprintf("Batch operation failed for %d/%d %ss: ", len(errs), len(items), c.resourceType)
		for i, err := range errs {
			if i > 0 {
				errMsg += "; "
			}
			errMsg += err.Error()
		}
		return errors.NewClientError(
			errors.ErrorCodeServerInternal,
			"Batch operation failed",
			errMsg,
		)
	}

	return nil
}

// ResourceNotFoundError creates a standardized not found error
func (c *CRUDManager) ResourceNotFoundError(identifier interface{}) error {
	return errors.NewClientError(
		errors.ErrorCodeResourceNotFound,
		fmt.Sprintf("%s not found", c.resourceType),
		fmt.Sprintf("%s '%v' not found", c.resourceType, identifier),
	)
}

// convertToSlice converts an interface{} to []interface{} using reflection
func convertToSlice(items interface{}) ([]interface{}, error) {
	if items == nil {
		return []interface{}{}, nil
	}

	v := reflect.ValueOf(items)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected slice, got %s", v.Kind())
	}

	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}

	return result, nil
}