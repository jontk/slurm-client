# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          | SLURM API Version |
| ------- | ------------------ | ----------------- |
| 0.1.x   | :white_check_mark: | v0.0.40 - v0.0.44 |
| < 0.1.0 | :x:                | -                 |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Please report security vulnerabilities to **security@jontk.com**.

### What to Include

Please include the following information in your report:

- Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 5 business days
- **Resolution Target**: Based on severity
  - Critical: 7 days
  - High: 14 days
  - Medium: 30 days
  - Low: 90 days

## Preferred Languages

We prefer all communications to be in English.

## Disclosure Policy

- Security issues will be disclosed after patches are available
- We'll coordinate disclosure with the reporter
- Credit will be given to reporters (unless they prefer to remain anonymous)

## Security Best Practices

When using slurm-client, we recommend:

1. **Keep Updated**: Always use the latest version to get security updates
2. **Secure Configuration**:
   - Use HTTPS for SLURM REST API connections
   - Store API tokens securely (environment variables, secret management systems)
   - Never commit credentials to version control
3. **Network Security**:
   - Use firewalls to restrict access to SLURM REST API endpoints
   - Consider using VPN or private networks for SLURM access
4. **Authentication**:
   - Use strong authentication tokens
   - Rotate tokens regularly
   - Use separate tokens for different environments (dev/staging/prod)
5. **Monitoring**:
   - Monitor for unusual API activity
   - Log all API interactions
   - Set up alerts for authentication failures

## Known Security Considerations

### Third-Party Dependencies

This project depends on several third-party libraries. We:

- Regularly update dependencies to get security patches
- Use automated tools (Dependabot, govulncheck) to detect vulnerabilities
- Review security advisories for our dependencies

Recent security updates:
- Updated kin-openapi from v0.128.0 to v0.133.0 (addresses CVE-2025-30153)

### SLURM REST API Security

When using slurm-client with SLURM REST API:

1. **Authentication**: Ensure proper authentication is configured on your SLURM REST API server
2. **Authorization**: Configure appropriate user permissions in SLURM
3. **Network Exposure**: Limit network access to the SLURM REST API
4. **TLS/SSL**: Always use HTTPS connections to the SLURM REST API
5. **API Rate Limiting**: Consider implementing rate limiting to prevent abuse

## Comments on this Policy

If you have suggestions on how this process could be improved, please submit a pull request.