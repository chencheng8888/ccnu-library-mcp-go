package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpServer interface {
	Run(context.Context) error
}

type LocalMcpServer struct {
	server *mcp.Server
}

func NewLocalMcpServer(server *mcp.Server) *LocalMcpServer {
	return &LocalMcpServer{server: server}
}
func (s *LocalMcpServer) Run(ctx context.Context) error {
	return s.server.Run(ctx, &mcp.StdioTransport{})
}

type RemoteMcpServer struct {
	server *mcp.Server
	port   int
}

func NewRemoteMcpServer(server *mcp.Server, port int) *RemoteMcpServer {
	return &RemoteMcpServer{server: server, port: port}
}
func (s *RemoteMcpServer) Run(ctx context.Context) error {
	// 1. 设置信号监听
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT (Ctrl+C) 和 SIGTERM (终止信号)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 2. 构造 HTTP 路由器和处理器
	mux := http.NewServeMux()
	mux.Handle("/mcp", mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return s.server
	}, &mcp.StreamableHTTPOptions{
		JSONResponse: true,
	}))

	// 3. 构造 http.Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	// 4. 启动服务器的协程
	go func() {
		// ListenAndServe 返回错误，除非是 http.ErrServerClosed
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("MCP Server failed to start: %v", err))
		}
	}()

	// 5. 阻塞，直到接收到终止信号
	<-quit

	// 6. 执行优雅关闭
	// 创建一个带有超时限制的上下文，给服务器一定的时间来完成现有请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Server forced to shutdown: %v", err))
	}
	return nil
}
