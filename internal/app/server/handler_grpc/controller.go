package handler_grpc

import (
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handler/pvz"
	pb "gitlab.ozon.dev/zlatoivan4/homework/internal/pkg/pb"
)

type Controller struct {
	pvzService   pvz.Service
	orderService order.Service
	pb.UnimplementedApiServer
}

func New(pvzService pvz.Service, orderService order.Service) Controller {
	controller := Controller{
		pvzService:   pvzService,
		orderService: orderService,
	}
	return controller
}
