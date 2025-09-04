# CCNU Library MCP Server

华中师范大学图书馆座位预约 MCP (Model Context Protocol) 服务器

## 项目简介

这是一个基于 Go 语言开发的 MCP 服务器，提供华中师范大学图书馆座位预约相关的工具功能。通过 MCP 协议，可以与各种 AI 助手集成，实现图书馆座位的查询和预约功能。

## 功能特性

- **学生注册**: 注册学生信息，保存学号和密码
- **座位查询**: 查询指定时间段内的座位占用情况
- **座位预约**: 预约指定座位和时间段
- **支持区域**: 南湖分馆一楼开敞座位区、一楼中庭开敞座位区、二楼开敞座位区

## 安装要求

- Go 1.24.4 或更高版本
- 华中师范大学有效学号和密码

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

### 3. 构建项目(根据你的系统,以windows为例)

```bash
go build -o ccnu-library-mcp-go.exe
```

### 4. 在mcp client中配置

```json
{
  "mcpServers": {
    "ccnu-library": {
      "disabled": false,
      "timeout": 60,
      "type": "stdio",
      "command": "xxx\\ccnu-library-mcp-go\\ccnu-library-mcp-go.exe",
      "args": []
    }
  }
}

```

## MCP 工具说明

### 1. 注册工具 (register)

注册学生信息，必须先注册才能使用其他功能。

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

## 项目结构

```
ccnu-library-mcp-go/
├── internal/
│   ├── auther/          # 认证模块
│   │   ├── auther.go    # 认证实现
│   │   └── auther_test.go
│   └── reverser/        # 预约模块
│       ├── reverser.go  # 预约实现
│       └── reverser_test.go
├── pkg/
│   └── tool.go          # 工具函数和常量
├── main.go              # 主程序入口
├── go.mod              # Go 模块定义
├── go.sum              # 依赖校验和
├── LICENSE             # Apache 2.0 许可证
└── README.md           # 项目说明文档
```

## 技术栈

- **语言**: Go 1.24.4
- **MCP SDK**: github.com/modelcontextprotocol/go-sdk v0.3.1
- **HTTP 客户端**: 标准库 net/http
- **HTML 解析**: github.com/PuerkitoBio/goquery v1.10.3

## 开发说明

### 构建测试

```bash
go test ./...
```

### 代码格式化

```bash
go fmt ./...
```

### 依赖检查

```bash
go mod tidy
```

## 许可证

本项目采用 Apache License 2.0 开源协议，详见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 注意事项

1. 使用前请确保拥有华中师范大学的有效学号和密码
2. 时间参数中的分钟必须是5的倍数（如 10:00, 10:05, 10:10 等）
3. 请遵守图书馆的使用规定，合理使用预约功能
4. 本项目仅用于学习和研究目的

## 联系方式

如有问题或建议，请通过 GitHub Issues 提交反馈。
