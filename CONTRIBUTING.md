# Contributing to JarCraft

First of all, thank you for considering contributing to JarCraft! Every contribution is welcome, whether it is a bug fix, a new feature, documentation improvement, or suggestion.

## Getting Started

Before contributing, make sure you have the following installed:

- [Go](https://go.dev/dl/) (version specified in `go.mod`).
- Git.
- Java JDK (required to test Java project builds).
- Maven or Gradle (optional if the project uses a wrapper).

### 1. Fork the repository

Fork the repository to your GitHub account and clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/jarcraft.git
cd jarcraft
```

### 2. Create a branch

Create a descriptive branch for your contribution:

```bash
git checkout -b feat/your-feature
```

Suggested branch prefixes:

- `feat/`: New features.
- `fix/`: Bug fixes.
- `docs/`: Documentation changes.
- `test/`: Test improvements.
- `refactor/`: Code refactoring.
- `chore/`: Maintenance tasks.

### 3. Run the tests

Before making changes, run the existing test suite:

```bash
go test ./...
```

You can also run:

```bash
go vet ./...
```

Make sure all tests pass before submitting your changes.

## Development Guidelines

Please follow these guidelines when contributing:

- Keep the code simple, readable, and maintainable.
- Follow standard Go conventions and formatting.
- Use `gofmt` to format Go files.
- Add or update tests for every behavior change.
- Avoid introducing unnecessary dependencies.
- Keep configuration backward-compatible whenever possible.
- Handle errors explicitly instead of silently ignoring them.
- Do not commit generated binaries, build artifacts, or local configuration files containing machine-specific paths.

Format your code before submitting:

```bash
gofmt -w .
```

## Commit Messages

JarCraft follows the **Conventional Commits** specification.

Use the following format:

```text
<type>(optional scope): <description>
```

Common types:

| Type | Description |
|---|---|
| `feat` | Introduces a new feature |
| `fix` | Fixes a bug |
| `docs` | Documentation changes |
| `test` | Adds or updates tests |
| `refactor` | Code changes without changing behavior |
| `chore` | Maintenance tasks |
| `ci` | CI/CD changes |

Examples:

```text
feat(builder): add support for custom Java versions
fix(scanner): prevent duplicate project detection
docs: improve installation instructions
test(engine): add tests for parallel builds
ci: update GitHub Actions workflow
```

## Pull Request Process

1. Make sure your branch is up to date with the main branch.
2. Run all tests and `go vet`.
3. Commit your changes using Conventional Commits.
4. Push your branch to your fork.
5. Open a Pull Request against the original repository's `main` branch.
6. Describe the problem your contribution solves and the changes introduced.
7. Include relevant test results and screenshots or logs when applicable.

### Pull Request Checklist

Before submitting your PR, verify that:

- [ ] The code follows the project's Go conventions.
- [ ] All tests pass (`go test ./...`).
- [ ] `go vet ./...` completes successfully.
- [ ] New functionality includes tests.
- [ ] Documentation has been updated when necessary.
- [ ] No generated binaries or machine-specific files are included.
- [ ] The PR description explains the changes clearly.

## Reporting Bugs

If you find a bug, please open a GitHub Issue and include:

- A clear description of the problem.
- Steps to reproduce it.
- Expected behavior.
- Actual behavior.
- Your operating system.
- Go version.
- Relevant logs or error messages.

## Suggesting Features

Feature suggestions are welcome!

When proposing a feature, explain:

- The problem it solves.
- Why it would be useful.
- How you imagine it working.
- Any relevant examples or use cases.

## Questions and Discussions

If you have questions, ideas, or need clarification about the project, feel free to open a GitHub Discussion or an Issue.

Thank you for helping make JarCraft better!

---

**Happy coding!** 🚀