#!/bin/bash

# Download OpenAPI specifications from SLURM REST API
# This script connects to a running slurmrestd instance and downloads the OpenAPI specs

set -e

# Configuration
SLURM_REST_URL="${SLURM_REST_URL:-http://localhost:6820}"
SPECS_DIR="openapi-specs"

# Create specs directory if it doesn't exist
mkdir -p "$SPECS_DIR"

echo "Downloading OpenAPI specifications from $SLURM_REST_URL..."

# List of API versions to download
VERSIONS=("v0.0.40" "v0.0.41" "v0.0.42" "v0.0.43" "v0.0.44" "v0.0.45")

# Function to download a specific version
download_version() {
    local version=$1
    local output_file="$SPECS_DIR/slurm-$version.json"
    
    echo "Downloading $version..."
    
    # Try to download the OpenAPI spec for this version
    if curl -f -s -o "$output_file" "$SLURM_REST_URL/openapi/$version" || \
       curl -f -s -o "$output_file" "$SLURM_REST_URL/openapi/v3?version=$version"; then
        echo "✓ Downloaded $version to $output_file"
        
        # Validate that it's valid JSON
        if ! jq empty "$output_file" 2>/dev/null; then
            echo "⚠ Warning: $output_file doesn't contain valid JSON"
            rm -f "$output_file"
            return 1
        fi
    else
        echo "⚠ Failed to download $version (may not be available)"
        rm -f "$output_file"
        return 1
    fi
}

# Check if slurmrestd is running
if ! curl -f -s "$SLURM_REST_URL/openapi/v3" >/dev/null 2>&1; then
    echo "Error: Cannot connect to slurmrestd at $SLURM_REST_URL"
    echo "Please ensure:"
    echo "1. SLURM is installed and configured"
    echo "2. slurmrestd is running on $SLURM_REST_URL"
    echo "3. The REST API plugins are loaded"
    echo ""
    echo "To start slurmrestd with v0.0.44 support:"
    echo "  slurmrestd -s openapi/slurmctld,openapi/slurmdbd,openapi/util -s data_parser/v0.0.45 localhost:6820"
    exit 1
fi

# Download specifications for each version
success_count=0
for version in "${VERSIONS[@]}"; do
    if download_version "$version"; then
        ((success_count++))
    fi
done

echo ""
echo "Downloaded $success_count OpenAPI specifications successfully."

# Special handling for v0.0.44 - try alternative endpoints if standard ones fail
if [ ! -f "$SPECS_DIR/slurm-v0.0.44.json" ]; then
    echo "Attempting alternative download methods for v0.0.44..."
    
    # Try different potential endpoints
    ENDPOINTS=(
        "/openapi/v0.0.44"
        "/openapi/v3?version=v0.0.44"
        "/openapi/v3?plugin=v0.0.44"
        "/openapi?version=0.0.44"
    )
    
    for endpoint in "${ENDPOINTS[@]}"; do
        echo "Trying $SLURM_REST_URL$endpoint..."
        if curl -f -s -o "$SPECS_DIR/slurm-v0.0.44.json" "$SLURM_REST_URL$endpoint"; then
            if jq empty "$SPECS_DIR/slurm-v0.0.44.json" 2>/dev/null; then
                echo "✓ Successfully downloaded v0.0.44 from $endpoint"
                break
            else
                rm -f "$SPECS_DIR/slurm-v0.0.44.json"
            fi
        fi
    done
fi

# List downloaded specifications
echo ""
echo "Available OpenAPI specifications:"
ls -la "$SPECS_DIR"/ || true

# Instructions for building SLURM with v0.0.44 support
if [ ! -f "$SPECS_DIR/slurm-v0.0.44.json" ]; then
    echo ""
    echo "📋 To build SLURM with v0.0.44 REST API support:"
    echo ""
    echo "1. Get SLURM source (version 25.11.x or later):"
    echo "   git clone https://github.com/SchedMD/slurm.git"
    echo "   cd slurm"
    echo ""
    echo "2. Configure with REST API support:"
    echo "   ./configure --enable-rest --enable-openapi"
    echo ""
    echo "3. Build:"
    echo "   make -j\$(nproc)"
    echo ""
    echo "4. Start slurmrestd with v0.0.44 plugin:"
    echo "   ./src/slurmrestd/slurmrestd -s openapi/slurmctld,openapi/slurmdbd,openapi/util -s data_parser/v0.0.45 localhost:6820"
    echo ""
fi

echo "Done."