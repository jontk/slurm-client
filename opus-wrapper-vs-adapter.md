This is a great architectural question. For a multi-versioned API SDK generated from OpenAPI specs, I'd recommend a combination approach, but the adapter pattern is typically the better primary choice. Here's why:
Adapter Pattern (Recommended as Primary)
The adapter pattern works well here because:

Each API version has its own generated client from the OpenAPI spec
A common interface defines what consumers of your SDK interact with
Version-specific adapters translate between the common interface and each generated client

// Common interface your SDK exposes
type UserService interface {
    GetUser(ctx context.Context, id string) (*User, error)
    CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
}

func (s *userServiceV1) GetUser(ctx context.Context, id string) (*User, error) {
        resp, err := s.client.GetUserById(ctx, id)
            if err != nil {
                        return nil, err
                            }
                                return toUser(resp), nil // transform v1 response to common model
}

// Adapter for v2 (different API shape)
type userServiceV2 struct {
        client *generated.V2Client
}
```

## Why Not Pure Wrapper?

The wrapper pattern (simple delegation with minimal transformation) works when you're just adding functionality around a single client. But with multiple API versions:

- Response/request shapes often differ between versions
- You need transformation logic, not just passthrough
- Consumers want a stable interface regardless of 18:26This is a great architectural question. For a multi-versioned API SDK generated from OpenAPI specs, I'd recommend a combination approach, but the adapter pattern is typically the better primary choice. Here's why:
Adapter Pattern (Recommended as Primary)
The adapter pattern works well here because:

Each API version has its own generated client from the OpenAPI spec
A common interface defines what consumers of your SDK interact with
Version-specific adapters translate between the common interface and each generated client

go// Common interface your SDK exposes
type UserService interface {
    GetUser(ctx context.Context, id string) (*User, error)
    CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
}

// Adapter for v1
type uversion

## Practical Structure
```
sdk/
├── client.go           # Main SDK entry point
├── users.go            # UserService interface + common types
├── internal/
│   ├── v1/
│   │   ├── generated/  # OpenAPI-generated code
│   │   └── adapter.go  # v1 adapter implementation
│   └── v2/
│       ├── generated/
│       └── adapter.go
One Consideration
If the API versions are very similar and you just want to expose the generated clients with some convenience methods, a thin wrapper might suffice. But the moment field names change, endpoints move, or response structures differ, you'll want adapters.
Would you like me to sketch out a more complete example with the factory pattern for version selection?

