package configmanager

import (
	"os"
	"sync"

	toml "github.com/pelletier/go-toml/v2"
)

type Manager struct {
	mu       sync.RWMutex
	filePath string
	config   FrpsConfig
}

func NewManager(filePath string) (*Manager, error) {
	m := &Manager{filePath: filePath}
	if err := m.Load(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) FilePath() string {
	return m.filePath
}

func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = FrpsConfig{}
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			m.config = defaultConfig()
			return nil
		}
		return err
	}
	if err := toml.Unmarshal(data, &m.config); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Get() FrpsConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

func (m *Manager) Update(cfg FrpsConfig) error {
	m.mu.Lock()
	m.config = cfg
	m.mu.Unlock()
	return m.Save()
}

func (m *Manager) Save() error {
	m.mu.RLock()
	cfg := m.config
	m.mu.RUnlock()

	cfg = normalizeConfig(cfg)
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0644)
}

func isTLSEmpty(t *TransportTLS) bool {
	return t == nil || (!boolPtrVal(t.Force) && t.CertFile == "" && t.KeyFile == "" && t.TrustedCaFile == "")
}

func isQUICEmpty(q *TransportQUIC) bool {
	return q == nil || (q.KeepalivePeriod == nil && q.MaxIdleTimeout == nil && q.MaxIncomingStreams == nil)
}

func isOIDCEmpty(o *AuthOIDC) bool {
	return o == nil || (o.Issuer == "" && o.Audience == "" && !boolPtrVal(o.SkipExpiryCheck) && !boolPtrVal(o.SkipIssuerCheck))
}

func isWebServerTLSEmpty(t *WebServerTLS) bool {
	return t == nil || (t.CertFile == "" && t.KeyFile == "")
}

func isSSHGatewayEmpty(g *SSHTunnelGateway) bool {
	return g == nil || (g.BindPort == nil && g.PrivateKeyFile == "" && g.AutoGenPrivateKeyPath == "" && g.AuthorizedKeysFile == "")
}

func boolPtrVal(b *bool) bool {
	return b != nil && *b
}

// normalizeConfig drops empty sub-structs so the written TOML does not contain
// empty sections (e.g. [webServer.tls]) that would trigger validation errors
// when frps loads the file.
func normalizeConfig(cfg FrpsConfig) FrpsConfig {
	if isTLSEmpty(cfg.Transport.TLS) {
		cfg.Transport.TLS = nil
	}
	if isQUICEmpty(cfg.Transport.QUIC) {
		cfg.Transport.QUIC = nil
	}
	if isOIDCEmpty(cfg.Auth.OIDC) {
		cfg.Auth.OIDC = nil
	}
	// An empty tokenSource (type not set to 'file' or 'exec') must not be
	// persisted, otherwise frps rejects the config on next start.
	if cfg.Auth.TokenSource == nil || cfg.Auth.TokenSource.Type == "" {
		cfg.Auth.TokenSource = nil
	}
	if cfg.Auth.TokenSource != nil && cfg.Auth.TokenSource.File != nil && cfg.Auth.TokenSource.File.Path == "" {
		cfg.Auth.TokenSource.File = nil
	}
	if isWebServerTLSEmpty(cfg.WebServer.TLS) {
		cfg.WebServer.TLS = nil
	}
	if isSSHGatewayEmpty(cfg.SSHTunnelGateway) {
		cfg.SSHTunnelGateway = nil
	}
	return cfg
}

func intPtr(v int) *int    { return &v }
func boolPtr(v bool) *bool { return &v }

func defaultConfig() FrpsConfig {
	return FrpsConfig{
		BindAddr: "0.0.0.0",
		BindPort: intPtr(7000),
		Transport: Transport{
			MaxPoolCount: intPtr(5),
		},
		VhostHTTPPort:  intPtr(80),
		VhostHTTPSPort: intPtr(443),
		WebServer: WebServer{
			Addr:     "127.0.0.1",
			Port:     7500,
			User:     "admin",
			Password: "admin",
		},
		Log: Log{
			To:    "./frps.log",
			Level: "info",
		},
		DetailedErrorsToClient: boolPtr(true),
		Auth: Auth{
			Method: "token",
			Token:  "12345678",
		},
	}
}
