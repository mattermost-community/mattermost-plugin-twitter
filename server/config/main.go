package config

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

var (
	config atomic.Value
)

// Configuration captures the plugin's external configuration as exposed in the Mattermost server
// configuration, as well as values computed from the configuration. Any public fields will be
// deserialized from the Mattermost server configuration in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and the plugin
// configuration can change at any time, access to the configuration must be synchronized. The
// strategy used in this plugin is to guard a pointer to the configuration, and clone the entire
// struct whenever it changes. You may replace this with whatever strategy you choose.
//
// If you add non-reference types to your configuration struct, be sure to rewrite Clone as a deep
// copy appropriate for your types.
type Configuration struct {
	OAuthClientID     string
	OAuthClientSecret string
	EncryptionKey     string
}

// Clone shallow copies the configuration. Your implementation may require a deep copy if
// your configuration has reference types.
func (c *Configuration) Clone() *Configuration {
	var clone = *c
	return &clone
}

// GetConfig retrieves the active configuration.
func GetConfig() *Configuration {
	return config.Load().(*Configuration)
}

// SetConfig replaces the active configuration.
func SetConfig(c *Configuration) {
	config.Store(c)
}

// IsValid checks if all needed fields are set.
func (c *Configuration) IsValid() error {
	if c.OAuthClientID == "" {
		return errors.New("must have a twitter oauth client id")
	}

	if c.OAuthClientSecret == "" {
		return errors.New("must have a twitter oauth client secret")
	}

	if c.EncryptionKey == "" {
		return errors.New("must have an encryption key")
	}

	return nil
}

func (c *Configuration) Serialize() map[string]interface{} {
	out := make(map[string]interface{})
	b, _ := json.Marshal(c)
	_ = json.Unmarshal(b, &out)
	return out
}
