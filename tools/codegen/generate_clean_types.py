#!/usr/bin/env python3
"""
Generate clean Go types from OpenAPI spec that match OpenAPI structure
but with simplified, unwrapped types.

Usage: python3 generate_clean_types.py [--version VERSION] openapi-spec.json output-dir/

Examples:
    python3 generate_clean_types.py openapi-specs/slurm-v0.0.44.json ./types/
    python3 generate_clean_types.py --version 0.0.45 openapi-specs/slurm-v0.0.45.json ./types/
    python3 generate_clean_types.py --config custom_config.yaml openapi-specs/slurm-v0.0.44.json ./types/
    python3 generate_clean_types.py --validate-only openapi-specs/slurm-v0.0.44.json
    python3 generate_clean_types.py --format minimal openapi-specs/slurm-v0.0.44.json ./types/
"""

import argparse
import json
import re
import sys
from datetime import datetime
from enum import Enum
from pathlib import Path
from typing import Dict, List, Set, Optional, Tuple
from dataclasses import dataclass, field

# Try to import yaml, fall back to built-in config if not available
try:
    import yaml
    YAML_AVAILABLE = True
except ImportError:
    YAML_AVAILABLE = False


# =============================================================================
# Constants
# =============================================================================

VERSION_PREFIX_TEMPLATE = "v{version}_"
ENUM_TYPE_SUFFIX = "Value"
DEFAULT_LINE_LENGTH = 80
DEFAULT_DESCRIPTION_CUTOFF = 40

# Go naming convention acronyms (lowercase -> uppercase)
GO_ACRONYMS = {
    'id': 'ID',
    'ids': 'IDs',
    'url': 'URL',
    'urls': 'URLs',
    'uri': 'URI',
    'uris': 'URIs',
    'ip': 'IP',
    'ips': 'IPs',
    'os': 'OS',
    'cpu': 'CPU',
    'cpus': 'CPUs',
    'gpu': 'GPU',
    'gpus': 'GPUs',
    'api': 'API',
    'http': 'HTTP',
    'https': 'HTTPS',
    'ssh': 'SSH',
    'tls': 'TLS',
    'ssl': 'SSL',
    'tcp': 'TCP',
    'udp': 'UDP',
    'dns': 'DNS',
    'io': 'IO',
    'uid': 'UID',
    'gid': 'GID',
    'pid': 'PID',
    'ram': 'RAM',
    'sql': 'SQL',
    'json': 'JSON',
    'xml': 'XML',
    'html': 'HTML',
    'css': 'CSS',
    'uuid': 'UUID',
    'ascii': 'ASCII',
    'utf': 'UTF',
    # SLURM-specific
    'qos': 'QoS',
    'tres': 'TRES',
    'mcs': 'MCS',
    'oci': 'OCI',
    'gres': 'GRES',
}

# Default configuration values (single source of truth)
DEFAULT_TIMESTAMP_FIELDS = {
    'boot_time', 'last_busy', 'eligible_time', 'end_time', 'start_time',
    'submit_time', 'deadline', 'preempt_time', 'suspend_time', 'resume_time',
    'reason_changed_at', 'slurmd_start_time', 'tls_cert_last_renewal',
    'accrue_time', 'resize_time', 'last_sched_evaluation',
    'created_time', 'modified_time', 'creation_time', 'preemptable_time',
}

DEFAULT_DURATION_FIELDS = {
    'time', 'time_limit', 'time_minimum', 'grace_time', 'pre_sus_time',
    'maximum_switch_wait_time', 'total_time', 'average_time',
    'gettimeofday_latency', 'resume_timeout', 'suspend_timeout',
}

DEFAULT_FRIENDLY_OVERRIDES = {
    'assoc': 'Association',
    'cluster_rec': 'Cluster',
    'job_info': 'Job',
    'partition_info': 'Partition',
    'reservation_info': 'Reservation',
    'acct_gather_energy': 'NodeEnergy',
    'process_exit_code_verbose': 'ExitCode',
    'slurm_step_id': 'StepID',
    'part_prio': 'JobPartitionPriority',
    'job_res': 'JobResources',
    'wckey': 'WCKey',
    'assoc_short': 'AssocShort',
    'tres': 'TRES',
    'job_res_node': 'JobResNode',
    'job_res_socket': 'JobResSocket',
    'job_res_core': 'JobResCore',
}

DEFAULT_PRIMITIVE_UNWRAP_PATTERNS = [
    ('uint32_no_val', 'uint32'),
    ('uint64_no_val', 'uint64'),
    ('uint16_no_val', 'uint16'),
    ('int32_no_val', 'int32'),
    ('float64_no_val', 'float64'),
    ('csv_string', '[]string'),
    ('string_array', '[]string'),
    ('string_list', '[]string'),
    ('hostlist_string', '[]string'),
    ('hostlist', '[]string'),
]

DEFAULT_BASE_ENTITIES = {
    'job_info': 'Job',
    'node': 'Node',
    'account': 'Account',
    'user': 'User',
    'partition_info': 'Partition',
    'qos': 'QoS',
    'reservation_info': 'Reservation',
    'assoc': 'Association',
    'cluster_rec': 'Cluster',
}

DEFAULT_AUXILIARY_TYPES = {
    'accounting': 'Accounting',
    'assoc_short': 'AssocShort',
    'coord': 'Coord',
    'tres': 'TRES',
    'job_res': 'JobResources',
    'acct_gather_energy': 'NodeEnergy',
    'process_exit_code_verbose': 'ExitCode',
    'slurm_step_id': 'StepID',
    'reservation_core_spec': 'ReservationCoreSpec',
    'part_prio': 'JobPartitionPriority',
    'wckey': 'WCKey',
    'job_res_node': 'JobResNode',
    'job_res_socket': 'JobResSocket',
    'job_res_core': 'JobResCore',
}

# Semantic enum overrides: (type_name, field_name) -> enum_type_name
# This allows multiple fields to share the same logical enum type
DEFAULT_ENUM_TYPE_OVERRIDES = {
    ('node', 'state'): 'NodeState',
    ('node', 'next_state_after_reboot'): 'NodeState',
    ('job', 'job_state'): 'JobState',
    ('partition', 'state'): 'PartitionState',
}


# =============================================================================
# Output Format
# =============================================================================

class OutputFormat(Enum):
    """Output format options for generated code."""
    FULL = "full"           # Full output with comments and descriptions
    MINIMAL = "minimal"     # No comments or descriptions
    COMPACT = "compact"     # Comments but truncated descriptions


# =============================================================================
# Data Classes
# =============================================================================

@dataclass
class Config:
    """Configuration for type generation, loaded from YAML or defaults."""
    timestamp_fields: Set[str] = field(default_factory=set)
    duration_fields: Set[str] = field(default_factory=set)
    friendly_overrides: Dict[str, str] = field(default_factory=dict)
    primitive_unwrap_patterns: List[Tuple[str, str]] = field(default_factory=list)
    base_entities: Dict[str, str] = field(default_factory=dict)
    type_unwrap: Dict[str, str] = field(default_factory=dict)
    auxiliary_types: Dict[str, str] = field(default_factory=dict)
    enum_type_overrides: Dict[Tuple[str, str], str] = field(default_factory=dict)
    write_entities: Dict[str, str] = field(default_factory=dict)
    write_auxiliary_types: Dict[str, str] = field(default_factory=dict)

    @classmethod
    def load(cls, config_path: Optional[Path], version: str) -> 'Config':
        """Load configuration from YAML file or use defaults."""
        if config_path and config_path.exists() and YAML_AVAILABLE:
            return cls._load_from_yaml(config_path, version)
        return cls._load_defaults(version)

    @classmethod
    def _load_from_yaml(cls, config_path: Path, version: str) -> 'Config':
        """Load configuration from YAML file."""
        with open(config_path) as f:
            data = yaml.safe_load(f)

        defaults = data.get('defaults', {})
        version_config = data.get('versions', {}).get(version, {})

        # Merge defaults with version-specific config
        config = cls()
        config.timestamp_fields = set(defaults.get('timestamp_fields', DEFAULT_TIMESTAMP_FIELDS))
        config.duration_fields = set(defaults.get('duration_fields', DEFAULT_DURATION_FIELDS))
        config.friendly_overrides = {**DEFAULT_FRIENDLY_OVERRIDES, **defaults.get('friendly_overrides', {})}

        # Convert primitive_unwrap_patterns from dict to list of tuples
        patterns_dict = defaults.get('primitive_unwrap_patterns', {})
        if patterns_dict:
            config.primitive_unwrap_patterns = list(patterns_dict.items())
        else:
            config.primitive_unwrap_patterns = DEFAULT_PRIMITIVE_UNWRAP_PATTERNS.copy()

        config.base_entities = {**DEFAULT_BASE_ENTITIES, **defaults.get('base_entities', {})}

        # Type unwrap from version-specific config
        config.type_unwrap = version_config.get('type_unwrap', {})
        config.auxiliary_types = {**DEFAULT_AUXILIARY_TYPES, **version_config.get('auxiliary_types', {})}

        # Enum type overrides
        enum_overrides = defaults.get('enum_type_overrides', {})
        config.enum_type_overrides = dict(DEFAULT_ENUM_TYPE_OVERRIDES)
        for key, value in enum_overrides.items():
            # Key format: "type.field" -> ("type", "field")
            parts = key.split('.')
            if len(parts) == 2:
                config.enum_type_overrides[(parts[0], parts[1])] = value

        # Write entities: merge defaults with version-specific
        default_write_entities = defaults.get('write_entities', {})
        version_write_entities = version_config.get('write_entities', {})
        config.write_entities = {**default_write_entities, **version_write_entities}

        # Write auxiliary types: merge defaults with version-specific
        default_write_aux = defaults.get('write_auxiliary_types', {})
        version_write_aux = version_config.get('write_auxiliary_types', {})
        config.write_auxiliary_types = {**default_write_aux, **version_write_aux}

        return config

    @classmethod
    def _load_defaults(cls, version: str) -> 'Config':
        """Load default configuration (fallback when YAML not available)."""
        config = cls()

        config.timestamp_fields = set(DEFAULT_TIMESTAMP_FIELDS)
        config.duration_fields = set(DEFAULT_DURATION_FIELDS)
        config.friendly_overrides = dict(DEFAULT_FRIENDLY_OVERRIDES)
        config.primitive_unwrap_patterns = list(DEFAULT_PRIMITIVE_UNWRAP_PATTERNS)
        config.base_entities = dict(DEFAULT_BASE_ENTITIES)
        config.auxiliary_types = dict(DEFAULT_AUXILIARY_TYPES)
        config.enum_type_overrides = dict(DEFAULT_ENUM_TYPE_OVERRIDES)
        config.write_entities = {'job_desc_msg': 'JobCreate'}  # Default write entities
        config.write_auxiliary_types = {}

        # Build type_unwrap with version prefix
        prefix = VERSION_PREFIX_TEMPLATE.format(version=version)
        config.type_unwrap = {
            f'{prefix}uint32_no_val_struct': 'uint32',
            f'{prefix}uint64_no_val_struct': 'uint64',
            f'{prefix}uint16_no_val_struct': 'uint16',
            f'{prefix}int32_no_val_struct': 'int32',
            f'{prefix}float64_no_val_struct': 'float64',
            f'{prefix}csv_string': '[]string',
            f'{prefix}string_array': '[]string',
            f'{prefix}accounting_list': '[]Accounting',
            f'{prefix}assoc_short_list': '[]AssocShort',
            f'{prefix}coord_list': '[]Coord',
            f'{prefix}qos_string_id_list': '[]string',
            f'{prefix}tres_list': '[]TRES',
            f'{prefix}job_info_gres_detail': '[]string',
            f'{prefix}priority_by_partition': '[]JobPartitionPriority',
            f'{prefix}qos_preempt_list': '[]string',
            f'{prefix}job_res_nodes': '[]JobResNode',
            f'{prefix}reservation_info_core_spec': '[]ReservationCoreSpec',
            f'{prefix}wckey_list': '[]WCKey',
            f'{prefix}job_res_socket_array': '[]JobResSocket',
            f'{prefix}job_res_core_array': '[]JobResCore',
        }

        return config


@dataclass
class GoField:
    """Represents a Go struct field."""
    name: str
    go_type: str
    json_tag: str
    description: str
    required: bool

    def to_string(self, output_format: OutputFormat = OutputFormat.FULL) -> str:
        """Convert to Go field declaration string."""
        # Don't add pointer for arrays, maps, required fields, or time.Time
        needs_pointer = (
            not self.required
            and not self.go_type.startswith("[]")
            and not self.go_type.startswith("map[")
            and self.go_type != "time.Time"
        )
        ptr = "*" if needs_pointer else ""
        omit = "" if self.required else ",omitempty"

        # Handle comments based on format
        if output_format == OutputFormat.MINIMAL:
            comment = ""
        elif output_format == OutputFormat.COMPACT:
            # Truncate at 40 chars
            desc = self.description[:40] + "..." if len(self.description) > 40 else self.description
            comment = f" // {desc}" if desc else ""
        else:
            comment = f" // {self.description}" if self.description else ""

        return f'\t{self.name} {ptr}{self.go_type} `json:"{self.json_tag}{omit}"`{comment}'

    def __str__(self):
        return self.to_string(OutputFormat.FULL)


@dataclass
class GeneratedType:
    """Result of generating a Go type."""
    filename: str
    code: str
    additional_types: List[str]
    field_count: int
    type_name: str
    schema_name: str


# =============================================================================
# Utility Functions
# =============================================================================

def get_version_prefix(version: str) -> str:
    """Get the version prefix for schema names."""
    return VERSION_PREFIX_TEMPLATE.format(version=version)


def to_go_field_name(name: str) -> str:
    """Convert JSON field name to Go field name (snake_case -> PascalCase)."""
    parts = name.split('_')
    result = []
    for part in parts:
        lower = part.lower()
        if lower in GO_ACRONYMS:
            result.append(GO_ACRONYMS[lower])
        else:
            result.append(part.capitalize())
    return ''.join(result)


def schema_to_friendly_name(schema_name: str, prefix: str,
                           config: Optional['Config'] = None) -> str:
    """Convert schema name to friendly Go type name."""
    # Remove version prefix
    clean = schema_name.replace(prefix, '')

    # Get friendly name overrides
    friendly_overrides = config.friendly_overrides if config else DEFAULT_FRIENDLY_OVERRIDES

    if clean in friendly_overrides:
        return friendly_overrides[clean]

    # Convert snake_case to PascalCase using GO_ACRONYMS
    parts = clean.split('_')
    result = []
    for part in parts:
        lower = part.lower()
        if lower in GO_ACRONYMS:
            result.append(GO_ACRONYMS[lower])
        elif lower == 'info':
            continue  # Skip 'info' suffix
        else:
            result.append(part.capitalize())

    return ''.join(result) if result else 'Unknown'


def find_all_refs(obj, found=None) -> Set[str]:
    """Recursively find all $ref values in a schema."""
    if found is None:
        found = set()
    if isinstance(obj, dict):
        if '$ref' in obj:
            ref = obj['$ref'].split('/')[-1]
            found.add(ref)
        for v in obj.values():
            find_all_refs(v, found)
    elif isinstance(obj, list):
        for item in obj:
            find_all_refs(item, found)
    return found


def truncate_description(description: str, max_length: int = DEFAULT_LINE_LENGTH) -> str:
    """Truncate description to max length, breaking at word boundary."""
    if len(description) <= max_length:
        return description

    # Find last space before max_length
    cutoff = description.rfind(' ', 0, max_length)
    if cutoff > DEFAULT_DESCRIPTION_CUTOFF:
        return description[:cutoff] + "..."
    return description[:max_length] + "..."


# =============================================================================
# OpenAPI Validation
# =============================================================================

class OpenAPIValidationError(Exception):
    """Raised when OpenAPI spec validation fails."""
    pass


def validate_openapi_spec(spec: dict, spec_file: Path) -> List[str]:
    """
    Validate that the OpenAPI spec has the required structure.

    Returns:
        List of warning messages (empty if valid)

    Raises:
        OpenAPIValidationError: If validation fails with errors
    """
    errors = []
    warnings = []

    # Check top-level structure
    if not isinstance(spec, dict):
        raise OpenAPIValidationError(f"{spec_file}: Root must be a JSON object")

    # Check for components section
    if 'components' not in spec:
        errors.append("Missing 'components' section")
    elif not isinstance(spec['components'], dict):
        errors.append("'components' must be an object")
    elif 'schemas' not in spec['components']:
        errors.append("Missing 'components.schemas' section")
    elif not isinstance(spec['components']['schemas'], dict):
        errors.append("'components.schemas' must be an object")
    elif len(spec['components']['schemas']) == 0:
        errors.append("'components.schemas' is empty")

    # Check for info section (optional but useful)
    if 'info' not in spec:
        warnings.append("Missing 'info' section (version detection may fail)")

    # Check that we have SLURM-style schemas
    if not errors:
        schemas = spec['components']['schemas']
        slurm_schemas = [k for k in schemas.keys() if k.startswith('v0.0.')]
        if not slurm_schemas:
            errors.append("No SLURM versioned schemas found (expected v0.0.XX_ prefix)")

    if errors:
        error_msg = f"{spec_file}: OpenAPI validation failed:\n" + "\n".join(f"  - {e}" for e in errors)
        raise OpenAPIValidationError(error_msg)

    return warnings


def load_openapi_spec(spec_file: Path) -> dict:
    """Load OpenAPI specification from JSON file."""
    with open(spec_file) as f:
        return json.load(f)


# =============================================================================
# Schema Discovery
# =============================================================================

def get_type_unwrap_map(version: str, config: Optional[Config] = None) -> Dict[str, str]:
    """
    Returns version-specific type unwrapping map.
    OpenAPI wrapper types -> Clean Go types
    """
    if config and config.type_unwrap:
        # Use config, but add version prefix if not present
        prefix = get_version_prefix(version)
        result = {}
        for key, value in config.type_unwrap.items():
            if not key.startswith('v'):
                key = prefix + key
            result[key] = value
        return result

    # Fallback to defaults
    return Config._load_defaults(version).type_unwrap


def get_entity_schemas(version: str, config: Optional[Config] = None) -> Dict[str, str]:
    """
    Returns version-specific entity schemas to generate.
    Maps OpenAPI schema name -> friendly Go type name.
    """
    prefix = get_version_prefix(version)

    if config:
        result = {}
        # Add base entities
        for key, value in config.base_entities.items():
            result[prefix + key] = value
        # Add auxiliary types
        for key, value in config.auxiliary_types.items():
            result[prefix + key] = value
        return result

    # Fallback to defaults
    default_config = Config._load_defaults(version)
    result = {}
    for key, value in default_config.base_entities.items():
        result[prefix + key] = value
    for key, value in default_config.auxiliary_types.items():
        result[prefix + key] = value
    return result


def get_base_entities(version: str, config: Optional[Config] = None) -> Dict[str, str]:
    """
    Returns just the base/main entities to generate.
    Used with --discover to auto-discover auxiliary types.
    """
    prefix = get_version_prefix(version)

    if config and config.base_entities:
        return {prefix + k: v for k, v in config.base_entities.items()}

    # Fallback to defaults
    return {prefix + k: v for k, v in DEFAULT_BASE_ENTITIES.items()}


def get_write_entity_schemas(version: str, config: Optional[Config] = None) -> Dict[str, str]:
    """
    Returns write entity schemas to generate (e.g., JobCreate from job_desc_msg).
    Maps OpenAPI schema name -> friendly Go type name.
    """
    prefix = get_version_prefix(version)

    if config and config.write_entities:
        result = {}
        # Add write entities
        for key, value in config.write_entities.items():
            result[prefix + key] = value
        # Add write auxiliary types
        for key, value in config.write_auxiliary_types.items():
            result[prefix + key] = value
        return result

    # Fallback to default write entities
    return {prefix + 'job_desc_msg': 'JobCreate'}


def discover_schemas(spec: dict, version: str, base_schemas: Dict[str, str],
                     config: Optional[Config] = None) -> Tuple[Dict[str, str], Dict[str, str]]:
    """
    Auto-discover all schemas needed by base_schemas.

    Returns:
        (schemas_to_generate, array_unwrap_map)
    """
    prefix = get_version_prefix(version)
    schemas = spec.get('components', {}).get('schemas', {})

    # Get primitive unwrap patterns from config or defaults
    primitive_patterns = config.primitive_unwrap_patterns if config else DEFAULT_PRIMITIVE_UNWRAP_PATTERNS

    # Start with base schemas
    to_generate = dict(base_schemas)
    array_unwrap = {}

    # Track what we've processed
    processed = set()
    queue = list(base_schemas.keys())

    while queue:
        schema_name = queue.pop(0)
        if schema_name in processed:
            continue
        processed.add(schema_name)

        schema = schemas.get(schema_name, {})

        # Find all $ref in this schema
        refs = find_all_refs(schema)

        for ref_name in refs:
            if ref_name in processed:
                continue

            ref_schema = schemas.get(ref_name, {})

            # Check if it's a primitive unwrap pattern
            clean_name = ref_name.replace(prefix, '')
            is_primitive = False
            for pattern, _ in primitive_patterns:
                if pattern in clean_name:
                    is_primitive = True
                    break

            if is_primitive:
                continue  # Already handled by get_type_unwrap_map

            # Skip deprecated types with no properties (convert to interface{})
            if (ref_schema.get('deprecated', False) and
                not ref_schema.get('properties', {})):
                array_unwrap[ref_name] = 'interface{}'
                continue

            # Check if it's an array type
            if ref_schema.get('type') == 'array':
                items = ref_schema.get('items', {})
                if '$ref' in items:
                    item_ref = items['$ref'].split('/')[-1]
                    item_name = schema_to_friendly_name(item_ref, prefix, config)
                    array_unwrap[ref_name] = f'[]{item_name}'
                    # Also need to generate the item type
                    if item_ref not in to_generate and item_ref not in processed:
                        queue.append(item_ref)
                        to_generate[item_ref] = item_name
                elif items.get('type'):
                    # Primitive array
                    item_type = items.get('type', 'interface{}')
                    go_type = {'string': 'string', 'integer': 'int32', 'boolean': 'bool'}.get(item_type, 'interface{}')
                    array_unwrap[ref_name] = f'[]{go_type}'
                else:
                    # Empty items: {} - infer from schema name
                    clean = ref_name.replace(prefix, '')
                    # Try to make a reasonable type name
                    item_name = schema_to_friendly_name(ref_name, prefix, config) + 'Item'
                    array_unwrap[ref_name] = f'[]interface{{}}'  # Safe fallback
            else:
                # It's a struct type - add to generation queue
                if ref_name not in to_generate:
                    friendly = schema_to_friendly_name(ref_name, prefix, config)
                    to_generate[ref_name] = friendly
                    queue.append(ref_name)

    return to_generate, array_unwrap


# =============================================================================
# Type Generator
# =============================================================================

class TypeGenerator:
    """Generates Go types from OpenAPI schemas."""

    # Class-level tracking of globally generated enum types (shared across instances)
    _global_enum_types: Dict[str, List[str]] = {}

    def __init__(self, spec: dict, version: str, friendly_names: Dict[str, str],
                 config: Optional[Config] = None,
                 output_format: OutputFormat = OutputFormat.FULL,
                 known_enum_types: Optional[Set[str]] = None):
        self.spec = spec
        self.version = version
        self.version_prefix = get_version_prefix(version)
        self.config = config or Config._load_defaults(version)
        self.type_unwrap = get_type_unwrap_map(version, config)
        self.friendly_names = friendly_names  # schema_name -> Go type name
        self.generated_type_names: Set[str] = set()
        self.enum_types: Dict[str, List[str]] = {}  # type_name -> enum values (per-file)
        self.current_type_name: str = ''  # Track current type being generated
        self.current_type_name_lower: str = ''  # Lowercase for enum override lookup
        self.output_format = output_format
        # Known enum types that already exist (skip generating them)
        self.known_enum_types = known_enum_types or set()

    def get_friendly_name(self, schema_name: str) -> str:
        """Get friendly Go type name for a schema, using explicit mapping or auto-generation."""
        # Check if we have an explicit mapping
        if schema_name in self.friendly_names:
            return self.friendly_names[schema_name]
        # Use module-level function for auto-generation
        return schema_to_friendly_name(schema_name, self.version_prefix, self.config)

    def is_timestamp_field(self, field_name: str) -> bool:
        """Check if a field represents a Unix timestamp (not a duration)."""
        return (field_name in self.config.timestamp_fields and
                field_name not in self.config.duration_fields)

    def _get_enum_type_name(self, json_name: str, enum_values: List[str]) -> str:
        """
        Get a unique enum type name and register it for generation.
        Reuses existing type if same values already registered.
        Skips registration if type is in known_enum_types.
        """
        # Check for semantic enum override
        override_key = (self.current_type_name_lower, json_name)
        if override_key in self.config.enum_type_overrides:
            semantic_name = self.config.enum_type_overrides[override_key]
            # Skip if this is a known type (already exists in other generated files)
            if semantic_name in self.known_enum_types:
                return semantic_name
            # Check if already registered with same values
            if semantic_name in TypeGenerator._global_enum_types:
                existing_values = TypeGenerator._global_enum_types[semantic_name]
                if set(existing_values) == set(enum_values):
                    return semantic_name
            # Register new semantic enum
            TypeGenerator._global_enum_types[semantic_name] = enum_values
            self.enum_types[semantic_name] = enum_values
            return semantic_name

        # Create a base type name from field name
        base_name = to_go_field_name(json_name)
        enum_type_name = base_name + ENUM_TYPE_SUFFIX

        # Skip if this is a known type (already exists in other generated files)
        if enum_type_name in self.known_enum_types:
            return enum_type_name

        # Create a key for the enum values (frozen set for comparison)
        values_key = tuple(sorted(enum_values))

        # Check if this exact enum type already exists globally
        for existing_name, existing_values in TypeGenerator._global_enum_types.items():
            if tuple(sorted(existing_values)) == values_key:
                # Reuse existing type
                return existing_name

        # Check for name collision - make unique if needed
        if enum_type_name in TypeGenerator._global_enum_types:
            # Name exists with different values - prefix with parent type
            enum_type_name = self.current_type_name + base_name + ENUM_TYPE_SUFFIX

        # Register globally and locally
        TypeGenerator._global_enum_types[enum_type_name] = enum_values
        self.enum_types[enum_type_name] = enum_values
        return enum_type_name

    def unwrap_type(self, ref_name: str, json_field: str) -> Optional[str]:
        """Unwrap OpenAPI special types to clean Go types."""
        ref_lower = ref_name.lower()

        # Check if it's a timestamp field (but not a duration)
        if self.is_timestamp_field(json_field):
            if 'uint64' in ref_lower or 'no_val' in ref_lower:
                return 'time.Time'

        # Check unwrap map
        if ref_lower in self.type_unwrap:
            return self.type_unwrap[ref_lower]

        return None

    def resolve_type(self, prop: dict, json_name: str) -> str:
        """Resolve property type to clean Go type."""
        # Handle $ref
        if '$ref' in prop:
            ref = prop['$ref']
            ref_name = ref.split('/')[-1]

            # Try to unwrap special types
            unwrapped = self.unwrap_type(ref_name, json_name)
            if unwrapped:
                return unwrapped

            # It's a nested type - use friendly name if available
            return self.get_friendly_name(ref_name)

        # Handle arrays
        if prop.get('type') == 'array':
            if 'items' in prop:
                items = prop['items']
                # Check if items have enum - generate enum type
                if 'enum' in items:
                    enum_type_name = self._get_enum_type_name(json_name, items['enum'])
                    return f'[]{enum_type_name}'
                # Check for $ref in items
                if '$ref' in items:
                    item_type = self.resolve_type(items, json_name)
                    return f'[]{item_type}'
                # Check for type in items
                if items.get('type'):
                    item_type = self.resolve_type(items, json_name)
                    return f'[]{item_type}'
                # Empty items: {} - use interface{}
                return '[]interface{}'
            return '[]interface{}'

        # Handle objects
        if prop.get('type') == 'object':
            if 'properties' in prop and prop['properties']:
                return to_go_field_name(json_name)  # Will be generated as nested type
            # Check for additionalProperties
            if 'additionalProperties' in prop:
                add_props = prop['additionalProperties']
                if isinstance(add_props, dict) and add_props.get('type'):
                    value_type = self.resolve_type(add_props, json_name)
                    return f'map[string]{value_type}'
            return 'map[string]interface{}'

        # Handle enums (non-array)
        if 'enum' in prop:
            return self._get_enum_type_name(json_name, prop['enum'])

        # Handle primitives
        type_map = {
            'string': 'string',
            'integer': 'int32',
            'number': 'float64',
            'boolean': 'bool',
        }

        prop_type = prop.get('type', '')
        go_type = type_map.get(prop_type, 'interface{}')

        # Handle format overrides
        fmt = prop.get('format', '')
        if fmt == 'int64':
            go_type = 'int64'
        elif fmt == 'uint32':
            go_type = 'uint32'
        elif fmt == 'uint16':
            go_type = 'uint16'
        elif fmt == 'uint64':
            go_type = 'time.Time' if self.is_timestamp_field(json_name) else 'uint64'
        elif fmt == 'date-time':
            go_type = 'time.Time'

        return go_type

    def generate_type(
        self,
        schema_name: str,
        schema: dict,
        type_name: str,
        include_header: bool = True
    ) -> GeneratedType:
        """
        Generate Go type definition from OpenAPI schema.

        Returns: GeneratedType with all generation results
        """
        file_name = type_name.lower() + '.gen.go'

        # Track current type for unique enum naming
        self.current_type_name = type_name
        self.current_type_name_lower = type_name.lower()

        # Get required fields
        required_fields = set(schema.get('required', []))

        # Generate fields
        fields = []
        additional_types = []
        props = schema.get('properties', {})

        for json_name in sorted(props.keys()):
            prop = props[json_name]
            field_type = None

            # Check for inline object that needs nested type
            if prop.get('type') == 'object' and prop.get('properties'):
                field_name = to_go_field_name(json_name)

                # Build unique nested type name
                nested_type_name = f"{type_name}{field_name}"

                # Deduplicate if needed
                base_name = nested_type_name
                counter = 2
                while nested_type_name in self.generated_type_names:
                    nested_type_name = f"{base_name}{counter}"
                    counter += 1

                self.generated_type_names.add(nested_type_name)

                nested_schema = {
                    'type': 'object',
                    'properties': prop['properties'],
                    'required': prop.get('required', [])
                }

                # Recursively generate nested type
                nested_result = self.generate_type(
                    f"{schema_name}_{json_name}",
                    nested_schema,
                    nested_type_name,
                    include_header=False
                )
                additional_types.append(nested_result.code)
                additional_types.extend(nested_result.additional_types)

                field_type = nested_type_name

            # Get field type
            if field_type is None:
                field_type = self.resolve_type(prop, json_name)

            # Clean up description
            description = prop.get('description', '').split('\n')[0]
            description = truncate_description(description)

            field = GoField(
                name=to_go_field_name(json_name),
                go_type=field_type,
                json_tag=json_name,
                description=description,
                required=json_name in required_fields
            )
            fields.append(field)

        # Generate Go code
        if include_header:
            needs_time = any('time.Time' in f.go_type for f in fields)
            time_import = 'import "time"\n\n' if needs_time else ''

            if self.output_format == OutputFormat.MINIMAL:
                type_comment = ""
            else:
                type_comment = f"// {type_name} represents a SLURM {type_name}.\n"

            code = f'''// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: {datetime.now().year} Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

{time_import}{type_comment}type {type_name} struct {{
{chr(10).join(f.to_string(self.output_format) for f in fields)}
}}
'''
        else:
            if self.output_format == OutputFormat.MINIMAL:
                nested_comment = ""
            else:
                nested_comment = f"// {type_name} is a nested type within its parent.\n"

            code = f'''{nested_comment}type {type_name} struct {{
{chr(10).join(f.to_string(self.output_format) for f in fields)}
}}
'''

        # Generate enum types collected during field resolution
        enum_code = self._generate_enum_types()
        if enum_code:
            additional_types.append(enum_code)

        return GeneratedType(
            filename=file_name,
            code=code,
            additional_types=additional_types,
            field_count=len(props),
            type_name=type_name,
            schema_name=schema_name
        )

    def _generate_enum_types(self) -> str:
        """Generate Go type aliases and constants for collected enum types."""
        if not self.enum_types:
            return ''

        lines = []
        for enum_name, values in sorted(self.enum_types.items()):
            # Generate type alias
            if self.output_format != OutputFormat.MINIMAL:
                field_hint = enum_name.replace(ENUM_TYPE_SUFFIX, '')
                lines.append(f'// {enum_name} represents possible values for {field_hint} field.')
            lines.append(f'type {enum_name} string')
            lines.append('')

            # Generate constants
            if self.output_format != OutputFormat.MINIMAL:
                lines.append(f'// {enum_name} constants.')
            lines.append('const (')
            for value in values:
                # Convert enum value to Go constant name
                const_name = self._enum_value_to_const(enum_name, value)
                lines.append(f'\t{const_name} {enum_name} = "{value}"')
            lines.append(')')
            lines.append('')

        # Clear enum types after generation (for next file)
        self.enum_types.clear()

        return '\n'.join(lines)

    def _enum_value_to_const(self, enum_name: str, value: str) -> str:
        """Convert an enum value to a Go constant name."""
        # Remove 'Value' suffix from enum_name for the prefix
        prefix = enum_name.replace(ENUM_TYPE_SUFFIX, '')
        # Replace special characters with underscores, then convert to PascalCase
        clean_value = re.sub(r'[^a-zA-Z0-9_]', '_', value)
        parts = clean_value.split('_')
        # Filter out empty parts and use GO_ACRONYMS for proper casing
        result_parts = []
        for p in parts:
            if not p:
                continue
            lower = p.lower()
            if lower in GO_ACRONYMS:
                result_parts.append(GO_ACRONYMS[lower])
            else:
                result_parts.append(p.capitalize())
        pascal_value = ''.join(result_parts)
        return f'{prefix}{pascal_value}'


# =============================================================================
# Main
# =============================================================================

def main():
    # Reset global state
    TypeGenerator._global_enum_types.clear()

    # Determine default config path (relative to script)
    script_dir = Path(__file__).resolve().parent
    default_config = script_dir / 'type_config.yaml'

    parser = argparse.ArgumentParser(
        description='Generate Go types from OpenAPI spec',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog='''
Examples:
  %(prog)s openapi-specs/slurm-v0.0.44.json ./types/
  %(prog)s --version 0.0.45 openapi-specs/slurm-v0.0.45.json ./types/
  %(prog)s --discover openapi-specs/slurm-v0.0.44.json ./types/  # Auto-discover needed types
  %(prog)s --config custom_config.yaml openapi-specs/slurm-v0.0.44.json ./types/
  %(prog)s --validate-only openapi-specs/slurm-v0.0.44.json      # Just validate spec
  %(prog)s --format minimal openapi-specs/slurm-v0.0.44.json ./types/  # No comments
        '''
    )
    parser.add_argument('spec_file', type=Path, help='Path to OpenAPI spec JSON file')
    parser.add_argument('output_dir', type=Path, nargs='?', default=None,
                        help='Output directory for generated files (not required with --validate-only)')
    parser.add_argument('--version', '-v', default='0.0.44',
                        help='API version (default: 0.0.44, auto-detected from spec if possible)')
    parser.add_argument('--config', '-c', type=Path, default=default_config,
                        help=f'Path to config YAML file (default: {default_config})')
    parser.add_argument('--dry-run', action='store_true',
                        help='Print what would be generated without writing files')
    parser.add_argument('--discover', action='store_true',
                        help='Auto-discover auxiliary types from base entities')
    parser.add_argument('--validate-only', action='store_true',
                        help='Only validate the OpenAPI spec, do not generate files')
    parser.add_argument('--format', '-f', choices=['full', 'minimal', 'compact'],
                        default='full',
                        help='Output format: full (default), minimal (no comments), compact (short comments)')
    parser.add_argument('--write-types', action='store_true',
                        help='Generate write types (e.g., JobCreate from job_desc_msg) instead of read types')

    args = parser.parse_args()

    # Validate arguments
    if not args.validate_only and args.output_dir is None:
        parser.error("output_dir is required unless --validate-only is specified")

    # Load and validate spec
    print(f"Loading OpenAPI spec from {args.spec_file}...")
    try:
        spec = load_openapi_spec(args.spec_file)
        warnings = validate_openapi_spec(spec, args.spec_file)
        for warning in warnings:
            print(f"  Warning: {warning}")
        print("  Validation passed")
    except (json.JSONDecodeError, OpenAPIValidationError) as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

    # If validate-only, we're done
    if args.validate_only:
        print("Validation completed successfully.")
        sys.exit(0)

    # Try to auto-detect version from spec
    version = args.version
    if 'info' in spec and 'version' in spec['info']:
        detected = spec['info']['version']
        if detected != version:
            print(f"Note: Detected version {detected} in spec (using --version {version})")

    # Load configuration
    config = Config.load(args.config, version)
    if args.config.exists() and YAML_AVAILABLE:
        print(f"Loaded config from {args.config}")
    elif not YAML_AVAILABLE and args.config.exists():
        print(f"Note: PyYAML not installed, using built-in defaults (pip install pyyaml to use config)")

    # Parse output format
    output_format = OutputFormat(args.format)

    # Get entity schemas for this version
    extra_unwrap: Dict[str, str] = {}  # Additional type unwraps from discovery
    if args.write_types:
        # Generate write types (e.g., JobCreate)
        entity_schemas = get_write_entity_schemas(version, config)
        print(f"Generating {len(entity_schemas)} write types...")
    elif args.discover:
        base_entities = get_base_entities(version, config)
        print(f"Auto-discovering types from {len(base_entities)} base entities...")
        entity_schemas, extra_unwrap = discover_schemas(spec, version, base_entities, config)
        print(f"  Discovered {len(entity_schemas)} total types, {len(extra_unwrap)} array unwraps")
    else:
        entity_schemas = get_entity_schemas(version, config)

    # Create output directory
    if not args.dry_run:
        args.output_dir.mkdir(parents=True, exist_ok=True)

    schemas = spec.get('components', {}).get('schemas', {})

    # Create generator (with any extra unwrap mappings from discovery)
    # When generating write types, provide known enum types to avoid regenerating them
    known_enums: Set[str] = set()
    if args.write_types:
        # These enums are already defined in the read types (job.gen.go, etc.)
        # We only want to generate NEW enums that are specific to write types
        known_enums = {
            'FlagsValue', 'MailTypeValue', 'ProfileValue', 'SharedValue',  # From Job
            'JobState', 'NodeState', 'PartitionState',  # Semantic enums
            # Add more as needed from other generated types
        }
    generator = TypeGenerator(spec, version, entity_schemas, config, output_format, known_enums)
    if args.discover:
        generator.type_unwrap.update(extra_unwrap)

    generated_files = []
    for schema_name, friendly_name in entity_schemas.items():
        if schema_name not in schemas:
            print(f"  Warning: Schema {schema_name} not found in spec")
            continue

        schema = schemas[schema_name]
        result = generator.generate_type(schema_name, schema, friendly_name)

        # Combine main type and nested types
        if result.additional_types:
            full_code = '\n\n'.join([result.code] + result.additional_types)
        else:
            full_code = result.code

        output_file = args.output_dir / result.filename

        if args.dry_run:
            print(f"Would generate: {output_file} ({result.field_count} fields)")
        else:
            with open(output_file, 'w') as f:
                f.write(full_code)

            print(f"  Generated {friendly_name} ({result.field_count} fields) -> {result.filename}")
            generated_files.append(output_file)

    print(f"\n{'Would generate' if args.dry_run else 'Successfully generated'} "
          f"{len(generated_files)} type files in {args.output_dir}")

    if not args.dry_run:
        print("\nNext steps:")
        print("1. Review generated files")
        print("2. Run: go build ./...")
        print("3. Run tests to verify compatibility")


if __name__ == '__main__':
    main()
