package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/storage/repo"
)

// createOrder creates order
func (s Server) createOrder(w http.ResponseWriter, req *http.Request) {
	newOrder, err := GetOrderWithoutIDFromReq(req)
	if err != nil {
		log.Printf("[createOrder] GetOrderWithoutIDFromReq: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	id, err := s.orderService.CreateOrder(req.Context(), newOrder)
	if err != nil {
		log.Printf("[createOrder] s.orderService.createOrder: %v\n", err)
		if errors.Is(err, repo.ErrorAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			writeComment(w, "ID already exists")
		}
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	log.Printf("Order created! id = %s", id)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ResponseID{ID: id})
	if err != nil {
		log.Printf("[createOrder] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// listOrders gets list of orders
func (s Server) listOrders(w http.ResponseWriter, req *http.Request) {
	list, err := s.orderService.ListOrders(req.Context())
	if err != nil {
		log.Printf("[listOrders] s.orderService.ListOrders: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	log.Println("Got list of orders!")

	w.Header().Set("Content-Type", "application/json")
	if len(list) == 0 {
		w.WriteHeader(http.StatusOK)
		writeComment(w, "No orders in database")
		return
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Printf("[listOrders] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getOrderByID gets Order by ID
func (s Server) getOrderByID(w http.ResponseWriter, req *http.Request) {
	id, err := GetOrderIDFromURL(req)
	if err != nil {
		log.Printf("[getOrderByID] GetOrderIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	order, err := s.orderService.GetOrderByID(req.Context(), id)
	if err != nil {
		log.Printf("[getOrderByID] s.orderService.GetOrderByID: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			writeComment(w, "Order not found by this ID")
		}
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	log.Printf("Order by ID:\n" + PrepToPrintOrder(order))

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		log.Printf("[getOrderByID] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// updateOrder updates Order
func (s Server) updateOrder(w http.ResponseWriter, req *http.Request) {
	updOrder, err := GetOrderFromReq(req)
	if err != nil {
		log.Printf("[updateOrder] GetOrderFromReq: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.orderService.UpdateOrder(req.Context(), updOrder)
	if err != nil {
		log.Printf("[updateOrder] s.orderService.UpdateOrder: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			writeComment(w, "Order not found by this ID")
		}
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	log.Println("Order updated!")

	w.WriteHeader(http.StatusOK)
}

// deleteOrder deletes Order
func (s Server) deleteOrder(w http.ResponseWriter, req *http.Request) {
	id, err := GetOrderIDFromURL(req)
	if err != nil {
		log.Printf("[deleteOrder] GetOrderIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.orderService.DeleteOrder(req.Context(), id)
	if err != nil {
		log.Printf("[deleteOrder] s.orderService.DeleteOrder: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			writeComment(w, "Order not found by this ID")
		}
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	log.Println("Order deleted!")

	w.WriteHeader(http.StatusOK)
}

// listClientOrders gets list of client orders
func (s Server) listClientOrders(w http.ResponseWriter, req *http.Request) {
	id, err := GetClientIDFromURL(req)
	if err != nil {
		log.Printf("[listClientOrders] GetClientIDFromURL: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	list, err := s.orderService.ListClientOrders(req.Context(), id)
	if err != nil {
		log.Printf("[listClientOrders] s.orderService.ListClientOrders: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	log.Printf("Got list of clients orders! Length = %d.\n", len(list))

	w.Header().Set("Content-Type", "application/json")
	if len(list) == 0 {
		w.WriteHeader(http.StatusOK)
		writeComment(w, "No orders in database")
		return
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Printf("[listClientOrders] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// giveOutOrders gives out a list of orders
func (s Server) giveOutOrders(w http.ResponseWriter, req *http.Request) {
	clientID, ids, err := GetDataForGiveOut(req)
	if err != nil {
		log.Printf("[giveOutOrders] GetDataForGiveOut: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.orderService.GiveOutOrders(req.Context(), clientID, ids)
	if err != nil {
		log.Printf("[giveOutOrders] s.orderService.GiveOutOrders: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	log.Println("Orders are given out")

	w.Header().Set("Content-Type", "application/json")
	writeComment(w, "Orders are given out")

	w.WriteHeader(http.StatusOK)
}

// returnOrder returns order
func (s Server) returnOrder(w http.ResponseWriter, req *http.Request) {
	clientID, id, err := GetDataForReturnOrder(req)
	if err != nil {
		log.Printf("[returnOrder] GetDataForGiveOut: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeComment(w, "Invalid data: "+err.Error())
		return
	}

	err = s.orderService.ReturnOrder(req.Context(), clientID, id)
	if err != nil {
		log.Printf("[returnOrder] s.orderService.ReturnOrder: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	log.Println("Order is returned")

	w.Header().Set("Content-Type", "application/json")
	writeComment(w, "Order is returned")

	w.WriteHeader(http.StatusOK)
}

// returnOrder returns order
func (s Server) listReturnedOrders(w http.ResponseWriter, req *http.Request) {
	list, err := s.orderService.ListReturnedOrders(req.Context())
	if err != nil {
		log.Printf("[listReturnedOrders] s.orderService.ListReturnedOrders: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	log.Printf("Got list of returned orders! Length = %d.\n", len(list))

	w.Header().Set("Content-Type", "application/json")
	if len(list) == 0 {
		w.WriteHeader(http.StatusOK)
		writeComment(w, "No returned orders in database")
		return
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Printf("[listReturnedOrders] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeComment(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
