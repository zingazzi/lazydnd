// panels/spells.go
package panels

// GetSpellsContent returns the content for the spells panel
func GetSpellsContent() string {
	return `✨ SPELLS ✨

Cantrips (At Will):
• Firebolt (1d10 fire damage)
  Range: 120 ft, Action
• Mage Hand
  Range: 30 ft, Action
• Prestidigitation
  Range: 10 ft, Action

Level 1 Spells (3/3 slots):
• Magic Missile
  3 darts, 1d4+1 force each
• Shield
  +5 AC until start of next turn
• Detect Magic
  Sense magic within 30 ft

Level 2 Spells (2/2 slots):
• Misty Step
  Teleport 30 ft as bonus action
• Scorching Ray
  3 ranged spell attacks, 2d6 fire each

Level 3 Spells (1/1 slots):
• Fireball
  8d6 fire damage, 20 ft radius
• Counterspell
  Stop a spell being cast

Spell Attack Bonus: +7
Spell Save DC: 15
Spellcasting Ability: Intelligence`
}
