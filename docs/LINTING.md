# Linting Governance and Workflow

This document describes how linting is managed in the slurm-client project, including how to run linters locally, propose changes, and the code review process.

## Quick Start

### Install Pre-commit Hooks

Pre-commit hooks automatically run linters before each commit:

```bash
pre-commit install
```

This ensures your code is properly formatted and linted before it reaches the review process.

### Run Linters Locally

#### Full Lint Check

Run the complete linting suite against all files:

```bash
make lint
```

Or directly with golangci-lint:

```bash
golangci-lint run ./...
```

#### Lint Only Staged Changes

Check only the files you've staged for commit:

```bash
make lint-staged
```

#### Format Code

Automatically format and fix simple issues:

```bash
make fmt
```

Or individually:

```bash
gofmt -w ./...
goimports -w ./...
```

## Linting Configuration

### Configuration File

All linting configuration is centralized in `.golangci.yml`:

```yaml
version: "2"
linters:
  default: all
  disable:
    # Comprehensive list of disabled linters with rationale
    # See docs/linting/DISABLED_LINTERS.md for details
```

### Understanding the Configuration

- **`default: all`** - Start with all available linters enabled
- **`disable`** - Selectively disable linters that don't fit the project's philosophy
- **`settings`** - Configure behavior for enabled linters
- **`exclusions`** - Exclude specific rules from specific paths (e.g., tests, generated code)

### Disabled Linters

For a complete list of disabled linters and the rationale behind each decision, see [docs/linting/DISABLED_LINTERS.md](./linting/DISABLED_LINTERS.md).

## Development Workflow

### Before Committing

1. **Run linters locally**:
   ```bash
   make lint
   ```

2. **Fix issues automatically** (when possible):
   ```bash
   make fmt
   ```

3. **Verify no issues remain**:
   ```bash
   make lint
   ```

4. **Stage and commit** - Pre-commit hooks will run automatically

### During Code Review

The continuous integration (CI) pipeline runs a comprehensive lint check on every pull request.

**What CI checks:**

- **Full golangci-lint** - All linters enabled as configured in `.golangci.yml`
- **Commit message linting** - Enforces [Conventional Commits](https://www.conventionalcommits.org/) format
- **Pre-commit hooks** - Standard file hygiene (trailing whitespace, file endings, etc.)

### CI Lint Failures

If CI reports a lint failure on your PR:

1. **Read the error message** - It includes the file, line number, and specific issue
2. **Run locally to reproduce**:
   ```bash
   make lint
   ```
3. **Fix the issue**:
   - Some fixes are automatic: `make fmt`
   - Others require code changes - follow the linter's suggestion
4. **Verify the fix**:
   ```bash
   make lint
   ```
5. **Push the changes** - CI will re-run automatically

## Proposing Linter Changes

### Adding a New Linter

To enable a currently disabled linter or enable a new linter:

1. **Research the linter**:
   - Check what it does and why it was disabled (see [docs/linting/DISABLED_LINTERS.md](./linting/DISABLED_LINTERS.md))
   - Run it locally to see violations: `golangci-lint run -E <linter> ./...`
   - Count violations and assess signal-to-noise ratio

2. **Evaluate the impact**:
   - How many violations exist?
   - Are they legitimate issues or false positives?
   - What refactoring effort is required?
   - Is the benefit worth the effort?

3. **Propose the change**:
   - Create an issue or discussion describing why the linter should be enabled
   - Include the number and types of violations
   - Explain the expected benefits
   - Estimate the refactoring effort

4. **Implementation**:
   - Enable the linter in `.golangci.yml`
   - Fix all violations in a dedicated PR
   - Include a clear commit message (e.g., "feat(linters): enable <linter_name>")
   - Update [docs/linting/DISABLED_LINTERS.md](./linting/DISABLED_LINTERS.md) if rationale changes

### Disabling a Linter

To disable a currently enabled linter:

1. **Document the rationale** - Update `.golangci.yml` with a comment explaining why
2. **Update documentation** - Add an entry to [docs/linting/DISABLED_LINTERS.md](./linting/DISABLED_LINTERS.md)
3. **Create a PR** - Include the change and documentation update

### Adjusting Linter Settings

To modify thresholds or rules for an existing linter:

1. **Update `.golangci.yml`** in the appropriate `settings` section
2. **Test locally** - Verify the new settings work as intended
3. **Document the change** - Update relevant documentation
4. **Create a PR** - Include rationale in the commit message

## Linting Phases

The project uses a phased approach to gradually introduce stricter linting rules through the `revive` linter:

### Current Phases (Complete)

- **Phase 1**: `empty-block`, `context-as-argument`, `context-keys-type`
- **Phase 2**: `exported`, `unused-parameter`, `unreachable-code`
- **Phase 3**: `error-naming`, `error-return`, `receiver-naming`, `indent-error-flow`, `blank-imports`, `var-declaration`, `unnecessary-stmt`

### Future Phases

Additional revive rules may be enabled in future phases based on code quality improvements and developer consensus.

## CI/CD Integration

### GitHub Actions Workflow

The linting check runs on:
- Every push to feature branches
- Every pull request
- On `main` branch commits (to prevent regressions)

### Linting Job Details

**File**: `.github/workflows/ci.yml`

**What runs**:
```bash
golangci-lint run ./...
```

**Exit behavior**:
- ✅ 0 issues → Pipeline succeeds
- ❌ 1+ issues → Pipeline fails, blocking merge

**Viewing Results**:
1. Go to the PR page on GitHub
2. Click the "Checks" tab
3. Expand the "lint" check to see details
4. Click on specific violations for line-specific feedback

## Best Practices

### General Guidelines

1. **Fix lint issues promptly** - Don't accumulate technical debt
2. **Understand the rules** - Read [docs/linting/DISABLED_LINTERS.md](./linting/DISABLED_LINTERS.md) to understand philosophy
3. **Use comments when needed** - Use `// nolint: <linter>` sparingly when justified
4. **Keep it clean** - Use pre-commit hooks to catch issues before review

### When Using `// nolint`

Only use `// nolint` comments when:
- The linter produces a false positive
- There's a legitimate technical reason to suppress the warning
- The suppression is justified by a comment

Example:
```go
// nolint: gosec - G402 is acceptable here for testing
config := &tls.Config{InsecureSkipVerify: true}
```

### Common Patterns

#### Error Handling

Always wrap errors with context:
```go
if err != nil {
    return fmt.Errorf("failed to connect: %w", err)
}
```

#### Interfaces

Return interfaces when appropriate for flexibility:
```go
// Acceptable - allows flexible implementations
type Manager interface {
    CreateJob(ctx context.Context, ...) (*Job, error)
    UpdateJob(ctx context.Context, ...) (*Job, error)
}
```

#### Named Returns

Use named returns for documentation:
```go
func (c *Client) GetJob(ctx context.Context, id string) (job *Job, err error) {
    // Clear what's being returned
}
```

## Getting Help

- **Linter documentation**: Run `golangci-lint help linters` for details
- **Configuration questions**: Check `.golangci.yml` comments or [golangci-lint docs](https://golangci-lint.run/)
- **Project philosophy**: See [docs/linting/DISABLED_LINTERS.md](./linting/DISABLED_LINTERS.md)
- **Coding standards**: Refer to [CONTRIBUTING.md](../CONTRIBUTING.md)

## Related Documentation

- [DISABLED_LINTERS.md](./linting/DISABLED_LINTERS.md) - Rationale for disabled linters
- [CONTRIBUTING.md](../CONTRIBUTING.md) - Overall contribution guidelines
- [golangci-lint documentation](https://golangci-lint.run/) - External reference
- [Conventional Commits](https://www.conventionalcommits.org/) - Commit message format
