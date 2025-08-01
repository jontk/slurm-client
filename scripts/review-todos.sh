#!/bin/bash
# SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
# SPDX-License-Identifier: Apache-2.0

# Script to review and categorize TODO comments
set -euo pipefail

echo "📝 TODO Comments Review"
echo "======================"
echo

# Colors for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Categories
CRITICAL_COUNT=0
NORMAL_COUNT=0
ENHANCEMENT_COUNT=0

echo -e "${RED}Critical TODOs (FIXME, XXX, HACK):${NC}"
echo "--------------------------------"
grep -r -n "FIXME\|XXX\|HACK" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . 2>/dev/null | while IFS= read -r line; do
    echo "  $line"
    ((CRITICAL_COUNT++))
done || true
echo

echo -e "${YELLOW}Security-related TODOs:${NC}"
echo "----------------------"
grep -r -n -i "TODO.*\(security\|auth\|encrypt\|password\|token\|credential\)" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . 2>/dev/null | while IFS= read -r line; do
    echo "  $line"
done || true
echo

echo -e "${BLUE}Normal TODOs:${NC}"
echo "------------"
grep -r -n "TODO" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . 2>/dev/null | grep -v -i "\(FIXME\|XXX\|HACK\|security\|auth\|encrypt\|password\|token\|credential\)" | head -20 | while IFS= read -r line; do
    echo "  $line"
    ((NORMAL_COUNT++))
done || true
echo

# Summary by file
echo -e "${GREEN}Summary by File:${NC}"
echo "---------------"
echo "Files with most TODOs:"
grep -r "TODO\|FIXME\|XXX\|HACK" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . 2>/dev/null | cut -d: -f1 | sort | uniq -c | sort -nr | head -10
echo

# Total count
TOTAL_COUNT=$(grep -r "TODO\|FIXME\|XXX\|HACK" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . 2>/dev/null | wc -l | tr -d ' ')
echo "Total TODO comments: $TOTAL_COUNT"

# Generate TODO report
REPORT_FILE="todo-report-$(date +%Y%m%d-%H%M%S).md"
echo
echo "Generating detailed report: $REPORT_FILE"

cat > "$REPORT_FILE" << EOF
# TODO Comments Report
Generated on: $(date)

## Summary
- Total TODO comments: $TOTAL_COUNT
- Critical (FIXME/XXX/HACK): $(grep -r "FIXME\|XXX\|HACK" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . 2>/dev/null | wc -l | tr -d ' ')
- Security-related: $(grep -r -i "TODO.*\(security\|auth\|encrypt\|password\|token\|credential\)" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . 2>/dev/null | wc -l | tr -d ' ')

## Critical Items Requiring Immediate Attention

$(grep -r -n "FIXME\|XXX\|HACK" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . 2>/dev/null || echo "None found")

## Security-Related TODOs

$(grep -r -n -i "TODO.*\(security\|auth\|encrypt\|password\|token\|credential\)" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . 2>/dev/null || echo "None found")

## All TODO Comments

$(grep -r -n "TODO\|FIXME\|XXX\|HACK" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . 2>/dev/null)
EOF

echo "Report generated: $REPORT_FILE"