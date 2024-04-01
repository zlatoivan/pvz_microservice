package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/server/delivery"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/repo"
)

// createOrder creates order
func (s Server) createOrder(w http.ResponseWriter, req *http.Request) {
	newOrder, packagingType, err := delivery.GetOrderWithoutIDFromReq(req)
	if err != nil {
		log.Printf("[createOrder] GetOrderWithoutIDFromReq: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	id, err := s.OrderService.CreateOrder(req.Context(), packagingType, newOrder)
	if err != nil {
		log.Printf("[createOrder] s.OrderService.createOrder: %v\n", err)
		if errors.Is(err, repo.ErrorAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			WriteComment(w, "ID already exists")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	log.Printf("Order created! id = %s", id)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(delivery.ResponseID{ID: id})
	if err != nil {
		log.Printf("[createOrder] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// listOrders gets list of orders
func (s Server) listOrders(w http.ResponseWriter, req *http.Request) {
	list, err := s.OrderService.ListOrders(req.Context())
	if err != nil {
		log.Printf("[listOrders] s.OrderService.ListOrders: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	log.Println("Got list of orders!")

	w.Header().Set("Content-Type", "application/json")
	if len(list) == 0 {
		w.WriteHeader(http.StatusOK)
		WriteComment(w, "No orders in database")
		return
	}

	respList := delivery.MakeRespList(list)
	err = json.NewEncoder(w).Encode(respList)
	if err != nil {
		log.Printf("[listOrders] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getOrderByID gets Order by ID
func (s Server) getOrderByID(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetOrderIDFromURL(req)
	if err != nil {
		log.Printf("[getOrderByID] GetOrderIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	order, err := s.OrderService.GetOrderByID(req.Context(), id)
	if err != nil {
		log.Printf("[getOrderByID] s.OrderService.GetOrderByID: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			WriteComment(w, "Order not found by this ID")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	log.Printf("Order by ID:\n" + delivery.PrepToPrintOrder(order))

	w.Header().Set("Content-Type", "application/json")
	respOrder := delivery.MakeRespOrder(order)
	err = json.NewEncoder(w).Encode(respOrder)
	if err != nil {
		log.Printf("[getOrderByID] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// updateOrder updates Order
func (s Server) updateOrder(w http.ResponseWriter, req *http.Request) {
	updOrder, err := delivery.GetOrderFromReq(req)
	if err != nil {
		log.Printf("[updateOrder] GetOrderFromReq: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.OrderService.UpdateOrder(req.Context(), updOrder)
	if err != nil {
		log.Printf("[updateOrder] s.OrderService.UpdateOrder: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			WriteComment(w, "Order not found by this ID")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	log.Println("Order updated!")
	WriteComment(w, "Order updated!")

	w.WriteHeader(http.StatusOK)
}

// deleteOrder deletes Order
func (s Server) deleteOrder(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetOrderIDFromURL(req)
	if err != nil {
		log.Printf("[deleteOrder] GetOrderIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.OrderService.DeleteOrder(req.Context(), id)
	if err != nil {
		log.Printf("[deleteOrder] s.OrderService.DeleteOrder: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			WriteComment(w, "Order not found by this ID")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	log.Println("Order deleted!")
	WriteComment(w, "Order deleted!")

	w.WriteHeader(http.StatusOK)
}

// listClientOrders gets list of client orders
func (s Server) listClientOrders(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetClientIDFromURL(req)
	if err != nil {
		log.Printf("[listClientOrders] GetClientIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	list, err := s.OrderService.ListClientOrders(req.Context(), id)
	if err != nil {
		log.Printf("[listClientOrders] s.OrderService.ListClientOrders: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	log.Printf("Got list of clients orders! Length = %d.\n", len(list))

	w.Header().Set("Content-Type", "application/json")
	if len(list) == 0 {
		w.WriteHeader(http.StatusOK)
		WriteComment(w, "No orders in database")
		return
	}

	respList := delivery.MakeRespList(list)
	err = json.NewEncoder(w).Encode(respList)
	if err != nil {
		log.Printf("[listClientOrders] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// giveOutOrders gives out a list of orders
func (s Server) giveOutOrders(w http.ResponseWriter, req *http.Request) {
	clientID, ids, err := delivery.GetDataForGiveOut(req)
	if err != nil {
		log.Printf("[giveOutOrders] GetDataForGiveOut: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.OrderService.GiveOutOrders(req.Context(), clientID, ids)
	if err != nil {
		log.Printf("[giveOutOrders] s.OrderService.GiveOutOrders: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	log.Println("Orders are given out")

	w.Header().Set("Content-Type", "application/json")
	WriteComment(w, "Orders are given out")

	w.WriteHeader(http.StatusOK)
}

// returnOrder returns order
func (s Server) returnOrder(w http.ResponseWriter, req *http.Request) {
	clientID, id, err := delivery.GetDataForReturnOrder(req)
	if err != nil {
		log.Printf("[returnOrder] GetDataForGiveOut: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.OrderService.ReturnOrder(req.Context(), clientID, id)
	if err != nil {
		log.Printf("[returnOrder] s.OrderService.ReturnOrder: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	log.Println("Order is returned")

	w.Header().Set("Content-Type", "application/json")
	WriteComment(w, "Order is returned")

	w.WriteHeader(http.StatusOK)
}

// returnOrder returns order
func (s Server) listReturnedOrders(w http.ResponseWriter, req *http.Request) {
	list, err := s.OrderService.ListReturnedOrders(req.Context())
	if err != nil {
		log.Printf("[listReturnedOrders] s.OrderService.ListReturnedOrders: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	log.Printf("Got list of returned orders! Length = %d.\n", len(list))

	w.Header().Set("Content-Type", "application/json")
	if len(list) == 0 {
		w.WriteHeader(http.StatusOK)
		WriteComment(w, "No returned orders in database")
		return
	}

	respList := delivery.MakeRespList(list)
	err = json.NewEncoder(w).Encode(respList)
	if err != nil {
		log.Printf("[listReturnedOrders] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteComment(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
