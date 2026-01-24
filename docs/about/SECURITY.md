# Security

This document outlines the security features and practices for the SLURM REST API Client Library.

## Supply Chain Security

### Software Bill of Materials (SBOM)

Starting with releases, each release includes a Software Bill of Materials (SBOM) that documents all dependencies and components included in the release artifacts.

**SBOM Files:**
- Format: SPDX JSON
- Location: Attached to each GitHub release
- Naming: `slurm-client_<version>_<os>_<arch>.sbom.json`

**What's Included:**
- All Go module dependencies
- Dependency versions and licenses
- Transitive dependencies
- Build-time information

**Usage:**
Download the SBOM from the release page to:
- Audit dependencies for security vulnerabilities
- Ensure license compliance
- Track supply chain components
- Integrate with vulnerability scanning tools

### Code Signing

All release artifacts are signed using [Cosign](https://github.com/sigstore/cosign) with keyless signing via [Sigstore](https://www.sigstore.dev/).

**What's Signed:**
- Release checksums file (`checksums.txt`)
- Provides cryptographic proof of authenticity

**Keyless Signing:**
- Uses OIDC identity from GitHub Actions
- No private keys to manage or secure
- Signatures stored in public transparency log (Rekor)
- Certificates issued by Fulcio CA

**Verifying Signatures:**

1. **Install cosign:**
   ```bash
   # macOS
   brew install cosign

   # Linux
   wget https://github.com/sigstore/cosign/releases/download/v2.4.1/cosign-linux-amd64
   chmod +x cosign-linux-amd64
   sudo mv cosign-linux-amd64 /usr/local/bin/cosign
   ```

2. **Download release artifacts:**
   ```bash
   # Download checksum file and signature
   VERSION="v1.0.0"
   wget https://github.com/jontk/slurm-client/releases/download/${VERSION}/checksums.txt
   wget https://github.com/jontk/slurm-client/releases/download/${VERSION}/checksums.txt.sig
   wget https://github.com/jontk/slurm-client/releases/download/${VERSION}/checksums.txt.pem
   ```

3. **Verify signature:**
   ```bash
   cosign verify-blob \
     --certificate checksums.txt.pem \
     --signature checksums.txt.sig \
     --certificate-identity-regexp "https://github.com/jontk/slurm-client" \
     --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
     checksums.txt
   ```

4. **Verify artifact checksums:**
   ```bash
   # After verification, check artifact integrity
   sha256sum --check checksums.txt
   ```

**Transparency:**
- All signatures are recorded in Rekor transparency log
- View certificate details at https://rekor.sigstore.dev
- Provides tamper-evident audit trail

## Vulnerability Disclosure

### Reporting Security Issues

**Please DO NOT report security vulnerabilities through public GitHub issues.**

Instead, report security vulnerabilities via GitHub's private vulnerability reporting feature:

1. Navigate to the [Security tab](https://github.com/jontk/slurm-client/security)
2. Click "Report a vulnerability"
3. Provide detailed information about the vulnerability

Alternatively, email security concerns to: [security contact needed]

### What to Include

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact
- Suggested fix (if you have one)
- Your contact information

### Response Timeline

- **Initial Response:** Within 48 hours
- **Status Update:** Within 7 days
- **Fix Timeline:** Depends on severity (critical issues prioritized)

### Disclosure Policy

- Security vulnerabilities will be fixed before public disclosure
- Coordinated disclosure with security researchers
- Credit given to security researchers (if desired)

## Security Scanning

The project uses multiple automated security scanners:

### GitHub Security Features
- **Dependabot:** Automated dependency updates for security patches
- **CodeQL:** Static analysis for security vulnerabilities
- **Secret Scanning:** Prevents accidental credential commits

### Third-Party Scanners
- **Gosec:** Go security checker in CI/CD
- **Trivy:** Container and dependency vulnerability scanner
- **Nancy:** Dependency vulnerability checker

### CI/CD Security Checks
All pull requests are scanned for:
- Known vulnerabilities in dependencies
- Security issues in code patterns
- License compliance
- SPDX header compliance

## Dependency Management

### Regular Updates
- Dependabot checks daily for updates
- Security patches applied within 48 hours
- Major version updates reviewed manually

### Minimal Dependencies
- Project minimizes external dependencies
- Direct dependencies are carefully vetted
- Unused dependencies are removed

## Best Practices for Users

### When Using the Library

1. **Keep Updated:**
   ```bash
   go get -u github.com/jontk/slurm-client@latest
   ```

2. **Scan Your Dependencies:**
   ```bash
   go list -m all | nancy sleuth
   ```

3. **Review SBOMs:**
   - Download SBOM from releases
   - Scan with your vulnerability scanner
   - Verify licenses match your requirements

4. **Verify Downloads:**
   - Always verify signatures before use
   - Check checksums match
   - Download from official GitHub releases only

### Secure Development

1. **API Credentials:**
   - Never commit credentials to version control
   - Use environment variables or secret managers
   - Rotate credentials regularly

2. **HTTPS Only:**
   - Always use HTTPS for SLURM API connections
   - Verify TLS certificates
   - Consider certificate pinning for production

3. **Error Handling:**
   - Don't log sensitive information
   - Sanitize error messages
   - Use appropriate log levels

4. **Input Validation:**
   - Validate all user inputs
   - Sanitize data before API calls
   - Use type-safe API methods

## Security Updates

Subscribe to security updates:
- Watch the repository for security advisories
- Enable GitHub security alerts
- Follow release notes for security fixes

## Audit Trail

All releases include:
- SBOM documenting all components
- Cryptographic signatures for verification
- Transparency log entries (Rekor)
- Detailed changelog

## Questions?

For security-related questions (non-vulnerabilities):
- Check [existing issues](https://github.com/jontk/slurm-client/issues)
- Review [documentation](./README.md)
- Open a [new issue](https://github.com/jontk/slurm-client/issues/new) for questions
