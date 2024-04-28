package handler_grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

// GetPVZByID creates PVZ
func (h Controller) GetPVZByID(ctx context.Context, in *pb.GetPVZByIDReq) (*pb.GetPVZByIDResp, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		log.Printf("[GetPVZByID] uuid.Parse: %v\n", err)
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}
	pvz, err := h.pvzService.GetPVZByID(ctx, id)
	if err != nil {
		log.Printf("[GetPVZByID] h.pvzService.GetPVZByID: %v\n", err)
		return nil, fmt.Errorf("h.pvzService.GetPVZByID: %w", err)
	}

	log.Printf("[GetPVZByID] Got PVZ by ID\n")
	resp := &pb.GetPVZByIDResp{
		Id:       pvz.ID.String(),
		Name:     pvz.Name,
		Address:  pvz.Address,
		Contacts: pvz.Contacts,
	}
	return resp, nil
}
