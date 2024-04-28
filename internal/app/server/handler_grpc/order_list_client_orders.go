package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// ListClientOrders creates Order
func (h Controller) ListClientOrders(ctx context.Context, in *pb.ListClientOrdersReq) (*pb.ListClientOrdersResp, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		log.Printf("[ListClientOrders] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	list, err := h.orderService.ListClientOrders(ctx, id)
	if err != nil {
		log.Printf("[ListClientOrders] h.orderService.ListClientOrders: %v\n", err)
		return nil, fmt.Errorf("h.orderService.ListClientOrders: %w", err)
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
	log.Printf("[ListClientOrders] Got list of clients orders. Length = %d.\n", len(list))
	resp := &pb.ListClientOrdersResp{
		Orders: orders,
	}
	return resp, nil
}
