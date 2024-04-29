package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/metrics"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// DeletePVZ creates PVZ
func (h Controller) DeletePVZ(ctx context.Context, in *pb.DeletePVZReq) (*pb.DeletePVZResp, error) {
	commonAttrs := []attribute.KeyValue{
		attribute.String("ID", in.Id),
	}
	_, span := h.Tracer.Start(
		ctx,
		"CollectorExporter-DeletePVZ",
		trace.WithAttributes(commonAttrs...))
	defer span.End()

	id, err := uuid.Parse(in.Id)
	if err != nil {
		log.Printf("[DeletePVZ] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	err = h.pvzService.DeletePVZ(ctx, id)
	if err != nil {
		log.Printf("[DeletePVZ] h.pvzService.DeletePVZ: %v\n", err)
		return nil, fmt.Errorf("h.pvzService.DeletePVZ: %w", err)
	}

	metrics.DeletedPVZsCounterMetric.Inc()

	log.Printf("[DeletePVZ] PVZ deleted\n")
	resp := &pb.DeletePVZResp{
		Comment: "PVZ deleted",
	}
	return resp, nil
}
