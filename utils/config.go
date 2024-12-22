// utils/config.go
package utils

import (
	"encoding/json"
	"fmt"
	"listarr-backend/models"
	"os"
	"strings"
	"sync"

	"github.com/knadh/koanf/parsers/dotenv"
	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	config     *models.Configuration
	configLock sync.RWMutex
	k          *koanf.Koanf
)

func InitConfig() error {
	k = koanf.New(".")

	// 1. Load defaults
	if err := k.Load(confmap.Provider(defaultConfig, "."), nil); err != nil {
		return fmt.Errorf("error loading defaults: %w", err)
	}

	if err := os.MkdirAll("./config", 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// 2. Load app.config.json
	f := file.Provider("config/app.config.json")

	if err := k.Load(f, kjson.Parser()); err != nil {
		// Only create default config if file doesn't exist
		if os.IsNotExist(err) {
			if err := saveDefaultConfig(); err != nil {
				return fmt.Errorf("error saving default config: %w", err)
			}
		} else {
			return fmt.Errorf("error loading config file: %w", err)
		}
	}

	fmt.Printf("After app.config.json - Database User: %s\n", k.String("db.user"))

	f.Watch(func(event interface{}, err error) {
		if err != nil {
			fmt.Printf("watch error: %v\n", err)
			return
		}

		configLock.Lock()
		defer configLock.Unlock()

		// Create a new instance and reload
		k = koanf.New(".")

		// Reload in the correct order
		k.Load(confmap.Provider(defaultConfig, "."), nil)
		k.Load(f, kjson.Parser())
		k.Load(file.Provider(".env"), dotenv.Parser())
		k.Load(env.Provider("LISTARR_", ".", envKeyReplacer), nil)

		// Update the config struct
		if err := k.Unmarshal("", config); err != nil {
			fmt.Printf("error unmarshaling config: %v\n", err)
			return
		}

		fmt.Println("Configuration reloaded due to file change")
	})

	// 4. Load environment variables (highest priority)
	fmt.Println("Loading environment variables")
	if err := k.Load(env.Provider("LISTARR_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "LISTARR_")), "_", ".")
	}), nil); err != nil {
		return fmt.Errorf("error loading environment variables: %w", err)
	}

	// Load the final config into our struct
	configLock.Lock()
	defer configLock.Unlock()

	config = &models.Configuration{}
	if err := k.UnmarshalWithConf("", config, koanf.UnmarshalConf{
		Tag: "json", // Use json tags from your struct
	}); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

func envKeyReplacer(s string) string {
	return strings.ReplaceAll(
		strings.ToLower(
			strings.TrimPrefix(s, "LISTARR_")),
		"_",
		".",
	)
}

var defaultConfig = map[string]interface{}{
	// App defaults
	"app.name":        "Listarr",
	"app.environment": "development",
	"app.appURL":      "http://localhost:3000",
	"app.apiBaseURL":  "http://localhost:8080",
	"app.logLevel":    "info",
	"app.maxPageSize": 100,

	// Database defaults
	"db.host":     "localhost",
	"db.port":     "5432",
	"db.name":     "listarr",
	"db.user":     "postgres_user",
	"db.password": "yourpassword",
	"db.maxConns": 20,
	"db.timeout":  30,

	// HTTP defaults
	"http.port":             "8080",
	"http.readTimeout":      30,
	"http.writeTimeout":     30,
	"http.idleTimeout":      60,
	"http.enableSSL":        false,
	"http.rateLimitEnabled": true,
	"http.requestsPerMin":   100,

	// Auth defaults
	"auth.enableLocal":     true,
	"auth.sessionTimeout":  60,
	"auth.enable2FA":       false,
	"auth.tokenExpiration": 24,
	"auth.allowedOrigins":  []string{"http://localhost:3000"},

	// Sync defaults
	"sync.enabled":          true,
	"sync.interval":         "0 */12 * * *",
	"sync.conflictStrategy": "skip",
	"sync.playlists": map[string]interface{}{
		"enableSync":   true,
		"syncInterval": "0 */6 * * *",
		"allowedTypes": []string{"music", "media"},
		"maxItems":     1000,
	},
	"sync.collections": map[string]interface{}{
		"enableSync":   true,
		"syncInterval": "0 */12 * * *",
		"allowedTypes": []string{"series", "movies", "music"},
		"maxItems":     5000,
	},

	// SpotDL defaults
	"spotdl.enabled":          false,
	"spotdl.downloadDir":      "./downloads",
	"spotdl.fileFormat":       "mp3",
	"spotdl.qualityPreset":    "high",
	"spotdl.namingTemplate":   "{artist} - {title}",
	"spotdl.maxRetries":       3,
	"spotdl.concurrentLimit":  2,
	"spotdl.notifyOnComplete": true,

	// Integrations defaults
	"integrations.emby": map[string]interface{}{
		"enabled": false,
		"host":    "localhost",
		"port":    8096,
		"ssl":     false,
	},
	"integrations.jellyfin": map[string]interface{}{
		"enabled": false,
		"host":    "localhost",
		"port":    8096,
		"ssl":     false,
	},
	"integrations.plex": map[string]interface{}{
		"enabled": false,
		"host":    "localhost",
		"port":    32400,
		"ssl":     false,
	},
	"integrations.navidrome": map[string]interface{}{
		"enabled": false,
		"host":    "localhost",
		"port":    4533,
		"ssl":     false,
	},
	"integrations.spotify": map[string]interface{}{
		"enabled": false,
		"scopes":  "user-library-read playlist-read-private",
	},
}

func saveDefaultConfig() error {
	if err := os.MkdirAll("./config", 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(k.Raw(), "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.WriteFile("config/app.config.json", data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func GetConfig() *models.Configuration {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func SaveConfig(cfg models.Configuration) error {
	configLock.Lock()
	defer configLock.Unlock()

	// Convert config struct to map
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &configMap); err != nil {
		return fmt.Errorf("error unmarshaling to map: %w", err)
	}

	// Load the new config
	if err := k.Load(confmap.Provider(configMap, "."), nil); err != nil {
		return fmt.Errorf("error loading new config: %w", err)
	}

	// Save to file
	if err := saveDefaultConfig(); err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}

	config = &cfg
	return nil
}

// GetFileConfig returns only the configuration from app.config.json
func GetFileConfig() *models.Configuration {
	// Create a new koanf instance just for the file
	k := koanf.New(".")

	// Load only the file configuration
	if err := k.Load(file.Provider("config/app.config.json"), kjson.Parser()); err != nil {
		return nil
	}

	config := &models.Configuration{}
	if err := k.Unmarshal("", config); err != nil {
		return nil
	}

	return config
}

// utils/config.go
func SaveFileConfig(cfg models.Configuration) error {
	configLock.Lock()
	defer configLock.Unlock()

	// Convert config struct to map for Koanf
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &configMap); err != nil {
		return fmt.Errorf("error unmarshaling to map: %w", err)
	}

	// Create a new Koanf instance just for the file config
	fileK := koanf.New(".")
	if err := fileK.Load(confmap.Provider(configMap, "."), nil); err != nil {
		return fmt.Errorf("error loading new config: %w", err)
	}

	// Write to file
	data, err := json.MarshalIndent(fileK.Raw(), "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.WriteFile("config/app.config.json", data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	// The watch functionality will automatically reload the config
	return nil
}

// ResetFileConfig resets app.config.json to default values
func ResetFileConfig() error {
	// Create a new koanf instance with just defaults
	k := koanf.New(".")
	if err := k.Load(confmap.Provider(defaultConfig, "."), nil); err != nil {
		return fmt.Errorf("error loading defaults: %w", err)
	}

	// Save defaults to file
	data, err := json.MarshalIndent(k.Raw(), "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.WriteFile("config/app.config.json", data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	// Reload the main configuration
	return InitConfig()
}
