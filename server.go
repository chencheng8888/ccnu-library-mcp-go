package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)


type McpServer interface{
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
	port  int
}

func NewRemoteMcpServer(server *mcp.Server,port int) *RemoteMcpServer {
	return &RemoteMcpServer{server: server, port: port}
}
func (s *RemoteMcpServer) Run(ctx context.Context) error {
	handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
		return s.server
	})
	return http.ListenAndServe(fmt.Sprintf(":%d",s.port), handler)
}
