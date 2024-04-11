package delivery

import (
	"time"

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
	st, _ := time.Parse(time.RFC3339, reqOrder.StoresTill)
	return model.Order{
		ID:            reqOrder.ID,
		ClientID:      reqOrder.ClientID,
		Weight:        reqOrder.Weight,
		Cost:          reqOrder.Cost,
		StoresTill:    st,
		PackagingType: reqOrder.PackagingType,
	}
}
