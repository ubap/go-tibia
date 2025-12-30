## 🛠 Development Guidelines

### Commit Message Convention

This project follows the [Conventional Commits](https://www.conventionalcommits.org/) specification. Structured commit messages make the project history readable and allow for automated changelog generation.

#### Format
`type(scope): description`

- **type**: See the list below for allowed types.
- **scope** (optional): The specific part of the codebase (e.g., `api`, `db`, `internal/storage`).
- **description**: A short, present-tense summary of the change.

#### Allowed Types
| Type | Description |
| :--- | :--- |
| `feat` | A new feature for the user |
| `fix` | A bug fix |
| `docs` | Documentation changes only |
| `refactor` | A code change that neither fixes a bug nor adds a feature |
| `perf` | A code change that improves performance |
| `test` | Adding missing tests or correcting existing tests |
| `chore` | Changes to the build process, auxiliary tools, or libraries (e.g., updating dependencies) |
| `ci` | Changes to CI configuration files and scripts (e.g., GitHub Actions) |

#### Examples

- **New Feature:** `feat(auth): add JWT middleware for user sessions`
- **Bug Fix:** `fix(server): resolve memory leak in websocket handler`
- **Refactor:** `refactor(internal): simplify interface for storage provider`
- **Docs:** `docs: update installation instructions in README`

> [!TIP]
> **Breaking Changes:** If a change breaks backward compatibility, append a `!` after the type, e.g., `feat!: change default database driver to Postgres`.