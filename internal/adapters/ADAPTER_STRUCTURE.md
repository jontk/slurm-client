# Adapter Structure Guide

## Directory Structure per Version

Each version directory (v0_0_40, v0_0_41, v0_0_42) needs the following files:

```
internal/adapters/v0_0_XX/
├── adapter.go                    ✅ Created (main adapter aggregator)
├── qos_adapter.go               ⏳ Pending
├── qos_converters.go            ⏳ Pending
├── job_adapter.go               ⏳ Pending
├── job_converters.go            ⏳ Pending
├── partition_adapter.go         ⏳ Pending
├── partition_converters.go      ⏳ Pending
├── node_adapter.go              ⏳ Pending
├── node_converters.go           ⏳ Pending
├── account_adapter.go           ⏳ Pending
├── account_converters.go        ⏳ Pending
├── user_adapter.go              ⏳ Pending
├── reservation_adapter.go       ⏳ Pending
└── association_adapter.go       ⏳ Pending
```

## Implementation Pattern

Each adapter follows the same pattern as v0_0_43:

### 1. Adapter Structure
```go
type XxxAdapter struct {
    *base.XxxBaseManager
    client  *api.ClientWithResponses
}
```

### 2. Constructor
```go
func NewXxxAdapter(client *api.ClientWithResponses) *XxxAdapter {
    return &XxxAdapter{
        XxxBaseManager: base.NewXxxBaseManager("v0.0.XX"),
        client:         client,
    }
}
```

### 3. Required Methods
Each adapter must implement its interface methods from `internal/adapters/common/interfaces.go`

### 4. Converter Functions
Each converter file provides:
- `convertFromAPIXxx` - Converts API type to common type
- `convertToAPIXxx` - Converts common type to API type
- Helper functions for nested structures

## Version-Specific Considerations

### v0_0_40
- Oldest version - may lack some fields
- Need to handle missing fields gracefully
- Some features might not be supported

### v0_0_41
- Intermediate version
- May have partial feature support
- Bridge between v0_0_40 and v0_0_42

### v0_0_42
- More complete than earlier versions
- Closer to v0_0_43 structure
- May still lack association manager in API

## Implementation Order

1. **QoS** - Already well-tested pattern
2. **Job** - Complex but essential
3. **Partition** - Simpler structure
4. **Node** - Straightforward
5. **Account** - Medium complexity
6. **User** - Simple structure
7. **Reservation** - Less commonly used
8. **Association** - May not exist in older versions

## Notes for Implementers

- Use base managers from `internal/adapters/base/`
- Follow error handling patterns from v0_0_43
- Include validation using base manager methods
- Handle nil pointers carefully in converters
- Test with actual API responses when possible