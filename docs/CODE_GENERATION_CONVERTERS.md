# Code Generation Guide - Type Converters

This document describes the automated type converter generation system.

## Overview

The slurm-client automatically generates type converters that translate between API types and common types for all supported SLURM REST API versions.

## Quick Start

Generate converters for all API versions:

```bash
make generate-converters
```

Verify converters are up to date:

```bash
make verify-converters
```

## Architecture

```
tools/codegen/
├── generate_converters_v2.go      # Main generator (600+ LOC)
├── converter_patterns.go          # Pattern detection & templates  
├── converter_helpers.go           # Shared helper functions
└── converter_config_enhanced.yaml # Configuration & field mappings
```

## Generated Output

**Total**: 41 files, ~4,400 lines across 5 API versions

```
internal/adapters/v0_0_*/
├── converter_helpers.gen.go    # Shared helpers
├── account_converters.gen.go   # Account converters
├── job_converters.gen.go       # Job converters
├── node_converters.gen.go      # Node converters
└── ... (9 entity types total)
```

## Conversion Patterns

The generator auto-detects 11 patterns:

- **SimpleCopy** - Direct assignment for non-pointer fields
- **SimplePointerCopy** - Pointer fields with nil check
- **PointerDereference** - *T → T conversion
- **SliceCast** - Element-wise slice conversion
- **NoValStructUnwrap** - Extract from SLURM's Set/Number pattern
- **TimeConversion** - Unix timestamp → time.Time
- **Custom** - Complex logic via helper functions

## CI Integration

GitHub Actions automatically verifies generated files are up to date:

```yaml
# .github/workflows/codegen-check.yml
- Triggers on PR or push to main
- Regenerates converters
- Fails if diffs detected
```

Fix CI failures:
```bash
make generate-converters
git add internal/adapters/*/*.gen.go
git commit -m "chore: update generated converters"
```

## See Also

- Full documentation: `docs/CONVERTER_GENERATOR_STATUS.md`
- Architecture details: `docs/PHASE_5A_ANALYSIS.md`
- Contributing guide: `CONTRIBUTING.md`
