# frp 编译指南

本文档记录了 frp v0.71.0 在本机（Ubuntu 24.04 / amd64 / Go 1.25.0）上的完整编译过程，
包括环境依赖、工具依赖、编译步骤与实际操作记录。

## 1. 项目概述

- 项目地址：`/root/frp`
- 版本：v0.71.0（commit `4a23aa18`）
- 源码结构：
  - `cmd/frps`：服务端入口
  - `cmd/frpc`：客户端入口
  - `client/`：客户端逻辑
  - `server/`：服务端逻辑
  - `pkg/`：公共库
  - `web/frps`、`web/frpc`：Web 管理面板（Vue3 + Vite）
- 编译产物：
  - `bin/frps`：服务端二进制（约 20.3 MB）
  - `bin/frpc`：客户端二进制（约 16.6 MB）

## 2. 环境与依赖

### 2.1 操作系统

| 项 | 值 |
| --- | --- |
| 系统 | Ubuntu 24.04.4 LTS (Noble Numbat) |
| 架构 | x86_64 (linux/amd64) |

### 2.2 Go 工具链（必需）

| 项 | 值 |
| --- | --- |
| Go 版本 | go1.25.0 linux/amd64 |
| 安装路径 | `/usr/local/go/bin/go` |
| GOPATH | `/root/go` |
| 模块缓存 | `/root/go/pkg/mod` |
| GOMODCACHE | `/root/go/pkg/mod` |

本项目 `go.mod` 中要求 `go 1.25.0`，因此 Go 工具链版本必须 ≥ 1.25.0。

### 2.3 编译工具（必需）

| 工具 | 用途 | 版本 | 安装方式 |
| --- | --- | --- | --- |
| `make` | 执行 Makefile | GNU Make 4.3 | `apt-get install make` |
| `gcc` | 部分 Go 依赖编译需要 | 系统自带 | `apt-get install gcc` |
| `node` | 构建 Web 管理面板 | v24.20.0 | 预装于 `/usr/local/node/bin` |
| `npm` | 安装 Web 前端依赖 | 11.19.0 | 随 node 自带 |
| `git` | 拉取源码/版本信息 | 系统自带 | 预装 |
| `curl` | 网络连通性测试 | 系统自带 | 预装 |

> 说明：
> - 核心 Go 二进制使用 `CGO_ENABLED=0` 静态编译，实际不依赖 gcc，但保守起见建议安装。
> - 若只编译不含 Web 面板的核心二进制，node/npm 不是必需的（见 4.3 节）。
> - Web 面板需要 Node ≥ 18（本项目使用 Vite 7 / Vue 3）。

### 2.4 网络依赖

Go 模块下载需要访问 Go module 代理。本机曾遇到 `proxy.golang.org` 超时，
可使用国内镜像解决：

```bash
export GOPROXY="https://goproxy.cn,direct"
```

连通性测试结果：

| 镜像 | 结果 |
| --- | --- |
| https://proxy.golang.org | 超时（不可用） |
| https://goproxy.cn | 200，延迟 0.07s（推荐） |
| https://goproxy.io | 200，延迟 1.2s |

## 3. 编译步骤

### 3.1 一次性准备

```bash
# 安装 make 与 gcc（root 环境）
apt-get update && apt-get install -y make gcc

# 设置 Go 模块代理（可选，网络不通时使用）
export GOPROXY="https://goproxy.cn,direct"

# 进入项目
cd /root/frp
```

### 3.2 快速编译核心二进制（不含 Web 面板）

```bash
make build
```

等价命令（Makefile 内部）：

```bash
# 服务端：加 noweb tag，不嵌入 Web 面板
env CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -tags "frps,noweb" -o bin/frps ./cmd/frps
# 客户端
env CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -tags "frpc,noweb" -o bin/frpc ./cmd/frpc
```

> `NOWEB_TAG` 由 Makefile 自动判断：当 `web/frps/dist` 与 `web/frpc/dist`
> 目录不存在时自动追加 `,noweb` tag，编译出的二进制不包含 Web 面板。

### 3.3 完整编译（含 Web 管理面板）

```bash
# 1) 安装前端依赖（在 web 目录）
cd /root/frp/web && npm install

# 2) 构建前端资源
cd /root/frp && make web
# 等价于：(make -C web/frps build) + (make -C web/frpc build)
# 产物：web/frps/dist 与 web/frpc/dist

# 3) 重新编译二进制（此时 dist 已存在，不再追加 noweb tag，自动嵌入面板）
make build
```

### 3.4 常用 Makefile 目标

| 目标 | 说明 |
| --- | --- |
| `make all` | 环境检查 + 格式化 + Web 构建 + 核心二进制 |
| `make build` | 仅编译 frps 与 frpc |
| `make frps` | 仅编译服务端 |
| `make frpc` | 仅编译客户端 |
| `make web` | 仅构建 Web 面板 |
| `make test` | 运行单元测试 |
| `make vet` | 运行 go vet（自动带 noweb tag） |
| `make clean` | 清理 bin/ 与临时目录 |

### 3.5 编译参数说明

| 参数 | 含义 |
| --- | --- |
| `CGO_ENABLED=0` | 使用纯静态编译，产物不依赖 glibc，可跨环境运行 |
| `-trimpath` | 去除构建路径信息，保证可复现构建 |
| `-ldflags "-s -w"` | 去掉符号表与 DWARF 调试信息，减小体积 |
| `-tags "frps"` / `"frpc"` | 区分服务端/客户端构建分支 |
| `-o bin/frps` | 输出文件名 |

## 4. 实际操作记录

### 4.1 环境检查

```bash
$ go version
go1.25.0 linux/amd64

$ make --version
GNU Make 4.3

$ node --version
v24.20.0

$ npm --version
11.19.0
```

`make` 初始未安装，已通过 `apt-get install -y make gcc` 补齐。

### 4.2 核心编译日志（片段）

```text
$ make build
env CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -tags "frpc,noweb" -o bin/frpc ./cmd/frpc
env CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -tags "frps,noweb" -o bin/frps ./cmd/frps
```

首次编译会下载大量依赖（约 40+ 个 module），视网络情况耗时数分钟。

出现的网络故障及解决：

```text
ERROR: ... Get "https://proxy.golang.org/...": dial tcp ...: i/o timeout
解决：export GOPROXY="https://goproxy.cn,direct" 后重试
```

### 4.3 Web 面板编译日志（片段）

```text
$ cd /root/frp/web && npm install
added 180+ packages   # 实际数量以安装为准
1 moderate severity vulnerability   # 来自传递依赖，不影响构建

$ make web
> frpc-dashboard@0.0.1 build-only
> vite build
✓ built in 6.65s
dist/index.html            0.38 kB │ gzip:   0.27 kB
dist/index-KasTiy2b.css   94.62 kB │ gzip:  14.99 kB
dist/index-Bjhz7A9c.js   383.65 kB │ gzip: 122.44 kB
```

### 4.4 编译后完整构建 Web 面板并重编二进制

```bash
# 生成 dist 后重新编译，构建命令变为（已无 noweb tag）：
env CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -tags "frps" -o bin/frps ./cmd/frps
env CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -tags "frpc" -o bin/frpc ./cmd/frpc
```

## 5. 产物验证

```bash
$ ls -la bin/
-rwxr-xr-x  1 root root 20283576  bin/frps
-rwxr-xr-x  1 root root 16552120  bin/frpc

$ file bin/frps
ELF 64-bit LSB executable, x86-64, statically linked, stripped
```

### 5.1 版本验证

```bash
$ ./bin/frps --version
0.71.0
$ ./bin/frpc --version
0.71.0
```

### 5.2 运行冒烟测试（服务端 + Web 面板）

使用最小配置启动 `frps`：

```toml
bindPort = 7000
webServer.addr = "127.0.0.1"
webServer.port = 7500
webServer.user = "admin"
webServer.password = "admin"
```

日志确认：

```text
[I] frps uses config file: frps.toml
[I] frps tcp listen on 0.0.0.0:7000
[I] frps started successfully
[I] dashboard listen on 127.0.0.1:7500
```

Web 面板访问：

```bash
$ curl -s -o /dev/null -w "HTTP %{http_code}\n" http://127.0.0.1:7500/
HTTP 401     # 未认证，401 说明面板已正常嵌入并对外服务
```

## 6. 常见问题排查

| 现象 | 原因 | 解决 |
| --- | --- | --- |
| `dial tcp ...: i/o timeout` | `proxy.golang.org` 不可达 | 换用 `export GOPROXY="https://goproxy.cn,direct"` |
| `make: command not found` | 未安装 make | `apt-get install -y make` |
| `No such file or directory: web/frps/dist` | 尚未构建 Web 面板 | 需先 `npm install && make web`，或使用 noweb 版本 |
| `CGO_ENABLED` 相关报错 | 缺少 C 编译器 | `apt-get install -y gcc` |
| 端口被占用 | 上次进程未退出 | 先 `pkill -f frps` 再启动 |

## 7. 参考

- 编译入口：`/root/frp/Makefile`
- 模块信息：`/root/frp/go.mod`
- Web 构建脚本：`/root/frp/web/frps/Makefile`、`/root/frp/web/frpc/Makefile`
- 项目文档：`/root/frp/README_zh.md`、`/root/frp/README.md`