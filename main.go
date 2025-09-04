package main

import (
	"ccnu-library-mcp-go/internal/auther"
	"ccnu-library-mcp-go/internal/reverser"
	"context"
	"flag"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var mcpServerType string

var port int

func init() {
	// 注册命令行参数
	flag.StringVar(&mcpServerType, "type", "stdio", "MCP server type: stdio or sse")
	flag.IntVar(&port, "port", 8080, "Port for SSE server")
}

func main() {
	a := auther.NewAuther()
	r := reverser.NewReverser(a)
	h := NewCCNULibHandler(r)

	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "ccnu-library-mcp", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "register", Description: "注册学生信息，只有这样才能进行其他操作"}, h.Register)
	mcp.AddTool(server, &mcp.Tool{Name: "get seat info", Description: "获取华中师范大学的图书馆的座位占用情况"}, h.GetSeats)
	mcp.AddTool(server, &mcp.Tool{Name: "reverse seat", Description: "预约座位"}, h.Reverse)

	var mcpServer McpServer

	switch mcpServerType {
	case "stdio":
		mcpServer = NewLocalMcpServer(server)
	case "sse":
		mcpServer = NewRemoteMcpServer(server, port)
	default:
		log.Fatalf("unknown mcp server type: %s", mcpServerType)
	}

	// Run the server over stdin/stdout, until the client disconnects
	if err := mcpServer.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
