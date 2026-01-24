# Permanently Disabled Linters

This document explains why certain golangci-lint linters are permanently disabled for the slurm-client project.

## Disabled Linters and Rationale

### Style and Formatting (Opinion-based)

#### `embeddedstructfieldcheck`
**Status**: Disabled
**Reason**: Embedded struct formatting is a stylistic choice. The linter enforces specific formatting for embedded fields, but our style preferences differ from its defaults.

#### `funcorder`
**Status**: Disabled
**Reason**: Function ordering in files is a developer's choice. Go doesn't enforce any particular order, and our team prefers to organize functions by logical grouping rather than export status or alphabetical order.

#### `nlreturn`
**Status**: Disabled
**Reason**: Enforces blank lines before return statements. This is stylistic and can reduce code density without improving readability.

#### `testpackage`
**Status**: Disabled
**Reason**: Enforces that tests are in separate `_test` packages (e.g., `package foo_test` instead of `package foo`). We use internal test packages for accessing unexported functions when necessary.

#### `wsl` / `wsl_v5`
**Status**: Disabled
**Reason**: Whitespace linting is extremely opinionated. These linters enforce specific blank line rules that conflict with our readability preferences. `wsl` is deprecated, and `wsl_v5` is its modern replacement - both are too strict for this project.

#### `godot`
**Status**: Disabled
**Reason**: Enforces punctuation in doc comments. We prefer flexibility in documentation formatting.

#### `godoclint`
**Status**: Disabled
**Reason**: Doc comment style enforcement (available in golangci-lint v2.8+). Similar to godot, we prefer flexible documentation.

### Complexity Metrics (Too Strict for Domain)

#### `cyclop`
**Status**: Disabled
**Reason**: Cyclomatic complexity linter. Our domain (job management, resource scheduling) involves inherent complexity that would require artificial refactoring to satisfy this linter.

#### `gocognit`
**Status**: Disabled
**Reason**: Cognitive complexity metric. Similar to cyclop, API layer business logic requires complex conditional logic for proper validation and state management.

#### `funlen`
**Status**: Disabled
**Reason**: Enforces maximum function length. Some functions legitimately need to be longer when breaking them down would reduce clarity (e.g., initialization sequences, comprehensive switch statements).

#### `maintidx`
**Status**: Disabled
**Reason**: Maintainability index scoring. Similar to complexity metrics, this is too strict for domain-specific business logic.

### Type and Interface Design (Intentional)

#### `ireturn`
**Status**: Disabled
**Reason**: Warns about returning interfaces. Our API intentionally returns interfaces to allow flexible implementations and mocking in tests.

#### `exhaustruct`
**Status**: Disabled
**Reason**: Requires all struct fields to be explicitly initialized. We use named field initialization where relevant but allow zero-value fields for optional configuration.

#### `interfacebloat`
**Status**: Disabled
**Reason**: Warns about interfaces with too many methods. Our manager interfaces are intentionally comprehensive to provide a complete API for resources.

#### `nonamedreturns`
**Status**: Disabled
**Reason**: Prohibits named return values. Named returns are useful for documenting return values and can improve code clarity in appropriate contexts.

### Error Handling Philosophy (Idiomatic)

#### `noinlineerr`
**Status**: Disabled
**Reason**: Prohibits inline error creation like `fmt.Errorf()`. Inline errors are idiomatic Go and often clearer than storing errors in variables first.

#### `wrapcheck`
**Status**: Disabled
**Reason**: Enforces error wrapping with `%w`. While error wrapping is good practice, this linter is too strict and incorrectly flags legitimate error handling patterns.

#### `err113`
**Status**: Disabled
**Reason**: Prohibits `fmt.Errorf` with constant error strings (suggests creating error variables). While constants can be useful, inline error messages are often more readable and contextual.

### Project Scope (Not Applicable)

#### `gosmopolitan`
**Status**: Disabled
**Reason**: Internationalization linter for i18n. This project is not designed for multiple language/locale support.

#### `godox`
**Status**: Disabled
**Reason**: Flags TODO, FIXME, and NOTE comments. We use these comments for development tracking; warnings would be noise.

#### `forbidigo`
**Status**: Disabled
**Reason**: Forbids specific patterns in code. Not configured for this project and generally too restrictive without specific use cases.

### Deprecated or Problematic (Technical Reasons)

#### `copyloopvar`
**Status**: Disabled
**Reason**: Flags loop variable copies. Fixed in Go 1.22+ with improved loop variable scoping. No longer needed for modern Go versions.

#### `dupl`
**Status**: Disabled
**Reason**: Detects duplicate code blocks. While useful for refactoring, it's too strict for adapter implementations where boilerplate patterns and interface implementations are intentionally similar across versions. The 95+ issues are mostly in legitimate adapter patterns for handling multiple SLURM API versions.

#### `canonicalheader`
**Status**: Disabled
**Reason**: HTTP header canonicalization enforcement. Not applicable to our SLURM client architecture.

#### `containedctx`
**Status**: Disabled
**Reason**: Flags embedded context.Context fields. While generally discouraged, there are legitimate use cases in our middleware and configuration structures.

#### `gochecknoinits` / `gochecknoglobals`
**Status**: Disabled
**Reason**: Discourages init functions and global variables. We use both in limited, justified cases (e.g., test helpers, package initialization).

#### `paralleltest`
**Status**: Disabled
**Reason**: Suggests running subtests in parallel. Not all tests can run in parallel (e.g., database tests, integration tests with shared resources).

#### `tagliatelle`
**Status**: Disabled
**Reason**: Enforces consistent struct tag naming conventions. Our structs follow SLURM API conventions which may differ from Go defaults.

#### `tparallel`
**Status**: Disabled
**Reason**: Flags tests that call `t.Parallel()` incorrectly. Over-aggressive for our test suite structure.

#### `varnamelen`
**Status**: Disabled
**Reason**: Enforces minimum variable name length. Short, meaningful names like `i`, `v`, `id` are idiomatic in Go.

#### `mnd`
**Status**: Disabled
**Reason**: Magic number detection. Constants like `1024` (for memory calculations) are clear without named constants in every context.

#### `lll`
**Status**: Disabled
**Reason**: Line length enforcement. Long lines are acceptable when they improve readability over breaking unnecessarily.

#### `modernize`
**Status**: Disabled
**Reason**: Modernization suggestions (v2.8+ only). Not all suggestions are appropriate; we evaluate improvements case-by-case.

## Linters Intentionally Enabled

These linters were recently enabled as part of major enhancements:
- ✅ `contextcheck` - Ensures context is properly propagated
- ✅ `usetesting` - Uses idiomatic test APIs (`t.Setenv` instead of `os.Setenv`)
- ✅ `dupword` - Catches duplicate words in comments
- ✅ `exhaustive` - Ensures switch statements handle all cases
- ✅ `noctx` - Ensures context is used in network calls
- ✅ `forcetypeassert` - Enforces checked type assertions
- ✅ `unparam` - Detects unused parameters and always-nil returns
- ✅ `intrange` - Suggests Go 1.22+ range loop syntax
- ✅ `revive` (Phase 1) - Comprehensive Go style checking

## Future Enablement Candidates

The following disabled linters may be candidates for future enablement:
- Future phases of `revive` - Additional style rules as code is refactored

## Adding New Linters

When considering enabling a new linter:
1. Check this document for rationale on why it's disabled
2. Run linting with the candidate linter to see issues
3. Assess whether issues are legitimate problems or false positives
4. Consider the signal-to-noise ratio
5. Document the decision in this file

## Configuration Reference

All disabled linters are configured in `.golangci.yml`:
```yaml
linters:
  default: all
  disable:
    # See above for rationale on each disabled linter
```
