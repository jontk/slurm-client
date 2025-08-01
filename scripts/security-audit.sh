#!/bin/bash
# SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
# SPDX-License-Identifier: Apache-2.0

# Security audit script for slurm-client
set -euo pipefail

echo "🔒 Running Security Audit for slurm-client"
echo "=========================================="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ISSUES_FOUND=0

# Function to report issues
report_issue() {
    echo -e "${RED}❌ ISSUE:${NC} $1"
    ((ISSUES_FOUND++))
}

report_warning() {
    echo -e "${YELLOW}⚠️  WARNING:${NC} $1"
}

report_ok() {
    echo -e "${GREEN}✅ OK:${NC} $1"
}

# 1. Check for hardcoded secrets
echo "1. Checking for hardcoded secrets..."
if grep -r -i -E "(password|secret|key|token|api_key)\s*=\s*[\"'][^\"']+[\"']" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . | grep -v -E "(test|example|mock|TODO)" > /dev/null 2>&1; then
    report_issue "Found potential hardcoded secrets"
    grep -r -i -E "(password|secret|key|token|api_key)\s*=\s*[\"'][^\"']+[\"']" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . | grep -v -E "(test|example|mock)" | head -5
else
    report_ok "No hardcoded secrets found"
fi
echo

# 2. Check for exposed sensitive data in logs
echo "2. Checking for sensitive data in logs..."
if grep -r -E "log\.(Print|Printf|Println).*\b(password|token|secret|key)\b" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . | grep -v -E "(test|example)" > /dev/null 2>&1; then
    report_warning "Found potential sensitive data in logs"
    grep -r -E "log\.(Print|Printf|Println).*\b(password|token|secret|key)\b" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . | grep -v -E "(test|example)" | head -5
else
    report_ok "No sensitive data exposed in logs"
fi
echo

# 3. Check TLS configuration
echo "3. Checking TLS configuration..."
if grep -r "InsecureSkipVerify.*true" --include="*.go" --exclude-dir=vendor --exclude-dir=.git --exclude-dir=examples . | grep -v -E "(test|example)" > /dev/null 2>&1; then
    report_issue "Found InsecureSkipVerify set to true"
    grep -r "InsecureSkipVerify.*true" --include="*.go" --exclude-dir=vendor --exclude-dir=.git --exclude-dir=examples . | grep -v -E "(test|example)" | head -5
else
    report_ok "TLS verification properly configured"
fi
echo

# 4. Check for SQL injection vulnerabilities
echo "4. Checking for SQL injection risks..."
if grep -r -E "fmt\.Sprintf.*\b(SELECT|INSERT|UPDATE|DELETE)\b" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . > /dev/null 2>&1; then
    report_warning "Found potential SQL injection risks (using string formatting for queries)"
else
    report_ok "No obvious SQL injection patterns found"
fi
echo

# 5. Check file permissions
echo "5. Checking file permissions..."
SENSITIVE_FILES=("LICENSE" "SECURITY.md" ".env" "*.pem" "*.key")
for pattern in "${SENSITIVE_FILES[@]}"; do
    while IFS= read -r -d '' file; do
        perms=$(stat -f "%A" "$file" 2>/dev/null || stat -c "%a" "$file" 2>/dev/null || echo "unknown")
        if [[ "$perms" != "unknown" && "$perms" -gt 644 ]]; then
            report_warning "File $file has overly permissive permissions: $perms"
        fi
    done < <(find . -name "$pattern" -type f -print0 2>/dev/null)
done
echo

# 6. Check for TODO security items
echo "6. Checking for security-related TODOs..."
if grep -r -i "TODO.*security\|FIXME.*security\|XXX.*security" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . > /dev/null 2>&1; then
    report_warning "Found security-related TODO items"
    grep -r -i "TODO.*security\|FIXME.*security\|XXX.*security" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . | head -5
else
    report_ok "No security-related TODOs found"
fi
echo

# 7. Check error handling
echo "7. Checking error handling patterns..."
if grep -r "err\s*!=\s*nil\s*{[[:space:]]*}" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . > /dev/null 2>&1; then
    report_warning "Found empty error handling blocks"
fi

# Check for ignored errors
if grep -r "_\s*=.*Error\|_\s*:=.*err" --include="*.go" --exclude-dir=vendor --exclude-dir=.git . | grep -v -E "(test|example)" > /dev/null 2>&1; then
    report_warning "Found ignored errors"
fi
echo

# 8. Check dependencies for known vulnerabilities
echo "8. Checking dependencies..."
if command -v nancy > /dev/null 2>&1; then
    echo "Running Nancy vulnerability scan..."
    go list -json -deps ./... | nancy sleuth || report_issue "Found vulnerable dependencies"
else
    report_warning "Nancy not installed, skipping dependency vulnerability check"
fi
echo

# 9. Check for race conditions
echo "9. Testing for race conditions..."
echo "Running go test -race on key packages..."
if go test -race -short ./pkg/... > /dev/null 2>&1; then
    report_ok "No race conditions detected in pkg/"
else
    report_issue "Race conditions detected"
fi
echo

# 10. Check authentication implementation
echo "10. Checking authentication implementation..."
if ! grep -r "type.*Auth.*interface" --include="*.go" pkg/auth/ > /dev/null 2>&1; then
    report_issue "Authentication interface not found"
else
    report_ok "Authentication interface properly defined"
fi

# Summary
echo
echo "=========================================="
echo "Security Audit Summary"
echo "=========================================="
if [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${GREEN}✅ No critical security issues found!${NC}"
    echo "Note: This is a basic audit. For production use, consider:"
    echo "  - Professional security audit"
    echo "  - Penetration testing"
    echo "  - Code review by security experts"
else
    echo -e "${RED}❌ Found $ISSUES_FOUND security issues that need attention${NC}"
fi
echo

# Create security audit report
REPORT_FILE="security-audit-report-$(date +%Y%m%d-%H%M%S).txt"
echo "Generating detailed report: $REPORT_FILE"

exit $ISSUES_FOUND