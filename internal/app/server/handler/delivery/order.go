package delivery

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func GetIDFromReq(req *http.Request) (uuid.UUID, error) {
	var reqID RequestID
	err := render.DecodeJSON(req.Body, &reqID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("render.DecodeJSON: %w", err)
	}

	if reqID.ID == uuid.Nil {
		return uuid.UUID{}, fmt.Errorf("id is nil")
	}

	return reqID.ID, nil
}

func GetOrderFromReq(req *http.Request) (model.Order, error) {
	var reqOrder RequestOrder
	err := render.DecodeJSON(req.Body, &reqOrder)
	if err != nil {
		return model.Order{}, fmt.Errorf("render.DecodeJSON: %w", err)
	}

	if reqOrder.ClientID == uuid.Nil {
		return model.Order{}, fmt.Errorf("client id is nil")
	}

	order := model.Order{
		ID:            reqOrder.ID,
		ClientID:      reqOrder.ClientID,
		StoresTill:    reqOrder.StoresTill,
		Weight:        reqOrder.Weight,
		Cost:          reqOrder.Cost,
		PackagingType: reqOrder.PackagingType,
	}

	return order, nil
}

func GetDataForGiveOutFromReq(req *http.Request) (uuid.UUID, []uuid.UUID, error) {
	var reqGiveOut RequestGiveOut
	err := render.DecodeJSON(req.Body, &reqGiveOut)
	if err != nil {
		return uuid.UUID{}, nil, fmt.Errorf("render.DecodeJSON: %w", err)
	}

	return reqGiveOut.ClientID, reqGiveOut.IDs, nil
}

func GetDataForReturnFromReq(req *http.Request) (uuid.UUID, uuid.UUID, error) {
	var reqReturn RequestReturn
	err := render.DecodeJSON(req.Body, &reqReturn)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, fmt.Errorf("render.DecodeJSON: %w", err)
	}

	if reqReturn.ClientID == uuid.Nil {
		return uuid.UUID{}, uuid.UUID{}, fmt.Errorf("client id is nil")
	}

	if reqReturn.ID == uuid.Nil {
		return uuid.UUID{}, uuid.UUID{}, fmt.Errorf("id is nil")
	}

	return reqReturn.ClientID, reqReturn.ID, nil
}
