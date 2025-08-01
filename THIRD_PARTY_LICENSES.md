# Third-Party Licenses

This project uses the following third-party dependencies:

## Direct Dependencies

### github.com/getkin/kin-openapi
- **Version**: v0.125.0
- **License**: MIT
- **Purpose**: OpenAPI 3.0 implementation
- **URL**: https://github.com/getkin/kin-openapi

### github.com/google/uuid
- **Version**: v1.5.0
- **License**: BSD-3-Clause
- **Purpose**: UUID generation
- **URL**: https://github.com/google/uuid

### github.com/gorilla/mux
- **Version**: v1.8.1
- **License**: BSD-3-Clause
- **Purpose**: HTTP router and URL matcher
- **URL**: https://github.com/gorilla/mux

### github.com/oapi-codegen/runtime
- **Version**: v1.1.1
- **License**: MIT
- **Purpose**: Runtime support for oapi-codegen
- **URL**: https://github.com/oapi-codegen/runtime

### github.com/stretchr/testify
- **Version**: v1.8.4
- **License**: MIT
- **Purpose**: Testing toolkit
- **URL**: https://github.com/stretchr/testify

## Indirect Dependencies

### github.com/apapsch/go-jsonmerge/v2
- **Version**: v2.0.0
- **License**: MIT

### github.com/davecgh/go-spew
- **Version**: v1.1.1
- **License**: ISC

### github.com/go-openapi/jsonpointer
- **Version**: v0.20.2
- **License**: Apache-2.0

### github.com/go-openapi/swag
- **Version**: v0.22.9
- **License**: Apache-2.0

### github.com/invopop/yaml
- **Version**: v0.2.0
- **License**: MIT

### github.com/josharian/intern
- **Version**: v1.0.0
- **License**: MIT

### github.com/mailru/easyjson
- **Version**: v0.7.7
- **License**: MIT

### github.com/mohae/deepcopy
- **Version**: v0.0.0-20170929034955-c48cc78d4826
- **License**: MIT

### github.com/perimeterx/marshmallow
- **Version**: v1.1.5
- **License**: MIT

### github.com/pmezard/go-difflib
- **Version**: v1.0.0
- **License**: BSD-3-Clause

### github.com/stretchr/objx
- **Version**: v0.5.0
- **License**: MIT

### gopkg.in/yaml.v2
- **Version**: v2.4.0
- **License**: MIT

### gopkg.in/yaml.v3
- **Version**: v3.0.1
- **License**: MIT

---

## License Compatibility

All dependencies are compatible with the Apache License 2.0 used by this project:

- **MIT License**: Compatible (can be included in Apache 2.0 projects)
- **BSD-3-Clause**: Compatible (can be included in Apache 2.0 projects)
- **Apache-2.0**: Compatible (same license)
- **ISC**: Compatible (similar to MIT/BSD)

## Updating This File

When adding new dependencies:
1. Add the dependency information to the appropriate section
2. Verify the license compatibility
3. Include the version and purpose
4. Update the go.mod file accordingly

Last updated: 2025-01-31