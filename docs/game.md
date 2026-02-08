# Game Package

The `game` package implements Zork I on top of the generic engine. It defines every room, item, NPC, verb handler, timed event, and puzzle in the game. The package dot-imports `engine` so game code reads like the original ZIL source.

## Files

| File | Responsibility |
|------|---------------|
| `zork_data.go` | `ZorkData` struct (game-specific mutable state), `gD()` accessor |
| `init.go` | `InitGame()`, `Run()`, clock function registration, well-known object wiring |
| `globals.go` | Global objects (it, me, hands, adventurer, etc.), sentinel objects, NPC action funcs |
| `items.go` | All game objects — rooms, items, NPCs — declared as package-level `Object` vars |
| `rooms_surface.go` | Surface world rooms (white house, forest, clearing, etc.) |
| `rooms_underground.go` | Underground rooms (cellar, troll room, temple, mine, etc.) |
| `rooms_maze.go` | Maze rooms and dead ends |
| `dungeon.go` | Lookup tables, random-selection pools, string data, navigation maps |
| `syntax_data.go` | `gameCommands` (all syntax definitions), `buzzWords`, `synonyms` map |
| `verbs.go` | Default verb handlers (vTake, vDrop, vOpen, vLook, etc.) |
| `verbs_movement.go` | Movement verbs (vWalk, vBoard, VClimb, vDisembark, vLaunch, etc.) |
| `verbs_meta.go` | Meta verbs (vScore, vInventory, vSave, vRestore, vRestart, vVersion, etc.) |
| `actions.go` | Object action functions, room action functions, puzzle logic, interrupt routines |
| `combat.go` | Combat system: melee tables, blow resolution, villain/hero attack logic |
| `thief.go` | thief NPC: stealing, treasure room, iThief daemon |
| `river.go` | river/boat system, dam/reservoir mechanics, water-level interrupts |
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

Accessed everywhere via `gD()`:

```go
func gD() *ZorkData {
    return G.GameData.(*ZorkData)
}
```

## Initialization (`init.go`)

`InitGame()` wires the game into the engine:

1. Creates `GameState` and `ZorkData`
2. Registers well-known objects (`RoomsObj`, `GlobalObj`, `NotHereObj`, `MeObj`, etc.)
3. Sets the implicit-take handler (`iTake`)
4. Calls `finalizeGameObjects()` to attach action functions and set item properties
5. Builds the object tree and vocabulary
6. Schedules initial clock events (thief daemon, combat daemon, lantern timer, candle timer)
7. Places the player at West of House

`Run()` calls `InitGame()`, shows the version banner and initial room description, then enters `MainLoop()`.

## World Model

### Rooms

Rooms are `Object` values with `Exits` maps and `Action` functions. They are children of the `rooms` sentinel object and are defined across three files by region:

- **Surface** (`rooms_surface.go`): White house perimeter, forest, clearing, canyon view, up a tree
- **Underground** (`rooms_underground.go`): cellar, troll room, maze areas, temple, mine, treasure room, dam, reservoir, loud room, cyclops room, and many more
- **Maze** (`rooms_maze.go`): The twisty maze passages, dead ends

Each room's `Action` function handles `ActLook` (room description), `ActBegin` (before command), `ActEnter` (on entry), and `ActEnd` (after command).

### Items

Items are `Object` values with `Synonyms`, `Adjectives`, `Desc`, flags, and optional `Action` functions. Items are defined in `items.go` and wired to their starting locations via `In` pointers.

Key item categories:
- **Treasures**: egg, canary, painting, chalice, pot of gold, jeweled scarab, trunk of jewels, etc. Each has a `TValue` that contributes to score when placed in the trophy case.
- **Tools**: lamp, sword, knife, shovel, screwdriver, wrench, keys, matches, rope
- **Containers**: mailbox, trophy case, machine, bottle, bag, basket, coffin
- **Vehicles**: Inflatable/inflated/punctured boat

### NPCs

NPCs are objects with `FlgPerson` or `FlgActor` and an `Action` function that handles multiple modes:

- **Normal verb handling**: responds to examine, talk, give, etc.
- **Combat callbacks**: `ActBusy`, `ActDead`, `ActUnconscious`, `ActConscious`, `ActFirst` — called by the combat system

Major NPCs: troll, thief, cyclops, bat, ghosts

## Verb System

### Syntax Definitions (`syntax_data.go`)

Every recognizable command is a `Syntax` entry in `gameCommands`. Examples:

| Input Pattern | Verb | VrbPrep | Obj1 | ObjPrep | Obj2 |
|--------------|------|---------|------|---------|------|
| `take X` | take | — | takeable | — | — |
| `put X in Y` | put | — | held | in | container |
| `attack X with Y` | attack | — | person | with | weapon |
| `look` | look | — | — | — | — |
| `turn X with Y` | turn | — | turnable | with | tool |

### Verb Handlers (`verbs.go`)

Default handlers for ~80 verbs. Each is a `func(ActionArg) bool`. Key verbs:

| Handler | Verbs | Behavior |
|---------|-------|----------|
| `vTake` | take, get | Pick up object, check weight/capacity |
| `vDrop` | drop | Place in room or vehicle |
| `vOpen` / `vClose` | open, close | Toggle FlgOpen, handle locked doors |
| `vLook` / `vFirstLook` | look | Describe room, list contents |
| `vExamine` | examine, x | Describe object in detail |
| `vRead` | read | Show object's Text field |
| `vAttack` | attack, kill | Initiate combat via `heroBlow()` |
| `vGive` | give | Transfer object to NPC |
| `vThrow` | throw | Throw at target |
| `vInventory` | inventory, i | List carried items |
| `vWalk` | walk, go | Move through exits |
| `vLampOn` / `vLampOff` | turn on/off | Toggle FlgOn, update lighting |
| `vEat` / `vDrink` | eat, drink | Consume food/water |

### Movement (`verbs_movement.go`)

`doWalk(dir)` is the core movement function:

1. Gets the exit for the direction from `G.Here.Exits`
2. Evaluates exit type (unconditional, conditional, door, function-based)
3. Calls `MoveToRoom(destination, true)` if passage is allowed
4. `MoveToRoom()` moves the player, updates `G.Here`, checks lighting, and optionally shows the room

Special movement: climbing, boarding/disembarking vehicles, launching the boat, entering/exiting through doors and windows.

### Meta Verbs (`verbs_meta.go`)

| Handler | Verb | Effect |
|---------|------|--------|
| `vScore` | score | Display current score and rank |
| `vInventory` | inventory | List carried objects |
| `vBrief` / `vVerbose` / `vSuperBrief` | brief/verbose/superbrief | Set description verbosity |
| `vSave` / `vRestore` | save/restore | Serialize/deserialize game state |
| `vRestart` | restart | Reset to initial state |
| `vVersion` | version | Show game version |
| `vQuit` | quit | End the game |
| `vDiagnose` | diagnose | Show health status |

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

Three villains: **troll** (best weapon: sword, +2), **thief** (best weapon: knife, +1), **cyclops** (no best weapon).

### Combat Flow

**Each turn**, `iFight()` (a daemon) checks if any villain is in the room and fighting:

```
iFight()
  ├─ For each registered villain in the room:
  │    ├─ Unconscious? → maybe wake up (probability increases each turn)
  │    ├─ Staggered? → recover
  │    ├─ Fighting? → villainBlow() — attack the player
  │    └─ Not in room? → call ActBusy handler (villain acts alone)
  └─ doFight() orchestrates multi-villain rounds
```

**Player attacks** are handled by `heroBlow()`:

```
heroBlow()
  ├─ Check if player is staggered (skip turn)
  ├─ Calculate attack strength (base + score-derived bonus + wounds)
  ├─ Roll on defense table (def1/def2A/def2B/def3A/def3B/def3C)
  ├─ Apply result (miss, stagger, light/heavy wound, unconscious, kill, disarm)
  ├─ Print melee message from heroMelee table
  └─ villainResult() — update villain state
```

### Blow Results

| Result | Effect on Defender |
|--------|-------------------|
| `blowMissed` | No effect |
| `blowStag` | Skip next attack turn |
| `blowLightWnd` | Strength -1 |
| `blowHeavyWnd` | Strength -2 |
| `blowLoseWpn` | Weapon drops to floor |
| `blowUncon` | Negative strength (unconscious) |
| `blowKill` | Strength 0, removed from game |
| `blowHesitate` | Villain pauses (player unconscious variant) |
| `blowSitDuck` | Villain finishes unconscious player |

### Melee Messages

Each villain has a `MeleeTable` — 9 sets of alternative messages (one set per blow result). Messages are composed of `MeleePart` fragments that can include literal text, weapon name markers (`fWep`), or defender name markers (`fDef`).

### Player Strength

Player strength is derived from score progress:

```
fightStrength = strengthMin + Score / (scoreMax / (strengthMax - strengthMin))
```

Wounds reduce effective strength. The `iCure` interrupt gradually heals the player (one point per 30 turns).

## Thief System (`thief.go`)

The thief is the most complex NPC, running as a daemon (`iThief`) every turn:

### Thief Daemon Behavior

```
iThief()
  ├─ In treasure room (not player's room)?
  │    → Deposit stolen treasures, hide them with FlgInvis
  ├─ In player's room?
  │    → thiefVsAdventurer():
  │       ├─ 30% chance to appear if hidden
  │       ├─ Retreats if losing combat
  │       └─ Steals from player (treasures, visible items)
  ├─ In another room?
  │    → rob room of treasures (75% chance per item)
  │    → Steal junk from room (worthless items, 10% chance)
  │    → In maze? → robMaze (taunt the player)
  └─ Move to a random non-sacred, reachable room
      → dropJunk() (drop worthless items, 30% chance each)
```

### Key Functions

- `rob(what, where, prob)` — steal treasures from a container/room
- `depositBooty(room)` — move stolen treasures to a room (treasure room)
- `thiefInTreasure()` — hide all room treasures when player enters treasure room
- `hackTreasures()` — make thief invisible, reveal all treasure room items
- `infested(room)` — check if room has visible actors (used by sword glow)

## River and Dam System (`river.go`)

### Boat Mechanics

Three boat states: `inflatableBoat` (deflated), `inflatedBoat` (usable), `puncturedBoat` (damaged).

- Inflate with pump, deflate manually
- Sharp objects puncture the boat when boarding
- Boat acts as a vehicle (`FlgVeh`) — player `MoveTo(boat)`, boat is in the room
- `rBoatFcn` handles boat-specific verbs (launch, deflate, navigate)

### River Navigation

The river flows through 5 rooms (river1–river5). When launched:

1. `iRiver` daemon activates with a speed based on the current room
2. Each tick, the boat moves downstream to the next river room
3. Speed increases as the boat approaches the falls
4. Reaching river5 with no exit → death (waterfall)

### Dam and Reservoir

The Flood Control Dam #3 has a control panel with colored buttons and a bolt:

- **Yellow button** → enables bolt turning
- **Brown button** → disables bolt turning
- **Red button** → toggles room lights
- **Blue button** → triggers maintenance room leak
- **bolt** (turned with wrench) → opens/closes sluice gates

Gate state drives two interrupts:
- `iRempty` — drains the reservoir over 8 turns, enabling crossing
- `iRfill` — refills the reservoir over 8 turns, flooding it again

The `loudRoom` becomes dangerous when gates are open (water roaring), and the `echo` command quiets it permanently.

## Death and Resurrection (`actions.go`)

`jigsUp(desc, isPlyr)` handles player death:

1. Print death message, deduct 10 points
2. First death with South Temple visited → sent to Entrance to Hades as a ghost
   - Player becomes "dead" — most actions are blocked by `deadFunction`
   - Praying at South Temple resurrects the player
3. First death without temple → respawn in forest1 with a second chance
4. Second death → permanent game over
5. `randomizeObjects()` scatters carried items across above-ground rooms
6. `killInterrupts()` disables most active clock events

## Scoring (`actions.go`)

- **Max score**: 350 points
- **Trophy case**: `otvalFrob()` sums `TValue` of all objects in the trophy case
- **Score = BaseScore + trophy case value**
- `scoreUpd(delta)` adjusts `BaseScore` and recalculates total
- Ranks: Beginner → Amateur → Novice → Junior → adventurer → Master → Wizard → Master adventurer

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
| troll bridge | Defeat troll in combat or throw him food |
| cyclops | Feed hot peppers → give water → falls asleep → say "odysseus" → opens wall |
| egg | thief opens it if stolen; player attempts damage it |
| Loud room | Say "echo" to quiet the deafening sound |
| dam/reservoir | Open sluice gates → drain reservoir → cross → explore |
| Entrance to Hades | Ring bell → light candles → read book (exorcism ceremony) |
| Mirror rooms | Rub mirror → teleport between Mirror Room 1 and 2 |
| rainbow | Wave sceptre at Aragain Falls → solidify rainbow → cross |
| coal mine | Put coal in machine → turn switch with screwdriver → diamond |
| Dome/torch room | Tie rope to railing → climb down |
| thief's treasure room | Kill thief → all hidden treasures revealed |

## Clock Events

| Key | Type | Purpose |
|-----|------|---------|
| `iFight` | Daemon | Runs villain combat each turn |
| `iSword` | Interrupt | sword glow detection (monsters nearby) |
| `iThief` | Daemon | thief movement, stealing, treasure management |
| `iLantern` | Interrupt | Lantern fuel countdown and warnings |
| `iCandles` | Interrupt | Candle burn-down countdown |
| `iCure` | Interrupt | Gradual wound healing (30-turn intervals) |
| `iCyclops` | Interrupt | cyclops anger escalation |
| `iForestRandom` | Interrupt | Random songbird chirp in forest |
| `iMaintRoom` | Interrupt | Maintenance room flooding |
| `iMatch` | Interrupt | match burns out (2 turns) |
| `iRempty` / `iRfill` | Interrupt | reservoir drain/fill |
| `iRiver` | Interrupt | river current moves boat downstream |
| `iXb` / `iXbh` / `iXc` | Interrupt | Hades exorcism ceremony timing |
