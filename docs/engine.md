# Engine Package

The `engine` package is a reusable, game-agnostic text-adventure engine. It provides everything needed to run a parser-driven interactive fiction game: an object model, a natural-language parser, an action dispatch system, a timed-event clock, save/restore infrastructure, and formatted I/O.

## Files

| File | Responsibility |
|------|---------------|
| `state.go` | `GameState` struct, `NewGameState()`, global `G` pointer, `RNG` interface |
| `object.go` | `Object` struct, flags, directions, facets, object-tree operations |
| `lexer.go` | Vocabulary, tokenizer, lexer, `LexItm` and `WordItm` types |
| `parser.go` | Full natural-language parser: `Parse()`, clause handling, orphan merging, object resolution |
| `syntax.go` | Syntax definitions (`Syntx`), vocabulary builder, `Commands` registry |
| `game.go` | `MainLoop()`, `Perform()`, `ResetGameState()`, `PerfRet` result type |
| `clock.go` | `ClockEvent`, `Queue()`, `Clocker()` — timed events and daemons |
| `helpers.go` | `IsInGlobal()`, `IsHeld()` |
| `utils.go` | `PickOne()`, `Random()`, `Prob()`, `IsFlaming()`, `IsOpenable()` |
| `output.go` | `Printf()` — formatted output to `G.GameOutput` |
| `save.go` | `ObjIndex`, `ObjToIdx()`, `IdxToObj()` — object serialization helpers |

## Game State (`state.go`)

All mutable state is centralized in `GameState`, accessed through the package-level pointer `G`:

```go
var G *GameState
```

Key field groups:

- **Core references** — `Here` (current room), `Winner` (acting character), `Player`.
- **Current command** — `DirObj`, `IndirObj`, `ActVerb`, resolved object lists, detected syntax.
- **Progress** — `Moves`, `Score`, `BaseScore`, `Lit`.
- **Settings** — `SuperBrief`, `Verbose`, `Lucky`.
- **Parser internals** — `ParsedSyntx`, `OrphanedSyntx`, `Params`, `LexRes`, continuation state.
- **Clock queue** — Fixed-size array of 30 `ClockEvent` slots plus daemon/interrupt boundary indices.
- **I/O** — `GameOutput` (writer), `GameInput` (reader), `Reader` (buffered). Defaults to stdout/stdin; tests swap these for buffers.
- **RNG** — `Rand` field implements the `RNG` interface (single method `Intn(n)`). Defaults to a time-seeded source; tests inject a deterministic one.
- **Well-known objects** — Pointers to sentinel objects (`RoomsObj`, `GlobalObj`, `NotHereObj`, `PseudoObj`, `ItPronounObj`, `MeObj`, `HandsObj`) set by the game during init.
- **Registries** — `Actions`, `PreActions`, `NormVerbs` maps (verb string → handler), `ClockFuncs` (key → tick function).
- **Game extension** — `GameData interface{}` holds a game-specific struct (e.g. `*ZorkData`).

### State Lifecycle

1. `NewGameState()` creates a fresh state with defaults.
2. `ResetGameState()` preserves I/O, RNG, and registered data, then rebuilds the state and object tree from snapshots. Used by restart and tests.

## Object Model (`object.go`)

Everything in the game world — rooms, items, NPCs, containers, vehicles, the player — is an `Object`.

### Core Fields

```go
type Object struct {
    Flags      Flags        // bitfield properties (FlgTake, FlgOpen, FlgLight, ...)
    In         *Object      // parent in the containment tree
    Children   []*Object    // child objects
    synonyms   []string     // nouns the parser matches
    Adjectives []string     // adjectives for disambiguation
    Desc       string       // short description ("brass lantern")
    LongDesc   string       // room/first-visit description
    FirstDesc  string       // description before first interaction
    Text       string       // readable text

    Action     Action       // handler called during Perform
    ContFcn    Action       // container function (intercepts child actions)
    DescFcn    Action       // custom description function

    Global     []*Object    // objects globally visible from this room
    Pseudo     []PseudoObj  // pseudo-objects (synonym + action, no real object)
    Exits      map[Direction]DirProps  // room exits
}
```

### Optional Facets

Role-specific data is stored in nil-able facet pointers to keep the core struct lean:

| Facet | Fields | Used By |
|-------|--------|---------|
| `*ItemData` | `Size`, `Value`, `TValue`, `Capacity` | Takeable items, containers, treasures |
| `*CombatData` | `Strength` | NPCs that participate in combat |
| `*VehicleData` | `Type` (flags like `FlgNonLand`) | Rideable vehicles (boat) |

Nil-safe getters/setters (e.g. `GetSize()`, `SetStrength()`) allocate the facet on first write.

### Containment Tree

Objects form a tree via `In` (parent) and `Children` (kids):

```
rooms
  ├─ westOfHouse
  │    ├─ mailbox
  │    │    └─ Leaflet
  │    └─ adventurer (player)
  │         ├─ sword
  │         └─ lamp
  └─ kitchen
       └─ kitchenTable
            ├─ bottle
            └─ Sack
globalObjects
  ├─ it (pronoun sentinel)
  ├─ hands
  ├─ me
  └─ grue
```

Key tree operations:

- `MoveTo(dest)` — detaches from current parent, attaches to dest.
- `Remove()` — detaches from parent, sets `In` to nil.
- `IsIn(loc)` — checks if direct parent is `loc`.
- `Location()` — returns `In` (the parent).
- `Has(flag)` / `Give(flag)` / `Take(flag)` — bitfield operations on `Flags`.

### Object Snapshots

`BuildObjectTree()` records the initial state of every object (parent, flags, combat strength, values, text). `ResetObjectTree()` restores objects to that snapshot, enabling clean test resets and game restart.

### Flags

Flags is a `uint64` bitfield with ~45 named constants covering:

- **Item properties**: `FlgTake`, `FlgTryTake`, `FlgCont`, `FlgOpen`, `FlgLock`, `FlgSurf`, `FlgTrans`
- **Light/fire**: `FlgLight`, `FlgOn`, `FlgFlame`, `FlgBurn`
- **Description**: `FlgNoDesc`, `FlgInvis`, `FlgTouch`, `FlgRead`
- **Character**: `FlgPerson`, `FlgFemale`, `FlgActor`
- **Grammar**: `FlgVowel`, `FlgNoArt`, `FlgPlural`, `FlgKludge`
- **Room type**: `FlgLand`, `FlgWater`, `FlgAir`, `FlgOut`, `FlgNonLand`, `FlgMaze`
- **Combat**: `FlgFight`, `FlgStagg`, `FlgWeapon`
- **Special**: `FlgSacred`, `FlgTool`, `FlgIntegral`, `FlgBodyPart`, `FlgNotAll`, `FlgDrop`, `FlgIn`, `FlgClimb`, `FlgDrink`, `FlgFood`, `FlgTurn`, `FlgRMung`, `FlgRLand`, `FlgWear`, `FlgWorn`, `FlgVeh`, `FlgSearch`, `FlgDoor`

### Directions and Exits

13 directions are defined (`North` through `Land`). Room exits are stored in `Exits map[Direction]DirProps`:

```go
type DirProps struct {
    NExit    string   // message for blocked exit
    UExit    bool     // unconditional exit flag
    RExit    *Object  // destination room
    FExit    FDir     // function-based exit (computed destination)
    CExit    CDir     // conditional check (returns bool)
    CExitStr string   // message when condition fails
    DExit    *Object  // door object that must be open
    DExitStr string   // message when door is closed
}
```

## Lexer (`lexer.go`)

The lexer converts raw player input into tagged tokens.

### Pipeline

```
Raw input string
  → Tokenize()    Split on whitespace, separate letters/digits/punctuation
  → Lex()         Look up each token in Vocabulary, tag with WordTypes
  → []LexItm      Tagged token list for the parser
```

### Word Types

| Type | Meaning | Examples |
|------|---------|---------|
| `WordDir` | Compass direction | north, south, up, down |
| `WordVerb` | Action verb | take, open, look, attack |
| `WordPrep` | Preposition | in, on, with, from |
| `WordAdj` | Adjective | brass, large, rusty |
| `WordObj` | Noun / object name | lamp, sword, door |
| `WordBuzz` | Filler (ignored) | the, a, an |

A single word can have multiple types (e.g. "light" can be both verb and adjective). The parser uses context to disambiguate.

### Vocabulary

`Vocabulary` is a `map[string]WordItm` built at init time by `BuildVocabulary()`. It merges:

1. Buzz words (articles, fillers)
2. Verbs and prepositions from syntax definitions
3. Directions
4. Object synonyms and adjectives
5. Word synonyms (e.g. "n" → "north")

## Parser (`parser.go`)

The parser is the most complex engine subsystem. It interprets player input and resolves it into a verb, direct object, and indirect object.

### Parse Flow

```
Parse()
  ├─ Handle OOPS correction
  ├─ Handle AGAIN / G repetition
  ├─ Lex loop: classify each word
  │    ├─ Direction? → set walk direction, done
  │    ├─ Verb? → record verb
  │    ├─ Prep/Adj/Obj? → Clause() — parse object clause
  │    ├─ "then" / "." → set continuation index, break
  │    └─ Unknown? → UnknownWord() error
  ├─ OrphanMerge() — merge with incomplete previous command
  ├─ SyntaxCheck() — match against Commands registry
  ├─ SnarfObjects() — resolve token clauses to Object lists
  ├─ ManyCheck() — reject multiple objects when syntax disallows
  └─ TakeCheck() — implicit-take check
```

### Disambiguation

When the parser matches multiple objects for a noun (e.g. "take key" when there are two keys), it:

1. Asks "Which key do you mean, the brass key or the iron key?"
2. Sets `ShldOrphan = true` and saves state in `OrphanedSyntx`
3. On the next input, `OrphanMerge()` combines the player's answer with the orphaned command

### Multi-command Input

Commands separated by "then", ".", or "," are handled via the continuation system:

- `G.Params.Continue` stores the index into `LexRes` where the next command begins.
- `G.Reserv` buffers commands separated by "and" or "," for sequential execution.

### Object Resolution

`SnarfObjects()` resolves each object clause:

1. `Snarfem()` walks the token list, handling "all", "but/except", "and", adjective-noun pairs
2. `GetObject()` searches the object tree using `SearchList()` and `GlobalCheck()`
3. Search scope is controlled by `LocFlags` from the syntax definition (held, carried, in-room, on-ground, etc.)
4. `IsThisIt()` checks if an object matches the current search criteria (synonym, adjective, flags)

## Syntax System (`syntax.go`)

Each recognized command pattern is a `Syntx`:

```go
type Syntx struct {
    NormVerb  string     // normalized verb for action lookup
    Verb      string     // primary verb word
    VrbPrep   string     // verb preposition ("look" + "at")
    Obj1      ObjProp    // first object slot constraints
    ObjPrep   string     // object preposition ("put X" + "in" + Y)
    Obj2      ObjProp    // second object slot constraints
    Action    VrbAction  // default verb handler
    PreAction VrbAction  // pre-action check
}
```

`ObjProp` specifies what kind of object a slot expects:

- `ObjFlags` — required object flags (e.g. `FlgPerson` for "attack")
- `LocFlags` — where to search (held, in room, take implicitly, allow many)
- `HasObj` — whether this slot is used at all

`SyntaxCheck()` iterates `Commands` looking for a syntax that matches the parsed verb, preposition(s), and object count. If an object is missing but a syntax expects one, it calls `FindWhatIMean()` to infer from context (printing e.g. "(with the key)").

## Action Dispatch (`game.go`)

### Perform

`Perform(verb, directObj, indirectObj)` dispatches an action through the handler chain (see Architecture doc). It:

1. Saves and restores `ActVerb`, `DirObj`, `IndirObj` (allowing nested `Perform` calls)
2. Resolves "it" pronoun references
3. Walks the handler chain in priority order
4. Returns `PerfRet`: `PerfNotHndld`, `PerfHndld`, `PerfFatal`, or `PerfQuit`

### MainLoop

`MainLoop()` runs the parse-perform-clock cycle until `InputExhausted` or `QuitRequested`. For multi-object commands ("take all"), it iterates over `DirObjPossibles`, calling `Perform` for each and printing the object name as a prefix.

## Clock System (`clock.go`)

The clock provides timed events (interrupts) and always-ticking daemons.

### ClockEvent

```go
type ClockEvent struct {
    Key  string       // unique identifier (e.g. "iLantern")
    Run  bool         // whether active
    Tick int          // countdown; fires at 0 (-1 = daemon, always fires)
    Fn   func() bool  // handler
}
```

### Queue Management

- `Queue(key, tick)` — find or allocate a slot, set its tick value.
- `QueueInt(key, dmn)` — find or allocate; if `dmn`, the event ticks even on bad parses.
- Slots are allocated from the end of a 30-element fixed array, growing downward. `QueueInts` and `QueueDmns` track the boundary.

### Clocker

`Clocker()` runs once per turn after `Perform`:

1. Skips if `ClockWait` is set (one-turn delay).
2. Iterates active events. If `ParserOk` is false, only daemons run.
3. Decrements `Tick`. When `Tick` reaches 0, calls `Fn()`.
4. Increments `G.Moves`.

## I/O (`output.go`, `state.go`)

- `Printf(format, args...)` writes to `G.GameOutput` (default: `os.Stdout`).
- `Read()` reads a line from `G.GameInput` via a buffered reader, lowercases it, tokenizes, and lexes.
- Tests replace `GameOutput` with `bytes.Buffer` and `GameInput` with `strings.Reader`.

## Save Infrastructure (`save.go`)

The engine provides object-index mapping for serialization:

- `BuildObjIndex()` builds a `map[*Object]int` from `G.AllObjects`.
- `ObjToIdx(o)` / `IdxToObj(idx)` convert between pointers and stable indices.

The actual save/restore/restart logic is provided by the game package via function hooks (`G.Save`, `G.Restore`, `G.Restart`). Each hook returns `error` (`nil` on success) so that callers can report the specific failure reason to the player.
