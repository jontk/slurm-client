#!/usr/bin/env python3
import json
import sys

def extract_version_spec(input_file, output_file, version="v0.0.41"):
    """Extract specific version endpoints and schemas from combined OpenAPI spec"""
    
    with open(input_file, 'r') as f:
        spec = json.load(f)
    
    # Create new spec with basic info
    new_spec = {
        "openapi": spec.get("openapi", "3.0.0"),
        "info": {
            "title": f"Slurm REST API - {version}",
            "version": version,
            "description": f"Extracted {version} endpoints from Slurm REST API"
        },
        "servers": spec.get("servers", []),
        "paths": {},
        "components": {
            "schemas": {},
            "securitySchemes": spec.get("components", {}).get("securitySchemes", {})
        }
    }
    
    # Extract paths for this version
    for path, methods in spec.get("paths", {}).items():
        if f"/{version}/" in path or f"/v0041/" in path:
            new_spec["paths"][path] = methods
    
    # Extract schemas referenced by these paths
    schemas_to_extract = set()
    
    def find_refs(obj):
        """Recursively find all $ref references"""
        if isinstance(obj, dict):
            for key, value in obj.items():
                if key == "$ref" and isinstance(value, str):
                    if "#/components/schemas/" in value:
                        schema_name = value.split("/")[-1]
                        schemas_to_extract.add(schema_name)
                else:
                    find_refs(value)
        elif isinstance(obj, list):
            for item in obj:
                find_refs(item)
    
    # Find all schema references in the paths
    find_refs(new_spec["paths"])
    
    # Extract schemas for v0.0.41
    all_schemas = spec.get("components", {}).get("schemas", {})
    
    # First pass: get v0.0.41 specific schemas
    for schema_name in list(schemas_to_extract):
        if schema_name in all_schemas:
            new_spec["components"]["schemas"][schema_name] = all_schemas[schema_name]
            # Find nested references
            find_refs(all_schemas[schema_name])
    
    # Also include schemas that start with v0_0_41 or v0.0.41
    for schema_name, schema_def in all_schemas.items():
        if schema_name.startswith("v0_0_41") or schema_name.startswith("v0.0.41"):
            new_spec["components"]["schemas"][schema_name] = schema_def
            schemas_to_extract.add(schema_name)
            find_refs(schema_def)
    
    # Second pass: get all referenced schemas
    max_iterations = 10  # Prevent infinite loops
    for _ in range(max_iterations):
        initial_count = len(schemas_to_extract)
        for schema_name in list(schemas_to_extract):
            if schema_name in all_schemas and schema_name not in new_spec["components"]["schemas"]:
                new_spec["components"]["schemas"][schema_name] = all_schemas[schema_name]
                find_refs(all_schemas[schema_name])
        if len(schemas_to_extract) == initial_count:
            break
    
    # Write the extracted spec
    with open(output_file, 'w') as f:
        json.dump(new_spec, f, indent=2)
    
    print(f"Extracted {len(new_spec['paths'])} paths and {len(new_spec['components']['schemas'])} schemas for {version}")

if __name__ == "__main__":
    extract_version_spec("slurm-openapi-v3.json", "slurm-v0.0.41-extracted.json", "v0.0.41")