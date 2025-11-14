# Lipgloss Removal Status

## Completed
- ✅ Phase 1: Styles system converted to TView colors
- ✅ Phase 2: Status bar updated (removed lipgloss)
- ✅ Phase 3: Layout rendering functions deprecated

## In Progress
- 🔄 Phase 5: Panel files - dice_roller.go partially done (needs syntax fixes)

## Remaining Work

### Panel Files (Phase 5)
- `panels/dice_roller.go` - Needs syntax fixes for style.Render() replacements
- `panels/spells.go` - Remove lipgloss imports and styles
- `panels/monsters.go` - Remove lipgloss imports and styles
- `panels/initiative_tracker.go` - Remove lipgloss imports and styles
- `panels/encounter_builder.go` - Remove lipgloss imports and styles
- `panels/search_utils.go` - Remove lipgloss imports and styles

### Popup Files (Phase 4)
- `ui/save_popup.go`
- `ui/help_popup.go`
- `ui/action_popup.go`
- `ui/condition_popup.go`
- `ui/cast_spell_popup.go`
- `ui/saving_throw_popup.go`
- `ui/quick_hp_popup.go`
- `ui/multi_target_popup.go`
- `ui/encounter_prompt_popup.go`
- `ui/encounter_generator_popup.go`

### Other
- Phase 6: Update main entry points
- Phase 7: Update TView app with modal management
- Phase 8: Remove dependencies from go.mod
- Phase 9: Update utility functions

## Notes
- All `.Render()` calls need to be replaced with plain text
- Style variables should be removed
- TView widgets handle all styling
