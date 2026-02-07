#!/usr/bin/env python3
"""
Unit tests for generate_clean_types.py

Run with: python3 -m pytest test_generate_clean_types.py -v
Or:       python3 test_generate_clean_types.py
"""

import unittest
from pathlib import Path

# Import the module under test
from generate_clean_types import (
    to_go_field_name,
    schema_to_friendly_name,
    find_all_refs,
    truncate_description,
    GoField,
    GeneratedType,
    Config,
    TypeGenerator,
    OutputFormat,
    discover_schemas,
    validate_openapi_spec,
    OpenAPIValidationError,
    GO_ACRONYMS,
    DEFAULT_ENUM_TYPE_OVERRIDES,
)


class TestToGoFieldName(unittest.TestCase):
    """Tests for to_go_field_name function."""

    def test_simple_snake_case(self):
        """Convert simple snake_case to PascalCase."""
        self.assertEqual(to_go_field_name('account'), 'Account')
        self.assertEqual(to_go_field_name('submit_time'), 'SubmitTime')

    def test_id_acronym(self):
        """ID should be uppercase per Go conventions."""
        self.assertEqual(to_go_field_name('job_id'), 'JobID')
        self.assertEqual(to_go_field_name('user_id'), 'UserID')
        self.assertEqual(to_go_field_name('id'), 'ID')
        self.assertEqual(to_go_field_name('array_job_id'), 'ArrayJobID')

    def test_common_acronyms(self):
        """Handle common acronyms correctly."""
        self.assertEqual(to_go_field_name('qos'), 'QoS')
        self.assertEqual(to_go_field_name('tres'), 'TRES')
        self.assertEqual(to_go_field_name('cpu'), 'CPU')
        self.assertEqual(to_go_field_name('cpus'), 'CPUs')
        self.assertEqual(to_go_field_name('gpu'), 'GPU')
        self.assertEqual(to_go_field_name('url'), 'URL')
        self.assertEqual(to_go_field_name('ip'), 'IP')
        self.assertEqual(to_go_field_name('http'), 'HTTP')
        self.assertEqual(to_go_field_name('api'), 'API')

    def test_special_cases(self):
        """Handle special cases like os, uid, gid."""
        self.assertEqual(to_go_field_name('os'), 'OS')
        self.assertEqual(to_go_field_name('uid'), 'UID')
        self.assertEqual(to_go_field_name('gid'), 'GID')
        self.assertEqual(to_go_field_name('pid'), 'PID')

    def test_mixed_case(self):
        """Handle multiple parts with acronyms."""
        self.assertEqual(to_go_field_name('tres_per_job'), 'TRESPerJob')
        self.assertEqual(to_go_field_name('qos_name'), 'QoSName')
        self.assertEqual(to_go_field_name('cpu_frequency'), 'CPUFrequency')


class TestSchemaToFriendlyName(unittest.TestCase):
    """Tests for schema_to_friendly_name function."""

    def test_with_overrides(self):
        """Use friendly overrides when available."""
        prefix = 'v0.0.44_'
        self.assertEqual(schema_to_friendly_name('v0.0.44_assoc', prefix), 'Association')
        self.assertEqual(schema_to_friendly_name('v0.0.44_job_info', prefix), 'Job')
        self.assertEqual(schema_to_friendly_name('v0.0.44_cluster_rec', prefix), 'Cluster')

    def test_auto_generation(self):
        """Auto-generate PascalCase from schema name."""
        prefix = 'v0.0.44_'
        self.assertEqual(schema_to_friendly_name('v0.0.44_node', prefix), 'Node')
        self.assertEqual(schema_to_friendly_name('v0.0.44_account', prefix), 'Account')

    def test_skip_info_suffix(self):
        """Skip 'info' suffix in generated names."""
        prefix = 'v0.0.44_'
        # Note: 'job_info' is overridden to 'Job', but 'some_info' would be 'Some'
        self.assertEqual(schema_to_friendly_name('v0.0.44_some_other_info', prefix), 'SomeOther')

    def test_step_id_override(self):
        """StepID should use proper ID casing."""
        prefix = 'v0.0.44_'
        self.assertEqual(schema_to_friendly_name('v0.0.44_slurm_step_id', prefix), 'StepID')


class TestFindAllRefs(unittest.TestCase):
    """Tests for find_all_refs function."""

    def test_simple_ref(self):
        """Find simple $ref."""
        obj = {'$ref': '#/components/schemas/v0.0.44_node'}
        refs = find_all_refs(obj)
        self.assertEqual(refs, {'v0.0.44_node'})

    def test_nested_refs(self):
        """Find refs in nested objects."""
        obj = {
            'properties': {
                'account': {'$ref': '#/components/schemas/v0.0.44_account'},
                'user': {'$ref': '#/components/schemas/v0.0.44_user'},
            }
        }
        refs = find_all_refs(obj)
        self.assertEqual(refs, {'v0.0.44_account', 'v0.0.44_user'})

    def test_refs_in_arrays(self):
        """Find refs in array items."""
        obj = {
            'type': 'array',
            'items': {'$ref': '#/components/schemas/v0.0.44_job'}
        }
        refs = find_all_refs(obj)
        self.assertEqual(refs, {'v0.0.44_job'})

    def test_no_refs(self):
        """Return empty set when no refs present."""
        obj = {'type': 'string', 'description': 'A simple string'}
        refs = find_all_refs(obj)
        self.assertEqual(refs, set())


class TestTruncateDescription(unittest.TestCase):
    """Tests for truncate_description function."""

    def test_short_description(self):
        """Short descriptions are not truncated."""
        desc = "A short description"
        self.assertEqual(truncate_description(desc), desc)

    def test_long_description(self):
        """Long descriptions are truncated at word boundary."""
        desc = "This is a very long description that exceeds the maximum length and should be truncated at a word boundary to make it readable"
        result = truncate_description(desc, max_length=80)
        self.assertTrue(len(result) <= 83)  # 80 + "..."
        self.assertTrue(result.endswith("..."))


class TestGoField(unittest.TestCase):
    """Tests for GoField dataclass."""

    def test_required_field(self):
        """Required fields don't have pointer or omitempty."""
        field = GoField(
            name='Name',
            go_type='string',
            json_tag='name',
            description='The name',
            required=True
        )
        self.assertIn('Name string `json:"name"`', str(field))
        self.assertNotIn('omitempty', str(field))
        self.assertNotIn('*', str(field))

    def test_optional_field(self):
        """Optional fields have pointer and omitempty."""
        field = GoField(
            name='Description',
            go_type='string',
            json_tag='description',
            description='Optional desc',
            required=False
        )
        result = str(field)
        self.assertIn('*string', result)
        self.assertIn('omitempty', result)

    def test_array_field_no_pointer(self):
        """Array fields never have pointer."""
        field = GoField(
            name='Items',
            go_type='[]string',
            json_tag='items',
            description='List of items',
            required=False
        )
        result = str(field)
        self.assertIn('[]string', result)
        self.assertNotIn('*[]string', result)

    def test_time_field_no_pointer(self):
        """time.Time fields never have pointer."""
        field = GoField(
            name='StartTime',
            go_type='time.Time',
            json_tag='start_time',
            description='Start time',
            required=False
        )
        result = str(field)
        self.assertIn('time.Time', result)
        self.assertNotIn('*time.Time', result)

    def test_output_format_minimal(self):
        """Minimal format removes comments."""
        field = GoField(
            name='Name',
            go_type='string',
            json_tag='name',
            description='The name',
            required=True
        )
        result = field.to_string(OutputFormat.MINIMAL)
        self.assertNotIn('//', result)
        self.assertNotIn('The name', result)

    def test_output_format_compact(self):
        """Compact format truncates long comments."""
        field = GoField(
            name='Name',
            go_type='string',
            json_tag='name',
            description='This is a very long description that should be truncated in compact mode',
            required=True
        )
        result = field.to_string(OutputFormat.COMPACT)
        self.assertIn('...', result)


class TestGeneratedType(unittest.TestCase):
    """Tests for GeneratedType dataclass."""

    def test_fields(self):
        """GeneratedType has all required fields."""
        gt = GeneratedType(
            filename='test.go',
            code='type Test struct {}',
            additional_types=[],
            field_count=5,
            type_name='Test',
            schema_name='v0.0.44_test'
        )
        self.assertEqual(gt.filename, 'test.go')
        self.assertEqual(gt.field_count, 5)
        self.assertEqual(gt.type_name, 'Test')


class TestConfig(unittest.TestCase):
    """Tests for Config class."""

    def test_load_defaults(self):
        """Load default config without YAML file."""
        config = Config._load_defaults('0.0.44')
        self.assertIn('boot_time', config.timestamp_fields)
        self.assertIn('time', config.duration_fields)
        self.assertIn('assoc', config.friendly_overrides)
        self.assertGreater(len(config.primitive_unwrap_patterns), 0)

    def test_is_timestamp_vs_duration(self):
        """Verify timestamp/duration field classification."""
        config = Config._load_defaults('0.0.44')
        # boot_time is a timestamp
        self.assertIn('boot_time', config.timestamp_fields)
        self.assertNotIn('boot_time', config.duration_fields)
        # time is a duration
        self.assertIn('time', config.duration_fields)

    def test_enum_type_overrides(self):
        """Config includes enum type overrides."""
        config = Config._load_defaults('0.0.44')
        self.assertIn(('node', 'state'), config.enum_type_overrides)
        self.assertEqual(config.enum_type_overrides[('node', 'state')], 'NodeState')

    def test_step_id_friendly_override(self):
        """StepID uses proper ID casing in friendly overrides."""
        config = Config._load_defaults('0.0.44')
        self.assertEqual(config.friendly_overrides.get('slurm_step_id'), 'StepID')


class TestTypeGenerator(unittest.TestCase):
    """Tests for TypeGenerator class."""

    def setUp(self):
        """Set up test fixtures."""
        # Reset global enum state
        TypeGenerator._global_enum_types.clear()

        self.spec = {
            'components': {
                'schemas': {
                    'v0.0.44_simple_type': {
                        'type': 'object',
                        'properties': {
                            'name': {'type': 'string', 'description': 'The name'},
                            'count': {'type': 'integer', 'description': 'Count'},
                        },
                        'required': ['name']
                    }
                }
            }
        }
        self.generator = TypeGenerator(
            self.spec,
            '0.0.44',
            {'v0.0.44_simple_type': 'SimpleType'}
        )

    def test_is_timestamp_field(self):
        """Check timestamp field detection."""
        self.assertTrue(self.generator.is_timestamp_field('boot_time'))
        self.assertTrue(self.generator.is_timestamp_field('submit_time'))
        self.assertFalse(self.generator.is_timestamp_field('time'))  # Duration
        self.assertFalse(self.generator.is_timestamp_field('time_limit'))  # Duration
        self.assertFalse(self.generator.is_timestamp_field('random_field'))

    def test_resolve_type_primitives(self):
        """Resolve primitive types correctly."""
        self.assertEqual(self.generator.resolve_type({'type': 'string'}, 'name'), 'string')
        self.assertEqual(self.generator.resolve_type({'type': 'integer'}, 'count'), 'int32')
        self.assertEqual(self.generator.resolve_type({'type': 'boolean'}, 'flag'), 'bool')
        self.assertEqual(self.generator.resolve_type({'type': 'number'}, 'amount'), 'float64')

    def test_resolve_type_arrays(self):
        """Resolve array types correctly."""
        self.assertEqual(
            self.generator.resolve_type({'type': 'array', 'items': {'type': 'string'}}, 'tags'),
            '[]string'
        )

    def test_generate_type_returns_generated_type(self):
        """Generate type returns GeneratedType dataclass."""
        result = self.generator.generate_type(
            'v0.0.44_simple_type',
            self.spec['components']['schemas']['v0.0.44_simple_type'],
            'SimpleType'
        )
        self.assertIsInstance(result, GeneratedType)
        self.assertEqual(result.filename, 'simpletype.go')
        self.assertIn('type SimpleType struct', result.code)
        self.assertIn('Name string `json:"name"`', result.code)  # Required
        self.assertIn('*int32 `json:"count,omitempty"`', result.code)  # Optional
        self.assertEqual(result.field_count, 2)
        self.assertEqual(result.type_name, 'SimpleType')

    def test_output_format_minimal(self):
        """Minimal output format removes comments."""
        generator = TypeGenerator(
            self.spec,
            '0.0.44',
            {'v0.0.44_simple_type': 'SimpleType'},
            output_format=OutputFormat.MINIMAL
        )
        result = generator.generate_type(
            'v0.0.44_simple_type',
            self.spec['components']['schemas']['v0.0.44_simple_type'],
            'SimpleType'
        )
        # Should not have type comment
        self.assertNotIn('// SimpleType represents', result.code)
        # But should have struct definition
        self.assertIn('type SimpleType struct', result.code)


class TestEnumTypeOverrides(unittest.TestCase):
    """Tests for semantic enum type overrides."""

    def setUp(self):
        """Reset global enum state."""
        TypeGenerator._global_enum_types.clear()

    def test_default_enum_overrides_exist(self):
        """Default enum overrides are defined."""
        self.assertIn(('node', 'state'), DEFAULT_ENUM_TYPE_OVERRIDES)
        self.assertIn(('node', 'next_state_after_reboot'), DEFAULT_ENUM_TYPE_OVERRIDES)
        self.assertEqual(
            DEFAULT_ENUM_TYPE_OVERRIDES[('node', 'state')],
            DEFAULT_ENUM_TYPE_OVERRIDES[('node', 'next_state_after_reboot')]
        )


class TestGoAcronyms(unittest.TestCase):
    """Tests for GO_ACRONYMS constant."""

    def test_common_acronyms(self):
        """Common acronyms are defined."""
        self.assertEqual(GO_ACRONYMS['id'], 'ID')
        self.assertEqual(GO_ACRONYMS['url'], 'URL')
        self.assertEqual(GO_ACRONYMS['api'], 'API')
        self.assertEqual(GO_ACRONYMS['http'], 'HTTP')

    def test_slurm_specific_acronyms(self):
        """SLURM-specific acronyms are defined."""
        self.assertEqual(GO_ACRONYMS['qos'], 'QoS')
        self.assertEqual(GO_ACRONYMS['tres'], 'TRES')
        self.assertEqual(GO_ACRONYMS['gres'], 'GRES')
        self.assertEqual(GO_ACRONYMS['mcs'], 'MCS')


class TestDiscoverSchemas(unittest.TestCase):
    """Tests for discover_schemas function."""

    def test_discover_from_refs(self):
        """Discover schemas referenced by base schemas."""
        spec = {
            'components': {
                'schemas': {
                    'v0.0.44_job': {
                        'type': 'object',
                        'properties': {
                            'account': {'$ref': '#/components/schemas/v0.0.44_account'}
                        }
                    },
                    'v0.0.44_account': {
                        'type': 'object',
                        'properties': {
                            'name': {'type': 'string'}
                        }
                    }
                }
            }
        }
        base = {'v0.0.44_job': 'Job'}
        schemas, array_unwrap = discover_schemas(spec, '0.0.44', base)

        self.assertIn('v0.0.44_job', schemas)
        self.assertIn('v0.0.44_account', schemas)

    def test_skip_primitive_patterns(self):
        """Skip schemas that match primitive unwrap patterns."""
        spec = {
            'components': {
                'schemas': {
                    'v0.0.44_job': {
                        'type': 'object',
                        'properties': {
                            'time': {'$ref': '#/components/schemas/v0.0.44_uint64_no_val_struct'}
                        }
                    },
                    'v0.0.44_uint64_no_val_struct': {
                        'type': 'object',
                        'properties': {'value': {'type': 'integer'}}
                    }
                }
            }
        }
        base = {'v0.0.44_job': 'Job'}
        schemas, _ = discover_schemas(spec, '0.0.44', base)

        # Should not include the primitive wrapper
        self.assertIn('v0.0.44_job', schemas)
        self.assertNotIn('v0.0.44_uint64_no_val_struct', schemas)


class TestValidateOpenAPISpec(unittest.TestCase):
    """Tests for validate_openapi_spec function."""

    def test_valid_spec(self):
        """Valid spec passes validation."""
        spec = {
            'info': {'version': '1.0'},
            'components': {
                'schemas': {
                    'v0.0.44_test': {'type': 'object'}
                }
            }
        }
        # Should not raise, returns warnings list
        warnings = validate_openapi_spec(spec, Path('test.json'))
        self.assertIsInstance(warnings, list)

    def test_missing_components(self):
        """Missing components section raises error."""
        spec = {'info': {'version': '1.0'}}
        with self.assertRaises(OpenAPIValidationError) as ctx:
            validate_openapi_spec(spec, Path('test.json'))
        self.assertIn('components', str(ctx.exception))

    def test_empty_schemas(self):
        """Empty schemas section raises error."""
        spec = {
            'info': {'version': '1.0'},
            'components': {'schemas': {}}
        }
        with self.assertRaises(OpenAPIValidationError) as ctx:
            validate_openapi_spec(spec, Path('test.json'))
        self.assertIn('empty', str(ctx.exception))

    def test_missing_info_warning(self):
        """Missing info section produces warning."""
        spec = {
            'components': {
                'schemas': {
                    'v0.0.44_test': {'type': 'object'}
                }
            }
        }
        warnings = validate_openapi_spec(spec, Path('test.json'))
        self.assertTrue(any('info' in w for w in warnings))


if __name__ == '__main__':
    unittest.main()
