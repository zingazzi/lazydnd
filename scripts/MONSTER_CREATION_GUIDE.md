# Monster Creation Guide

Quick reference for creating custom monsters for LazyDnD.

## Quick Start

```bash
./scripts/create_monster.py
```

Follow the interactive prompts!

## Field Guide

### Required Fields

| Field | Description | Example |
|-------|-------------|---------|
| `name` | Monster name | "Dire Wolf", "Ancient Red Dragon" |
| `meta` | Size, type, alignment | "Large beast, unaligned" |
| `Armor Class` | AC with source | "15 (Natural Armor)", "10" |
| `Hit Points` | HP with dice | "45 (6d8 + 18)" |
| `Speed` | Movement speeds | "30 ft.", "30 ft., fly 60 ft." |
| `STR` - `CHA` | Ability scores (3-30) | "10", "18", "7" |
| `Challenge` | CR with XP | "1 (200 XP)", "1/4 (50 XP)" |

### Optional Fields

| Field | Description | Example |
|-------|-------------|---------|
| `Saving Throws` | Saving throw bonuses | "STR +5, CON +3" |
| `Skills` | Skill bonuses | "Perception +4, Stealth +6" |
| `Senses` | Special senses | "Darkvision 60 ft., Passive Perception 13" |
| `Languages` | Languages known | "Common, Draconic" |
| `Traits` | Special abilities | See HTML format below |
| `Actions` | Available actions | See HTML format below |
| `Legendary Actions` | Legendary options | See HTML format below |
| `img_url` | Monster image | "https://..." |

## Challenge Ratings

Common CR values and their XP:

| CR | XP | CR | XP | CR | XP |
|----|----|----|----|----|-----|
| 0 | 10 | 4 | 1,100 | 13 | 10,000 |
| 1/8 | 25 | 5 | 1,800 | 14 | 11,500 |
| 1/4 | 50 | 6 | 2,300 | 15 | 13,000 |
| 1/2 | 100 | 7 | 2,900 | 16 | 15,000 |
| 1 | 200 | 8 | 3,900 | 17 | 18,000 |
| 2 | 450 | 9 | 5,000 | 18 | 20,000 |
| 3 | 700 | 10 | 5,900 | 19 | 22,000 |
| | | | | 20 | 25,000 |

## HTML Formatting

### Traits Example

```html
<p><em><strong>Pack Tactics.</strong></em> The wolf has advantage on attack rolls against a creature if at least one ally is within 5 feet of the creature and the ally isn't incapacitated.</p>
```

### Actions Example

```html
<p><em><strong>Bite.</strong></em> <em>Melee Weapon Attack:</em> +5 to hit, reach 5 ft., one target. <em>Hit:</em> 10 (2d6 + 3) piercing damage.</p>
```

### Multiple Traits/Actions

Concatenate multiple `<p>` tags:

```html
<p><em><strong>First Trait.</strong></em> Description here.</p><p><em><strong>Second Trait.</strong></em> Another description.</p>
```

## Action Types

The tool supports three action types:

### 1. Melee Attack
- Attack bonus (e.g., `+5`)
- Reach (e.g., `5ft`, `10ft`)
- Damage dice (e.g., `1d8 + 3`)
- Damage type (e.g., `piercing`, `slashing`)

### 2. Ranged Attack
- Attack bonus (e.g., `+5`)
- Range (e.g., `30/120 ft.`)
- Damage dice (e.g., `1d6 + 2`)
- Damage type (e.g., `piercing`)

### 3. Other
- Spells, special abilities, etc.
- Can include save DC and type

## Adding Monsters to LazyDnD

After creating your monster JSON file:

### Option 1: Direct Save (Easiest - Recommended)

The monster creator tool will ask if you want to save directly to `~/.lazydnd/monsters/`:
- Answer **Yes** (default)
- File is saved to `~/.lazydnd/monsters/your_monster.json`
- Restart LazyDnD
- Your monster appears in the Monsters panel! 🎉

### Option 2: Manual Copy

If you saved the file locally:

```bash
# Create the monsters directory if it doesn't exist
mkdir -p ~/.lazydnd/monsters

# Copy your monster file
cp my_monster.json ~/.lazydnd/monsters/

# Restart LazyDnD
```

### How It Works

- LazyDnD automatically scans `~/.lazydnd/monsters/` on startup
- All `.json` files are loaded and added to the monster list
- One file per monster (easier to manage!)
- Custom monsters override default monsters with the same name

### Managing Custom Monsters

**Add a monster:**
```bash
cp new_monster.json ~/.lazydnd/monsters/
```

**Remove a monster:**
```bash
rm ~/.lazydnd/monsters/unwanted_monster.json
```

**Edit a monster:**
```bash
nano ~/.lazydnd/monsters/my_monster.json
# Or use your favorite editor
```

**List all custom monsters:**
```bash
ls -lh ~/.lazydnd/monsters/
```

## Tips

1. **Ability Modifiers**: The tool auto-calculates these - no need to do math!
2. **HTML Formatting**: You can copy-paste from existing monsters in `assets/monsters.json`
3. **Preview First**: Always preview before saving to catch mistakes
4. **Use Defaults**: Press Enter to accept defaults for quick creation
5. **Cancel Anytime**: Press Ctrl+C to exit without saving

## Example: Simple Monster

```json
{
  "name": "Giant Rat",
  "meta": "Small beast, unaligned",
  "Armor Class": "12",
  "Hit Points": "7 (2d6)",
  "Speed": "30 ft.",
  "STR": "7",
  "STR_mod": "(-2)",
  "DEX": "15",
  "DEX_mod": "(+2)",
  "CON": "11",
  "CON_mod": "(+0)",
  "INT": "2",
  "INT_mod": "(-4)",
  "WIS": "10",
  "WIS_mod": "(+0)",
  "CHA": "4",
  "CHA_mod": "(-3)",
  "Senses": "Darkvision 60 ft., Passive Perception 10",
  "Challenge": "1/8 (25 XP)",
  "Traits": "<p><em><strong>Keen Smell.</strong></em> The rat has advantage on Wisdom (Perception) checks that rely on smell.</p>",
  "Actions": "<p><em><strong>Bite.</strong></em> <em>Melee Weapon Attack:</em> +4 to hit, reach 5 ft., one target. <em>Hit:</em> 4 (1d4 + 2) piercing damage.</p>",
  "ActionNumber": 1,
  "ActionList": [
    {
      "name": "Bite",
      "type": "melee",
      "description": "Bite. Melee Weapon Attack: +4 to hit, reach 5 ft., one target. Hit: 4 (1d4 + 2) piercing damage.",
      "roll": "+4",
      "reach": "5ft",
      "damage": "1d4 + 2",
      "damage_type": "piercing"
    }
  ]
}
```

## Troubleshooting

**Q: My monster doesn't appear in LazyDnD**
- A: Make sure the JSON is valid (use a JSON validator)
- A: Check that `custom_monsters.json` is in `~/.lazydnd/`
- A: Restart LazyDnD after adding monsters

**Q: Actions aren't parsed correctly**
- A: Ensure ActionList format matches the examples
- A: Check that action types are: "melee", "ranged", or "other"

**Q: Can I edit an existing monster?**
- A: Yes! Copy the monster object from `assets/monsters.json`, edit it, and save as a custom monster

**Q: How do I delete a custom monster?**
- A: Edit `~/.lazydnd/custom_monsters.json` and remove the monster object

---

Happy monster creating! 🐉
