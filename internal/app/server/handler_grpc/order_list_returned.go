package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// ListReturnedOrders creates Order
func (h Controller) ListReturnedOrders(ctx context.Context, in *pb.ListReturnedOrdersReq) (*pb.ListReturnedOrdersResp, error) {
	list, err := h.orderService.ListReturnedOrders(ctx)
	if err != nil {
		log.Printf("[ListReturnedOrders] h.orderService.ListReturnedOrders: %v\n", err)
		return nil, fmt.Errorf("h.orderService.ListReturnedOrders: %w", err)
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
	log.Printf("[ListReturnedOrders] Got list of returned orders. Length = %d.\n", len(list))
	resp := &pb.ListReturnedOrdersResp{
		Orders: orders,
	}
	return resp, nil
}
