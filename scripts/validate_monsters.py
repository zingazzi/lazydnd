#!/usr/bin/env python3
"""Validate the updated monsters.json file"""

import json

# Load monsters
with open('assets/monsters.json', 'r', encoding='utf-8') as f:
    monsters = json.load(f)

print(f"Total monsters: {len(monsters)}")
print(f"\nSample of monsters with their action counts:")
for m in monsters[:10]:
    print(f"  {m['name']}: {m['ActionNumber']} actions")

print(f"\nMonsters with most actions:")
sorted_monsters = sorted(monsters, key=lambda m: m.get('ActionNumber', 0), reverse=True)
for m in sorted_monsters[:5]:
    print(f"  {m['name']}: {m['ActionNumber']} actions")

print(f"\nValidation:")
print(f"  - All monsters have ActionNumber: {all('ActionNumber' in m for m in monsters)}")
print(f"  - All monsters have ActionList: {all('ActionList' in m for m in monsters)}")
print(f"  - Monsters with actions: {sum(1 for m in monsters if m.get('ActionNumber', 0) > 0)}/{len(monsters)}")

# Check for parsing issues
print(f"\nChecking action types:")
action_types = {}
for m in monsters:
    for action in m.get('ActionList', []):
        action_type = action.get('type', 'unknown')
        action_types[action_type] = action_types.get(action_type, 0) + 1

for action_type, count in sorted(action_types.items(), key=lambda x: x[1], reverse=True):
    print(f"  {action_type}: {count}")

print("\n✅ Validation complete!")
