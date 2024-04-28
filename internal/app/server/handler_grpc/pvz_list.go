package handler_grpc

import (
	"context"
	"fmt"
	"log"

	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// ListPVZs creates PVZ
func (h Controller) ListPVZs(ctx context.Context, in *pb.ListPVZsReq) (*pb.ListPVZsResp, error) {
	list, err := h.pvzService.ListPVZs(ctx)
	if err != nil {
		log.Printf("[ListPVZs] h.pvzService.ListPVZs: %v\n", err)
		return nil, fmt.Errorf("h.pvzService.ListPVZs: %w", err)
	}

	pvzs := make([]*pb.ModelPVZ, 0)
	for _, p := range list {
		pvz := &pb.ModelPVZ{
			Id:       p.ID.String(),
			Name:     p.Name,
			Address:  p.Address,
			Contacts: p.Contacts,
		}
		pvzs = append(pvzs, pvz)
	}
	log.Printf("[ListPVZs] Got list of PVZs. Length = %d.\n", len(list))
	resp := &pb.ListPVZsResp{
		Pvzs: pvzs,
	}
	return resp, nil
}
