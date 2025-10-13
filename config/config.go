// config/config.go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds all application configuration
type Config struct {
	// UI Settings
	Theme             Theme            `json:"theme"`
	AutoSave          AutoSaveConfig   `json:"auto_save"`
	DiceRoller        DiceRollerConfig `json:"dice_roller"`
	InitiativeTracker InitiativeConfig `json:"initiative_tracker"`
	Display           DisplayConfig    `json:"display"`
	Paths             PathsConfig      `json:"paths"`
}

// Theme configuration
type Theme struct {
	PrimaryColor   string `json:"primary_color"`   // Main accent color
	BorderColor    string `json:"border_color"`    // Panel border color
	HighlightColor string `json:"highlight_color"` // Selected item color
	ErrorColor     string `json:"error_color"`     // Error message color
	SuccessColor   string `json:"success_color"`   // Success message color
}

// AutoSaveConfig for campaign auto-save settings
type AutoSaveConfig struct {
	Enabled         bool `json:"enabled"`          // Enable auto-save
	IntervalMinutes int  `json:"interval_minutes"` // Auto-save interval in minutes
}

// DiceRollerConfig for dice roller settings
type DiceRollerConfig struct {
	HistorySize        int    `json:"history_size"`        // Number of rolls to keep in history
	ShowIndividual     bool   `json:"show_individual"`     // Show individual die results
	MinimumValue       int    `json:"minimum_value"`       // Minimum roll value (D&D standard: 1)
	CriticalHitEnabled bool   `json:"critical_hit_enabled"` // Enable critical hit detection
	CriticalHitMode    string `json:"critical_hit_mode"`    // "double" or "max" - how to calculate crit damage
}

// InitiativeConfig for initiative tracker settings
type InitiativeConfig struct {
	AutoSort        bool `json:"auto_sort"`        // Auto-sort by initiative
	ShowHP          bool `json:"show_hp"`          // Show HP in list
	ShowAC          bool `json:"show_ac"`          // Show AC in list
	HighlightActive bool `json:"highlight_active"` // Highlight active turn
	RoundCounter    bool `json:"round_counter"`    // Show round counter
}

// DisplayConfig for display settings
type DisplayConfig struct {
	ShowHelpHints      bool `json:"show_help_hints"`     // Show help hints at bottom
	CompactMode        bool `json:"compact_mode"`        // Use compact display
	AnimateTransitions bool `json:"animate_transitions"` // Animate panel transitions
	LineWrap           bool `json:"line_wrap"`           // Wrap long text lines
	MaxLineLength      int  `json:"max_line_length"`     // Maximum line length before wrap
}

// PathsConfig for file and directory paths
type PathsConfig struct {
	SaveDirectory   string `json:"save_directory"`   // Directory for campaign saves (empty = default ~/.lazydnd)
	ConfigDirectory string `json:"config_directory"` // Config directory (empty = default ~/.config/lazydnd)
	BackupEnabled   bool   `json:"backup_enabled"`   // Enable automatic backups
	BackupDirectory string `json:"backup_directory"` // Directory for backups (empty = saves/.backups)
	MaxBackups      int    `json:"max_backups"`      // Maximum number of backups to keep per campaign
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Theme: Theme{
			PrimaryColor:   "#7D56F4",
			BorderColor:    "#444444",
			HighlightColor: "#00FF00",
			ErrorColor:     "#FF0000",
			SuccessColor:   "#00FF00",
		},
		AutoSave: AutoSaveConfig{
			Enabled:         true,
			IntervalMinutes: 5,
		},
		DiceRoller: DiceRollerConfig{
			HistorySize:        15,
			ShowIndividual:     true,
			MinimumValue:       1,
			CriticalHitEnabled: true,
			CriticalHitMode:    "double", // "double" or "max"
		},
		InitiativeTracker: InitiativeConfig{
			AutoSort:        true,
			ShowHP:          true,
			ShowAC:          true,
			HighlightActive: true,
			RoundCounter:    true,
		},
		Display: DisplayConfig{
			ShowHelpHints:      true,
			CompactMode:        false,
			AnimateTransitions: false,
			LineWrap:           true,
			MaxLineLength:      50,
		},
		Paths: PathsConfig{
			SaveDirectory:   "", // Empty = use default ~/.lazydnd
			ConfigDirectory: "", // Empty = use default ~/.config/lazydnd
			BackupEnabled:   true,
			BackupDirectory: "", // Empty = use saves/.backups
			MaxBackups:      10,
		},
	}
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "lazydnd")
	return configDir, nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// Load loads the configuration from file, or returns default if not found
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return default and save it
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := Default()
		if err := cfg.Save(); err != nil {
			return cfg, fmt.Errorf("failed to save default config: %w", err)
		}
		return cfg, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Parse JSON
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Validate checks if the configuration values are valid
func (c *Config) Validate() error {
	// Validate auto-save interval
	if c.AutoSave.IntervalMinutes < 1 {
		return fmt.Errorf("auto_save.interval_minutes must be at least 1")
	}

	// Validate dice roller history size
	if c.DiceRoller.HistorySize < 1 || c.DiceRoller.HistorySize > 100 {
		return fmt.Errorf("dice_roller.history_size must be between 1 and 100")
	}

	// Validate minimum value
	if c.DiceRoller.MinimumValue < 0 {
		return fmt.Errorf("dice_roller.minimum_value cannot be negative")
	}

	// Validate max line length
	if c.Display.MaxLineLength < 20 {
		return fmt.Errorf("display.max_line_length must be at least 20")
	}

	// Validate max backups
	if c.Paths.MaxBackups < 1 {
		return fmt.Errorf("paths.max_backups must be at least 1")
	}

	// Validate save directory if specified
	if c.Paths.SaveDirectory != "" {
		expandedPath := expandPath(c.Paths.SaveDirectory)
		if !filepath.IsAbs(expandedPath) {
			return fmt.Errorf("paths.save_directory must be an absolute path")
		}
	}

	// Validate backup directory if specified
	if c.Paths.BackupDirectory != "" {
		expandedPath := expandPath(c.Paths.BackupDirectory)
		if !filepath.IsAbs(expandedPath) {
			return fmt.Errorf("paths.backup_directory must be an absolute path")
		}
	}

	return nil
}

// Reset resets configuration to defaults
func (c *Config) Reset() {
	*c = *Default()
}

// GetSaveDirectory returns the save directory path (either custom or default)
func (c *Config) GetSaveDirectory() (string, error) {
	if c.Paths.SaveDirectory != "" {
		return expandPath(c.Paths.SaveDirectory), nil
	}

	// Use default ~/.lazydnd
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".lazydnd"), nil
}

// GetBackupDirectory returns the backup directory path
func (c *Config) GetBackupDirectory() (string, error) {
	if c.Paths.BackupDirectory != "" {
		return expandPath(c.Paths.BackupDirectory), nil
	}

	// Use default: saves/.backups
	saveDir, err := c.GetSaveDirectory()
	if err != nil {
		return "", err
	}

	return filepath.Join(saveDir, ".backups"), nil
}

// expandPath expands ~ to home directory and environment variables
func expandPath(path string) string {
	// Expand ~ to home directory
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			if len(path) == 1 {
				return homeDir
			}
			return filepath.Join(homeDir, path[1:])
		}
	}

	// Expand environment variables
	return os.ExpandEnv(path)
}

// EnsureDirectoriesExist creates all configured directories if they don't exist
func (c *Config) EnsureDirectoriesExist() error {
	// Create save directory
	saveDir, err := c.GetSaveDirectory()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	// Create backup directory if backups are enabled
	if c.Paths.BackupEnabled {
		backupDir, err := c.GetBackupDirectory()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}
	}

	return nil
}
