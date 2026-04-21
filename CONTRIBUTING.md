# Contributing to GoZork

## Prerequisites

- Go 1.25+
- golangci-lint v2 (for linting)

## Getting started

```bash
git clone https://github.com/ajdnik/gozork.git
cd gozork
make build
```

## Workflow

### Branching

Branch from `main`. Use this naming convention:

| Type        | Pattern                    | Example                      |
|-------------|----------------------------|------------------------------|
| Feature     | `feat/<short-description>` | `feat/add-save-slots`        |
| Bug fix     | `fix/<short-description>`  | `fix/combat-damage-calc`     |
| Chore/docs  | `chore/<short-description>`| `chore/update-dependencies`  |

```bash
git checkout main
git pull origin main
git checkout -b feat/<short-description>
```

### Making changes

1. Make changes on your branch
2. Run `make check` — must pass before submitting
3. Commit using [Conventional Commits](#commit-style)

### Opening a pull request

Target `main`. Title must follow Conventional Commits format.

```bash
# Push branch
git push -u origin feat/<short-description>

# Open PR with gh CLI
gh pr create --base main --title "feat: <short description>" --body "$(cat <<'EOF'
## Description

<!-- What does this PR do? -->

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Refactoring
- [ ] Documentation
- [ ] CI/Build

## Checklist

- [ ] Tests pass (`make test`)
- [ ] Linter clean (`make lint`)
- [ ] `go vet` clean (`make vet`)
EOF
)"
```

Opening via GitHub UI after pushing will pre-fill the template automatically.

CI runs lint, vet, and tests automatically. All checks must pass before merge. Keep PRs focused — one feature or fix per PR.

## Make targets

| Target       | What it does                                      |
|--------------|---------------------------------------------------|
| `make check` | fmt + vet + lint + test (run this before PRs)     |
| `make build` | Compile the `gozork` binary                       |
| `make run`   | Build and launch the game                         |
| `make test`  | Run all tests (verbose, no cache)                 |
| `make cover` | Tests + per-package coverage summary              |
| `make fuzz`  | Fuzz `FuzzTokenize` for 30 seconds                |
| `make lint`  | golangci-lint v2                                  |
| `make vet`   | go vet static analysis                            |
| `make fmt`   | gofmt all source files                            |
| `make clean` | Remove build artifacts                            |

## Testing

Tests live in `engine/` and `game/`. Two kinds:

- **Unit tests** — test individual engine functions
- **Playthrough tests** — feed scripted command sequences, assert on captured output

When adding game logic, add a playthrough test that exercises the new behavior. Use the existing test harness in `game/game_harness_test.go`.

Engine tests use I/O substitution:
- `G.GameInput` → `strings.Reader` with scripted commands
- `G.GameOutput` → `bytes.Buffer` to capture output
- `G.Rand` → deterministic RNG for reproducible results

Coverage threshold: 80%+ total. Check with `make cover`.

### Fuzz testing

`FuzzTokenize` in `engine/fuzz_test.go` verifies the tokenizer never panics on arbitrary input.

Run for 30 seconds (default):

```bash
make fuzz
```

Run longer or with custom duration:

```bash
go test -fuzz=FuzzTokenize -fuzztime=5m ./engine/
```

If fuzzing finds a failure, Go saves the input to `testdata/fuzz/FuzzTokenize/` automatically. Reproduce it:

```bash
go test -run=FuzzTokenize/testdata/fuzz/FuzzTokenize/<id> ./engine/
```

Commit any new corpus entries in `testdata/fuzz/` that expose real bugs.

## Code structure

```
engine/   # Game-agnostic engine — parser, objects, clock, I/O
game/     # Zork I content — rooms, items, NPCs, verbs, actions
main.go   # Entry point
```

Engine must stay game-agnostic. No Zork-specific logic in `engine/`. Game content is injected at init time through registries.

## Commit style

Conventional Commits:

```
feat: add <thing>
fix: correct <thing>
chore: <maintenance>
test: <test changes>
docs: <docs changes>
```

Subject ≤ 72 characters. No period at end.

## Pull requests

- Target `main`
- CI runs lint, vet, and tests automatically
- All checks must pass
- Keep PRs focused — one feature or fix per PR

## Reporting bugs

Open a [bug report issue](https://github.com/ajdnik/gozork/issues/new?template=bug_report.md). Fill in all sections:

- **Steps to reproduce** — exact command sequence that triggers the bug
- **Expected vs actual behavior** — what should happen vs what does
- **Environment** — Go version (`go version`), OS, GoZork version or commit

Example:

```
## Description
Typing "take all" in a room with no items crashes the game.

## Steps to reproduce
1. Run ./gozork
2. Navigate to an empty room (e.g. "go east" from West of House)
3. Type "take all"

## Expected behavior
Prints "There is nothing here to take."

## Actual behavior
Panic: runtime error: index out of range [0] with length 0

## Environment
- Go version: go1.25.0 darwin/arm64
- OS: macOS 15.0
- GoZork version: v1.2.0
```

## Requesting features

Open a [feature request issue](https://github.com/ajdnik/gozork/issues/new?template=feature_request.md). Fill in all sections:

- **Summary** — one sentence describing the feature
- **Motivation** — problem it solves or improvement it makes
- **Proposed solution** — how it should work
- **Alternatives considered** — other approaches and why they were ruled out

Example:

```
## Summary
Add a --no-color flag to disable ANSI color output.

## Motivation
Terminal emulators that don't support ANSI codes display raw escape sequences,
making output unreadable.

## Proposed solution
Parse a --no-color flag in main.go and pass it to the engine's output writer
to suppress escape sequences.

## Alternatives considered
Detecting $NO_COLOR env var (https://no-color.org) — could do both.
```
