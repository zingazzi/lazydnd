# LazyDnD Scripts

This directory contains utility scripts for maintaining and extending LazyDnD.

## Monster Creator Tool

### `create_monster.py`

An interactive CLI tool to create custom monsters in the correct JSON format for LazyDnD.

#### Usage

```bash
# From the scripts directory
python3 create_monster.py

# Or run directly (if executable)
./create_monster.py
```

#### Features

- ✅ **Interactive prompts** - Step-by-step guidance through all monster fields
- ✅ **Auto-calculate modifiers** - Automatically calculates ability score modifiers
- ✅ **Flexible actions** - Support for melee, ranged, and special actions
- ✅ **Validation** - Ensures required fields are filled
- ✅ **Default values** - Sensible defaults for quick creation
- ✅ **Preview before save** - Review the JSON before saving
- ✅ **HTML formatting** - Supports HTML formatting for traits and actions

#### Example Workflow

1. **Run the tool**:
   ```bash
   ./scripts/create_monster.py
   ```

2. **Follow the prompts**:
   - Basic info (name, size, type, alignment)
   - Combat stats (AC, HP, Speed)
   - Ability scores (STR, DEX, CON, INT, WIS, CHA)
   - Challenge rating
   - Optional: Skills, senses, languages, traits
   - Actions (melee, ranged, or special)
   - Optional: Legendary actions, image URL

3. **Preview and save**:
   - Review the generated JSON
   - Save to a file (e.g., `my_monster.json`)

4. **Add to LazyDnD**:

   The tool will ask if you want to save directly to `~/.lazydnd/monsters/` (recommended).

   **Option 1: Direct save (easiest)**
   - Select "Yes" when asked to save to ~/.lazydnd/monsters/
   - Restart LazyDnD
   - Your monster is ready!

   **Option 2: Manual copy**
   ```bash
   mkdir -p ~/.lazydnd/monsters
   mv my_monster.json ~/.lazydnd/monsters/
   ```

   **Note:** LazyDnD automatically loads all `.json` files from `~/.lazydnd/monsters/`

#### Example Monster Creation

```
Monster name: Dire Wolf
Size, type, and alignment: Large beast, unaligned
Armor Class: 14 (Natural Armor)
Hit Points: 37 (5d10 + 10)
Speed: 50 ft.

STR: 17
DEX: 15
CON: 15
INT: 3
WIS: 12
CHA: 7

Challenge rating: 1 (200 XP)

Add skills? [y/N]: y
Skills: Perception +3, Stealth +4

Add senses? [y/N]: y
Senses: Darkvision 60 ft., Passive Perception 13

Add traits? [y/N]: y
Traits: <p><em><strong>Keen Hearing and Smell.</strong></em> The wolf has advantage on Wisdom (Perception) checks that rely on hearing or smell.</p><p><em><strong>Pack Tactics.</strong></em> The wolf has advantage on an attack roll against a creature if at least one of the wolf's allies is within 5 feet of the creature and the ally isn't incapacitated.</p>

--- Creating Action ---
Action name: Bite
Action description: Melee Weapon Attack: +5 to hit, reach 5 ft., one target. Hit: 10 (2d6 + 3) piercing damage. If the target is a creature, it must succeed on a DC 13 Strength saving throw or be knocked prone.
Select type (1-3): 1
Attack bonus: +5
Reach: 5ft
Damage dice: 2d6 + 3
Damage type: piercing
Does this action require a saving throw? [y/N]: y
Save DC: DC 13
Save type: Strength
```

## Other Scripts

### `parse_monster_actions.py`
Parses monster action text into structured ActionList format.

### `validate_monsters.py`
Validates monster JSON files for correct structure and required fields.

### `extract_changelog.sh`
Extracts changelog entries for a specific version.

### `extract_version.sh`
Extracts the current version from main.go.

---

## Contributing

If you create useful scripts for LazyDnD, feel free to add them to this directory with appropriate documentation!
