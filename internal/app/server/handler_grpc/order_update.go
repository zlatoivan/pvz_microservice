package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

// UpdateOrder creates Order
func (h Controller) UpdateOrder(ctx context.Context, in *pb.UpdateOrderReq) (*pb.UpdateOrderResp, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		log.Printf("[UpdateOrder] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	clientID, err := uuid.Parse(in.ClientId)
	if err != nil {
		log.Printf("[UpdateOrder] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	updOrder := model.Order{
		ID:            id,
		ClientID:      clientID,
		Weight:        int(in.Weight),
		Cost:          int(in.Cost),
		StoresTill:    in.StoresTill.AsTime(),
		PackagingType: in.PackagingType,
	}
	err = h.orderService.UpdateOrder(ctx, updOrder)
	if err != nil {
		log.Printf("[UpdateOrder] h.orderService.UpdateOrder: %v\n", err)
		return nil, fmt.Errorf("h.orderService.UpdateOrder: %w", err)
	}

	log.Printf("[UpdateOrder] Order updated\n")
	resp := &pb.UpdateOrderResp{
		Comment: "Order updated",
	}
	return resp, nil
}
