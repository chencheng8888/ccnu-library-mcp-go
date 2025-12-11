package main

import (
	"ccnu-library-mcp-go/pkg"
	"context"
	"fmt"
	libraryreservations "github.com/chencheng8888/ccnu-library-reservations"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CCNULibHandler struct {
	au libraryreservations.Auther
	r  libraryreservations.Reverser
}

func NewCCNULibHandler(au libraryreservations.Auther) *CCNULibHandler {
	return &CCNULibHandler{
		au: au,
		r:  libraryreservations.NewReverser(au),
	}
}

type RegisterParams struct {
	StuID string `json:"stu_id" jsonschema:"学号"`
	Pwd   string `json:"pwd" jsonschema:"密码"`
}

func (h *CCNULibHandler) Register(ctx context.Context, req *mcp.CallToolRequest, args RegisterParams) (*mcp.CallToolResult, any, error) {
	err := h.au.StoreStuInfo(ctx, args.StuID, args.Pwd)
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

type GetSeatsParams struct {
	StuID         string `json:"stu_id" jsonschema:"学号"`
	RoomName      string `json:"room_name" jsonschema:"楼层名,目前只有n1,n1m,n2; n1代表图书馆南湖分馆一楼开敞座位区,n1m代表图书馆南湖分馆一楼中庭开敞座位区,n2代表图书馆南湖分馆二楼开敞座位区"`
	StartTime     string `json:"start_time" jsonschema:"开始时间,格式为2025-06-01 10:00 分钟必须是5的倍数"`
	EndTime       string `json:"end_time" jsonschema:"结束时间,格式为2025-06-01 10:00 分钟必须是5的倍数"`
	OnlyAvailable bool   `json:"only_available" jsonschema:"是否只返回空闲座位"`
}

func (h *CCNULibHandler) GetSeats(ctx context.Context, req *mcp.CallToolRequest, args GetSeatsParams) (*mcp.CallToolResult, any, error) {
	startTime, _ := pkg.TransferStringToTime(args.StartTime, pkg.FORMAT2)
	endTime, _ := pkg.TransferStringToTime(args.EndTime, pkg.FORMAT2)

	seats, err := h.r.GetSeatsByTime(ctx, args.StuID, pkg.Rooms[args.RoomName],
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

type ReverseParams struct {
	StuID     string `json:"stu_id" jsonschema:"学号"`
	SeatID    string `json:"seat_id" jsonschema:"座位号"`
	StartTime string `json:"start_time" jsonschema:"开始时间,格式为2025-06-01 10:00 分钟必须是5的倍数"`
	EndTime   string `json:"end_time" jsonschema:"结束时间,格式为2025-06-01 10:00 分钟必须是5的倍数"`
}

func (h *CCNULibHandler) Reverse(ctx context.Context, req *mcp.CallToolRequest, args ReverseParams) (*mcp.CallToolResult, any, error) {

	startTime, _ := pkg.TransferStringToTime(args.StartTime, pkg.FORMAT2)
	endTime, _ := pkg.TransferStringToTime(args.EndTime, pkg.FORMAT2)

	err := h.r.Reverse(ctx, args.StuID, args.SeatID, startTime, endTime)
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
