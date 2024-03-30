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
		return
	}

	log.Printf("Order created! id = %s", id)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ResponseID{ID: id})
	if err != nil {
		log.Printf("[createOrder] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
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
		return
	}

	log.Println("Got order list!")

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
		return
	}

	log.Printf("Order by ID:\n" + PrepToPrintOrder(order))

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		log.Printf("[getOrderByID] json.NewEncoder().Encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
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
		return
	}

	log.Println("Order deleted!")

	w.WriteHeader(http.StatusOK)
}
