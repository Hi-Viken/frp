export interface AllowPort {
  start?: number
  end?: number
  single?: number
}

export interface TransportTLS {
  force?: boolean
  certFile?: string
  keyFile?: string
  trustedCaFile?: string
}

export interface TransportQUIC {
  keepalivePeriod?: number
  maxIdleTimeout?: number
  maxIncomingStreams?: number
}

export interface Transport {
  tls: TransportTLS
  quic: TransportQUIC
  maxPoolCount?: number
  heartbeatTimeout?: number
  tcpMux?: boolean
  tcpMuxKeepaliveInterval?: number
  tcpKeepalive?: number
}

export interface AuthOIDC {
  issuer?: string
  audience?: string
  skipExpiryCheck?: boolean
  skipIssuerCheck?: boolean
}

export interface AuthTokenFile {
  path: string
}

export interface AuthTokenSource {
  type: string
  file: AuthTokenFile
}

export interface Auth {
  method: string
  token?: string
  tokenSource?: AuthTokenSource
  additionalScopes?: string[]
  oidc: AuthOIDC
}

export interface WebServerTLS {
  certFile?: string
  keyFile?: string
}

export interface WebServer {
  addr: string
  port: number
  user: string
  password: string
  tls: WebServerTLS
  assetsDir?: string
  pprofEnable?: boolean
}

export interface Log {
  to: string
  level: string
  maxDays?: number
  disablePrintColor?: boolean
}

export interface SSHTunnelGateway {
  bindPort?: number
  privateKeyFile?: string
  autoGenPrivateKeyPath?: string
  authorizedKeysFile?: string
}

export interface HTTPPlugin {
  name: string
  addr: string
  path: string
  ops: string[]
}

export interface FrpsConfig {
  bindAddr: string
  bindPort?: number
  kcpBindPort?: number
  quicBindPort?: number
  proxyBindAddr?: string
  transport: Transport
  vhostHTTPPort?: number
  vhostHTTPSPort?: number
  vhostHTTPTimeout?: number
  tcpmuxHTTPConnectPort?: number
  tcpmuxPassthrough?: boolean
  webServer: WebServer
  enablePrometheus?: boolean
  log: Log
  detailedErrorsToClient?: boolean
  auth: Auth
  userConnTimeout?: number
  allowPorts?: AllowPort[]
  maxPortsPerClient?: number
  subDomainHost?: string
  custom404Page?: string
  udpPacketSize?: number
  natholeAnalysisDataReserveHours?: number
  sshTunnelGateway: SSHTunnelGateway
  httpPlugins: HTTPPlugin[]
}

export interface ConfigAPIResponse {
  success: boolean
  data?: FrpsConfig
  error?: string
}

export function emptyConfig(): FrpsConfig {
  return {
    bindAddr: '0.0.0.0',
    bindPort: 7000,
    transport: {
      tls: {},
      quic: {},
      maxPoolCount: 5,
    },
    webServer: { addr: '127.0.0.1', port: 7500, user: '', password: '', tls: {} },
    log: { to: '', level: 'info' },
    auth: { method: 'token', oidc: {} },
    sshTunnelGateway: {},
    httpPlugins: [],
  }
}
