package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	libraryreservations "github.com/chencheng8888/ccnu-library-reservations"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var mcpServerType string

var port int


var student_config_file string

var init_students  map[string]string

func init() {
	// 注册命令行参数
	flag.StringVar(&mcpServerType, "type", "stdio", "MCP server type: stdio or remote")
	flag.IntVar(&port, "port", 8080, "Port for SSE server")
	flag.StringVar(&student_config_file, "conf", "", "Path to student config file,the file should be a json file, like {\"stu_id1\":\"pwd1\",\"stu_id2\":\"pwd2\"}")
}

func main() {
	flag.Parse()

	if len(student_config_file) > 0 {
		init_students = readStudentConfig(student_config_file)
	}

	a := libraryreservations.NewAuther()

	if len(init_students)>0 {
		for stuId,pwd := range init_students {
			// fmt.Println("Initializing student:",stuId,pwd)
			a.StoreStuInfo(context.Background(),stuId,pwd)
		}
	}

	h := NewCCNULibHandler(a)

	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "ccnu-library-mcp", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "register", Description: "注册学生信息，其他操作都需要至少执行过或曾经执行过一次这个接口"}, h.Register)
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


func readStudentConfig(configFile string) map[string]string {
	// 读取文件内容
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil
	}

	// 解析 JSON 内容到 map
	var students map[string]string
	err = json.Unmarshal(data, &students)
	if err != nil {
		return nil
	}

	return students
}