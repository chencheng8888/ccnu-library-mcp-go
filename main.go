package main

import (
	"context"
	"flag"
	libraryreservations "github.com/chencheng8888/ccnu-library-reservations"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var mcpServerType string

var port int

func init() {
	// 注册命令行参数
	flag.StringVar(&mcpServerType, "type", "stdio", "MCP server type: stdio or remote")
	flag.IntVar(&port, "port", 8080, "Port for SSE server")
}

func main() {
	flag.Parse()
	a := libraryreservations.NewAuther()
	h := NewCCNULibHandler(a)

	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "ccnu-library-mcp", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "register", Description: "注册学生信息"}, h.Register)
	mcp.AddTool(server, &mcp.Tool{Name: "get seat info", Description: "获取华中师范大学的图书馆的座位占用情况"}, h.GetSeats)
	mcp.AddTool(server, &mcp.Tool{Name: "reverse seat", Description: "预约座位"}, h.Reverse)

	var mcpServer McpServer

	switch mcpServerType {
	case "stdio":
		mcpServer = NewLocalMcpServer(server)
	case "remote":
		mcpServer = NewRemoteMcpServer(server, port)
	default:
		log.Fatalf("unknown mcp server type: %s", mcpServerType)
	}

	// Run the server over stdin/stdout, until the client disconnects
	if err := mcpServer.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
