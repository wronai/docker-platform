# 👋 Welcome Contributors!

Thank you for your interest in contributing to Media Vault! We're excited to have you on board. This guide will help you get started with contributing to our project.

## 📋 Table of Contents

- [Code of Conduct](#-code-of-conduct)
- [Getting Started](#-getting-started)
- [Development Workflow](#-development-workflow)
- [Code Style](#-code-style)
- [Testing](#-testing)
- [Pull Request Process](#-pull-request-process)
- [Reporting Bugs](#-reporting-bugs)
- [Feature Requests](#-feature-requests)
- [Documentation](#-documentation)
- [Community](#-community)

## 🤝 Code of Conduct

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing. We are committed to fostering a welcoming and inclusive community.

## 🚀 Getting Started

### Prerequisites

- Git
- Docker 20.10+ and Docker Compose
- Node.js 16+ (for frontend development)
- Go 1.19+ (for backend development)
- Make (optional but recommended)

### Setting Up the Development Environment

1. **Fork the repository**
   ```bash
   git clone https://github.com/wronai/docker-platform.git
   cd docker-platform
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start the development environment**
   ```bash
   make dev
   ```

4. **Access the applications**
   - Web UI: http://localhost:3000
   - API: http://localhost:8080
   - Documentation: http://localhost:8080/docs
   - Monitoring: http://localhost:9090

## 🔄 Development Workflow

1. **Create a new branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

2. **Make your changes**
   - Follow the code style guidelines
   - Write tests for new features
   - Update documentation as needed

3. **Run tests**
   ```bash
   make test
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

5. **Push your changes**
   ```bash
   git push origin your-branch-name
   ```

6. **Create a Pull Request**
   - Go to the [repository](https://github.com/wronai/docker-platform)
   - Click "New Pull Request"
   - Follow the PR template
   - Request reviews from maintainers

## 🎨 Code Style

### Backend (Go)
- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for formatting
- Run `golint` and `go vet` before committing
- Write unit tests for new functionality

### Frontend (Flutter/Dart)
- Follow the [Dart Style Guide](https://dart.dev/guides/language/effective-dart/style)
- Use `dart format` for formatting
- Follow BLoC pattern for state management
- Write widget tests for UI components

### Git Commit Messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code changes that neither fixes a bug nor adds a feature
- `perf`: Performance improvements
- `test`: Adding tests
- `chore`: Changes to the build process or auxiliary tools

Example:
```
feat(auth): add two-factor authentication

- Add TOTP support
- Add recovery codes
- Update login flow

Closes #123
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run backend tests
make test-backend

# Run frontend tests
make test-frontend

# Run integration tests
make test-integration

# Run linters
make lint
```

### Writing Tests
- Write unit tests for all new functionality
- Include integration tests for critical paths
- Use descriptive test names
- Test edge cases and error conditions

## 🔄 Pull Request Process

1. Ensure your code passes all tests
2. Update documentation if needed
3. Ensure your branch is up to date with `main`
4. Run linters and fix any issues
5. Request reviews from at least one maintainer
6. Address all review comments
7. Wait for CI to pass
8. Squash and merge when approved

## 🐛 Reporting Bugs

Found a bug? Please let us know by [opening an issue](https://github.com/wronai/docker-platform/issues/new?template=bug_report.md).

### Before Submitting a Bug Report
1. Check if the issue has already been reported
2. Try to reproduce the issue with the latest version
3. Collect as much information as possible

### Bug Report Template
```markdown
## Describe the bug
A clear and concise description of what the bug is.

## To Reproduce
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

## Expected behavior
A clear and concise description of what you expected to happen.

## Screenshots
If applicable, add screenshots to help explain your problem.

## Additional context
Add any other context about the problem here.
```

## 💡 Feature Requests

Have an idea for a new feature? [Open a feature request](https://github.com/wronai/docker-platform/issues/new?template=feature_request.md).

### Feature Request Template
```markdown
## Is your feature request related to a problem? Please describe.
A clear and concise description of what the problem is.

## Describe the solution you'd like
A clear and concise description of what you want to happen.

## Describe alternatives you've considered
A clear and concise description of any alternative solutions or features you've considered.

## Additional context
Add any other context or screenshots about the feature request here.
```

## 📚 Documentation

Good documentation is crucial for the success of any project. We appreciate contributions that improve our documentation.

### Documentation Guidelines
- Keep it clear and concise
- Use proper formatting
- Include examples where helpful
- Keep it up to date

### Building Documentation
```bash
# Build documentation
make docs

# Serve documentation locally
make docs-serve
```

## 🌍 Community

### Getting Help
- [GitHub Discussions](https://github.com/wronai/docker-platform/discussions)
- [Discord Server](https://discord.gg/your-invite-link)
- [Community Forum](https://community.mediavault.example.com)

### Stay Updated
- Follow us on [Twitter](https://twitter.com/yourhandle)
- Subscribe to our [blog](https://blog.mediavault.example.com)
- Join our [newsletter](https://mediavault.example.com/newsletter)

## 🙏 Acknowledgments

- Thank you to all our contributors
- Special thanks to our core maintainers
- Shoutout to all our beta testers

---

<div align="center">
  <p>Happy coding! 🚀</p>
  <p>If you have any questions, feel free to ask in our <a href="https://github.com/wronai/docker-platform/discussions">Discussions</a>.</p>
</div>
