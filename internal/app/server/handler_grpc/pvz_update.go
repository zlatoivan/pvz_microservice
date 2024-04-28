package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// UpdatePVZ creates PVZ
func (h Controller) UpdatePVZ(ctx context.Context, in *pb.UpdatePVZReq) (*pb.UpdatePVZResp, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		log.Printf("[UpdatePVZ] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	updPVZ := model.PVZ{
		ID:       id,
		Name:     in.Name,
		Address:  in.Address,
		Contacts: in.Contacts,
	}
	err = h.pvzService.UpdatePVZ(ctx, updPVZ)
	if err != nil {
		log.Printf("[UpdatePVZ] h.pvzService.UpdatePVZ: %v\n", err)
		return nil, fmt.Errorf("h.pvzService.UpdatePVZ: %w", err)
	}

	log.Printf("[UpdatePVZ] PVZ updated\n")
	resp := &pb.UpdatePVZResp{
		Comment: "PVZ updated",
	}
	return resp, nil
}
