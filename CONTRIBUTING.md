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

1. Fork the repo and create a branch from `main`
2. Make changes
3. Run `make check` — must pass before submitting
4. Open a pull request targeting `main`

## Make targets

| Target       | What it does                                      |
|--------------|---------------------------------------------------|
| `make check` | fmt + vet + lint + test (run this before PRs)     |
| `make test`  | Run all tests (verbose, no cache)                 |
| `make cover` | Tests + per-package coverage summary              |
| `make lint`  | golangci-lint v2                                  |
| `make vet`   | go vet static analysis                            |
| `make fmt`   | gofmt all source files                            |

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

Open a GitHub issue with:
- Go version (`go version`)
- Steps to reproduce
- Expected vs actual output
