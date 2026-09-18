# Changelog

## 2026-09-18 — v0.71.0-custom (d3018b00)

### Added
- **Web 面板登录**：配置 `webServer.user` / `webServer.password` 后自动启用 Web 页登录（替代浏览器 Basic Auth 弹窗），支持「记住我」（30 天 httpOnly cookie）。
- **配置管理面板**：新增 `/config` 页面，支持在线查看、编辑、保存 frps TOML 配置，并可一键 reload/restart（配置写回文件后自动重启服务端）。
- **配置热重启**：frps 支持运行中通过 `/api/config/restart` 接口或 Web 面板触发重启，自动重新加载配置文件。
- **多语言支持（i18n）**：集成 vue-i18n，支持英文/中文切换（语言偏好存储在 localStorage）。
- **404 页面重新设计**：暗色主题卡片式 404 页面，支持渐变动画。
- **systemd 服务文件**：新增 `packaging/frpc.service` 和 `packaging/frps.service`，可直接用于 systemd 部署。
- **编译文档**：新增 `doc/build.md`（Linux 编译指南）和 `doc/build-windows.md`（Windows 交叉编译指南）。
- **测试配置示例**：新增 `conf/frps_test_example.toml`，包含所有配置项的中文注释。

### Changed
- `server/api_router.go`：静态资源路由（`/static/`、`/`、favicon）移至鉴权之外，仅保护 API 路由。
- `pkg/util/net/http.go`：HTTP 认证中间件新增 token 校验能力（cookie `frp-auth-token` / Bearer token），保留 BasicAuth 兜底。
- `server/http/controller.go`：构造函数增加 `sessionMgr` 参数。
- `server/service.go`：新增 `sessionMgr`、`configMgr`、vhost 监听器生命周期管理、重启原子标记。
- `cmd/frps/root.go`：启动入口改为循环模式，支持热重启后重新创建服务。
- `web/frps/src/views/`：全部硬编码文案替换为 `$t()` 国际化调用。
- `web/frps/src/App.vue`：增加登录页/主布局切换、语言下拉、退出按钮、config 侧边栏。

### Fixed
- `server/http/controller_v2_test.go`：适配 `NewController` 新增参数，所有测试通过。

---

## 2026-09-18 — v0.71.0-custom（d3018b00）

### 新增
- **Web 面板登录**：配置 `webServer.user` / `webServer.password` 后自动启用 Web 页登录（替代浏览器 Basic Auth 弹窗），支持「记住我」（30 天 httpOnly cookie）。
- **配置管理面板**：新增 `/config` 页面，支持在线查看、编辑、保存 frps TOML 配置，并可一键 reload/restart（配置写回文件后自动重启服务端）。
- **配置热重启**：frps 支持运行中通过 `/api/config/restart` 接口或 Web 面板触发重启，自动重新加载配置文件。
- **多语言支持（i18n）**：集成 vue-i18n，支持英文/中文切换（语言偏好存储在 localStorage）。
- **404 页面重新设计**：暗色主题卡片式 404 页面，支持渐变动画。
- **systemd 服务文件**：新增 `packaging/frpc.service` 和 `packaging/frps.service`，可直接用于 systemd 部署。
- **编译文档**：新增 `doc/build.md`（Linux 编译指南）和 `doc/build-windows.md`（Windows 交叉编译指南）。
- **测试配置示例**：新增 `conf/frps_test_example.toml`，包含所有配置项的中文注释。

### 变更
- `server/api_router.go`：静态资源路由（`/static/`、`/`、favicon）移至鉴权之外，仅保护 API 路由。
- `pkg/util/net/http.go`：HTTP 认证中间件新增 token 校验能力（cookie `frp-auth-token` / Bearer token），保留 BasicAuth 兜底。
- `server/http/controller.go`：构造函数增加 `sessionMgr` 参数。
- `server/service.go`：新增 `sessionMgr`、`configMgr`、vhost 监听器生命周期管理、重启原子标记。
- `cmd/frps/root.go`：启动入口改为循环模式，支持热重启后重新创建服务。
- `web/frps/src/views/`：全部硬编码文案替换为 `$t()` 国际化调用。
- `web/frps/src/App.vue`：增加登录页/主布局切换、语言下拉、退出按钮、config 侧边栏。

### 修复
- `server/http/controller_v2_test.go`：适配 `NewController` 新增参数，所有测试通过。
