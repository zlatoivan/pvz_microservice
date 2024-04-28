package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// GetOrderByID creates Order
func (h Controller) GetOrderByID(ctx context.Context, in *pb.GetOrderByIDReq) (*pb.GetOrderByIDResp, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		log.Printf("[GetOrderByID] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("h.orderService.GetOrderByID: %w", err)
	}
	order, err := h.orderService.GetOrderByID(ctx, id)
	if err != nil {
		log.Printf("[GetOrderByID] h.orderService.GetOrderByID: %v\n", err)
		return nil, fmt.Errorf("h.orderService.GetOrderByID: %w", err)
	}

	log.Printf("[GetOrderByID] Got Order by ID\n")
	resp := &pb.GetOrderByIDResp{
		Id:            in.Id,
		ClientId:      order.ClientID.String(),
		Weight:        int64(order.Weight),
		Cost:          int64(order.Cost),
		StoresTill:    timestamppb.New(order.StoresTill),
		GiveOutTime:   timestamppb.New(order.GiveOutTime),
		IsReturned:    order.IsReturned,
		PackagingType: order.PackagingType,
	}
	return resp, nil
}
