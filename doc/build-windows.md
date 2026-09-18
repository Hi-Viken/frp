# frp Windows 版交叉编译指南

本文档记录了 frp v0.71.0 在 Ubuntu 24.04（amd64）主机上交叉编译出
Windows 版（.exe）二进制及其文档的过程。

## 1. 说明

frp 使用 Go 编写，支持跨平台编译。本编译基于上一份《编译指南》
（`doc/build.md`）中已构建好的 Web 面板资源（`web/frps/dist`、
`web/frpc/dist`），在本机（Linux）上直接产出 Windows 可执行文件。

- 编译主机：Ubuntu 24.04.4 / x86_64
- 目标平台：Windows amd64(x86-64)、Windows arm64(AArch64)
- 编译方式：Go 交叉编译，`CGO_ENABLED=0` 纯静态，**无需 MinGW-w64 工具链**

## 2. 原理

Go 交叉编译通过在 `go build` 时指定目标平台环境变量实现：

| 环境变量 | 含义 |
| --- | --- |
| `GOOS=windows` | 目标操作系统为 Windows |
| `GOARCH=amd64` / `arm64` | 目标 CPU 架构 |
| `CGO_ENABLED=0` | 关闭 CGO，纯 Go 静态交叉编译，不依赖 MinGW 交叉编译器 |
| `GOPROXY` | 模块代理（本机需用 `https://goproxy.cn,direct`） |

> 关键点：`CGO_ENABLED=0` 时 Go 可以无交叉工具链地跨平台编译；
> 若开启 CGO（默认），Windows 目标需要安装 `mingw-w64`，本编译未使用。

## 3. 环境与依赖

### 3.1 必备工具（与 Linux 编译相同）

| 工具 | 版本 | 备注 |
| --- | --- | --- |
| Go | go1.25.0 linux/amd64 | 要求 ≥ 1.25.0 |
| make | GNU Make 4.3 | 可选（也可直接用 go build） |
| node / npm | v24.20.0 / 11.19.0 | 仅首次构建 Web 面板需要 |
| gcc | 系统自带 | 本编译未用到 |

### 3.2 网络依赖

Windows 目标编译会额外拉取 Windows 专用依赖
`golang.zx2c4.com/wintun`（wireguard-tun 在 Windows 下的实现），
该模块在 Linux 编译时不会被下载。`proxy.golang.org` 在本机不可达，
须使用镜像：

```bash
export GOPROXY="https://goproxy.cn,direct"
```

## 4. 编译步骤

### 4.1 前置准备

```bash
# 1) 确保 Web 面板已构建（否则二进制不含 Web 管理界面）
# 若 web/frps/dist 与 web/frpc/dist 不存在，先执行：
cd /root/frp/web && npm install
cd /root/frp && make web

# 2) 设置 Go 代理
export GOPROXY="https://goproxy.cn,direct"

# 3) 创建输出目录
mkdir -p /root/frp/release
```

### 4.2 编译 Windows amd64（x86-64）

```bash
cd /root/frp

# 服务端
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
  -trimpath -ldflags "-s -w" \
  -tags "frps" -o ./release/frps_windows_amd64.exe ./cmd/frps

# 客户端
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
  -trimpath -ldflags "-s -w" \
  -tags "frpc" -o ./release/frpc_windows_amd64.exe ./cmd/frpc
```

> 若未构建 Web 面板，请在 `-tags` 中追加 `,noweb`（`-tags "frps,noweb"`）。

### 4.3 编译 Windows arm64（ARM64 架构，如高通骁龙平台）

```bash
env CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build \
  -trimpath -ldflags "-s -w" \
  -tags "frps" -o ./release/frps_windows_arm64.exe ./cmd/frps

env CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build \
  -trimpath -ldflags "-s -w" \
  -tags "frpc" -o ./release/frpc_windows_arm64.exe ./cmd/frpc
```

### 4.4 批量编译（可选用官方脚本）

项目自带 `Makefile.cross-compiles`，可一键编译绝大多数平台
（含 windows amd64/arm64），产出到 `release/`：

```bash
cd /root/frp
make -f Makefile.cross-compiles
```

> 该脚本会编译 18 组平台目标，耗时较长；若只需 Windows 版，
> 推荐按 4.2 / 4.3 单独编译。

## 5. 编译参数说明

| 参数 | 含义 |
| --- | --- |
| `CGO_ENABLED=0` | 纯静态交叉编译，无需 MinGW，且产物不依赖 VC 运行库 |
| `GOOS=windows` | 目标系统 Windows，产出 PE 可执行文件 |
| `GOARCH=amd64/arm64` | 目标架构 |
| `-trimpath` | 去除构建路径，保证可复现 |
| `-ldflags "-s -w"` | 去除符号表与调试信息，减小体积 |
| `-tags "frps"/"frpc"` | 服务端/客户端构建分支 |
| `-o release/frpx_windows_xxx.exe` | 输出路径（.exe 后缀由命名给出） |

## 6. 实际操作记录

### 6.1 首次编译碰到的依赖下载问题

```text
$ env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ... ./cmd/frps
go: downloading golang.zx2c4.com/wintun ...
  Get "https://proxy.golang.org/golang.zx2c4.com/wintun/...": i/o timeout
```

处理：改用 `goproxy.cn` 镜像后成功下载：

```bash
GOPROXY="https://goproxy.cn,direct" go mod download golang.zx2c4.com/wintun
```

### 6.2 编译命令与输出

```text
S (amd64): env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ... -tags "frps" ...
C (amd64): env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ... -tags "frpc" ...
S (arm64): env CGO_ENABLED=0 GOOS=windows GOARCH=arm64 ... -tags "frps" ...
C (arm64): env CGO_ENABLED=0 GOOS=windows GOARCH=arm64 ... -tags "frpc" ...

全部输出 "amd64 build OK" / "arm64 build OK"
```

### 6.3 产物清单

```text
/root/frp/release/
  frps_windows_amd64.exe   20440064 bytes
  frpc_windows_amd64.exe   16669696 bytes
  frps_windows_arm64.exe   18073088 bytes
  frpc_windows_arm64.exe   14571520 bytes
```

## 7. 产物验证

### 7.1 文件格式（PE 可执行文件）

```text
$ file release/*.exe
frpc_windows_amd64.exe: PE32+ executable (console) x86-64, for MS Windows, 8 sections
frps_windows_amd64.exe: PE32+ executable (console) x86-64, for MS Windows, 8 sections
frpc_windows_arm64.exe: PE32+ executable (console) Aarch64,    for MS Windows, 6 sections
frps_windows_arm64.exe: PE32+ executable (console) Aarch64,    for MS Windows, 6 sections
```

> `PE32+` = 64 位 Windows 可执行格式，(console) = 控制台程序，
> `stripped`（已去除调试符号，配合 `-ldflags "-s -w"`）。

### 7.2 Go 构建信息

```text
$ go version release/frps_windows_amd64.exe
release/frps_windows_amd64.exe: go1.25.0
$ go version release/frpc_windows_arm64.exe
release/frpc_windows_arm64.exe: go1.25.0
```

### 7.3 Web 面板嵌入确认

```text
$ strings release/frps_windows_amd64.exe | grep -c "index-B"
2     # 匹配到嵌入的前端资源文件名，确认面板已内置
```

### 7.4 在 Windows 上运行验证（需在 Windows 机器上执行）

```bat
REM 查看版本
frps.exe --version
frpc.exe --version

REM 服务端启动（依赖 frps.toml 配置，可参考 conf/frps.toml）
frps.exe -c frps.toml
```

## 8. 常见问题排查

| 现象 | 原因 | 解决 |
| --- | --- | --- |
| `wintun ... i/o timeout` | `proxy.golang.org` 不可达，拉不下 Windows 专用依赖 | `export GOPROXY="https://goproxy.cn,direct"` 后重试 |
| 产物是 ELF 而不是 PE | 忘了设置 `GOOS=windows` | 编译时显式指定 `GOOS=windows GOARCH=amd64` |
| need MinGW / C compiler 报错 | CGO 未关闭 | 加 `CGO_ENABLED=0` |
| 运行时提示缺少 DLL（.dll） | 产物为动态链接版 | 必须用 `CGO_ENABLED=0` 静态编译 |
| Windows 下启动闪退 | 缺少配置文件 | 使用 `frps.exe -c xxx.toml` 指定配置或查看日志 |

## 9. 参考

- 官方交叉编译脚本：`/root/frp/Makefile.cross-compiles`
- 上一份 Linux 编译文档：`/root/frp/doc/build.md`
- Windows 专用依赖：`golang.zx2c4.com/wintun`（`go.mod` 间接依赖）