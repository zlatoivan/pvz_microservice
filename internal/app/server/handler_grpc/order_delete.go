package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// DeleteOrder creates Order
func (h Controller) DeleteOrder(ctx context.Context, in *pb.DeleteOrderReq) (*pb.DeleteOrderResp, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		log.Printf("[DeleteOrder] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	err = h.orderService.DeleteOrder(ctx, id)
	if err != nil {
		log.Printf("[DeleteOrder] h.orderService.DeleteOrder: %v\n", err)
		return nil, fmt.Errorf("h.orderService.DeleteOrder: %w", err)
	}

	log.Printf("[DeleteOrder] Order deleted\n")
	resp := &pb.DeleteOrderResp{
		Comment: "Order deleted",
	}
	return resp, nil
}
