package glide

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config is a sealed interface representing a parsed strategy configuration.
// Implementations are flatConfig and hierarchyConfig.
type Config interface {
	GetStrategy() Strategy
}

// flatConfig holds configuration for the flat placement strategy.
type flatConfig struct {
	Separator string `toml:"separator"`
}

// GetStrategy returns a flatStrategy configured with this config.
func (c flatConfig) GetStrategy() Strategy {
	sep := c.Separator
	if sep == "" {
		sep = "-"
	}
	return flatStrategy{separator: sep}
}

// hierarchyConfig holds configuration for the hierarchy placement strategy.
type hierarchyConfig struct {
	Root string `toml:"root"`
}

// GetStrategy returns a hierarchyStrategy configured with this config.
func (c hierarchyConfig) GetStrategy() Strategy {
	return hierarchyStrategy{root: c.Root}
}

// compositeConfig is the direct representation of a config file.
type compositeConfig struct {
	Strategy  string          `toml:"strategy"`
	Flat      flatConfig      `toml:"flat"`
	Hierarchy hierarchyConfig `toml:"hierarchy"`
}

// resolve converts compositeConfig into a typed Config based on the strategy field.
func (r compositeConfig) resolve() (Config, error) {
	switch r.Strategy {
	case "", "flat":
		return r.Flat, nil
	case "hierarchy":
		return r.Hierarchy, nil
	default:
		return nil, fmt.Errorf("Error: unknown strategy '%s'", r.Strategy)
	}
}

// LoadConfig loads the effective strategy configuration for the given repository root.
// It first checks for a local config at <root>/.glide/config, falling back to
// the global config at $XDG_CONFIG_HOME/glide/config.
// Returns a flat strategy with defaults if no config file is found.
func LoadConfig(root string) (Config, error) {
	localPath := filepath.Join(root, ".glide", "config")
	if _, err := os.Stat(localPath); err == nil {
		local, err := loadConfigFile(localPath)
		if err != nil {
			return nil, err
		}
		if local.Strategy != "" {
			return local.resolve()
		}
	}

	globalPath := filepath.Join(globalConfigRootDir(), "glide", "config")
	if _, err := os.Stat(globalPath); err == nil {
		global, err := loadConfigFile(globalPath)
		if err != nil {
			return nil, err
		}
		return global.resolve()
	}

	return flatConfig{}, nil
}

// globalConfigRootDir returns the root directory for global config files,
// respecting $XDG_CONFIG_HOME or falling back to ~/.config.
func globalConfigRootDir() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base != "" {
		return base
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}

// loadConfigFile parses a TOML config file at the given path into a compositeConfig.
func loadConfigFile(path string) (compositeConfig, error) {
	var c compositeConfig
	if _, err := toml.DecodeFile(path, &c); err != nil {
		return compositeConfig{}, err
	}
	return c, nil
}
