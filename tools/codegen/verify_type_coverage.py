#!/usr/bin/env python3
"""
Verify generated .gen.go types cover OpenAPI spec fields.

Usage: python3 verify_type_coverage.py [--spec SPEC] [--types-dir DIR]

Examples:
    python3 verify_type_coverage.py
    python3 verify_type_coverage.py --spec openapi-specs/slurm-v0.0.44.json
    python3 verify_type_coverage.py --types-dir internal/common/types
"""

import argparse
import json
import re
import sys
from pathlib import Path
from typing import Set, List, Tuple


def load_openapi_spec(spec_file: Path) -> dict:
    """Load OpenAPI specification from JSON file."""
    with open(spec_file) as f:
        return json.load(f)


def extract_json_tags(go_file: Path) -> Set[str]:
    """Extract JSON field names from Go struct."""
    content = go_file.read_text()
    pattern = r'`json:"([^",]+)'
    return set(re.findall(pattern, content))


def get_schema_fields(schema: dict) -> Set[str]:
    """Extract field names from OpenAPI schema."""
    return set(schema.get('properties', {}).keys())


def detect_version(spec: dict) -> str:
    """Detect API version from schema names."""
    schemas = spec.get('components', {}).get('schemas', {})
    for name in schemas.keys():
        match = re.match(r'v(\d+\.\d+\.\d+)_', name)
        if match:
            return match.group(1)
    return '0.0.44'


def get_type_mappings(version: str) -> dict:
    """Get mapping from Go file -> OpenAPI schema name."""
    prefix = f'v{version}_'
    return {
        'job.gen.go': f'{prefix}job_info',
        'node.gen.go': f'{prefix}node',
        'account.gen.go': f'{prefix}account',
        'user.gen.go': f'{prefix}user',
        'partition.gen.go': f'{prefix}partition_info',
        'qos.gen.go': f'{prefix}qos',
        'reservation.gen.go': f'{prefix}reservation_info',
        'association.gen.go': f'{prefix}assoc',
        'cluster.gen.go': f'{prefix}cluster_rec',
        'stepid.gen.go': f'{prefix}slurm_step_id',
        'exitcode.gen.go': f'{prefix}process_exit_code_verbose',
        'nodeenergy.gen.go': f'{prefix}acct_gather_energy',
        'jobresources.gen.go': f'{prefix}job_res',
        'assocshort.gen.go': f'{prefix}assoc_short',
        'coord.gen.go': f'{prefix}coord',
        'wckey.gen.go': f'{prefix}wckey',
        'tres.gen.go': f'{prefix}tres',
        'accounting.gen.go': f'{prefix}accounting',
        'jobpartitionpriority.gen.go': f'{prefix}part_prio',
        'reservationcorespec.gen.go': f'{prefix}reservation_core_spec',
        'jobresnode.gen.go': f'{prefix}job_res_node',
        'jobressocket.gen.go': f'{prefix}job_res_socket',
        'jobrescore.gen.go': f'{prefix}job_res_core',
    }


def verify_coverage(
    spec_file: Path,
    types_dir: Path,
    verbose: bool = False
) -> Tuple[int, int, List[str]]:
    """
    Verify type coverage against OpenAPI spec.

    Returns:
        (total_fields, covered_fields, issues)
    """
    spec = load_openapi_spec(spec_file)
    schemas = spec.get('components', {}).get('schemas', {})
    version = detect_version(spec)
    mappings = get_type_mappings(version)

    total_spec_fields = 0
    total_matched = 0
    results = []
    issues = []

    for go_file, schema_name in mappings.items():
        go_path = types_dir / go_file
        if not go_path.exists():
            if verbose:
                print(f"  {go_file}: NOT FOUND")
            continue

        if schema_name not in schemas:
            if verbose:
                print(f"  {go_file}: Schema {schema_name} not in spec")
            continue

        go_tags = extract_json_tags(go_path)
        spec_fields = get_schema_fields(schemas[schema_name])

        matched = go_tags & spec_fields
        missing = spec_fields - go_tags

        total_spec_fields += len(spec_fields)
        total_matched += len(matched)

        pct = (len(matched) / len(spec_fields) * 100) if spec_fields else 100

        results.append((go_file, schema_name, len(matched), len(spec_fields), pct, missing))

        if missing:
            issues.append(f"{go_file}: missing {sorted(missing)}")

    # Sort by coverage
    results.sort(key=lambda x: x[4], reverse=True)

    for go_file, schema_name, matched, total, pct, missing in results:
        status = "✅" if pct == 100 else ("⚠️ " if pct >= 90 else "❌")
        print(f"{status} {go_file}: {matched}/{total} fields ({pct:.1f}%)")
        if missing and verbose:
            print(f"   Missing: {sorted(missing)[:5]}{'...' if len(missing) > 5 else ''}")

    return total_spec_fields, total_matched, issues


def main():
    # Find repo root (where this script lives is tools/codegen/)
    script_dir = Path(__file__).resolve().parent
    repo_root = script_dir.parent.parent

    parser = argparse.ArgumentParser(
        description='Verify generated Go types cover OpenAPI spec fields',
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        '--spec', '-s',
        type=Path,
        default=repo_root / 'openapi-specs/slurm-v0.0.44.json',
        help='Path to OpenAPI spec JSON file'
    )
    parser.add_argument(
        '--types-dir', '-t',
        type=Path,
        default=repo_root / 'internal/common/types',
        help='Directory containing generated Go types'
    )
    parser.add_argument('--verbose', '-v', action='store_true', help='Show detailed output')
    parser.add_argument('--strict', action='store_true', help='Exit with error if any fields missing')

    args = parser.parse_args()

    if not args.spec.exists():
        print(f"Error: Spec file not found: {args.spec}", file=sys.stderr)
        sys.exit(1)

    if not args.types_dir.exists():
        print(f"Error: Types directory not found: {args.types_dir}", file=sys.stderr)
        sys.exit(1)

    print("=" * 90)
    print("Generated Types Coverage vs OpenAPI Spec")
    print("=" * 90)
    print(f"Spec: {args.spec}")
    print(f"Types: {args.types_dir}")
    print()

    total, matched, issues = verify_coverage(args.spec, args.types_dir, args.verbose)

    print()
    print("=" * 90)
    overall_pct = (matched / total * 100) if total else 100
    print(f"Overall: {matched}/{total} fields ({overall_pct:.1f}%)")
    print("=" * 90)

    if issues and args.strict:
        print(f"\n{len(issues)} issues found", file=sys.stderr)
        sys.exit(1)

    sys.exit(0 if overall_pct >= 95 else 1)


if __name__ == '__main__':
    main()
