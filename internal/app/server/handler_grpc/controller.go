package handler_grpc

import (
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/pvz"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

type Controller struct {
	pvzService   pvz.Service
	orderService order.Service
	pb.UnimplementedApiV1Server
}

func New(pvzService pvz.Service, orderService order.Service) Controller {
	controller := Controller{
		pvzService:   pvzService,
		orderService: orderService,
	}
	return controller
}
