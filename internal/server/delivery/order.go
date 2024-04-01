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
	"gitlab.ozon.dev/zlatoivan4/homework/internal/server"
)

func PrepToPrintOrder(order model.Order) string {
	if order.ID == uuid.Nil {
		return fmt.Sprintf("   ClientID: %s\n   StoresTill: %s\n   GiveOutTime: %s   IsReturned: %t\n", order.ClientID, order.StoresTill, order.GiveOutTime, order.IsReturned)
	}
	return fmt.Sprintf("   Id: %s\n   ClientID: %s\n   StoresTill: %s\n   GiveOutTime: %s   IsReturned: %t\n", order.ID, order.ClientID, order.StoresTill, order.GiveOutTime, order.IsReturned)
}

func GetOrderIDFromURL(req *http.Request) (uuid.UUID, error) {
	idStr := chi.URLParam(req, "orderID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("uuid.Parse: %w", err)
	}
	return id, nil
}

func GetClientIDFromURL(req *http.Request) (uuid.UUID, error) {
	idStr := chi.URLParam(req, "clientID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("uuid.Parse: %w", err)
	}
	return id, nil
}

const dateLayout = "02.01.2006 15:04"

func GetOrderWithoutIDFromReq(req *http.Request) (model.Order, string, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return model.Order{}, "", fmt.Errorf("io.ReadAll: %w", err)
	}
	err = req.Body.Close()
	if err != nil {
		return model.Order{}, "", fmt.Errorf("req.Body.Close: %w", err)
	}

	var order server.RequestOrder
	err = json.Unmarshal(data, &order)
	if err != nil {
		return model.Order{}, "", fmt.Errorf("json.Unmarshal: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	if order.ClientID == uuid.Nil {
		return model.Order{}, "", fmt.Errorf("ClientID is Nil")
	}
	storesTill, err := time.Parse(dateLayout, order.StoresTill)
	if err != nil {
		return model.Order{}, "", fmt.Errorf("time.Parse: %w", err)
	}
	newOrder := model.Order{
		ClientID:   order.ClientID,
		StoresTill: storesTill,
		Weight:     order.Weight,
		Cost:       order.Cost,
	}

	return newOrder, order.PackagingType, nil
}

func GetOrderFromReq(req *http.Request) (model.Order, error) {
	id, err := GetOrderIDFromURL(req)
	if err != nil {
		return model.Order{}, fmt.Errorf("GetOrderIDFromURL: %w", err)
	}
	order, _, err := GetOrderWithoutIDFromReq(req)
	if err != nil {
		return model.Order{}, fmt.Errorf("GetOrderWithoutIDFromReq: %w", err)
	}
	order.ID = id
	return order, nil
}

func GetClientOrdersIDsFromReq(req *http.Request) ([]uuid.UUID, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	err = req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("req.Body.Close: %w", err)
	}

	var reqClientOrders server.RequestClientOrders
	err = json.Unmarshal(data, &reqClientOrders)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	return reqClientOrders.IDs, nil
}

func GetDataForGiveOut(req *http.Request) (uuid.UUID, []uuid.UUID, error) {
	clientID, err := GetClientIDFromURL(req)
	if err != nil {
		return uuid.UUID{}, nil, fmt.Errorf("GetClientIDFromURL: %w", err)
	}

	ids, err := GetClientOrdersIDsFromReq(req)
	if err != nil {
		return uuid.UUID{}, nil, fmt.Errorf("GetClientOrdersIDsFromReq: %w", err)
	}

	return clientID, ids, nil
}

func GetOrderIDFromReq(req *http.Request) (uuid.UUID, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("io.ReadAll: %w", err)
	}
	err = req.Body.Close()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("req.Body.Close: %w", err)
	}

	var reqOrderID server.RequestOrderID
	err = json.Unmarshal(data, &reqOrderID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	return reqOrderID.ID, nil
}

func GetDataForReturnOrder(req *http.Request) (uuid.UUID, uuid.UUID, error) {
	clientID, err := GetClientIDFromURL(req)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, fmt.Errorf("GetClientIDFromURL: %w", err)
	}

	id, err := GetOrderIDFromReq(req)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, fmt.Errorf("GetClientOrdersIDsFromReq: %w", err)
	}

	return clientID, id, nil
}

func MakeRespList(list []model.Order) []server.ResponseOrder {
	respList := make([]server.ResponseOrder, 0)
	for _, order := range list {
		respOrder := server.ResponseOrder{
			ID:          order.ID,
			ClientID:    order.ClientID,
			Weight:      order.Weight,
			Cost:        order.Cost,
			StoresTill:  order.StoresTill,
			GiveOutTime: order.GiveOutTime,
			IsReturned:  order.IsReturned,
		}
		respList = append(respList, respOrder)
	}
	return respList
}

func MakeRespOrder(order model.Order) server.ResponseOrder {
	respOrder := server.ResponseOrder{
		ID:          order.ID,
		ClientID:    order.ClientID,
		Weight:      order.Weight,
		Cost:        order.Cost,
		StoresTill:  order.StoresTill,
		GiveOutTime: order.GiveOutTime,
		IsReturned:  order.IsReturned,
	}
	return respOrder
}
