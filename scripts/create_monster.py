#!/usr/bin/env python3
"""
LazyDnD Monster Creator
Interactive CLI tool to create custom monsters in the correct JSON format.
"""

import json
import sys
import os
from typing import List, Dict, Optional

def calculate_modifier(ability_score: int) -> str:
    """Calculate ability modifier from ability score."""
    mod = (ability_score - 10) // 2
    if mod >= 0:
        return f"(+{mod})"
    return f"({mod})"

def get_input(prompt: str, default: str = "", required: bool = True) -> str:
    """Get user input with optional default value."""
    if default:
        user_input = input(f"{prompt} [{default}]: ").strip()
        return user_input if user_input else default
    else:
        while True:
            user_input = input(f"{prompt}: ").strip()
            if user_input or not required:
                return user_input
            print("This field is required. Please provide a value.")

def get_int(prompt: str, default: Optional[int] = None) -> int:
    """Get integer input with validation."""
    while True:
        try:
            if default is not None:
                user_input = input(f"{prompt} [{default}]: ").strip()
                return int(user_input) if user_input else default
            else:
                return int(input(f"{prompt}: ").strip())
        except ValueError:
            print("Please enter a valid number.")

def get_yes_no(prompt: str, default: bool = False) -> bool:
    """Get yes/no input."""
    default_str = "Y/n" if default else "y/N"
    while True:
        response = input(f"{prompt} [{default_str}]: ").strip().lower()
        if not response:
            return default
        if response in ['y', 'yes']:
            return True
        if response in ['n', 'no']:
            return False
        print("Please enter 'y' or 'n'.")

def create_action() -> Dict:
    """Create a monster action interactively."""
    print("\n--- Creating Action ---")

    action = {
        "name": get_input("Action name (e.g., 'Bite', 'Claw')", required=True),
        "description": get_input("Action description", required=True)
    }

    print("\nAction type:")
    print("  1. Melee attack")
    print("  2. Ranged attack")
    print("  3. Other (spell, ability, etc.)")

    action_type = get_input("Select type (1-3)", "3")
    if action_type == "1":
        action["type"] = "melee"
        action["roll"] = get_input("Attack bonus (e.g., '+5')", required=False)
        action["reach"] = get_input("Reach (e.g., '5ft', '10ft')", "5ft")
        action["damage"] = get_input("Damage dice (e.g., '1d8 + 3')", required=False)
        action["damage_type"] = get_input("Damage type (e.g., 'piercing')", required=False)
    elif action_type == "2":
        action["type"] = "ranged"
        action["roll"] = get_input("Attack bonus (e.g., '+5')", required=False)
        action["range"] = get_input("Range (e.g., '30/120 ft.')", required=False)
        action["damage"] = get_input("Damage dice (e.g., '1d8 + 3')", required=False)
        action["damage_type"] = get_input("Damage type (e.g., 'piercing')", required=False)
    else:
        action["type"] = "other"

    # Optional save DC
    if get_yes_no("Does this action require a saving throw?", False):
        action["save_dc"] = get_input("Save DC (e.g., 'DC 13')", required=False)
        action["save_type"] = get_input("Save type (e.g., 'Dexterity')", required=False)

    return action

def create_monster() -> Dict:
    """Create a monster interactively."""
    print("\n" + "="*60)
    print("LazyDnD Monster Creator".center(60))
    print("="*60)
    print("\nCreate a custom D&D 5e monster for LazyDnD")
    print("Press Ctrl+C at any time to cancel\n")

    monster = {}

    # Basic info
    print("=== BASIC INFORMATION ===\n")
    monster["name"] = get_input("Monster name", required=True)

    # Meta (size, type, alignment)
    print("\nExamples: 'Medium humanoid, neutral evil' or 'Large dragon, chaotic good'")
    monster["meta"] = get_input("Size, type, and alignment", "Medium humanoid, unaligned")

    # Combat stats
    print("\n=== COMBAT STATS ===\n")
    monster["Armor Class"] = get_input("Armor Class (e.g., '15 (Natural Armor)')", "10")
    monster["Hit Points"] = get_input("Hit Points (e.g., '45 (6d8 + 18)')", "10 (3d8)")
    monster["Speed"] = get_input("Speed (e.g., '30 ft., fly 60 ft.')", "30 ft.")

    # Ability scores
    print("\n=== ABILITY SCORES ===\n")
    for ability in ["STR", "DEX", "CON", "INT", "WIS", "CHA"]:
        score = get_int(f"{ability}", 10)
        monster[ability] = str(score)
        monster[f"{ability}_mod"] = calculate_modifier(score)

    # Challenge rating
    print("\n=== CHALLENGE RATING ===\n")
    print("Examples: '0 (10 XP)', '1/4 (50 XP)', '1 (200 XP)', '5 (1,800 XP)'")
    monster["Challenge"] = get_input("Challenge rating", "0 (10 XP)")

    # Optional fields
    print("\n=== OPTIONAL INFORMATION ===\n")

    if get_yes_no("Add saving throw proficiencies?", False):
        print("Example: 'STR +5, CON +3'")
        monster["Saving Throws"] = get_input("Saving throws", required=False)

    if get_yes_no("Add skills?", False):
        print("Example: 'Perception +4, Stealth +6'")
        monster["Skills"] = get_input("Skills", required=False)

    if get_yes_no("Add senses?", False):
        print("Example: 'Darkvision 60 ft., Passive Perception 12'")
        monster["Senses"] = get_input("Senses", "Passive Perception 10")

    if get_yes_no("Add languages?", False):
        print("Example: 'Common, Elvish'")
        monster["Languages"] = get_input("Languages", required=False)

    # Traits
    if get_yes_no("Add traits/features?", False):
        print("\nExample: '<p><em><strong>Keen Smell.</strong></em> The wolf has advantage on Wisdom (Perception) checks that rely on smell.</p>'")
        print("Note: Use HTML tags for formatting (see monsters.json for examples)")
        monster["Traits"] = get_input("Traits (HTML formatted)", required=False)

    # Actions
    print("\n=== ACTIONS ===\n")
    actions_html = []
    action_list = []

    if get_yes_no("Add actions?", True):
        while True:
            action = create_action()
            action_list.append(action)

            # Create HTML representation
            html = f"<p><em><strong>{action['name']}</strong></em> {action['description']}</p>"
            actions_html.append(html)

            if not get_yes_no("Add another action?", False):
                break

    if actions_html:
        monster["Actions"] = "".join(actions_html)
        monster["ActionList"] = action_list
        monster["ActionNumber"] = len(action_list)
    else:
        monster["Actions"] = ""
        monster["ActionList"] = []
        monster["ActionNumber"] = 0

    # Legendary actions
    if get_yes_no("Add legendary actions?", False):
        print("\nExample: '<p>The dragon can take 3 legendary actions...</p>'")
        monster["Legendary Actions"] = get_input("Legendary actions (HTML formatted)", required=False)

    # Image URL
    if get_yes_no("Add image URL?", False):
        monster["img_url"] = get_input("Image URL", required=False)

    return monster

def save_monster(monster: Dict, output_file: str = None, save_to_lazydnd: bool = False):
    """Save monster to a JSON file."""
    if output_file is None:
        monster_name = monster["name"].lower().replace(" ", "_")
        output_file = f"{monster_name}.json"

    # If saving to LazyDnD directory
    if save_to_lazydnd:
        home_dir = os.path.expanduser("~")
        lazydnd_monsters_dir = os.path.join(home_dir, ".lazydnd", "monsters")

        # Create directory if it doesn't exist
        try:
            os.makedirs(lazydnd_monsters_dir, exist_ok=True)
        except Exception as e:
            print(f"\n⚠️  Could not create {lazydnd_monsters_dir}: {e}")
            print("Saving to current directory instead.")
            save_to_lazydnd = False

        if save_to_lazydnd:
            output_file = os.path.join(lazydnd_monsters_dir, output_file)

    # Check if file exists
    if os.path.exists(output_file):
        if not get_yes_no(f"\n'{output_file}' already exists. Overwrite?", False):
            monster_name = monster["name"].lower().replace(" ", "_")
            output_file = get_input("Enter new filename", f"{monster_name}_new.json")
            if save_to_lazydnd and not output_file.startswith(os.path.expanduser("~")):
                output_file = os.path.join(lazydnd_monsters_dir, os.path.basename(output_file))

    try:
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump([monster], f, indent=2, ensure_ascii=False)

        print(f"\n✅ Monster saved to: {output_file}")

        if save_to_lazydnd:
            print(f"\n🎉 Monster ready to use!")
            print(f"   Restart LazyDnD to see your new monster in the Monsters panel.")
        else:
            print(f"\nTo add this monster to LazyDnD:")
            print(f"1. Move or copy {output_file} to ~/.lazydnd/monsters/")
            print(f"2. Restart LazyDnD")
        return True
    except Exception as e:
        print(f"\n❌ Error saving file: {e}")
        return False

def preview_monster(monster: Dict):
    """Display a preview of the created monster."""
    print("\n" + "="*60)
    print("MONSTER PREVIEW".center(60))
    print("="*60)
    print(json.dumps(monster, indent=2))
    print("="*60)

def main():
    """Main function."""
    try:
        monster = create_monster()

        # Preview
        print("\n")
        if get_yes_no("Preview the monster JSON?", True):
            preview_monster(monster)

        # Save
        print("\n")
        if get_yes_no("Save this monster?", True):
            save_to_lazydnd = get_yes_no("Save directly to ~/.lazydnd/monsters/ (recommended)?", True)

            output_file = None
            if get_yes_no("Specify custom filename?", False):
                default_name = f"{monster['name'].lower().replace(' ', '_')}.json"
                output_file = get_input("Output filename", default_name)

            save_monster(monster, output_file, save_to_lazydnd)
        else:
            print("\n❌ Monster not saved.")

        print("\n✨ Done!\n")

    except KeyboardInterrupt:
        print("\n\n❌ Cancelled by user.\n")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ Error: {e}\n")
        sys.exit(1)

if __name__ == "__main__":
    main()
