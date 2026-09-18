package configmanager

type AllowPort struct {
	Start  *int `toml:"start,omitempty" json:"start,omitempty"`
	End    *int `toml:"end,omitempty" json:"end,omitempty"`
	Single *int `toml:"single,omitempty" json:"single,omitempty"`
}

type TransportTLS struct {
	Force         *bool  `toml:"force,omitempty" json:"force,omitempty"`
	CertFile      string `toml:"certFile,omitempty" json:"certFile,omitempty"`
	KeyFile       string `toml:"keyFile,omitempty" json:"keyFile,omitempty"`
	TrustedCaFile string `toml:"trustedCaFile,omitempty" json:"trustedCaFile,omitempty"`
}

type TransportQUIC struct {
	KeepalivePeriod    *int `toml:"keepalivePeriod,omitempty" json:"keepalivePeriod,omitempty"`
	MaxIdleTimeout     *int `toml:"maxIdleTimeout,omitempty" json:"maxIdleTimeout,omitempty"`
	MaxIncomingStreams *int `toml:"maxIncomingStreams,omitempty" json:"maxIncomingStreams,omitempty"`
}

type Transport struct {
	TLS                     *TransportTLS  `toml:"tls,omitempty" json:"tls,omitempty"`
	QUIC                    *TransportQUIC `toml:"quic,omitempty" json:"quic,omitempty"`
	MaxPoolCount            *int           `toml:"maxPoolCount,omitempty" json:"maxPoolCount,omitempty"`
	HeartbeatTimeout        *int           `toml:"heartbeatTimeout,omitempty" json:"heartbeatTimeout,omitempty"`
	TCPMux                  *bool          `toml:"tcpMux,omitempty" json:"tcpMux,omitempty"`
	TCPMuxKeepaliveInterval *int           `toml:"tcpMuxKeepaliveInterval,omitempty" json:"tcpMuxKeepaliveInterval,omitempty"`
	TCPKeepalive            *int           `toml:"tcpKeepalive,omitempty" json:"tcpKeepalive,omitempty"`
}

type AuthOIDC struct {
	Issuer          string `toml:"issuer,omitempty" json:"issuer,omitempty"`
	Audience        string `toml:"audience,omitempty" json:"audience,omitempty"`
	SkipExpiryCheck *bool  `toml:"skipExpiryCheck,omitempty" json:"skipExpiryCheck,omitempty"`
	SkipIssuerCheck *bool  `toml:"skipIssuerCheck,omitempty" json:"skipIssuerCheck,omitempty"`
}

type AuthTokenSource struct {
	Type string         `toml:"type" json:"type"`
	File *AuthTokenFile `toml:"file,omitempty" json:"file,omitempty"`
}

type AuthTokenFile struct {
	Path string `toml:"path" json:"path"`
}

type Auth struct {
	Method           string           `toml:"method" json:"method"`
	Token            string           `toml:"token,omitempty" json:"token,omitempty"`
	TokenSource      *AuthTokenSource `toml:"tokenSource,omitempty" json:"tokenSource,omitempty"`
	AdditionalScopes []string         `toml:"additionalScopes,omitempty" json:"additionalScopes,omitempty"`
	OIDC             *AuthOIDC        `toml:"oidc,omitempty" json:"oidc,omitempty"`
}

type WebServerTLS struct {
	CertFile string `toml:"certFile,omitempty" json:"certFile,omitempty"`
	KeyFile  string `toml:"keyFile,omitempty" json:"keyFile,omitempty"`
}

type WebServer struct {
	Addr        string        `toml:"addr" json:"addr"`
	Port        int           `toml:"port" json:"port"`
	User        string        `toml:"user" json:"user"`
	Password    string        `toml:"password" json:"password"`
	TLS         *WebServerTLS `toml:"tls,omitempty" json:"tls,omitempty"`
	AssetsDir   string        `toml:"assetsDir,omitempty" json:"assetsDir,omitempty"`
	PprofEnable *bool         `toml:"pprofEnable,omitempty" json:"pprofEnable,omitempty"`
}

type Log struct {
	To                string `toml:"to" json:"to"`
	Level             string `toml:"level" json:"level"`
	MaxDays           *int   `toml:"maxDays,omitempty" json:"maxDays,omitempty"`
	DisablePrintColor *bool  `toml:"disablePrintColor,omitempty" json:"disablePrintColor,omitempty"`
}

type SSHTunnelGateway struct {
	BindPort              *int   `toml:"bindPort,omitempty" json:"bindPort,omitempty"`
	PrivateKeyFile        string `toml:"privateKeyFile,omitempty" json:"privateKeyFile,omitempty"`
	AutoGenPrivateKeyPath string `toml:"autoGenPrivateKeyPath,omitempty" json:"autoGenPrivateKeyPath,omitempty"`
	AuthorizedKeysFile    string `toml:"authorizedKeysFile,omitempty" json:"authorizedKeysFile,omitempty"`
}

type HTTPPlugin struct {
	Name string   `toml:"name" json:"name"`
	Addr string   `toml:"addr" json:"addr"`
	Path string   `toml:"path" json:"path"`
	Ops  []string `toml:"ops" json:"ops"`
}

type FrpsConfig struct {
	BindAddr                        string            `toml:"bindAddr,omitempty" json:"bindAddr"`
	BindPort                        *int              `toml:"bindPort,omitempty" json:"bindPort"`
	KCPBindPort                     *int              `toml:"kcpBindPort,omitempty" json:"kcpBindPort"`
	QuicBindPort                    *int              `toml:"quicBindPort,omitempty" json:"quicBindPort,omitempty"`
	ProxyBindAddr                   string            `toml:"proxyBindAddr,omitempty" json:"proxyBindAddr,omitempty"`
	Transport                       Transport         `toml:"transport" json:"transport"`
	VhostHTTPPort                   *int              `toml:"vhostHTTPPort,omitempty" json:"vhostHTTPPort,omitempty"`
	VhostHTTPSPort                  *int              `toml:"vhostHTTPSPort,omitempty" json:"vhostHTTPSPort,omitempty"`
	VhostHTTPTimeout                *int              `toml:"vhostHTTPTimeout,omitempty" json:"vhostHTTPTimeout,omitempty"`
	TCPMuxHTTPConnectPort           *int              `toml:"tcpmuxHTTPConnectPort,omitempty" json:"tcpmuxHTTPConnectPort,omitempty"`
	TCPMuxPassthrough               *bool             `toml:"tcpmuxPassthrough,omitempty" json:"tcpmuxPassthrough,omitempty"`
	WebServer                       WebServer         `toml:"webServer" json:"webServer"`
	EnablePrometheus                *bool             `toml:"enablePrometheus,omitempty" json:"enablePrometheus,omitempty"`
	Log                             Log               `toml:"log" json:"log"`
	DetailedErrorsToClient          *bool             `toml:"detailedErrorsToClient,omitempty" json:"detailedErrorsToClient,omitempty"`
	Auth                            Auth              `toml:"auth" json:"auth"`
	UserConnTimeout                 *int              `toml:"userConnTimeout,omitempty" json:"userConnTimeout,omitempty"`
	AllowPorts                      []AllowPort       `toml:"allowPorts,omitempty" json:"allowPorts,omitempty"`
	MaxPortsPerClient               *int              `toml:"maxPortsPerClient,omitempty" json:"maxPortsPerClient,omitempty"`
	SubDomainHost                   string            `toml:"subDomainHost,omitempty" json:"subDomainHost,omitempty"`
	Custom404Page                   string            `toml:"custom404Page,omitempty" json:"custom404Page,omitempty"`
	UDPPacketSize                   *int              `toml:"udpPacketSize,omitempty" json:"udpPacketSize,omitempty"`
	NatholeAnalysisDataReserveHours *int              `toml:"natholeAnalysisDataReserveHours,omitempty" json:"natholeAnalysisDataReserveHours,omitempty"`
	SSHTunnelGateway                *SSHTunnelGateway `toml:"sshTunnelGateway,omitempty" json:"sshTunnelGateway,omitempty"`
	HTTPPlugins                     []HTTPPlugin      `toml:"httpPlugins" json:"httpPlugins"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
