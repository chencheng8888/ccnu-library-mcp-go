package main

import (
	"ccnu-library-mcp-go/internal/auther"
	"ccnu-library-mcp-go/internal/reverser"
	"ccnu-library-mcp-go/pkg"
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CCNULibMcpServer struct {
	r *reverser.Reverser
}

func NewCCNULibMcpServer(r *reverser.Reverser) *CCNULibMcpServer {
	return &CCNULibMcpServer{r: r}
}

var (
	server *CCNULibMcpServer
)

func init() {
	a := auther.NewAuther()
	r := reverser.NewReverser(a)
	server = NewCCNULibMcpServer(r)
}

type GetSeatsParams struct {
	StuID         string `json:"stu_id" jsonschema:"学号"`
	RoomName      string `json:"room_name" jsonschema:"楼层名,目前只有n1,n1m,n2; n1代表图书馆南湖分馆一楼开敞座位区,n1m代表图书馆南湖分馆一楼中庭开敞座位区,n2代表图书馆南湖分馆二楼开敞座位区"`
	StartTime     string `json:"start_time" jsonschema:"开始时间,格式为2025-06-01 10:00 分钟必须是5的倍数"`
	EndTime       string `json:"end_time" jsonschema:"结束时间,格式为2025-06-01 10:00 分钟必须是5的倍数"`
	OnlyAvailable bool   `json:"only_available" jsonschema:"是否只返回空闲座位"`
}

func GetSeats(ctx context.Context, req *mcp.CallToolRequest, args GetSeatsParams) (*mcp.CallToolResult, any, error) {
	startTime, _ := pkg.TransferStringToTime(args.StartTime, pkg.FORMAT2)
	endTime, _ := pkg.TransferStringToTime(args.EndTime, pkg.FORMAT2)

	seats, err := server.r.GetSeatsByTime(ctx, args.StuID, pkg.Rooms[args.RoomName],
		startTime, endTime, args.OnlyAvailable)
	if err != nil || seats == nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%+v", seats),
			},
		},
	}, nil, nil
}

type RegisterParams struct {
	StuID string `json:"stu_id" jsonschema:"学号"`
	Pwd   string `json:"pwd" jsonschema:"密码"`
}

func Register(ctx context.Context, req *mcp.CallToolRequest, args RegisterParams) (*mcp.CallToolResult, any, error) {
	err := server.r.StoreStuInfo(ctx, args.StuID, args.Pwd)
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "注册成功",
			},
		},
	}, nil, nil
}

type ReverseParams struct {
	StuID     string `json:"stu_id" jsonschema:"学号"`
	SeatID    string `json:"seat_id" jsonschema:"座位号"`
	StartTime string `json:"start_time" jsonschema:"开始时间,格式为2025-06-01 10:00 分钟必须是5的倍数"`
	EndTime   string `json:"end_time" jsonschema:"结束时间,格式为2025-06-01 10:00 分钟必须是5的倍数"`
}

func Reverse(ctx context.Context, req *mcp.CallToolRequest, args ReverseParams) (*mcp.CallToolResult, any, error) {

	startTime, _ := pkg.TransferStringToTime(args.StartTime, pkg.FORMAT2)
	endTime, _ := pkg.TransferStringToTime(args.EndTime, pkg.FORMAT2)

	err := server.r.Reverse(ctx, args.StuID, args.SeatID, startTime, endTime)
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "预约成功",
			},
		},
	}, nil, nil
}

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "ccnu-library-mcp", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "register", Description: "注册学生信息，只有这样才能进行其他操作"}, Register)
	mcp.AddTool(server, &mcp.Tool{Name: "get seat info", Description: "获取华中师范大学的图书馆的座位占用情况"}, GetSeats)
	mcp.AddTool(server, &mcp.Tool{Name: "reverse seat", Description: "预约座位"}, Reverse)
	// Run the server over stdin/stdout, until the client disconnects
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
