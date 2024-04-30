package handler_grpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/metric"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// GiveOutOrders creates Order
func (h Controller) GiveOutOrders(ctx context.Context, in *pb.GiveOutOrdersReq) (*pb.GiveOutOrdersResp, error) {
	commonAttrs := []attribute.KeyValue{
		attribute.String("ClientID", in.ClientId),
		attribute.String("IDs", strings.Join(in.Ids, ", ")),
	}
	_, span := h.Tracer.Start(
		ctx,
		"CollectorExporter-GiveOutOrders",
		trace.WithAttributes(commonAttrs...))
	defer span.End()

	clientID, err := uuid.Parse(in.ClientId)
	if err != nil {
		log.Printf("[GiveOutOrders] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	ids := make([]uuid.UUID, 0)
	for _, pbID := range in.Ids {
		id, err := uuid.Parse(pbID)
		if err != nil {
			log.Printf("[GiveOutOrders] uuid.Parse: %v\n", err)
			return nil, fmt.Errorf("uuid.Parse: %w", err)
		}
		ids = append(ids, id)
	}
	err = h.orderService.GiveOutOrders(ctx, clientID, ids)
	if err != nil {
		log.Printf("[GiveOutOrders] h.orderService.GiveOutOrders: %v\n", err)
		return nil, fmt.Errorf("h.orderService.GiveOutOrders: %w", err)
	}

	metric.GivenOutOrdersCounterMetric.Add(float64(len(ids)))
	metric.ClientGivenOutOrdersCounterMetric.WithLabelValues(clientID.String()).Inc()

	log.Printf("[GiveOutOrders] Orders are given out\n")
	resp := &pb.GiveOutOrdersResp{
		Comment: "Orders are given out",
	}
	return resp, nil
}
