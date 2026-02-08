# Architecture Overview

GoZork is a faithful Go port of Zork I, the classic interactive fiction game originally written in ZIL (Zork Implementation Language) for the Z-machine. The codebase is split into two packages that mirror the original ZIL separation between the interpreter (engine) and the story file (game).

## Package Layout

```
main.go          Entry point — calls game.Run()
engine/          Reusable text-adventure engine (generic, game-agnostic)
game/            Zork I game data and logic (depends on engine)
```

## Design Principles

- **Single mutable state pointer.** All mutable state lives in `engine.GameState` (accessed via the global `engine.G`). This makes save/restore, restart, and test isolation straightforward — swap the pointer, reset the object tree, and you have a fresh game.
- **Dot-import by design.** The `game` package dot-imports `engine` so that game code reads like the original ZIL — `G.Here`, `Printf(...)`, `Perform(...)` — with no package prefix.
- **Engine is game-agnostic.** The engine knows nothing about Zork-specific objects, rooms, or verbs. It provides the parser, object model, action dispatch, clock system, and I/O. All game content is injected at init time through registries (`Commands`, `Vocabulary`, `G.AllObjects`, etc.).
- **Game-specific state via interface field.** `GameState.GameData` is an `interface{}` that the game package fills with its own `*ZorkData` struct. The engine never touches it; the game accesses it through `GD()`.

## Startup Sequence

```
main.go
  └─ game.Run()
       └─ game.InitGame()
            ├─ engine.NewGameState()         // fresh engine state
            ├─ game.NewZorkData()            // fresh Zork-specific state
            ├─ registerWellKnownObjects()    // wire sentinel objects into engine
            ├─ engine.ResetGameState()       // snapshot + rebuild
            ├─ initClockFuncs()              // register timed event handlers
            ├─ FinalizeGameObjects()         // link actions, set item data, etc.
            ├─ engine.BuildObjectTree()      // populate Children from In pointers
            ├─ engine.BuildVocabulary(...)   // build word → type map + action map
            ├─ engine.InitReader()           // set up buffered stdin reader
            ├─ initSaveSystem()             // wire save/restore/restart hooks
            └─ (initial clock events, starting room, player placement)
       └─ engine.MainLoop()                  // parse → perform → clock cycle
```

## Main Loop

Each turn follows this cycle:

```
┌──────────────────────────────────────────┐
│  1. Parse()                              │
│     Read input → Tokenize → Lex → Clause │
│     → SyntaxCheck → SnarfObjects         │
│     → ManyCheck → TakeCheck              │
├──────────────────────────────────────────┤
│  2. Perform(verb, directObj, indirectObj)│
│     Walk the handler chain:              │
│       winner → room(M-BEG) → preAction   │
│       → indirectObj → container          │
│       → directObj → verbAction           │
├──────────────────────────────────────────┤
│  3. Room end handler (M-END)             │
├──────────────────────────────────────────┤
│  4. Clocker()                            │
│     Tick all active clock events         │
│     Increment move counter               │
└──────────────────────────────────────────┘
```

## Handler Chain (Perform)

`Perform` dispatches an action through a priority-ordered chain. The first handler to return `true` claims the action and stops the chain:

| Priority | Handler | Source |
|----------|---------|--------|
| 1 | Not-here object | `G.NotHereObj.Action` — if either object is the not-here sentinel |
| 2 | Winner (actor) | `G.Winner.Action` — lets the actor intercept any action |
| 3 | Room begin | `G.Here.Action(ActBegin)` — room-level interception |
| 4 | Pre-action | `G.PreActions[verb]` — verb-specific pre-check |
| 5 | Indirect object | `G.IndirObj.Action` — the indirect object reacts |
| 6 | Container | `G.DirObj.Location().ContFcn` — the container reacts |
| 7 | Direct object | `G.DirObj.Action` — the direct object reacts |
| 8 | Verb action | `G.Actions[verb]` — the default verb handler |

After the chain completes, the room's end handler (`ActEnd`) runs, followed by the clock system.

## Data Flow Between Packages

```
engine (generic)                    game (Zork I)
─────────────────                   ──────────────
GameState.GameData  ←── filled ───  *ZorkData
G.AllObjects        ←── filled ───  Objects slice (items.go)
G.RoomsObj          ←── filled ───  &Rooms (globals.go)
G.Actions map       ←── built ────  GameCommands (syntax_data.go)
G.PreActions map    ←── built ────  GameCommands (syntax_data.go)
Vocabulary map      ←── built ────  synonyms, objects, adjectives
Commands slice      ←── set ──────  GameCommands (syntax_data.go)
G.ClockFuncs map    ←── filled ───  initClockFuncs() (init.go)
G.Save/Restore/     ←── wired ───  doSave/doRestore/doRestart (save.go)
  Restart hooks                     (return error, nil on success)
```

## Testing Strategy

Tests live in `game/` and use the engine's I/O abstraction:

- `G.GameInput` is replaced with a `strings.Reader` containing scripted commands.
- `G.GameOutput` is replaced with a `bytes.Buffer` to capture output.
- `G.Rand` is replaced with a deterministic RNG for reproducible combat and random events.
- `ResetGameState()` + `ResetObjectTree()` restore the world to its initial state between tests.

Playthrough tests feed a full sequence of commands and assert on the captured output, ensuring the game logic matches the original Zork I behavior.
