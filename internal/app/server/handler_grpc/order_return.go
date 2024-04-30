package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/metric"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// ReturnOrder creates Order
func (h Controller) ReturnOrder(ctx context.Context, in *pb.ReturnOrderReq) (*pb.ReturnOrderResp, error) {
	commonAttrs := []attribute.KeyValue{
		attribute.String("ClientID", in.ClientId),
		attribute.String("ID", in.Id),
	}
	_, span := h.Tracer.Start(
		ctx,
		"CollectorExporter-ReturnOrder",
		trace.WithAttributes(commonAttrs...))
	defer span.End()

	clientID, err := uuid.Parse(in.ClientId)
	if err != nil {
		log.Printf("[ReturnOrder] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		log.Printf("[ReturnOrder] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	err = h.orderService.ReturnOrder(ctx, clientID, id)
	if err != nil {
		log.Printf("[ReturnOrder] h.orderService.ReturnOrder: %v\n", err)
		return nil, fmt.Errorf("h.orderService.ReturnOrder: %w", err)
	}

	metric.ReturnedOrdersCounterMetric.Inc()

	log.Println("[ReturnOrder] Order is returned")
	resp := &pb.ReturnOrderResp{
		Comment: "Order is returned",
	}
	return resp, nil
}
