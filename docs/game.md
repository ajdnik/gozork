# Game Package

The `game` package implements Zork I on top of the generic engine. It defines every room, item, NPC, verb handler, timed event, and puzzle in the game. The package dot-imports `engine` so game code reads like the original ZIL source.

## Files

| File | Responsibility |
|------|---------------|
| `zork_data.go` | `ZorkData` struct (game-specific mutable state), `GD()` accessor |
| `init.go` | `InitGame()`, `Run()`, clock function registration, well-known object wiring |
| `globals.go` | Global objects (It, Me, Hands, Adventurer, etc.), sentinel objects, NPC action funcs |
| `items.go` | All game objects — rooms, items, NPCs — declared as package-level `Object` vars |
| `rooms_surface.go` | Surface world rooms (white house, forest, clearing, etc.) |
| `rooms_underground.go` | Underground rooms (cellar, troll room, temple, mine, etc.) |
| `rooms_maze.go` | Maze rooms and dead ends |
| `dungeon.go` | Lookup tables, random-selection pools, string data, navigation maps |
| `syntax_data.go` | `GameCommands` (all syntax definitions), `BuzzWords`, `Synonyms` map |
| `verbs.go` | Default verb handlers (VTake, VDrop, VOpen, VLook, etc.) |
| `verbs_movement.go` | Movement verbs (VWalk, VBoard, VClimb, VDisembark, VLaunch, etc.) |
| `verbs_meta.go` | Meta verbs (VScore, VInventory, VSave, VRestore, VRestart, VVersion, etc.) |
| `actions.go` | Object action functions, room action functions, puzzle logic, interrupt routines |
| `combat.go` | Combat system: melee tables, blow resolution, villain/hero attack logic |
| `thief.go` | Thief NPC: stealing, treasure room, IThief daemon |
| `river.go` | River/boat system, dam/reservoir mechanics, water-level interrupts |
| `save.go` | Save/restore/restart implementation using `encoding/gob` |
| `game_test.go` | Test helpers and playthrough infrastructure |
| `playthrough_test.go` | Partial playthrough tests |
| `fullplaythrough_test.go` | Full game completion test |

## Game-Specific State (`zork_data.go`)

`ZorkData` holds all mutable state that is specific to Zork I and not part of the generic engine:

```go
type ZorkData struct {
    Dead, WonGame       bool
    Deaths              int
    HelloSailor         int
    TrollFlag           bool   // troll defeated
    CyclopsFlag         bool   // cyclops asleep
    MagicFlag           bool   // cyclops wall opened
    LowTide             bool   // reservoir drained
    RainbowFlag         bool   // rainbow solidified
    LLDFlag             bool   // spirits banished
    DomeFlag            bool   // rope tied to railing
    GrateRevealed       bool   // grate found under leaves
    RugMoved            bool   // rug moved, trap door visible
    MirrorMung          bool   // mirror broken (bad luck!)
    GatesOpen           bool   // dam sluice gates open
    EggSolve            bool   // thief opened the egg
    LoudFlag            bool   // echo room quieted
    CycloWrath          int    // cyclops anger level
    WaterLevel          int    // maintenance room flood level
    MatchCount          int    // matches remaining
    Villains            []*VillainEntry  // combat-registered NPCs
    // ... ~40 more flags and counters
}
```

Accessed everywhere via `GD()`:

```go
func GD() *ZorkData {
    return G.GameData.(*ZorkData)
}
```

## Initialization (`init.go`)

`InitGame()` wires the game into the engine:

1. Creates `GameState` and `ZorkData`
2. Registers well-known objects (`RoomsObj`, `GlobalObj`, `NotHereObj`, `MeObj`, etc.)
3. Sets the implicit-take handler (`ITake`)
4. Calls `FinalizeGameObjects()` to attach action functions and set item properties
5. Builds the object tree and vocabulary
6. Schedules initial clock events (thief daemon, combat daemon, lantern timer, candle timer)
7. Places the player at West of House

`Run()` calls `InitGame()`, shows the version banner and initial room description, then enters `MainLoop()`.

## World Model

### Rooms

Rooms are `Object` values with `Exits` maps and `Action` functions. They are children of the `Rooms` sentinel object and are defined across three files by region:

- **Surface** (`rooms_surface.go`): White house perimeter, forest, clearing, canyon view, up a tree
- **Underground** (`rooms_underground.go`): Cellar, troll room, maze areas, temple, mine, treasure room, dam, reservoir, loud room, cyclops room, and many more
- **Maze** (`rooms_maze.go`): The twisty maze passages, dead ends

Each room's `Action` function handles `ActLook` (room description), `ActBegin` (before command), `ActEnter` (on entry), and `ActEnd` (after command).

### Items

Items are `Object` values with `Synonyms`, `Adjectives`, `Desc`, flags, and optional `Action` functions. Items are defined in `items.go` and wired to their starting locations via `In` pointers.

Key item categories:
- **Treasures**: Egg, canary, painting, chalice, pot of gold, jeweled scarab, trunk of jewels, etc. Each has a `TValue` that contributes to score when placed in the trophy case.
- **Tools**: Lamp, sword, knife, shovel, screwdriver, wrench, keys, matches, rope
- **Containers**: Mailbox, trophy case, machine, bottle, bag, basket, coffin
- **Vehicles**: Inflatable/inflated/punctured boat

### NPCs

NPCs are objects with `FlgPerson` or `FlgActor` and an `Action` function that handles multiple modes:

- **Normal verb handling**: responds to examine, talk, give, etc.
- **Combat callbacks**: `ActBusy`, `ActDead`, `ActUnconscious`, `ActConscious`, `ActFirst` — called by the combat system

Major NPCs: Troll, Thief, Cyclops, Bat, Ghosts

## Verb System

### Syntax Definitions (`syntax_data.go`)

Every recognizable command is a `Syntx` entry in `GameCommands`. Examples:

| Input Pattern | Verb | VrbPrep | Obj1 | ObjPrep | Obj2 |
|--------------|------|---------|------|---------|------|
| `take X` | take | — | takeable | — | — |
| `put X in Y` | put | — | held | in | container |
| `attack X with Y` | attack | — | person | with | weapon |
| `look` | look | — | — | — | — |
| `turn X with Y` | turn | — | turnable | with | tool |

### Verb Handlers (`verbs.go`)

Default handlers for ~80 verbs. Each is a `func(ActArg) bool`. Key verbs:

| Handler | Verbs | Behavior |
|---------|-------|----------|
| `VTake` | take, get | Pick up object, check weight/capacity |
| `VDrop` | drop | Place in room or vehicle |
| `VOpen` / `VClose` | open, close | Toggle FlgOpen, handle locked doors |
| `VLook` / `VFirstLook` | look | Describe room, list contents |
| `VExamine` | examine, x | Describe object in detail |
| `VRead` | read | Show object's Text field |
| `VAttack` | attack, kill | Initiate combat via `HeroBlow()` |
| `VGive` | give | Transfer object to NPC |
| `VThrow` | throw | Throw at target |
| `VInventory` | inventory, i | List carried items |
| `VWalk` | walk, go | Move through exits |
| `VLampOn` / `VLampOff` | turn on/off | Toggle FlgOn, update lighting |
| `VEat` / `VDrink` | eat, drink | Consume food/water |

### Movement (`verbs_movement.go`)

`DoWalk(dir)` is the core movement function:

1. Gets the exit for the direction from `G.Here.Exits`
2. Evaluates exit type (unconditional, conditional, door, function-based)
3. Calls `Goto(destination, true)` if passage is allowed
4. `Goto()` moves the player, updates `G.Here`, checks lighting, and optionally shows the room

Special movement: climbing, boarding/disembarking vehicles, launching the boat, entering/exiting through doors and windows.

### Meta Verbs (`verbs_meta.go`)

| Handler | Verb | Effect |
|---------|------|--------|
| `VScore` | score | Display current score and rank |
| `VInventory` | inventory | List carried objects |
| `VBrief` / `VVerbose` / `VSuperBrief` | brief/verbose/superbrief | Set description verbosity |
| `VSave` / `VRestore` | save/restore | Serialize/deserialize game state |
| `VRestart` | restart | Reset to initial state |
| `VVersion` | version | Show game version |
| `VQuit` | quit | End the game |
| `VDiagnose` | diagnose | Show health status |

## Combat System (`combat.go`)

The combat system faithfully ports the original ZIL melee mechanics.

### Villain Registry

Each combat-capable NPC is registered as a `VillainEntry`:

```go
type VillainEntry struct {
    Villain *Object      // the NPC object
    Best    *Object      // best weapon against this villain
    BestAdv int          // advantage bonus when using best weapon
    Prob    int          // probability of waking from unconsciousness
    Msgs    *MeleeTable  // melee message table
}
```

Three villains: **Troll** (best weapon: sword, +2), **Thief** (best weapon: knife, +1), **Cyclops** (no best weapon).

### Combat Flow

**Each turn**, `IFight()` (a daemon) checks if any villain is in the room and fighting:

```
IFight()
  ├─ For each registered villain in the room:
  │    ├─ Unconscious? → maybe wake up (probability increases each turn)
  │    ├─ Staggered? → recover
  │    ├─ Fighting? → VillainBlow() — attack the player
  │    └─ Not in room? → call ActBusy handler (villain acts alone)
  └─ DoFight() orchestrates multi-villain rounds
```

**Player attacks** are handled by `HeroBlow()`:

```
HeroBlow()
  ├─ Check if player is staggered (skip turn)
  ├─ Calculate attack strength (base + score-derived bonus + wounds)
  ├─ Roll on defense table (Def1/Def2A/Def2B/Def3A/Def3B/Def3C)
  ├─ Apply result (miss, stagger, light/heavy wound, unconscious, kill, disarm)
  ├─ Print melee message from HeroMelee table
  └─ VillainResult() — update villain state
```

### Blow Results

| Result | Effect on Defender |
|--------|-------------------|
| `BlowMissed` | No effect |
| `BlowStag` | Skip next attack turn |
| `BlowLightWnd` | Strength -1 |
| `BlowHeavyWnd` | Strength -2 |
| `BlowLoseWpn` | Weapon drops to floor |
| `BlowUncon` | Negative strength (unconscious) |
| `BlowKill` | Strength 0, removed from game |
| `BlowHesitate` | Villain pauses (player unconscious variant) |
| `BlowSitDuck` | Villain finishes unconscious player |

### Melee Messages

Each villain has a `MeleeTable` — 9 sets of alternative messages (one set per blow result). Messages are composed of `MeleePart` fragments that can include literal text, weapon name markers (`FWep`), or defender name markers (`FDef`).

### Player Strength

Player strength is derived from score progress:

```
FightStrength = StrengthMin + Score / (ScoreMax / (StrengthMax - StrengthMin))
```

Wounds reduce effective strength. The `ICure` interrupt gradually heals the player (one point per 30 turns).

## Thief System (`thief.go`)

The thief is the most complex NPC, running as a daemon (`IThief`) every turn:

### Thief Daemon Behavior

```
IThief()
  ├─ In treasure room (not player's room)?
  │    → Deposit stolen treasures, hide them with FlgInvis
  ├─ In player's room?
  │    → ThiefVsAdventurer():
  │       ├─ 30% chance to appear if hidden
  │       ├─ Retreats if losing combat
  │       └─ Steals from player (treasures, visible items)
  ├─ In another room?
  │    → Rob room of treasures (75% chance per item)
  │    → Steal junk from room (worthless items, 10% chance)
  │    → In maze? → RobMaze (taunt the player)
  └─ Move to a random non-sacred, reachable room
      → DropJunk() (drop worthless items, 30% chance each)
```

### Key Functions

- `Rob(what, where, prob)` — steal treasures from a container/room
- `DepositBooty(room)` — move stolen treasures to a room (treasure room)
- `ThiefInTreasure()` — hide all room treasures when player enters treasure room
- `HackTreasures()` — make thief invisible, reveal all treasure room items
- `Infested(room)` — check if room has visible actors (used by sword glow)

## River and Dam System (`river.go`)

### Boat Mechanics

Three boat states: `InflatableBoat` (deflated), `InflatedBoat` (usable), `PuncturedBoat` (damaged).

- Inflate with pump, deflate manually
- Sharp objects puncture the boat when boarding
- Boat acts as a vehicle (`FlgVeh`) — player `MoveTo(boat)`, boat is in the room
- `RBoatFcn` handles boat-specific verbs (launch, deflate, navigate)

### River Navigation

The river flows through 5 rooms (River1–River5). When launched:

1. `IRiver` daemon activates with a speed based on the current room
2. Each tick, the boat moves downstream to the next river room
3. Speed increases as the boat approaches the falls
4. Reaching River5 with no exit → death (waterfall)

### Dam and Reservoir

The Flood Control Dam #3 has a control panel with colored buttons and a bolt:

- **Yellow button** → enables bolt turning
- **Brown button** → disables bolt turning
- **Red button** → toggles room lights
- **Blue button** → triggers maintenance room leak
- **Bolt** (turned with wrench) → opens/closes sluice gates

Gate state drives two interrupts:
- `IRempty` — drains the reservoir over 8 turns, enabling crossing
- `IRfill` — refills the reservoir over 8 turns, flooding it again

The `LoudRoom` becomes dangerous when gates are open (water roaring), and the `echo` command quiets it permanently.

## Death and Resurrection (`actions.go`)

`JigsUp(desc, isPlyr)` handles player death:

1. Print death message, deduct 10 points
2. First death with South Temple visited → sent to Entrance to Hades as a ghost
   - Player becomes "dead" — most actions are blocked by `DeadFunction`
   - Praying at South Temple resurrects the player
3. First death without temple → respawn in Forest1 with a second chance
4. Second death → permanent game over
5. `RandomizeObjects()` scatters carried items across above-ground rooms
6. `KillInterrupts()` disables most active clock events

## Scoring (`actions.go`)

- **Max score**: 350 points
- **Trophy case**: `OtvalFrob()` sums `TValue` of all objects in the trophy case
- **Score = BaseScore + trophy case value**
- `ScoreUpd(delta)` adjusts `BaseScore` and recalculates total
- Ranks: Beginner → Amateur → Novice → Junior → Adventurer → Master → Wizard → Master Adventurer

## Save/Restore (`save.go`)

Uses Go's `encoding/gob` for binary serialization.

### Saved State

```go
type gameState struct {
    ObjStates  []savedObject   // per-object: parent index, flags, strength, values, text, descs
    // Engine scalars: score, moves, lit, settings, flags
    // All ZorkData flags and counters
    // Object reference indices: Here, Winner, Player, ItObj, DescObj
    // Clock queue state: all 30 slots, boundaries
    // Villain probabilities
}
```

### Operations

All three functions return `error` (`nil` on success), propagating the underlying cause so verb handlers can display a meaningful message to the player (e.g. "Failed: open save file: no such file or directory").

- `doSave() error` — prompts for filename, captures state, encodes to file
- `doRestore() error` — prompts for filename, decodes, validates object count, applies state, rebuilds object tree
- `doRestart() error` — reinitializes everything from scratch (same as `InitGame()` but preserves I/O)

## Key Puzzles and Mechanics

| Puzzle | Mechanism |
|--------|-----------|
| Trap door | Move rug → reveal trap door → open it → cellar access |
| Troll bridge | Defeat troll in combat or throw him food |
| Cyclops | Feed hot peppers → give water → falls asleep → say "odysseus" → opens wall |
| Egg | Thief opens it if stolen; player attempts damage it |
| Loud room | Say "echo" to quiet the deafening sound |
| Dam/Reservoir | Open sluice gates → drain reservoir → cross → explore |
| Entrance to Hades | Ring bell → light candles → read book (exorcism ceremony) |
| Mirror rooms | Rub mirror → teleport between Mirror Room 1 and 2 |
| Rainbow | Wave sceptre at Aragain Falls → solidify rainbow → cross |
| Coal mine | Put coal in machine → turn switch with screwdriver → diamond |
| Dome/Torch room | Tie rope to railing → climb down |
| Thief's treasure room | Kill thief → all hidden treasures revealed |

## Clock Events

| Key | Type | Purpose |
|-----|------|---------|
| `IFight` | Daemon | Runs villain combat each turn |
| `ISword` | Interrupt | Sword glow detection (monsters nearby) |
| `IThief` | Daemon | Thief movement, stealing, treasure management |
| `ILantern` | Interrupt | Lantern fuel countdown and warnings |
| `ICandles` | Interrupt | Candle burn-down countdown |
| `ICure` | Interrupt | Gradual wound healing (30-turn intervals) |
| `ICyclops` | Interrupt | Cyclops anger escalation |
| `IForestRandom` | Interrupt | Random songbird chirp in forest |
| `IMaintRoom` | Interrupt | Maintenance room flooding |
| `IMatch` | Interrupt | Match burns out (2 turns) |
| `IRempty` / `IRfill` | Interrupt | Reservoir drain/fill |
| `IRiver` | Interrupt | River current moves boat downstream |
| `IXb` / `IXbh` / `IXc` | Interrupt | Hades exorcism ceremony timing |
