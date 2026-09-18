<template>
  <div class="server-overview">
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard
          :label="$t('nav.clients')"
          :value="data.status.clientCounts"
          type="clients"
          :subtitle="$t('overview.connectedClients')"
          to="/clients"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard
          :label="$t('nav.proxies')"
          :value="proxyCounts"
          type="proxies"
          :subtitle="$t('overview.activeProxies')"
          to="/proxies/tcp"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard
          :label="$t('clientDetail.connections')"
          :value="data.status.curConns"
          type="connections"
          :subtitle="$t('overview.currentConnections')"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard
          :label="$t('proxyDetail.traffic')"
          :value="formatTrafficTotal()"
          type="traffic"
          :subtitle="$t('overview.totalToday')"
        />
      </el-col>
    </el-row>

    <el-row :gutter="20" class="charts-row">
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">{{ $t('overview.networkTraffic') }}</span>
              <el-tag size="small" type="info">{{ $t('overview.today') }}</el-tag>
            </div>
          </template>
          <div class="traffic-summary">
            <div class="traffic-item in">
              <div class="traffic-icon">
                <el-icon><Download /></el-icon>
              </div>
              <div class="traffic-info">
                <div class="label">{{ $t('overview.inbound') }}</div>
                <div class="value">
                  {{ formatFileSize(data.status.totalTrafficIn) }}
                </div>
              </div>
            </div>
            <div class="traffic-divider"></div>
            <div class="traffic-item out">
              <div class="traffic-icon">
                <el-icon><Upload /></el-icon>
              </div>
              <div class="traffic-info">
                <div class="label">{{ $t('overview.outbound') }}</div>
                <div class="value">
                  {{ formatFileSize(data.status.totalTrafficOut) }}
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">{{ $t('overview.proxyTypes') }}</span>
              <el-tag size="small" type="info">{{ $t('overview.now') }}</el-tag>
            </div>
          </template>
          <div class="proxy-types-grid">
            <div
              v-for="(count, type) in data.status.proxyTypeCount"
              :key="type"
              class="proxy-type-item"
              v-show="count > 0"
            >
              <div class="proxy-type-name">{{ type.toUpperCase() }}</div>
              <div class="proxy-type-count">{{ count }}</div>
            </div>
            <div v-if="!hasActiveProxies" class="no-data">
              {{ $t('overview.noActiveProxies') }}
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="config-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">{{ $t('overview.serverConfiguration') }}</span>
          <el-tag size="small" type="success">v{{ data.version }}</el-tag>
        </div>
      </template>
      <div class="config-grid">
        <div class="config-item">
          <span class="config-label">{{ $t('config.bindPort') }}</span>
          <span class="config-value">{{ data.config.bindPort }}</span>
        </div>
        <div class="config-item" v-if="data.config.kcpBindPort != 0">
          <span class="config-label">{{ $t('config.kcpPort') }}</span>
          <span class="config-value">{{ data.config.kcpBindPort }}</span>
        </div>
        <div class="config-item" v-if="data.config.quicBindPort != 0">
          <span class="config-label">{{ $t('config.quicPort') }}</span>
          <span class="config-value">{{ data.config.quicBindPort }}</span>
        </div>
        <div class="config-item" v-if="data.config.vhostHTTPPort != 0">
          <span class="config-label">{{ $t('config.httpPort') }}</span>
          <span class="config-value">{{ data.config.vhostHTTPPort }}</span>
        </div>
        <div class="config-item" v-if="data.config.vhostHTTPSPort != 0">
          <span class="config-label">{{ $t('config.httpsPort') }}</span>
          <span class="config-value">{{ data.config.vhostHTTPSPort }}</span>
        </div>
        <div class="config-item" v-if="data.config.tcpmuxHTTPConnectPort != 0">
          <span class="config-label">{{ $t('config.tcpmuxPort') }}</span>
          <span class="config-value">{{ data.config.tcpmuxHTTPConnectPort }}</span>
        </div>
        <div class="config-item" v-if="data.config.subdomainHost != ''">
          <span class="config-label">{{ $t('config.subdomainHost') }}</span>
          <span class="config-value">{{ data.config.subdomainHost }}</span>
        </div>
        <div class="config-item">
          <span class="config-label">{{ $t('config.maxPoolCount') }}</span>
          <span class="config-value">{{ data.config.maxPoolCount }}</span>
        </div>
        <div class="config-item">
          <span class="config-label">{{ $t('config.maxPortsPerClient') }}</span>
          <span class="config-value">{{ maxPortsPerClientLabel }}</span>
        </div>
        <div class="config-item" v-if="data.config.allowPortsStr != ''">
          <span class="config-label">{{ $t('config.allowPorts') }}</span>
          <span class="config-value">{{ data.config.allowPortsStr }}</span>
        </div>
        <div class="config-item" v-if="data.config.tlsForce">
          <span class="config-label">{{ $t('config.tlsForce') }}</span>
          <el-tag size="small" type="warning">{{ $t('common.enabled') }}</el-tag>
        </div>
        <div class="config-item">
          <span class="config-label">{{ $t('config.heartbeatTimeout') }}</span>
          <span class="config-value">{{ data.config.heartbeatTimeout }}s</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { formatFileSize } from '../utils/format'
import { Download, Upload } from '@element-plus/icons-vue'
import StatCard from '../components/StatCard.vue'
import { getServerInfo } from '../api/server'
import type { ServerInfo } from '../types/server'

const { t } = useI18n()

const data = ref<ServerInfo>({
  version: '',
  config: {
    bindPort: 0,
    kcpBindPort: 0,
    quicBindPort: 0,
    vhostHTTPPort: 0,
    vhostHTTPSPort: 0,
    tcpmuxHTTPConnectPort: 0,
    subdomainHost: '',
    maxPoolCount: 0,
    maxPortsPerClient: 0,
    allowPortsStr: '',
    tlsForce: false,
    heartbeatTimeout: 0,
  },
  status: {
    clientCounts: 0,
    curConns: 0,
    totalTrafficIn: 0,
    totalTrafficOut: 0,
    proxyTypeCount: {},
  },
})

const hasActiveProxies = computed(() => {
  return Object.values(data.value.status.proxyTypeCount).some((c) => c > 0)
})

const proxyCounts = computed(() => {
  return Object.values(data.value.status.proxyTypeCount).reduce(
    (sum, count) => sum + (count || 0),
    0,
  )
})

const maxPortsPerClientLabel = computed(() => {
  const value = data.value.config.maxPortsPerClient
  return value === 0 ? t('common.noLimit') : String(value)
})

const formatTrafficTotal = () => {
  const total =
    data.value.status.totalTrafficIn + data.value.status.totalTrafficOut
  return formatFileSize(total)
}

const fetchData = async () => {
  try {
    const json = await getServerInfo()
    data.value = json
  } catch {
    ElMessage({
      showClose: true,
      message: t('overview.fetchError'),
      type: 'error',
    })
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.server-overview {
  padding: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.charts-row {
  margin-bottom: 20px;
}

.chart-card {
  border-radius: 12px;
  border: 1px solid #e4e7ed;
  height: 100%;
}

html.dark .chart-card {
  border-color: #3a3d5c;
  background: #27293d;
}

.config-card {
  border-radius: 12px;
  border: 1px solid #e4e7ed;
  margin-bottom: 20px;
}

html.dark .config-card {
  border-color: #3a3d5c;
  background: #27293d;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

html.dark .card-title {
  color: #e5e7eb;
}

.traffic-summary {
  display: flex;
  align-items: center;
  justify-content: space-around;
  min-height: 120px;
  padding: 10px 0;
}

.traffic-item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.traffic-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.traffic-item.in .traffic-icon {
  background: rgba(84, 112, 198, 0.1);
  color: #5470c6;
}

.traffic-item.out .traffic-icon {
  background: rgba(145, 204, 117, 0.1);
  color: #91cc75;
}

.traffic-info {
  display: flex;
  flex-direction: column;
}

.traffic-info .label {
  font-size: 14px;
  color: #909399;
}

.traffic-info .value {
  font-size: 24px;
  font-weight: 500;
  color: #303133;
}

html.dark .traffic-info .value {
  color: #e5e7eb;
}

.traffic-divider {
  width: 1px;
  height: 60px;
  background: #e4e7ed;
}

html.dark .traffic-divider {
  background: #3a3d5c;
}

.proxy-types-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 16px;
  min-height: 120px;
  align-content: center;
  padding: 10px 0;
}

.proxy-type-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
}

html.dark .proxy-type-item {
  background: #1e1e2d;
}

.proxy-type-name {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
  margin-bottom: 4px;
}

.proxy-type-count {
  font-size: 20px;
  font-weight: 500;
  color: #303133;
}

html.dark .proxy-type-count {
  color: #e5e7eb;
}

.no-data {
  grid-column: 1 / -1;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
  font-size: 14px;
}

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
}

.config-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  transition: background 0.2s;
}

html.dark .config-item {
  background: #1e1e2d;
}

.config-label {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
}

html.dark .config-label {
  color: #9ca3af;
}

.config-value {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
  word-break: break-all;
}

html.dark .config-value {
  color: #e5e7eb;
}

@media (max-width: 768px) {
  .chart-container {
    height: 250px;
  }

  .config-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
