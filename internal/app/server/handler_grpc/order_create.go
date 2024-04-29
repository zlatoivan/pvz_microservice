package handler_grpc

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

// CreateOrder creates Order
func (h Controller) CreateOrder(ctx context.Context, in *pb.CreateOrderReq) (*pb.CreateOrderResp, error) {
	commonAttrs := []attribute.KeyValue{
		attribute.String("ClientID", in.ClientId),
		attribute.String("Weight", strconv.Itoa(int(in.Weight))),
		attribute.String("Cost", strconv.Itoa(int(in.Cost))),
		attribute.String("StoresTill", in.StoresTill.String()),
		attribute.String("PackagingType", in.PackagingType),
	}
	_, span := h.Tracer.Start(
		ctx,
		"CollectorExporter-CreateOrder",
		trace.WithAttributes(commonAttrs...))
	defer span.End()

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
