# CCNU Library MCP Server

华中师范大学图书馆座位预约 MCP (Model Context Protocol) 服务器

## 项目简介

这是一个基于 Go 语言开发的 MCP 服务器，提供华中师范大学图书馆座位预约相关的工具功能。通过 MCP 协议，可以与各种 AI 助手（如 Claude Desktop、VSCode 等）集成，实现图书馆座位的查询和预约功能。

该项目支持两种运行模式：
- **stdio 模式**: 通过标准输入输出与客户端通信，适用于本地集成
- **SSE 模式**: 通过 HTTP 服务器提供 Server-Sent Events 接口，适用于远程调用

## 功能特性

- **学生注册**: 注册学生信息，保存学号和密码
- **座位查询**: 查询指定时间段内的座位占用情况
- **座位预约**: 预约指定座位和时间段
- **支持区域**: 南湖分馆一楼开敞座位区、一楼中庭开敞座位区、二楼开敞座位区

## 系统要求

- Go 1.24.4 或更高版本
- 华中师范大学有效学号和密码
- 支持的操作系统: Windows, Linux, macOS

## 安装和运行

### 1. 克隆项目

```bash
git clone git@github.com:chencheng8888/ccnu-library-mcp-go.git
cd ccnu-library-mcp-go
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 构建项目

#### Windows
```bash
go build -o ccnu-library-mcp-go.exe
```

#### Linux/macOS
```bash
go build -o ccnu-library-mcp-go
```

### 4. 运行服务器

#### Stdio 模式（默认,可不加type）
```bash
# Windows
./ccnu-library-mcp-go.exe -type stdio

# Linux/macOS  
./ccnu-library-mcp-go -type stdio
```

#### SSE 模式（用于 HTTP 接口,端口默认8080）
```bash
# 默认端口 8080
./ccnu-library-mcp-go -type sse

# 自定义端口
./ccnu-library-mcp-go -type sse -port 3000
```

### 5. MCP 客户端配置


**studio模式**
```json
{
  "mcpServers": {
    "ccnu-library": {
      "disabled": false,
      "timeout": 60,
      "type": "stdio",
      "command": "C:\\path\\to\\ccnu-library-mcp-go\\ccnu-library-mcp-go.exe",
      "args": ["-type", "stdio"]
    }
  }
}
```

**sse模式**
```json
{
  "mcpServers": {
    "ccnu-library-mcp-remote": {
      "url": "http://addr:port",
      "disabled": false
    }
  }
}
```

## MCP 工具说明

### 1. 注册工具 (register)

注册学生信息

**参数:**
- `stu_id`: 学号 (必填)
- `pwd`: 密码 (必填)

**示例:**
```json
{
  "stu_id": "2021xxxxxx",
  "pwd": "your_password"
}
```

### 2. 座位查询工具 (get seat info)

查询指定时间段内的座位占用情况。

**参数:**
- `stu_id`: 学号 (必填)
- `room_name`: 楼层名 (必填，可选值: n1, n1m, n2)
  - `n1`: 图书馆南湖分馆一楼开敞座位区
  - `n1m`: 图书馆南湖分馆一楼中庭开敞座位区  
  - `n2`: 图书馆南湖分馆二楼开敞座位区
- `start_time`: 开始时间 (必填，格式: 2025-06-01 10:00，分钟必须是5的倍数)
- `end_time`: 结束时间 (必填，格式: 2025-06-01 12:00，分钟必须是5的倍数)
- `only_available`: 是否只返回空闲座位 (必填)

**示例:**
```json
{
  "stu_id": "2021xxxxxx",
  "room_name": "n1",
  "start_time": "2025-06-01 10:00",
  "end_time": "2025-06-01 12:00",
  "only_available": true
}
```

### 3. 座位预约工具 (reverse seat)

预约指定座位和时间段。

**参数:**
- `stu_id`: 学号 (必填)
- `seat_id`: 座位号 (必填)
- `start_time`: 开始时间 (必填，格式: 2025-06-01 10:00，分钟必须是5的倍数)
- `end_time`: 结束时间 (必填，格式: 2025-06-01 12:00，分钟必须是5的倍数)

**示例:**
```json
{
  "stu_id": "2021xxxxxx",
  "seat_id": "123456",
  "start_time": "2025-06-01 10:00",
  "end_time": "2025-06-01 12:00"
}
```

## 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-type` | `stdio` | 服务器类型，可选值: `stdio`, `sse` |
| `-port` | `8080` | SSE 模式下的 HTTP 服务器端口 |

**示例:**
```bash
# 使用 stdio 模式（默认）
./ccnu-library-mcp-go

# 使用 SSE 模式，端口 3000
./ccnu-library-mcp-go -type sse -port 3000
```

## 项目结构

```
ccnu-library-mcp-go/
├── internal/
│   ├── auther/          # 认证模块
│   │   ├── auther.go    # 学生信息存储和认证
│   │   └── auther_test.go
│   └── reverser/        # 预约模块
│       ├── reverser.go  # 座位查询和预约核心逻辑
│       └── reverser_test.go
├── pkg/
│   └── tool.go          # 工具函数、时间处理、房间映射
├── main.go              # 主程序入口，MCP 服务器初始化
├── handler.go           # MCP 工具处理器实现
├── server.go            # 服务器接口，支持 stdio 和 SSE 模式
├── go.mod              # Go 模块定义
├── go.sum              # 依赖校验和
├── LICENSE             # Apache 2.0 许可证
└── README.md           # 项目说明文档
```


## 开发说明

### 运行测试
```bash
# 运行所有测试
go test ./...

# 运行指定模块测试
go test ./internal/auther
go test ./internal/reverser
```

### 代码质量
```bash
# 代码格式化
go fmt ./...

# 代码检查
go vet ./...

# 依赖整理
go mod tidy
```


## 许可证

本项目采用 Apache License 2.0 开源协议，详见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 注意事项

1. **认证要求**: 使用前请确保拥有华中师范大学的有效学号和密码
2. **时间格式**: 时间参数中的分钟必须是5的倍数（如 10:00, 10:05, 10:10 等）
3. **房间代码**: 
   - `n1`: 南湖分馆一楼开敞座位区 (ID: 101699179)
   - `n1m`: 南湖分馆一楼中庭开敞座位区 (ID: 101699187)  
   - `n2`: 南湖分馆二楼开敞座位区 (ID: 101699189)
4. **使用规范**: 请遵守图书馆的使用规定，合理使用预约功能
5. **时区设置**: 所有时间均使用上海时区 (Asia/Shanghai)
6. **用途声明**: 本项目仅用于学习和研究目的
