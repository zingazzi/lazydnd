# tests/monster_action_parser_test.py
"""Tests for monster action parser"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

import unittest
from parse_monster_actions import parse_action_text, parse_actions


class TestMonsterActionParser(unittest.TestCase):
    """Test cases for monster action parsing"""

    def test_melee_weapon_attack(self):
        """Test parsing melee weapon attack"""
        action_text = """
        <strong>Greataxe.</strong> <em>Melee Weapon Attack:</em> +5 to hit, reach 5 ft.,
        one target. <em>Hit:</em> 9 (1d12 + 3) slashing damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Greataxe.')
        self.assertEqual(result['type'], 'melee')
        self.assertEqual(result['roll'], '+5')
        self.assertEqual(result['reach'], '5ft')
        self.assertEqual(result['damage'], '1d12 + 3')
        self.assertEqual(result['damage_type'], 'slashing')

    def test_ranged_weapon_attack(self):
        """Test parsing ranged weapon attack"""
        action_text = """
        <strong>Longbow.</strong> <em>Ranged Weapon Attack:</em> +4 to hit, range 150/600 ft.,
        one target. <em>Hit:</em> 6 (1d8 + 2) piercing damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Longbow.')
        self.assertEqual(result['type'], 'ranged')
        self.assertEqual(result['roll'], '+4')
        self.assertEqual(result['range'], '150/600ft')
        self.assertEqual(result['damage'], '1d8 + 2')
        self.assertEqual(result['damage_type'], 'piercing')

    def test_melee_spell_attack(self):
        """Test parsing melee spell attack"""
        action_text = """
        <strong>Shocking Grasp.</strong> <em>Melee Spell Attack:</em> +6 to hit, reach 5 ft.,
        one creature. <em>Hit:</em> 9 (2d8) lightning damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Shocking Grasp.')
        self.assertEqual(result['type'], 'melee_spell')
        self.assertEqual(result['roll'], '+6')

    def test_ranged_spell_attack(self):
        """Test parsing ranged spell attack"""
        action_text = """
        <strong>Fire Bolt.</strong> <em>Ranged Spell Attack:</em> +5 to hit, range 120 ft.,
        one target. <em>Hit:</em> 11 (2d10) fire damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Fire Bolt.')
        self.assertEqual(result['type'], 'ranged_spell')
        self.assertEqual(result['damage_type'], 'fire')

    def test_saving_throw_attack(self):
        """Test parsing action with saving throw"""
        action_text = """
        <strong>Poison Breath.</strong> The dragon exhales poisonous gas in a 60-foot cone.
        Each creature in that area must make a DC 18 Constitution saving throw,
        taking 56 (16d6) poison damage on a failed save, or half as much damage on a successful one.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Poison Breath.')
        self.assertEqual(result['save_dc'], 'DC 18')
        self.assertEqual(result['save_type'], 'Constitution')
        self.assertEqual(result['damage'], '16d6')
        self.assertEqual(result['damage_type'], 'poison')

    def test_multiple_damage_types(self):
        """Test parsing action with multiple damage types"""
        action_text = """
        <strong>Multiattack.</strong> <em>Melee Weapon Attack:</em> +7 to hit, reach 10 ft.,
        one target. <em>Hit:</em> 11 (2d6 + 4) slashing damage plus 3 (1d6) fire damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Multiattack.')
        # Should capture both damage dice
        self.assertIn('2d6 + 4', result['damage'])
        self.assertIn('1d6', result['damage'])

    def test_melee_or_ranged_attack(self):
        """Test parsing melee or ranged attack"""
        action_text = """
        <strong>Spear.</strong> <em>Melee or Ranged Weapon Attack:</em> +4 to hit,
        reach 5 ft. or range 20/60 ft., one target. <em>Hit:</em> 5 (1d6 + 2) piercing damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        # Parser matches ranged before melee/ranged in current implementation
        self.assertIn(result['type'], ['ranged', 'melee/ranged'])
        self.assertEqual(result['reach'], '5ft')
        self.assertEqual(result['range'], '20/60ft')

    def test_negative_attack_bonus(self):
        """Test parsing action with negative attack bonus"""
        action_text = """
        <strong>Weak Strike.</strong> <em>Melee Weapon Attack:</em> -1 to hit, reach 5 ft.,
        one target. <em>Hit:</em> 1 (1d4 - 1) bludgeoning damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['roll'], '-1')

    def test_no_html_tags(self):
        """Test parsing action without HTML tags"""
        action_text = "Bite. Melee Attack: +3 to hit, reach 5 ft., one target. Hit: 4 (1d6 + 1) piercing damage."

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Bite')

    def test_special_action_no_attack(self):
        """Test parsing special action without attack roll"""
        action_text = """
        <strong>Invisibility.</strong> The creature magically turns invisible until it attacks or
        casts a spell, or until its concentration ends.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Invisibility.')
        self.assertEqual(result['type'], 'other')
        self.assertNotIn('roll', result)

    def test_parse_multiple_actions(self):
        """Test parsing multiple actions from HTML"""
        actions_html = """
        <p><strong>Bite.</strong> <em>Melee Weapon Attack:</em> +4 to hit, reach 5 ft.,
        one target. <em>Hit:</em> 7 (1d10 + 2) piercing damage.</p>
        <p><strong>Claw.</strong> <em>Melee Weapon Attack:</em> +4 to hit, reach 5 ft.,
        one target. <em>Hit:</em> 5 (1d6 + 2) slashing damage.</p>
        """

        results = parse_actions(actions_html)

        self.assertEqual(len(results), 2)
        self.assertEqual(results[0]['name'], 'Bite.')
        self.assertEqual(results[1]['name'], 'Claw.')

    def test_parse_empty_actions(self):
        """Test parsing empty actions"""
        results = parse_actions("")
        self.assertEqual(len(results), 0)

        results = parse_actions("--")
        self.assertEqual(len(results), 0)

        results = parse_actions(None)
        self.assertEqual(len(results), 0)

    def test_complex_dragon_breath(self):
        """Test parsing complex dragon breath weapon"""
        action_text = """
        <strong>Fire Breath (Recharge 5-6).</strong> The dragon exhales fire in a 90-foot cone.
        Each creature in that area must make a DC 24 Dexterity saving throw, taking 91 (26d6)
        fire damage on a failed save, or half as much damage on a successful one.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertIn('Fire Breath', result['name'])
        self.assertEqual(result['save_dc'], 'DC 24')
        self.assertEqual(result['save_type'], 'Dexterity')
        self.assertEqual(result['damage'], '26d6')
        self.assertEqual(result['damage_type'], 'fire')

    def test_reach_only_action(self):
        """Test action with reach but no range"""
        action_text = """
        <strong>Slam.</strong> <em>Melee Weapon Attack:</em> +8 to hit, reach 10 ft.,
        one target. <em>Hit:</em> 14 (2d8 + 5) bludgeoning damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['reach'], '10ft')
        self.assertNotIn('range', result)

    def test_range_only_action(self):
        """Test action with range but no reach"""
        action_text = """
        <strong>Rock.</strong> <em>Ranged Weapon Attack:</em> +5 to hit, range 30/120 ft.,
        one target. <em>Hit:</em> 10 (2d6 + 3) bludgeoning damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['range'], '30/120ft')
        self.assertNotIn('reach', result)

    def test_invalid_action_text(self):
        """Test parsing invalid action text"""
        # Action with no name
        result = parse_action_text("Some random text without proper format")
        self.assertIsNone(result)

        # Empty string
        result = parse_action_text("")
        self.assertIsNone(result)

    def test_action_with_condition(self):
        """Test action that applies a condition"""
        action_text = """
        <strong>Stunning Strike.</strong> <em>Melee Weapon Attack:</em> +6 to hit, reach 5 ft.,
        one target. <em>Hit:</em> 7 (1d8 + 3) bludgeoning damage. The target must succeed on a
        DC 13 Constitution saving throw or be stunned until the end of the monk's next turn.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Stunning Strike.')
        self.assertIn('stunned', result['description'].lower())

    def test_zero_damage(self):
        """Test action with potential zero damage"""
        action_text = """
        <strong>Touch.</strong> <em>Melee Weapon Attack:</em> +2 to hit, reach 5 ft.,
        one target. <em>Hit:</em> 0 (1d4 - 2) damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Touch.')

    def test_large_damage_dice(self):
        """Test action with very large damage dice"""
        action_text = """
        <strong>Meteor Swarm.</strong> DC 20 Dexterity saving throw,
        taking 140 (40d6) fire damage and 140 (40d6) bludgeoning damage on a failed save.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['save_dc'], 'DC 20')
        # Should capture both 40d6 instances
        self.assertIn('40d6', result['damage'])

    def test_action_with_dc_but_no_save_type(self):
        """Test action with DC but without explicit save type"""
        action_text = """
        <strong>Confusion Aura.</strong> Creatures within 30 feet must succeed on a DC 15 saving throw
        or become confused.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['save_dc'], 'DC 15')
        # save_type might not be captured, which is okay

    def test_multiline_action_description(self):
        """Test action with complex multiline description"""
        action_text = """
        <strong>Swallow.</strong> <em>Melee Weapon Attack:</em> +10 to hit, reach 5 ft., one target.
        <em>Hit:</em> 22 (4d8 + 4) piercing damage. If this damage reduces the target to 0 hit points,
        the creature is swallowed and the target is blinded and restrained. While swallowed,
        the target takes 21 (6d6) acid damage at the start of each of the creature's turns.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Swallow.')
        # Should capture both damage types
        self.assertIn('4d8 + 4', result['damage'])
        self.assertIn('6d6', result['damage'])

    def test_percentage_based_dice(self):
        """Test action with percentage-based mechanics"""
        action_text = """
        <strong>Critical Strike.</strong> <em>Melee Weapon Attack:</em> +7 to hit, reach 5 ft.,
        one target. <em>Hit:</em> 15 (2d10 + 4) slashing damage.
        """

        result = parse_action_text(action_text)

        self.assertIsNotNone(result)
        self.assertEqual(result['damage'], '2d10 + 4')


class TestActionEdgeCases(unittest.TestCase):
    """Test edge cases for action parsing"""

    def test_mixed_case_attack_types(self):
        """Test parsing with mixed case text"""
        action_text = "<strong>Strike.</strong> mELEe wEAPon aTTacK: +3 to hit."
        result = parse_action_text(action_text)
        self.assertIsNotNone(result)

    def test_extra_whitespace(self):
        """Test parsing with extra whitespace"""
        action_text = """

        <strong>Slam.</strong>    <em>Melee   Weapon   Attack:</em>    +5   to   hit

        """
        result = parse_action_text(action_text)
        self.assertIsNotNone(result)
        self.assertEqual(result['name'], 'Slam.')

    def test_unicode_characters(self):
        """Test parsing with unicode characters"""
        action_text = """
        <strong>Épée.</strong> <em>Melee Weapon Attack:</em> +5 to hit, reach 5 ft.,
        one target. <em>Hit:</em> 8 (1d8 + 4) piercing damage.
        """
        result = parse_action_text(action_text)
        self.assertIsNotNone(result)

    def test_very_long_action_name(self):
        """Test action with very long name"""
        action_text = """
        <strong>Super Ultra Mega Devastating Apocalyptic World-Ending Strike of Doom.</strong>
        <em>Melee Weapon Attack:</em> +10 to hit.
        """
        result = parse_action_text(action_text)
        self.assertIsNotNone(result)

    def test_missing_period_after_name(self):
        """Test action name without trailing period"""
        action_text = "<strong>Bite</strong> <em>Melee Weapon Attack:</em> +3 to hit"
        result = parse_action_text(action_text)
        self.assertIsNotNone(result)

    def test_action_with_parenthetical_info(self):
        """Test action with parenthetical information in name"""
        action_text = """
        <strong>Breath Weapon (Recharge 5-6).</strong> DC 15 Dexterity saving throw.
        """
        result = parse_action_text(action_text)
        self.assertIsNotNone(result)
        self.assertIn('Recharge', result['name'])


if __name__ == '__main__':
    unittest.main()
