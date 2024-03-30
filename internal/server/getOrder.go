package server

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

func GetOrderIDFromURL(req *http.Request) (uuid.UUID, error) {
	idStr := chi.URLParam(req, "orderID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("uuid.Parse: %w", err)
	}
	return id, nil
}

const dateLayout = "02.01.2006 15:04"

func GetOrderWithoutIDFromReq(req *http.Request) (model.Order, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return model.Order{}, fmt.Errorf("io.ReadAll: %w", err)
	}
	err = req.Body.Close()
	if err != nil {
		return model.Order{}, fmt.Errorf("req.Body.Close: %w", err)
	}

	var order requestOrder
	err = json.Unmarshal(data, &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	if order.ClientID == uuid.Nil {
		return model.Order{}, fmt.Errorf("ClientID is Nil")
	}
	storesTill, err := time.Parse(dateLayout, order.StoresTill)
	if err != nil {
		return model.Order{}, fmt.Errorf("time.Parse: %w", err)
	}
	newOrder := model.Order{
		ClientID:   order.ClientID,
		StoresTill: storesTill,
	}

	return newOrder, nil
}

func GetOrderFromReq(req *http.Request) (model.Order, error) {
	id, err := GetOrderIDFromURL(req)
	if err != nil {
		return model.Order{}, fmt.Errorf("GetOrderIDFromURL: %w", err)
	}
	order, err := GetOrderWithoutIDFromReq(req)
	if err != nil {
		return model.Order{}, fmt.Errorf("GetOrderWithoutIDFromReq: %w", err)
	}
	order.ID = id
	return order, nil
}

func PrepToPrintOrder(order model.Order) string {
	if order.ID == uuid.Nil {
		return fmt.Sprintf("   ClientID: %s\n   StoresTill: %s\n   IsDeleted: %t   GiveOutTime: %s   IsReturned: %t\n", order.ClientID, order.StoresTill, order.IsDeleted, order.GiveOutTime, order.IsReturned)
	}
	return fmt.Sprintf("   Id: %s\n   ClientID: %s\n   StoresTill: %s\n   IsDeleted: %t   GiveOutTime: %s   IsReturned: %t\n", order.ID, order.ClientID, order.StoresTill, order.IsDeleted, order.GiveOutTime, order.IsReturned)
}
