package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// ListOrders creates Order
func (h Controller) ListOrders(ctx context.Context, in *pb.ListOrdersReq) (*pb.ListOrdersResp, error) {
	list, err := h.orderService.ListOrders(ctx)
	if err != nil {
		log.Printf("[ListOrders] h.orderService.ListOrders: %v\n", err)
		return nil, fmt.Errorf("h.orderService.ListOrders: %w", err)
	}

	orders := make([]*pb.ModelOrder, 0)
	for _, order := range list {
		pbOrder := &pb.ModelOrder{
			Id:            order.ID.String(),
			ClientId:      order.ClientID.String(),
			Weight:        int64(order.Weight),
			Cost:          int64(order.Cost),
			StoresTill:    timestamppb.New(order.StoresTill),
			GiveOutTime:   timestamppb.New(order.GiveOutTime),
			IsReturned:    order.IsReturned,
			PackagingType: order.PackagingType,
		}
		orders = append(orders, pbOrder)
	}
	log.Printf("[ListOrders] Got list of Orders. Length = %d.\n", len(list))
	resp := &pb.ListOrdersResp{
		Orders: orders,
	}
	return resp, nil
}
