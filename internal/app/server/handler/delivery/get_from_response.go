package delivery

import (
	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func GetOrderFromRespOrder(respOrder ResponseOrder) model.Order {
	return model.Order{
		ID:            respOrder.ID,
		ClientID:      respOrder.ClientID,
		Weight:        respOrder.Weight,
		Cost:          respOrder.Cost,
		StoresTill:    respOrder.StoresTill,
		GiveOutTime:   respOrder.GiveOutTime,
		IsReturned:    respOrder.IsReturned,
		IsDeleted:     respOrder.IsDeleted,
		PackagingType: respOrder.PackagingType,
	}
}

func GetOrderFromReqOrder(reqOrder RequestOrder) model.Order {
	return model.Order{
		ID:            reqOrder.ID,
		ClientID:      reqOrder.ClientID,
		Weight:        reqOrder.Weight,
		Cost:          reqOrder.Cost,
		StoresTill:    reqOrder.StoresTill,
		PackagingType: reqOrder.PackagingType,
	}
}
