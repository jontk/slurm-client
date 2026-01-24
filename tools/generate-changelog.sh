#!/bin/bash
# SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
# SPDX-License-Identifier: Apache-2.0

# Generate CHANGELOG.md from conventional commits

set -e

CHANGELOG_FILE="${1:-docs/development/CHANGELOG.md}"

# Function to generate changelog section for a range
generate_section_for_range() {
    local TO=$1
    local FROM=$2

    local RANGE
    if [ -z "$FROM" ]; then
        RANGE="$TO"
    else
        RANGE="${FROM}..${TO}"
    fi

    # Get commits by type
    local FEATURES=$(git log --format="%s|||%b" "$RANGE" | grep "^feat" || true)
    local FIXES=$(git log --format="%s|||%b" "$RANGE" | grep "^fix" || true)
    local DOCS=$(git log --format="%s|||%b" "$RANGE" | grep "^docs" || true)
    local PERFORMANCE=$(git log --format="%s|||%b" "$RANGE" | grep "^perf" || true)
    local REFACTOR=$(git log --format="%s|||%b" "$RANGE" | grep "^refactor" || true)
    local TESTS=$(git log --format="%s|||%b" "$RANGE" | grep "^test" || true)
    local BUILD=$(git log --format="%s|||%b" "$RANGE" | grep "^build" || true)
    local CI=$(git log --format="%s|||%b" "$RANGE" | grep "^ci" || true)
    local CHORES=$(git log --format="%s|||%b" "$RANGE" | grep "^chore" || true)

    # Check for breaking changes
    local BREAKING=$(git log --format="%s|||%b" "$RANGE" | grep "BREAKING CHANGE" || true)

    if [ -n "$BREAKING" ]; then
        echo "### ⚠️ BREAKING CHANGES"
        echo ""
        echo "$BREAKING" | while IFS='|||' read -r subject body; do
            # Extract description after BREAKING CHANGE:
            local desc=$(echo "$body" | grep -A 100 "BREAKING CHANGE:" | tail -n +2 | head -n 1)
            if [ -z "$desc" ]; then
                desc=$(echo "$subject" | sed 's/^[^:]*: //')
            fi
            echo "- $desc"
        done
        echo ""
    fi

    if [ -n "$FEATURES" ]; then
        echo "### ✨ Features"
        echo ""
        echo "$FEATURES" | while IFS='|||' read -r subject body; do
            local msg=$(echo "$subject" | sed 's/^feat[(].*[)]: //; s/^feat: //')
            echo "- $msg"
        done
        echo ""
    fi

    if [ -n "$FIXES" ]; then
        echo "### 🐛 Bug Fixes"
        echo ""
        echo "$FIXES" | while IFS='|||' read -r subject body; do
            local msg=$(echo "$subject" | sed 's/^fix[(].*[)]: //; s/^fix: //')
            echo "- $msg"
        done
        echo ""
    fi

    if [ -n "$PERFORMANCE" ]; then
        echo "### ⚡ Performance"
        echo ""
        echo "$PERFORMANCE" | while IFS='|||' read -r subject body; do
            local msg=$(echo "$subject" | sed 's/^perf[(].*[)]: //; s/^perf: //')
            echo "- $msg"
        done
        echo ""
    fi

    if [ -n "$DOCS" ]; then
        echo "### 📚 Documentation"
        echo ""
        echo "$DOCS" | while IFS='|||' read -r subject body; do
            local msg=$(echo "$subject" | sed 's/^docs[(].*[)]: //; s/^docs: //')
            echo "- $msg"
        done
        echo ""
    fi

    if [ -n "$REFACTOR" ]; then
        echo "### ♻️ Code Refactoring"
        echo ""
        echo "$REFACTOR" | while IFS='|||' read -r subject body; do
            local msg=$(echo "$subject" | sed 's/^refactor[(].*[)]: //; s/^refactor: //')
            echo "- $msg"
        done
        echo ""
    fi

    if [ -n "$TESTS" ]; then
        echo "### ✅ Tests"
        echo ""
        echo "$TESTS" | while IFS='|||' read -r subject body; do
            local msg=$(echo "$subject" | sed 's/^test[(].*[)]: //; s/^test: //')
            echo "- $msg"
        done
        echo ""
    fi

    if [ -n "$BUILD" ] || [ -n "$CI" ] || [ -n "$CHORES" ]; then
        echo "### 🔧 Maintenance"
        echo ""

        if [ -n "$BUILD" ]; then
            echo "$BUILD" | while IFS='|||' read -r subject body; do
                local msg=$(echo "$subject" | sed 's/^build[(].*[)]: //; s/^build: //')
                echo "- Build: $msg"
            done
        fi

        if [ -n "$CI" ]; then
            echo "$CI" | while IFS='|||' read -r subject body; do
                local msg=$(echo "$subject" | sed 's/^ci[(].*[)]: //; s/^ci: //')
                echo "- CI: $msg"
            done
        fi

        if [ -n "$CHORES" ]; then
            echo "$CHORES" | while IFS='|||' read -r subject body; do
                local msg=$(echo "$subject" | sed 's/^chore[(].*[)]: //; s/^chore: //')
                echo "- Chore: $msg"
            done
        fi

        echo ""
    fi
}

echo "Generating CHANGELOG from git history..."

# Create temporary file for the changelog
TEMP_FILE=$(mktemp)

# Add header
cat > "$TEMP_FILE" <<'EOF'
# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
and uses [Conventional Commits](https://www.conventionalcommits.org/) for commit messages.

EOF

# Get all tags sorted by version
TAGS=$(git tag -l "v*" --sort=-version:refname)

# If no tags exist, generate from all commits
if [ -z "$TAGS" ]; then
    echo "## [Unreleased]" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    generate_section_for_range "HEAD" "" >> "$TEMP_FILE"
else
    # Generate unreleased section
    LATEST_TAG=$(echo "$TAGS" | head -n1)
    echo "## [Unreleased]" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    generate_section_for_range "HEAD" "$LATEST_TAG" >> "$TEMP_FILE"

    # Generate sections for each tag
    PREV_TAG=""
    for TAG in $TAGS; do
        # Get tag date
        TAG_DATE=$(git log -1 --format=%ai "$TAG" | cut -d' ' -f1)

        echo "## [$TAG] - $TAG_DATE" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"

        if [ -z "$PREV_TAG" ]; then
            generate_section_for_range "$TAG" "" >> "$TEMP_FILE"
        else
            generate_section_for_range "$TAG" "$PREV_TAG" >> "$TEMP_FILE"
        fi

        PREV_TAG="$TAG"
    done
fi

# Move temp file to final location
mv "$TEMP_FILE" "$CHANGELOG_FILE"

echo "✓ CHANGELOG generated successfully: $CHANGELOG_FILE"
