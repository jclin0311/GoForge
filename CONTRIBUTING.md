# Contributing to GoForge

Thank you for your interest in contributing to GoForge! This document provides guidelines and instructions for contributing.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/goforge.git
   cd goforge
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/goforge/ai-platform.git
   ```

## Development Setup

1. **Install prerequisites**:
   - Go 1.22+
   - Docker
   - Kubernetes (Minikube or similar)
   - Protocol Buffers compiler
   - grpcurl (for testing)

2. **Install Go dependencies**:
   ```bash
   go mod download
   ```

3. **Generate protobuf code**:
   ```bash
   make proto
   ```

4. **Run tests**:
   ```bash
   make test
   ```

## Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards

3. **Add tests** for new functionality

4. **Run tests and linting**:
   ```bash
   make test
   go vet ./...
   ```

5. **Commit your changes**:
   ```bash
   git commit -m "Add feature: description of your changes"
   ```
   
   Follow conventional commits format:
   - `feat: add new feature`
   - `fix: fix bug in component`
   - `docs: update documentation`
   - `test: add tests for feature`
   - `refactor: refactor code`

6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request** on GitHub

## Coding Standards

### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors explicitly, don't ignore them

### Example:

```go
// PredictInput processes the input data and returns a prediction
func (s *Server) PredictInput(ctx context.Context, data []byte) (*Prediction, error) {
    if len(data) == 0 {
        return nil, fmt.Errorf("input data is empty")
    }
    
    // Process data...
    
    return prediction, nil
}
```

### Protocol Buffers

- Use descriptive field names
- Add comments explaining message purpose
- Use appropriate field numbers (1-15 for frequently used fields)

### Kubernetes Manifests

- Include resource requests and limits
- Add appropriate labels and annotations
- Use ConfigMaps for configuration
- Use Secrets for sensitive data

## Testing

### Unit Tests

- Write tests for all new functions
- Use table-driven tests where appropriate
- Aim for >80% code coverage
- Mock external dependencies

Example:
```go
func TestPredict(t *testing.T) {
    tests := []struct {
        name    string
        input   []byte
        want    *Prediction
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   []byte("test data"),
            want:    &Prediction{Value: 0.95},
            wantErr: false,
        },
        // Add more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Predict(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Predict() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Assert results
        })
    }
}
```

### Integration Tests

- Test gRPC service interactions
- Test Kubernetes deployments
- Use test containers where possible

## Pull Request Process

1. **Update documentation** if you changed APIs or added features
2. **Update CHANGELOG.md** with your changes
3. **Ensure all tests pass** in CI
4. **Request review** from maintainers
5. **Address review comments**
6. **Squash commits** if requested
7. **Wait for approval** before merging

### Pull Request Checklist

- [ ] Tests pass locally
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Commit messages are clear
- [ ] No unnecessary files included

## Reporting Bugs

1. **Search existing issues** to avoid duplicates
2. **Create a new issue** with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, K8s version)
   - Logs and error messages

### Bug Report Template

```markdown
**Describe the bug**
A clear description of the bug.

**To Reproduce**
Steps to reproduce:
1. Deploy service X
2. Send request Y
3. See error Z

**Expected behavior**
What you expected to happen.

**Environment**
- OS: [e.g., macOS 14.0]
- Go version: [e.g., 1.22.0]
- Kubernetes: [e.g., 1.29.0]
- GoForge version: [e.g., v1.0.0]

**Logs**
```
Paste relevant logs here
```
```

## Feature Requests

1. **Check roadmap** to see if it's already planned
2. **Create an issue** describing:
   - Use case and motivation
   - Proposed solution
   - Alternatives considered
   - Impact on existing functionality

## Code Review Guidelines

### For Reviewers

- Be respectful and constructive
- Explain reasoning for suggestions
- Approve when code meets standards
- Request changes when necessary

### For Contributors

- Respond to all comments
- Ask questions if unclear
- Don't take feedback personally
- Update PR based on feedback

## Community

- **Slack**: #goforge-dev
- **Mailing List**: dev@goforge.io
- **Monthly Meeting**: First Tuesday of each month

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to ask questions in:
- GitHub Issues
- Slack channel
- Mailing list

Thank you for contributing to GoForge!
