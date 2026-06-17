# nginx-autoindex-tui

[![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Tag](https://img.shields.io/github/v/tag/cangyunye/nginx-autoindex-tui)](https://github.com/cangyunye/nginx-autoindex-tui/tags)

一个基于 **Go + Bubble Tea + Lip Gloss** 的终端文件浏览与下载工具，用于访问启用了 `autoindex on;` 的 nginx 静态文件服务器。同时支持 HTML 和 JSON 两种 autoindex 输出格式。

## 功能

- **双格式解析** — 自动识别 HTML（goquery）或 JSON 输出，无需用户配置
- **交互式表格** — 显示索引、类型、文件名、修改时间、文件大小
- **目录导航** — Enter 进入目录，`Esc` / `b` 返回上级，光标循环滚动
- **文件下载** — Enter 选择文件，调系统 `wget` 下载
- **覆盖确认** — 文件已存在时提示 `[Y]覆盖 [N]跳过 [A]始终覆盖`
- **强制覆盖** — `--force` / `-f` flag，或运行时按 `F` 切换
- **主题系统** — `--theme` 可选 `dark` / `light` / `mono`
- **边框风格** — `--border-style` 可选 `normal` / `rounded` / `ascii`
- **根路径约束** — 防止 `../` 逃逸出用户指定的根目录
- **SSL 跳过** — `--insecure` / `-k` 用于自签名证书测试
- **自定义 User-Agent** — `--user-agent` / `-u`

## 快速开始

### 前置要求

- Go 1.21+
- 系统安装 `wget`（用于文件下载）
- 一个启用了 `autoindex on;` 的 nginx 服务器

### 安装与运行

```bash
git clone https://github.com/cangyunye/nginx-autoindex-tui.git
cd nginx-autoindex-tui
make build
./nginx-autoindex-tui http://your-server/files/
```

### 命令行参数

```
./nginx-autoindex-tui [URL] [flags]

Flags:
  -f, --force                强制覆盖已存在文件，跳过确认
  -c, --concurrent int       同时下载的 wget 进程最大数 (default 3)
  -o, --output-dir string    下载文件保存目录（默认当前目录）
  -k, --insecure             跳过 HTTPS 的 SSL 证书验证
  -u, --user-agent string    HTTP 请求的 User-Agent 头
      --border-style string  表格边框风格：normal / rounded / ascii (default "normal")
  -t, --theme string         颜色方案：dark / light / mono (default "dark")
```

### Makefile 命令

```bash
make build    # 优化编译 (-ldflags "-w -s")
make test     # 运行全部测试
make tidy     # 清理依赖
make clean    # 删除二进制 + 构建缓存
```

## 操作说明

| 按键 | 功能 |
|------|------|
| `↑` / `k` | 上移（到顶后跳到最后） |
| `↓` / `j` | 下移（到底后跳到开头） |
| `Enter` | 进入目录 / 下载文件 |
| `Esc` / `b` | 返回上级目录 |
| `F` | 切换强制覆盖开关 |
| `R` | 刷新当前目录 |
| `q` / `Ctrl+C` | 退出 |

## 架构

```
main.go                        # CLI 入口（cobra flag 解析）
  └─ internal/tui/model.go     # Bubble Tea Elm 架构（Init/Update/View）
  └─ internal/tui/styles.go    # 主题系统（dark/light/mono + 边框）
  └─ internal/fetcher/fetcher.go  # HTTP 请求 + URL 解析
  └─ internal/parser/parser.go    # HTML / JSON 双格式解析
tests/
  ├─ parser_test.go            # 解析器表格驱动测试
  ├─ fetcher/fetcher_test.go   # URL 拼接测试
  └─ fixtures/                 # 测试夹具（HTML + JSON）
```

## 技术栈

| 领域 | 工具 |
|------|------|
| 语言 | Go 1.21+ |
| TUI 框架 | [Bubble Tea](https://github.com/charmbracelet/bubbletea)（Elm 架构） |
| 样式与表格 | [Lip Gloss](https://github.com/charmbracelet/lipgloss) |
| HTML 解析 | [goquery](https://github.com/PuerkitoBio/goquery) |
| CLI 框架 | [Cobra](https://github.com/spf13/cobra) |
| 下载 | 系统 `wget`（`os/exec`） |

## 开发

```bash
# 运行测试
make test

# 测试指定包
go test ./tests/... -v
go test ./tests/fetcher/... -v
```

## 许可证

MIT License
