// models/config.go
package models

// Configuration represents the complete application configuration
// @Description Complete application configuration settings
type Configuration struct {
	// App contains core application settings
	App struct {
		Name        string `json:"name" mapstructure:"name" example:"Listarr" binding:"required"`
		Environment string `json:"environment" mapstructure:"environment" example:"development" binding:"required,oneof=development staging production"`
		AppURL      string `json:"appURL" mapstructure:"appURL" example:"http://localhost:3000" binding:"required,url"`
		APIBaseURL  string `json:"apiBaseURL" mapstructure:"apiBaseURL" example:"http://localhost:8080" binding:"required,url"`
		LogLevel    string `json:"logLevel" mapstructure:"logLevel" example:"info" binding:"required,oneof=debug info warn error"`
		MaxPageSize int    `json:"maxPageSize" mapstructure:"maxPageSize" example:"100" binding:"required,min=1,max=1000"`
	} `json:"app"`

	// Database contains database connection settings
	Db struct {
		Host     string `json:"host" mapstructure:"url" example:"localhost" binding:"required"`
		Port     string `json:"port" mapstructure:"port" example:"5432" binding:"required"`
		Name     string `json:"name" mapstructure:"name" example:"listarr" binding:"required"`
		User     string `json:"user" mapstructure:"user" example:"postgres_user" binding:"required"`
		Password string `json:"password" mapstructure:"password" example:"yourpassword" binding:"required"`
		MaxConns int    `json:"maxConns" mapstructure:"maxConns" example:"20" binding:"required,min=1"`
		Timeout  int    `json:"timeout" mapstructure:"timeout" example:"30" binding:"required,min=1"`
	} `json:"db" mapstructure:"db"`

	// HTTP contains HTTP server configuration
	HTTP struct {
		Port             string `json:"port" mapstructure:"port" example:"8080" binding:"required"`
		ReadTimeout      int    `json:"readTimeout" mapstructure:"readTimeout" example:"30" binding:"required,min=1"`
		WriteTimeout     int    `json:"writeTimeout" mapstructure:"writeTimeout" example:"30" binding:"required,min=1"`
		IdleTimeout      int    `json:"idleTimeout" mapstructure:"idleTimeout" example:"60" binding:"required,min=1"`
		EnableSSL        bool   `json:"enableSSL" mapstructure:"enableSSL" example:"false"`
		SSLCert          string `json:"sslCert" mapstructure:"sslCert" example:"/path/to/cert.pem"`
		SSLKey           string `json:"sslKey" mapstructure:"sslKey" example:"/path/to/key.pem"`
		ProxyEnabled     bool   `json:"proxyEnabled" mapstructure:"proxyEnabled" example:"false"`
		ProxyURL         string `json:"proxyURL" mapstructure:"proxyURL" example:"http://proxy:8080"`
		RateLimitEnabled bool   `json:"rateLimitEnabled" mapstructure:"rateLimitEnabled" example:"true"`
		RequestsPerMin   int    `json:"requestsPerMin" mapstructure:"requestsPerMin" example:"100" binding:"min=0"`
	} `json:"http"`

	// Auth contains authentication settings
	Auth struct {
		EnableLocal     bool     `json:"enableLocal" mapstructure:"enableLocal" example:"true"`
		SessionTimeout  int      `json:"sessionTimeout" mapstructure:"sessionTimeout" example:"60" binding:"required,min=1"`
		Enable2FA       bool     `json:"enable2FA" mapstructure:"enable2FA" example:"false"`
		JWTSecret       string   `json:"jwtSecret" mapstructure:"jwtSecret" example:"your-secret-key" binding:"required"`
		TokenExpiration int      `json:"tokenExpiration" mapstructure:"tokenExpiration" example:"24" binding:"required,min=1"`
		AllowedOrigins  []string `json:"allowedOrigins" mapstructure:"allowedOrigins" example:"http://localhost:3000"`
	} `json:"auth"`

	// Integrations contains all third-party service configurations
	Integrations struct {
		Emby      EmbyConfig      `json:"emby" mapstructure:"emby"`
		Jellyfin  JellyfinConfig  `json:"jellyfin" mapstructure:"jellyfin"`
		Plex      PlexConfig      `json:"plex" mapstructure:"plex"`
		Trakt     TraktConfig     `json:"trakt" mapstructure:"trakt"`
		Navidrome NavidromeConfig `json:"navidrome" mapstructure:"navidrome"`
		Spotify   SpotifyConfig   `json:"spotify" mapstructure:"spotify"`
	} `json:"integrations"`

	// Sync contains synchronization settings
	Sync struct {
		Enabled          bool   `json:"enabled" mapstructure:"enabled" example:"true"`
		Interval         string `json:"interval" mapstructure:"interval" example:"0 */12 * * *" binding:"required"`
		ConflictStrategy string `json:"conflictStrategy" mapstructure:"conflictStrategy" example:"skip" binding:"required,oneof=overwrite skip merge"`

		Playlists struct {
			EnableSync   bool     `json:"enableSync" mapstructure:"enableSync" example:"true"`
			SyncInterval string   `json:"syncInterval" mapstructure:"syncInterval" example:"0 */6 * * *"`
			AllowedTypes []string `json:"allowedTypes" mapstructure:"allowedTypes" example:"music,media"`
			MaxItems     int      `json:"maxItems" mapstructure:"maxItems" example:"1000" binding:"required,min=1"`
		} `json:"playlists"`

		Collections struct {
			EnableSync   bool     `json:"enableSync" mapstructure:"enableSync" example:"true"`
			SyncInterval string   `json:"syncInterval" mapstructure:"syncInterval" example:"0 */12 * * *"`
			AllowedTypes []string `json:"allowedTypes" mapstructure:"allowedTypes" example:"series,movies,music"`
			MaxItems     int      `json:"maxItems" mapstructure:"maxItems" example:"5000" binding:"required,min=1"`
		} `json:"collections"`
	} `json:"sync"`

	// SpotDL contains Spotify download integration settings
	SpotDL struct {
		Enabled          bool   `json:"enabled" mapstructure:"enabled" example:"false"`
		DownloadDir      string `json:"downloadDirectory" mapstructure:"downloadDirectory" example:"./downloads" binding:"required"`
		FileFormat       string `json:"fileFormat" mapstructure:"fileFormat" example:"mp3" binding:"required,oneof=mp3 flac"`
		QualityPreset    string `json:"qualityPreset" mapstructure:"qualityPreset" example:"high" binding:"required,oneof=low medium high"`
		NamingTemplate   string `json:"namingTemplate" mapstructure:"namingTemplate" example:"{artist} - {title}" binding:"required"`
		MaxRetries       int    `json:"maxRetries" mapstructure:"maxRetries" example:"3" binding:"required,min=0"`
		ConcurrentLimit  int    `json:"concurrentDownloads" mapstructure:"concurrentDownloads" example:"2" binding:"required,min=1"`
		NotifyOnComplete bool   `json:"notifyOnComplete" mapstructure:"notifyOnComplete" example:"true"`
	} `json:"spotdl"`
}

// Integration config types
// @Description Emby media server configuration
type EmbyConfig struct {
	Enabled  bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	Host     string `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port     int    `json:"port" mapstructure:"port" example:"8096" binding:"required_if=Enabled true"`
	APIKey   string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Username string `json:"username" mapstructure:"username" example:"admin"`
	SSL      bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Jellyfin media server configuration
type JellyfinConfig struct {
	Enabled  bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	Host     string `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port     int    `json:"port" mapstructure:"port" example:"8096" binding:"required_if=Enabled true"`
	APIKey   string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Username string `json:"username" mapstructure:"username" example:"admin"`
	SSL      bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Plex media server configuration
type PlexConfig struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	Host    string `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port    int    `json:"port" mapstructure:"port" example:"32400" binding:"required_if=Enabled true"`
	Token   string `json:"token" mapstructure:"token" example:"your-plex-token" binding:"required_if=Enabled true"`
	SSL     bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Trakt.tv configuration
type TraktConfig struct {
	Enabled      bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	ClientID     string `json:"clientId" mapstructure:"clientId" example:"your-client-id" binding:"required_if=Enabled true"`
	ClientSecret string `json:"clientSecret" mapstructure:"clientSecret" example:"your-client-secret" binding:"required_if=Enabled true"`
	RedirectURI  string `json:"redirectUri" mapstructure:"redirectUri" example:"http://localhost:8080/callback" binding:"required_if=Enabled true"`
}

// @Description Navidrome music server configuration
type NavidromeConfig struct {
	Enabled  bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	Host     string `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port     int    `json:"port" mapstructure:"port" example:"4533" binding:"required_if=Enabled true"`
	Username string `json:"username" mapstructure:"username" example:"admin" binding:"required_if=Enabled true"`
	Password string `json:"password" mapstructure:"password" example:"your-password" binding:"required_if=Enabled true"`
	SSL      bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Spotify configuration
type SpotifyConfig struct {
	Enabled      bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	ClientID     string `json:"clientId" mapstructure:"clientId" example:"your-client-id" binding:"required_if=Enabled true"`
	ClientSecret string `json:"clientSecret" mapstructure:"clientSecret" example:"your-client-secret" binding:"required_if=Enabled true"`
	RedirectURI  string `json:"redirectUri" mapstructure:"redirectUri" example:"http://localhost:8080/callback" binding:"required_if=Enabled true"`
	Scopes       string `json:"scopes" mapstructure:"scopes" example:"user-library-read playlist-read-private"`
}

// ConfigResponse represents the response structure for configuration endpoints
// @Description Configuration response wrapper
type ConfigResponse struct {
	Data  *Configuration `json:"data,omitempty"`
	Error string         `json:"error,omitempty"`
}
