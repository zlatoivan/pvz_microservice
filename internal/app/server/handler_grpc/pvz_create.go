package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// CreatePVZ creates PVZ
func (h Controller) CreatePVZ(ctx context.Context, in *pb.CreatePVZReq) (*pb.CreatePVZResp, error) {
	newPVZ := model.PVZ{
		Name:     in.Name,
		Address:  in.Address,
		Contacts: in.Contacts,
	}
	id, err := h.pvzService.CreatePVZ(ctx, newPVZ)
	if err != nil {
		log.Printf("[CreatePVZ] h.Service.CreatePVZ: %v\n", err)
		return nil, fmt.Errorf("h.pvzService.CreatePVZ: %w", err)
	}

	log.Printf("[CreatePVZ] PVZ created. id = %s\n", id)
	resp := &pb.CreatePVZResp{
		Id: id.String(),
	}
	return resp, nil
}
