package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func GetClientIDFromURL(req *http.Request) (uuid.UUID, error) {
	idStr := chi.URLParam(req, "clientID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("uuid.Parse: %w", err)
	}
	return id, nil
}

func GetOrderFromReq(req *http.Request) (model.Order, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return model.Order{}, fmt.Errorf("io.ReadAll: %w", err)
	}
	err = req.Body.Close()
	if err != nil {
		return model.Order{}, fmt.Errorf("req.Body.Close: %w", err)
	}

	var order RequestOrder
	err = json.Unmarshal(data, &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	if order.ClientID == uuid.Nil {
		return model.Order{}, fmt.Errorf("client id is nil")
	}
	storesTill, err := time.Parse(time.RFC3339, order.StoresTill)
	if err != nil {
		return model.Order{}, fmt.Errorf("time.Parse: %w", err)
	}
	newOrder := model.Order{
		ID:            order.ID,
		ClientID:      order.ClientID,
		StoresTill:    storesTill,
		Weight:        order.Weight,
		Cost:          order.Cost,
		PackagingType: order.PackagingType,
	}

	return newOrder, nil
}

func GetGiveOutDataFromReq(req *http.Request) (uuid.UUID, []uuid.UUID, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return uuid.UUID{}, nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	err = req.Body.Close()
	if err != nil {
		return uuid.UUID{}, nil, fmt.Errorf("req.Body.Close: %w", err)
	}

	var reqGiveOut RequestGiveOut
	err = json.Unmarshal(data, &reqGiveOut)
	if err != nil {
		return uuid.UUID{}, nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	return reqGiveOut.ClientID, reqGiveOut.IDs, nil
}

func GetIDFromReq(req *http.Request) (uuid.UUID, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("io.ReadAll: %w", err)
	}
	err = req.Body.Close()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("req.Body.Close: %w", err)
	}

	var reqID RequestID
	err = json.Unmarshal(data, &reqID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	if reqID.ID == uuid.Nil {
		return uuid.UUID{}, fmt.Errorf("id is nil")
	}

	return reqID.ID, nil
}

func GetDataForReturnOrder(req *http.Request) (uuid.UUID, uuid.UUID, error) {
	clientID, err := GetClientIDFromURL(req)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, fmt.Errorf("GetClientIDFromURL: %w", err)
	}

	id, err := GetIDFromReq(req)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, fmt.Errorf("GetClientOrdersIDsFromReq: %w", err)
	}

	return clientID, id, nil
}
