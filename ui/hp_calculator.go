// ui/hp_calculator.go
package ui

import "math"

// HPCalculator provides safe HP calculation methods with edge case handling
type HPCalculator struct{}

// ApplyDamage applies damage to an entry, handling temp HP and edge cases
// Returns the new HP and new TempHP values
func (c *HPCalculator) ApplyDamage(currentHP, maxHP, currentTempHP, damage int) (newHP, newTempHP int) {
	// Validate inputs
	if damage < 0 {
		damage = 0 // Negative damage is treated as 0
	}
	if currentHP < 0 {
		currentHP = 0
	}
	if currentTempHP < 0 {
		currentTempHP = 0
	}
	if maxHP < 1 {
		maxHP = 1
	}

	// Apply damage to temp HP first
	newTempHP = currentTempHP
	newHP = currentHP

	if newTempHP > 0 {
		if damage <= newTempHP {
			// All damage absorbed by temp HP
			newTempHP -= damage
		} else {
			// Temp HP absorbs some, rest goes to real HP
			remainingDamage := damage - newTempHP
			newTempHP = 0
			newHP -= remainingDamage
			if newHP < 0 {
				newHP = 0
			}
		}
	} else {
		// No temp HP, damage goes directly to real HP
		newHP -= damage
		if newHP < 0 {
			newHP = 0
		}
	}

	return newHP, newTempHP
}

// ApplyHealing applies healing to an entry, capping at max HP
// Returns the new HP value and the amount actually healed
func (c *HPCalculator) ApplyHealing(currentHP, maxHP, healing int) (newHP, actualHealed int) {
	// Validate inputs
	if healing < 0 {
		healing = 0 // Negative healing is treated as 0
	}
	if currentHP < 0 {
		currentHP = 0
	}
	if maxHP < 1 {
		maxHP = 1
	}

	// Cap current HP if it somehow exceeds max
	if currentHP > maxHP {
		currentHP = maxHP
	}

	// Apply healing
	newHP = currentHP + healing

	// Cap at max HP
	if newHP > maxHP {
		actualHealed = maxHP - currentHP
		newHP = maxHP
	} else {
		actualHealed = healing
	}

	return newHP, actualHealed
}

// SetTempHP sets temporary HP, replacing any existing temp HP
// Returns the validated temp HP value
func (c *HPCalculator) SetTempHP(tempHP int) int {
	if tempHP < 0 {
		return 0
	}
	if tempHP > MaxHPValue {
		return MaxHPValue
	}
	return tempHP
}

// SetMaxHP sets maximum HP and adjusts current HP if needed
// Returns the new max HP and new current HP
func (c *HPCalculator) SetMaxHP(currentHP, newMaxHP int) (adjustedMaxHP, adjustedCurrentHP int) {
	// Ensure max HP is at least 1
	if newMaxHP < 1 {
		newMaxHP = 1
	}

	// Cap max HP at system limit
	if newMaxHP > MaxHPValue {
		newMaxHP = MaxHPValue
	}

	// Cap current HP if it exceeds new max
	adjustedCurrentHP = currentHP
	if adjustedCurrentHP > newMaxHP {
		adjustedCurrentHP = newMaxHP
	}
	if adjustedCurrentHP < 0 {
		adjustedCurrentHP = 0
	}

	return newMaxHP, adjustedCurrentHP
}

// ValidateHP ensures HP is within valid bounds
// Returns the validated HP value
func (c *HPCalculator) ValidateHP(hp, maxHP int) int {
	if hp < 0 {
		return 0
	}
	if hp > maxHP {
		return maxHP
	}
	if hp > MaxHPValue {
		return MaxHPValue
	}
	return hp
}

// CalculateHPChange calculates the result of applying an HP change
// Positive change = healing, negative change = damage
// Returns new HP, new temp HP, and a description of what happened
func (c *HPCalculator) CalculateHPChange(currentHP, maxHP, currentTempHP, change int) (newHP, newTempHP int, description string) {
	if change < 0 {
		// Damage
		damage := -change
		newHP, newTempHP = c.ApplyDamage(currentHP, maxHP, currentTempHP, damage)

		if currentTempHP > 0 {
			realHPLost := currentHP - newHP

			if realHPLost > 0 {
				description = "Temp HP absorbed some damage, real HP also reduced"
			} else {
				description = "All damage absorbed by temp HP"
			}
		} else {
			description = "Damage applied to HP"
		}
	} else if change > 0 {
		// Healing
		var actualHealed int
		newHP, actualHealed = c.ApplyHealing(currentHP, maxHP, change)
		newTempHP = currentTempHP

		if actualHealed < change {
			description = "Healed to maximum HP (some healing wasted)"
		} else {
			description = "Healed successfully"
		}
	} else {
		// No change
		newHP = currentHP
		newTempHP = currentTempHP
		description = "No change"
	}

	return newHP, newTempHP, description
}

// SafeAddHP adds to HP with overflow protection
// Returns the result, capped at MaxHPValue
func (c *HPCalculator) SafeAddHP(a, b int) int {
	// Check for overflow
	if b > 0 && a > math.MaxInt32-b {
		return MaxHPValue
	}
	if b < 0 && a < math.MinInt32-b {
		return 0
	}

	result := a + b
	if result > MaxHPValue {
		return MaxHPValue
	}
	if result < 0 {
		return 0
	}
	return result
}

// SafeSubtractHP subtracts from HP with underflow protection
// Returns the result, minimum 0
func (c *HPCalculator) SafeSubtractHP(a, b int) int {
	if b > a {
		return 0
	}
	result := a - b
	if result < 0 {
		return 0
	}
	return result
}

// GetHPPercentage calculates HP as a percentage of max HP
// Returns 0-100, where 100 = full HP, 0 = no HP
func (c *HPCalculator) GetHPPercentage(currentHP, maxHP int) int {
	if maxHP <= 0 {
		return 0
	}
	if currentHP <= 0 {
		return 0
	}
	if currentHP >= maxHP {
		return 100
	}

	// Calculate percentage
	percentage := (currentHP * 100) / maxHP
	if percentage > 100 {
		return 100
	}
	if percentage < 0 {
		return 0
	}
	return percentage
}

// IsConscious returns true if the creature is conscious (HP > 0)
func (c *HPCalculator) IsConscious(hp int) bool {
	return hp > 0
}

// IsBlooded returns true if HP is below 50%
func (c *HPCalculator) IsBlooded(currentHP, maxHP int) bool {
	percentage := c.GetHPPercentage(currentHP, maxHP)
	return percentage < 50
}

// IsCritical returns true if HP is below 25%
func (c *HPCalculator) IsCritical(currentHP, maxHP int) bool {
	percentage := c.GetHPPercentage(currentHP, maxHP)
	return percentage < 25
}

// Global HP calculator instance
var HPCalc = &HPCalculator{}

