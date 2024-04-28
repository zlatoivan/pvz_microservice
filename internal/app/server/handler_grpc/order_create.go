package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

// CreateOrder creates Order
func (h Controller) CreateOrder(ctx context.Context, in *pb.CreateOrderReq) (*pb.CreateOrderResp, error) {
	clientID, err := uuid.Parse(in.ClientId)
	if err != nil {
		log.Printf("[CreateOrder] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	newOrder := model.Order{
		ClientID:      clientID,
		Weight:        int(in.Weight),
		Cost:          int(in.Cost),
		StoresTill:    in.StoresTill.AsTime(),
		PackagingType: in.PackagingType,
	}
	id, err := h.orderService.CreateOrder(ctx, newOrder)
	if err != nil {
		log.Printf("[CreateOrder] h.orderService.CreateOrder: %v\n", err)
		return nil, fmt.Errorf("h.orderService.CreateOrder: %w", err)
	}

	log.Printf("[CreateOrder] Order created. id = %s\n", id)
	resp := &pb.CreateOrderResp{
		Id: id.String(),
	}
	return resp, nil
}
