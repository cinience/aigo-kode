# Go Anon Kode

Go Anon Kode 是一个基于终端的 AI 编码工具，可以使用任何支持 OpenAI 风格 API 的模型。这是对原始 [anon-kode](https://github.com/dnakov/anon-kode) 项目的 Golang 重新实现，并增加了 Web 界面。

## 功能特点

- 修复代码问题
- 解释函数功能
- 运行测试、Shell 命令等
- 提供终端和 Web 两种界面
- 支持多种 AI 模型（通过 OpenAI 兼容 API）

## 系统要求

- Go 1.18+
- Node.js 18+（用于 Web 界面）

## 安装

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/user/go-anon-kode.git
cd go-anon-kode

# 构建 CLI 应用
cd cmd/cli
go build -o ../../bin/go-anon-kode-cli

# 构建服务器应用
cd ../server
go build -o ../../bin/go-anon-kode-server

# 构建 Web 前端
cd ../../web
npm install
npm run build
```

## 使用方法

### CLI 模式

```bash
# 启动交互式终端
./bin/go-anon-kode-cli

# 直接执行命令
./bin/go-anon-kode-cli file-read /path/to/file.txt
./bin/go-anon-kode-cli bash "ls -la"
```

### Web 模式

```bash
# 启动 Web 服务器
./bin/go-anon-kode-server

# 然后在浏览器中访问
# http://localhost:8080
```

## 配置

首次运行时，应用会引导您完成配置过程。您需要提供 OpenAI API 密钥或其他兼容的 API 密钥。

配置文件存储在：
- 全局配置：`~/.aigo-kode/config.json`
- 项目配置：`<project-dir>/.aigo-kode.json`

## 主要功能

### 文件操作

- 读取文件内容
- 写入/编辑文件
- 搜索文件内容
- 查找匹配文件

### 代码分析

- 解释代码功能
- 提供代码改进建议
- 修复代码问题

### Shell 命令

- 执行 Shell 命令
- 查看命令输出

## 架构

该项目采用模块化架构，主要组件包括：

1. **核心库 (core)**: 包含共享的业务逻辑和接口
2. **CLI 应用 (cmd/cli)**: 终端界面应用
3. **Web 服务 (cmd/server)**: Web API 和前端服务
4. **工具实现 (tools)**: 各种工具的具体实现
5. **AI 服务 (ai)**: AI 模型集成
6. **配置 (config)**: 配置管理

## 开发

### 运行测试

```bash
# 运行单元测试
go test ./...

# 运行集成测试
./integration_test.sh
```

### 添加新工具

要添加新工具，需要实现 `core.Tool` 接口并将其注册到 `tools.ToolRegistry`。

## 许可证

与原始 anon-kode 项目相同的许可证。

## 致谢

- 原始 [anon-kode](https://github.com/dnakov/anon-kode) 项目
- 所有使用的开源库和框架
