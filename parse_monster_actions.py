# /Users/marcozingoni/Playgound/lazydnd/parse_monster_actions.py
#!/usr/bin/env python3
"""
Parse monster actions from HTML format and create structured ActionList
"""

import json
import re
from html.parser import HTMLParser
from typing import List, Dict, Optional

class ActionHTMLParser(HTMLParser):
    """Parse HTML action descriptions"""

    def __init__(self):
        super().__init__()
        self.current_text = []
        self.actions = []

    def handle_data(self, data):
        self.current_text.append(data.strip())

    def get_text(self):
        return ' '.join(self.current_text).strip()

def parse_action_text(action_text: str) -> Optional[Dict]:
    """Parse a single action text into structured format"""

    # Extract action name (usually in <strong> tags or at start)
    name_match = re.search(r'<strong>(.*?)</strong>', action_text, re.IGNORECASE)
    if not name_match:
        name_match = re.search(r'^([^.]+?)\.', action_text)

    if not name_match:
        return None

    name = name_match.group(1).strip()

    # Remove HTML tags for parsing
    clean_text = re.sub(r'<[^>]+>', '', action_text)

    # Determine attack type
    attack_type = "other"
    if "Melee Weapon Attack" in clean_text or "Melee Attack" in clean_text:
        attack_type = "melee"
    elif "Ranged Weapon Attack" in clean_text or "Ranged Attack" in clean_text:
        attack_type = "ranged"
    elif "Melee or Ranged Weapon Attack" in clean_text:
        attack_type = "melee/ranged"
    elif "Melee Spell Attack" in clean_text:
        attack_type = "melee_spell"
    elif "Ranged Spell Attack" in clean_text:
        attack_type = "ranged_spell"

    # Extract attack roll bonus
    roll_match = re.search(r'([+\-]\d+)\s+to hit', clean_text)
    attack_roll = roll_match.group(1) if roll_match else None

    # Extract reach
    reach_match = re.search(r'reach (\d+)\s*ft', clean_text, re.IGNORECASE)
    reach = f"{reach_match.group(1)}ft" if reach_match else None

    # Extract range
    range_match = re.search(r'range (\d+(?:/\d+)?)\s*ft', clean_text, re.IGNORECASE)
    attack_range = f"{range_match.group(1)}ft" if range_match else None

    # Extract damage (can be multiple)
    damage_list = []
    damage_type_list = []

    # Pattern for damage: "11 (2d6 + 4) piercing damage"
    damage_patterns = re.finditer(r'(\d+)\s*\(([^)]+)\)\s+(\w+)\s+damage', clean_text, re.IGNORECASE)
    for match in damage_patterns:
        dice_formula = match.group(2).strip()
        damage_type = match.group(3).strip()
        damage_list.append(dice_formula)
        damage_type_list.append(damage_type)

    # Also look for saving throw damage
    save_damage = re.finditer(r'taking\s+\d+\s*\(([^)]+)\)\s+(\w+)\s+damage', clean_text, re.IGNORECASE)
    for match in save_damage:
        dice_formula = match.group(1).strip()
        damage_type = match.group(2).strip()
        if dice_formula not in damage_list:
            damage_list.append(dice_formula)
            damage_type_list.append(damage_type)

    # Combine damage
    damage = ", ".join(damage_list) if damage_list else None
    damage_types = ", ".join(damage_type_list) if damage_type_list else None

    # Extract DC for saving throws
    dc_match = re.search(r'DC (\d+)', clean_text)
    save_dc = f"DC {dc_match.group(1)}" if dc_match else None

    # Extract saving throw type
    save_match = re.search(r'DC \d+ (\w+) saving throw', clean_text, re.IGNORECASE)
    save_type = save_match.group(1) if save_match else None

    # Build action dict
    action = {
        "name": name,
        "type": attack_type,
        "description": clean_text
    }

    if attack_roll:
        action["roll"] = attack_roll
    if reach:
        action["reach"] = reach
    if attack_range:
        action["range"] = attack_range
    if damage:
        action["damage"] = damage
    if damage_types:
        action["damage_type"] = damage_types
    if save_dc:
        action["save_dc"] = save_dc
    if save_type:
        action["save_type"] = save_type

    return action

def parse_actions(actions_html: str) -> List[Dict]:
    """Parse all actions from HTML string"""

    if not actions_html or actions_html == "--":
        return []

    # Split by <p> tags or <strong> tags to separate actions
    action_blocks = re.split(r'</p>\s*<p>', actions_html)

    parsed_actions = []
    for block in action_blocks:
        # Clean up the block
        block = block.replace('<p>', '').replace('</p>', '').strip()
        if not block:
            continue

        # Parse the action
        action = parse_action_text(block)
        if action:
            parsed_actions.append(action)

    return parsed_actions

def process_monsters_file(input_file: str, output_file: str):
    """Process the monsters JSON file and add ActionList to each monster"""

    print(f"Loading monsters from {input_file}...")
    with open(input_file, 'r', encoding='utf-8') as f:
        monsters = json.load(f)

    print(f"Processing {len(monsters)} monsters...")

    for i, monster in enumerate(monsters):
        if (i + 1) % 100 == 0:
            print(f"  Processed {i + 1}/{len(monsters)} monsters...")

        # Get actions HTML
        actions_html = monster.get('Actions', '')

        # Parse actions
        action_list = parse_actions(actions_html)

        # Add to monster
        monster['ActionNumber'] = len(action_list)
        monster['ActionList'] = action_list

    print(f"Writing updated monsters to {output_file}...")
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(monsters, f, indent=2, ensure_ascii=False)

    print("Done!")
    print(f"Processed {len(monsters)} monsters")

    # Print some statistics
    total_actions = sum(m.get('ActionNumber', 0) for m in monsters)
    monsters_with_actions = sum(1 for m in monsters if m.get('ActionNumber', 0) > 0)
    print(f"Total actions: {total_actions}")
    print(f"Monsters with actions: {monsters_with_actions}/{len(monsters)}")

if __name__ == "__main__":
    input_file = "assets/monsters.json"
    output_file = "assets/monsters.json"

    process_monsters_file(input_file, output_file)
