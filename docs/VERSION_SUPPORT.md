# Version Support Policy

This document outlines the support policy for different SLURM REST API versions in the slurm-client SDK.

## Supported API Versions

| API Version | SLURM Version | Support Tier | End of Life | Notes |
|-------------|---------------|--------------|-------------|-------|
| v0.0.44     | 24.11.x       | **Active**   | -           | Latest version, full feature support |
| v0.0.43     | 24.05.x       | **Active**   | -           | Full feature support |
| v0.0.42     | 23.11.x       | Maintenance  | 2026-06-01  | Security fixes only |
| v0.0.41     | 23.02.x       | Deprecated   | 2025-06-01  | Removed in next major version |
| v0.0.40     | 22.05.x       | Deprecated   | 2025-06-01  | Removed in next major version |

## Support Tier Definitions

### Active
- **Full feature support**: New features, enhancements, and bug fixes
- **Performance optimizations**: Ongoing improvements to efficiency and speed
- **Security updates**: Immediate security patches
- **Documentation**: Comprehensive and actively maintained
- **Community support**: Active support in discussions and issues

### Maintenance
- **Bug fixes**: Critical bugs only
- **Security updates**: Security patches as needed
- **No new features**: Feature development has ended
- **Limited documentation updates**: Only for critical corrections
- **Limited community support**: Best-effort support

### Deprecated
- **Security updates only**: Only critical security vulnerabilities
- **No bug fixes**: No general bug fixes
- **No new features**: No feature development
- **Documentation frozen**: No further updates
- **No community support**: Users encouraged to upgrade

## Version Selection

The SDK automatically negotiates the best available API version based on your SLURM server's capabilities. You can also explicitly specify a version:

```go
import (
    "github.com/jontk/slurm-client"
)

// Automatic version negotiation (recommended)
client, err := slurm.NewClient("https://slurm.example.com")

// Explicit version specification
client, err := slurm.NewClient("https://slurm.example.com",
    slurm.WithAPIVersion("v0.0.44"))
```

## Migration Guide

### Upgrading from Deprecated Versions

If you're using a deprecated API version (v0.0.40 or v0.0.41), we recommend upgrading to v0.0.43 or v0.0.44:

1. **Review breaking changes**: Check the [CHANGELOG](../CHANGELOG.md) for breaking changes between versions
2. **Test thoroughly**: Run your integration tests against the new version
3. **Update dependencies**: Update your go.mod to use the latest SDK version
4. **Monitor for issues**: Watch for any unexpected behavior after upgrade

### Key Differences Between Versions

#### v0.0.44 vs v0.0.43
- Enhanced cluster management features
- Additional association management capabilities
- Performance improvements in list operations

#### v0.0.43 vs v0.0.42
- Added AssociationManager and ClusterManager interfaces
- Improved error handling and reporting
- Enhanced job step tracking capabilities

#### v0.0.42 vs v0.0.41
- Improved node state handling
- Enhanced partition management
- Better pagination support

## Support Timeline

```
2025-06-01: v0.0.40 and v0.0.41 reach End of Life
2026-06-01: v0.0.42 reaches End of Life
```

## Getting Help

- **Documentation**: Check the [main README](../README.md) and [examples](../examples/)
- **Issues**: Report bugs on [GitHub Issues](https://github.com/jontk/slurm-client/issues)
- **Discussions**: Ask questions in [GitHub Discussions](https://github.com/jontk/slurm-client/discussions)

## Version-Specific Implementation Details

### Adapter-Based Architecture

All supported versions use an adapter-based architecture that:
- Converts between API-specific types and common SDK types
- Handles version-specific quirks and differences
- Provides a consistent interface across all versions

### Feature Availability

Some features are only available in specific API versions:

| Feature | v0.0.40 | v0.0.41 | v0.0.42 | v0.0.43 | v0.0.44 |
|---------|---------|---------|---------|---------|---------|
| Job Management | ✓ | ✓ | ✓ | ✓ | ✓ |
| Node Management | ✓ | ✓ | ✓ | ✓ | ✓ |
| Partition Management | ✓ | ✓ | ✓ | ✓ | ✓ |
| Association Management | ✗ | ✗ | ✗ | ✓ | ✓ |
| Cluster Management | ✗ | ✗ | ✗ | ✓ | ✓ |
| Enhanced Job Steps | ✗ | ✗ | ✓ | ✓ | ✓ |
| Job Analytics | ✗ | ✗ | ✗ | ✓ | ✓ |

## Feedback

We welcome feedback on our version support policy. Please open an issue or discussion to share your thoughts.
