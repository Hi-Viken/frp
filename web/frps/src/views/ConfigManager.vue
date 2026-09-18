<template>
  <div class="config-manager">
    <div class="config-header">
      <h2 class="config-title">{{ $t('configMgr.title', 'Configuration Manager') }}</h2>
      <div class="config-actions">
        <el-button type="primary" :loading="saving" @click="handleSave">
          {{ $t('configMgr.save', 'Save') }}
        </el-button>
        <el-button :loading="reloading" @click="handleReload">
          {{ $t('configMgr.refresh', 'Refresh') }}
        </el-button>
        <el-popconfirm
          :title="$t('configMgr.restartConfirm', 'Restart frps? All active connections will be disconnected.')"
          :confirm-button-text="$t('common.confirm', 'Confirm')"
          :cancel-button-text="$t('common.cancel', 'Cancel')"
          @confirm="handleRestart"
        >
          <template #reference>
            <el-button type="danger" :loading="restarting">
              {{ $t('configMgr.restart', 'Restart') }}
            </el-button>
          </template>
        </el-popconfirm>
      </div>
    </div>

    <el-alert
      v-if="errorMsg"
      :title="errorMsg"
      type="error"
      show-icon
      closable
      class="config-alert"
      @close="errorMsg = ''"
    />

    <el-tabs v-model="activeTab" class="config-tabs">
      <el-tab-pane :label="$t('configMgr.tabs.basic', 'Basic')" name="basic">
        <el-form label-width="180px" class="config-form">
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.bindAddr', 'Bind Address')">
              <el-input v-model="form.bindAddr" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.bindPort', 'Bind Port')">
              <el-input-number v-model="form.bindPort" :min="0" :max="65535" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.kcpBindPort', 'KCP Bind Port')">
              <el-input-number v-model="form.kcpBindPort" :min="0" :max="65535" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.quicBindPort', 'QUIC Bind Port')">
              <el-input-number v-model="form.quicBindPort" :min="0" :max="65535" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.proxyBindAddr', 'Proxy Bind Address')">
              <el-input v-model="form.proxyBindAddr" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.udpPacketSize', 'UDP Packet Size')">
              <el-input-number v-model="form.udpPacketSize" :min="0" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.natholeAnalysisHours', 'NAT Hole Analysis Hours')">
              <el-input-number v-model="form.natholeAnalysisDataReserveHours" :min="0" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.subDomainHost', 'Subdomain Host')">
              <el-input v-model="form.subDomainHost" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.custom404Page', 'Custom 404 Page')">
              <el-input v-model="form.custom404Page" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.detailedErrors', 'Detailed Errors to Client')">
              <el-switch v-model="form.detailedErrorsToClient" />
            </el-form-item>
          </div>
        </el-form>
      </el-tab-pane>

      <el-tab-pane :label="$t('configMgr.tabs.transport', 'Transport')" name="transport">
        <el-form label-width="220px" class="config-form">
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.maxPoolCount', 'Max Pool Count')">
              <el-input-number v-model="form.transport.maxPoolCount" :min="0" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.heartbeatTimeout', 'Heartbeat Timeout (s)')">
              <el-input-number v-model="form.transport.heartbeatTimeout" :min="0" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.tcpMux', 'TCP Mux')">
              <el-switch v-model="form.transport.tcpMux" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.tcpMuxKeepalive', 'TCP Mux Keepalive Interval')">
              <el-input-number v-model="form.transport.tcpMuxKeepaliveInterval" :min="0" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.tcpKeepalive', 'TCP Keepalive')">
              <el-input-number v-model="form.transport.tcpKeepalive" :min="0" />
            </el-form-item>
          </div>

          <el-divider>{{ $t('configMgr.tls', 'TLS Settings') }}</el-divider>
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.tlsForce', 'Force TLS')">
              <el-switch v-model="form.transport.tls.force" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.tlsCertFile', 'TLS Cert File')">
              <el-input v-model="form.transport.tls.certFile" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.tlsKeyFile', 'TLS Key File')">
              <el-input v-model="form.transport.tls.keyFile" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.tlsCaFile', 'TLS Trusted CA File')">
              <el-input v-model="form.transport.tls.trustedCaFile" />
            </el-form-item>
          </div>

          <el-divider>{{ $t('configMgr.quic', 'QUIC Settings') }}</el-divider>
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.quicKeepalive', 'QUIC Keepalive Period')">
              <el-input-number v-model="form.transport.quic.keepalivePeriod" :min="0" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.quicIdleTimeout', 'QUIC Max Idle Timeout')">
              <el-input-number v-model="form.transport.quic.maxIdleTimeout" :min="0" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.quicMaxStreams', 'QUIC Max Incoming Streams')">
              <el-input-number v-model="form.transport.quic.maxIncomingStreams" :min="0" />
            </el-form-item>
          </div>
        </el-form>
      </el-tab-pane>

      <el-tab-pane :label="$t('configMgr.tabs.vhost', 'Virtual Host')" name="vhost">
        <el-form label-width="220px" class="config-form">
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.vhostHTTPPort', 'VHost HTTP Port')">
              <el-input-number v-model="form.vhostHTTPPort" :min="0" :max="65535" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.vhostHTTPSPort', 'VHost HTTPS Port')">
              <el-input-number v-model="form.vhostHTTPSPort" :min="0" :max="65535" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.vhostHTTPTimeout', 'VHost HTTP Timeout (s)')">
              <el-input-number v-model="form.vhostHTTPTimeout" :min="0" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.tcpmuxPort', 'TCPMux HTTP Connect Port')">
              <el-input-number v-model="form.tcpmuxHTTPConnectPort" :min="0" :max="65535" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.tcpmuxPassthrough', 'TCPMux Passthrough')">
              <el-switch v-model="form.tcpmuxPassthrough" />
            </el-form-item>
          </div>
        </el-form>
      </el-tab-pane>

      <el-tab-pane :label="$t('configMgr.tabs.dashboard', 'Dashboard')" name="dashboard">
        <el-form label-width="200px" class="config-form">
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.webAddr', 'Web Server Address')">
              <el-input v-model="form.webServer.addr" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.webPort', 'Web Server Port')">
              <el-input-number v-model="form.webServer.port" :min="0" :max="65535" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.webUser', 'Web Server User')">
              <el-input v-model="form.webServer.user" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.webPassword', 'Web Server Password')">
              <el-input v-model="form.webServer.password" type="password" show-password />
            </el-form-item>
            <el-form-item :label="$t('configMgr.pprofEnable', 'Enable Pprof')">
              <el-switch v-model="form.webServer.pprofEnable" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.prometheus', 'Enable Prometheus')">
              <el-switch v-model="form.enablePrometheus" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.assetsDir', 'Assets Directory')">
              <el-input v-model="form.webServer.assetsDir" />
            </el-form-item>
          </div>

          <el-divider>{{ $t('configMgr.webTls', 'Web Server TLS') }}</el-divider>
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.webTlsCert', 'TLS Cert File')">
              <el-input v-model="form.webServer.tls.certFile" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.webTlsKey', 'TLS Key File')">
              <el-input v-model="form.webServer.tls.keyFile" />
            </el-form-item>
          </div>
        </el-form>
      </el-tab-pane>

      <el-tab-pane :label="$t('configMgr.tabs.auth', 'Auth')" name="auth">
        <el-form label-width="200px" class="config-form">
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.authMethod', 'Auth Method')">
              <el-select v-model="form.auth.method">
                <el-option label="Token" value="token" />
                <el-option label="OIDC" value="oidc" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('configMgr.authToken', 'Token')">
              <el-input v-model="form.auth.token" type="password" show-password />
            </el-form-item>
            <el-form-item
              v-if="form.auth.tokenSource"
              :label="$t('configMgr.tokenSourceType', 'Token Source Type')"
            >
              <el-input v-model="form.auth.tokenSource.type" :placeholder="$t('configMgr.tokenSourceTypeHint', 'file or exec')" />
            </el-form-item>
            <el-form-item
              v-if="form.auth.tokenSource"
              :label="$t('configMgr.tokenSourcePath', 'Token Source Path')"
            >
              <el-input v-model="form.auth.tokenSource.file.path" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.tokenSource', 'Token Source')">
              <el-switch
                :model-value="!!form.auth.tokenSource"
                @update:model-value="(v: boolean) => toggleTokenSource(v)"
              />
            </el-form-item>
            <el-form-item :label="$t('configMgr.additionalScopes', 'Additional Scopes')">
              <el-input v-model="additionalScopesText" :placeholder="$t('configMgr.scopesHint', 'comma-separated')" />
            </el-form-item>
          </div>

          <el-divider>{{ $t('configMgr.oidc', 'OIDC Settings') }}</el-divider>
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.oidcIssuer', 'OIDC Issuer')">
              <el-input v-model="form.auth.oidc.issuer" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.oidcAudience', 'OIDC Audience')">
              <el-input v-model="form.auth.oidc.audience" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.oidcSkipExpiry', 'Skip Expiry Check')">
              <el-switch v-model="form.auth.oidc.skipExpiryCheck" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.oidcSkipIssuer', 'Skip Issuer Check')">
              <el-switch v-model="form.auth.oidc.skipIssuerCheck" />
            </el-form-item>
          </div>
        </el-form>
      </el-tab-pane>

      <el-tab-pane :label="$t('configMgr.tabs.logging', 'Logging')" name="logging">
        <el-form label-width="200px" class="config-form">
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.logTo', 'Log To')">
              <el-input v-model="form.log.to" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.logLevel', 'Log Level')">
              <el-select v-model="form.log.level">
                <el-option label="Trace" value="trace" />
                <el-option label="Debug" value="debug" />
                <el-option label="Info" value="info" />
                <el-option label="Warn" value="warn" />
                <el-option label="Error" value="error" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('configMgr.logMaxDays', 'Max Days')">
              <el-input-number v-model="form.log.maxDays" :min="0" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.logDisableColor', 'Disable Print Color')">
              <el-switch v-model="form.log.disablePrintColor" />
            </el-form-item>
          </div>
        </el-form>
      </el-tab-pane>

      <el-tab-pane :label="$t('configMgr.tabs.ports', 'Ports')" name="ports">
        <el-form label-width="200px" class="config-form">
          <el-form-item :label="$t('configMgr.maxPortsPerClient', 'Max Ports Per Client')">
            <el-input-number v-model="form.maxPortsPerClient" :min="0" />
          </el-form-item>

          <el-form-item :label="$t('configMgr.allowPorts', 'Allow Ports')">
            <div class="dynamic-list">
              <div v-for="(port, index) in form.allowPorts" :key="index" class="dynamic-item">
                <el-select :model-value="portType(port)" @update:model-value="(val: string) => setPortType(port, val)">
                  <el-option :label="$t('configMgr.portSingle', 'Single')" value="single" />
                  <el-option :label="$t('configMgr.portRange', 'Range')" value="range" />
                </el-select>
                <el-input-number
                  v-if="portType(port) === 'single'"
                  v-model="port.single"
                  :min="1"
                  :max="65535"
                  :placeholder="$t('configMgr.port', 'Port')"
                />
                <template v-else>
                  <el-input-number v-model="port.start" :min="1" :max="65535" :placeholder="$t('configMgr.startPort', 'Start')" />
                  <span class="range-separator">-</span>
                  <el-input-number v-model="port.end" :min="1" :max="65535" :placeholder="$t('configMgr.endPort', 'End')" />
                </template>
                <el-button type="danger" :icon="Delete" circle @click="removeAllowPort(index)" />
              </div>
              <el-button type="primary" :icon="Plus" @click="addAllowPort">
                {{ $t('configMgr.addPort', 'Add Port') }}
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane :label="$t('configMgr.tabs.plugins', 'Plugins')" name="plugins">
        <el-form label-width="200px" class="config-form">
          <el-divider content-position="left">{{ $t('configMgr.httpPlugins', 'HTTP Plugins') }}</el-divider>
          <div class="plugins-list">
            <el-card v-for="(plugin, index) in form.httpPlugins" :key="index" class="plugin-card" shadow="hover">
              <template #header>
                <div class="plugin-header">
                  <span>{{ plugin.name || $t('configMgr.newPlugin', 'New Plugin') }}</span>
                  <el-button type="danger" :icon="Delete" circle @click="removePlugin(index)" />
                </div>
              </template>
              <div class="form-grid">
                <el-form-item :label="$t('configMgr.pluginName', 'Name')">
                  <el-input v-model="plugin.name" />
                </el-form-item>
                <el-form-item :label="$t('configMgr.pluginAddr', 'Address')">
                  <el-input v-model="plugin.addr" />
                </el-form-item>
                <el-form-item :label="$t('configMgr.pluginPath', 'Path')">
                  <el-input v-model="plugin.path" />
                </el-form-item>
                <el-form-item :label="$t('configMgr.pluginOps', 'Ops')">
                  <el-input :model-value="pluginOpsText[index]" @update:model-value="(val: string) => setPluginOps(index, val)" :placeholder="$t('configMgr.opsHint', 'comma-separated')" />
                </el-form-item>
              </div>
            </el-card>
            <el-button type="primary" :icon="Plus" @click="addPlugin">
              {{ $t('configMgr.addPlugin', 'Add Plugin') }}
            </el-button>
          </div>

          <el-divider content-position="left">{{ $t('configMgr.sshGateway', 'SSH Tunnel Gateway') }}</el-divider>
          <div class="form-grid">
            <el-form-item :label="$t('configMgr.sshBindPort', 'Bind Port')">
              <el-input-number v-model="form.sshTunnelGateway.bindPort" :min="0" :max="65535" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.sshPrivateKey', 'Private Key File')">
              <el-input v-model="form.sshTunnelGateway.privateKeyFile" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.sshAutoGenKey', 'Auto Gen Private Key Path')">
              <el-input v-model="form.sshTunnelGateway.autoGenPrivateKeyPath" />
            </el-form-item>
            <el-form-item :label="$t('configMgr.sshAuthorizedKeys', 'Authorized Keys File')">
              <el-input v-model="form.sshTunnelGateway.authorizedKeysFile" />
            </el-form-item>
          </div>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Delete, Plus } from '@element-plus/icons-vue'
import { configApi } from '../api/config'
import type { FrpsConfig, AllowPort } from '../types/config'
import { emptyConfig } from '../types/config'

const { t } = useI18n()

const activeTab = ref('basic')
const saving = ref(false)
const reloading = ref(false)
const restarting = ref(false)
const loading = ref(false)
const errorMsg = ref('')
const form = ref<FrpsConfig>(emptyConfig())

const additionalScopesText = computed({
  get: () => (form.value.auth.additionalScopes || []).join(','),
  set: (val: string) => {
    form.value.auth.additionalScopes = val.split(',').map((s) => s.trim()).filter(Boolean)
  },
})

const pluginOpsText = computed(() =>
  form.value.httpPlugins.map((p) => p.ops.join(',')),
)

function setPluginOps(index: number, val: string) {
  form.value.httpPlugins[index].ops = val.split(',').map((s) => s.trim()).filter(Boolean)
}

function portType(port: AllowPort): string {
  return port.single !== undefined ? 'single' : 'range'
}

function setPortType(port: AllowPort, type: string) {
  if (type === 'single') {
    port.single = port.start || 0
    delete port.start
    delete port.end
  } else {
    port.start = port.single || 0
    port.end = port.single || 0
    delete port.single
  }
}

function toggleTokenSource(v: boolean) {
  if (v) {
    form.value.auth.tokenSource = { type: 'file', file: { path: '' } }
  } else {
    delete form.value.auth.tokenSource
  }
}

function addAllowPort() {
  if (!form.value.allowPorts) form.value.allowPorts = []
  form.value.allowPorts.push({ single: 80 })
}

function removeAllowPort(index: number) {
  form.value.allowPorts?.splice(index, 1)
}

function addPlugin() {
  form.value.httpPlugins.push({ name: '', addr: '', path: '', ops: [] })
}

function removePlugin(index: number) {
  form.value.httpPlugins.splice(index, 1)
}

function normalizeForm(cfg: FrpsConfig) {
  cfg.transport = cfg.transport || {}
  if (!cfg.transport.tls) cfg.transport.tls = {}
  if (!cfg.transport.quic) cfg.transport.quic = {}
  cfg.auth = cfg.auth || { method: 'token' }
  if (!cfg.auth.oidc) cfg.auth.oidc = {}
  cfg.webServer = cfg.webServer || { addr: '0.0.0.0', port: 7500, user: '', password: '' }
  if (!cfg.webServer.tls) cfg.webServer.tls = {}
  if (!cfg.sshTunnelGateway) cfg.sshTunnelGateway = {}
  return cfg
}

async function loadConfig() {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await configApi.get()
    if (res.success && res.data) {
      form.value = normalizeForm(res.data)
      if (!form.value.allowPorts) form.value.allowPorts = []
      if (!form.value.httpPlugins) form.value.httpPlugins = []
    } else {
      errorMsg.value = res.error || t('configMgr.loadError', 'Failed to load config')
      ElMessage.error(errorMsg.value)
    }
  } catch (err) {
    errorMsg.value = err instanceof Error ? err.message : t('configMgr.loadError', 'Failed to load config')
    ElMessage.error(errorMsg.value)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  errorMsg.value = ''
  try {
    const payload: FrpsConfig = normalizeForm(JSON.parse(JSON.stringify(form.value)))
    if (!payload.allowPorts) payload.allowPorts = []
    if (!payload.httpPlugins) payload.httpPlugins = []
    const res = await configApi.save(payload)
    if (res.success) {
      ElMessage.success(t('configMgr.saveSuccess', 'Configuration saved successfully'))
      if (res.data) form.value = normalizeForm(res.data)
    } else {
      errorMsg.value = res.error || t('configMgr.saveError', 'Failed to save config')
      ElMessage.error(errorMsg.value)
    }
  } catch (err) {
    errorMsg.value = err instanceof Error ? err.message : t('configMgr.saveError', 'Failed to save config')
    ElMessage.error(errorMsg.value)
  } finally {
    saving.value = false
  }
}

async function handleReload() {
  reloading.value = true
  errorMsg.value = ''
  try {
    const res = await configApi.reload()
    if (res.success) {
      ElMessage.success(t('configMgr.reloadSuccess', 'Configuration reloaded successfully'))
      if (res.data) form.value = normalizeForm(res.data)
    } else {
      errorMsg.value = res.error || t('configMgr.reloadError', 'Failed to reload config')
      ElMessage.error(errorMsg.value)
    }
  } catch (err) {
    errorMsg.value = err instanceof Error ? err.message : t('configMgr.reloadError', 'Failed to reload config')
    ElMessage.error(errorMsg.value)
  } finally {
    reloading.value = false
  }
}

async function handleRestart() {
  restarting.value = true
  errorMsg.value = ''
  try {
    const res = await configApi.restart()
    if (res.success) {
      ElMessage.success(t('configMgr.restartSuccess', 'frps is restarting, page will reload in a few seconds...'))
      setTimeout(() => window.location.reload(), 3000)
    } else {
      errorMsg.value = res.error || t('configMgr.restartError', 'Failed to restart frps')
      ElMessage.error(errorMsg.value)
      restarting.value = false
    }
  } catch (err) {
    errorMsg.value = err instanceof Error ? err.message : t('configMgr.restartError', 'Failed to restart frps')
    ElMessage.error(errorMsg.value)
    restarting.value = false
  }
}

onMounted(() => {
  loadConfig()
})
</script>

<style scoped>
.config-manager {
  padding: 0;
}

.config-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.config-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.config-actions {
  display: flex;
  gap: 12px;
}

.config-alert {
  margin-bottom: 20px;
}

.config-tabs {
  background: var(--color-bg-surface);
  border-radius: 12px;
  padding: 20px;
  border: 1px solid var(--color-border-light);
}

.config-form {
  padding: 20px 0;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

.dynamic-list {
  width: 100%;
}

.dynamic-item {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.range-separator {
  color: var(--color-text-muted);
  font-weight: 500;
}

.plugins-list {
  margin-bottom: 20px;
}

.plugin-card {
  margin-bottom: 16px;
  border-radius: 8px;
  border: 1px solid var(--color-border-light);
}

.plugin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

html.dark .plugin-card {
  border-color: var(--color-border);
  background: var(--color-bg-tertiary);
}
</style>
