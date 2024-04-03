package server

import (
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
		log.Printf("[createOrder] GetOrderWithoutIDFromReq: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	id, err := s.OrderService.CreateOrder(req.Context(), packagingType, newOrder)
	if err != nil {
		log.Printf("[createOrder] s.OrderService.createOrder: %v\n", err)
		if errors.Is(err, repo.ErrorAlreadyExists) {
			delivery.RenderResponse(w, req, http.StatusConflict, delivery.MakeRespErrAlreadyExists(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("Order created. id = %s\n", id)
	delivery.RenderResponse(w, req, http.StatusCreated, delivery.MakeRespId(id))
}

// listOrders gets list of orders
func (s Server) listOrders(w http.ResponseWriter, req *http.Request) {
	list, err := s.OrderService.ListOrders(req.Context())
	if err != nil {
		log.Printf("[listOrders] s.OrderService.ListOrders: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("Got list of orders. Length = %d.\n", len(list))
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespOrderList(list))
}

// getOrderByID gets Order by ID
func (s Server) getOrderByID(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetOrderIDFromURL(req)
	if err != nil {
		log.Printf("[getOrderByID] GetOrderIDFromURL: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	order, err := s.OrderService.GetOrderByID(req.Context(), id)
	if err != nil {
		log.Printf("[getOrderByID] s.OrderService.GetOrderByID: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("Got order by ID")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespOrder(order))
}

// updateOrder updates Order
func (s Server) updateOrder(w http.ResponseWriter, req *http.Request) {
	updOrder, err := delivery.GetOrderFromReq(req)
	if err != nil {
		log.Printf("[updateOrder] GetOrderFromReq: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.OrderService.UpdateOrder(req.Context(), updOrder)
	if err != nil {
		log.Printf("[updateOrder] s.OrderService.UpdateOrder: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("Order updated")
	delivery.RenderResponse(w, req, http.StatusOK, "Order updated")
}

// deleteOrder deletes Order
func (s Server) deleteOrder(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetOrderIDFromURL(req)
	if err != nil {
		log.Printf("[deleteOrder] GetOrderIDFromURL: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.OrderService.DeleteOrder(req.Context(), id)
	if err != nil {
		log.Printf("[deleteOrder] s.OrderService.DeleteOrder: %v\n", err)
		if errors.Is(err, repo.ErrorNotFound) {
			delivery.RenderResponse(w, req, http.StatusNotFound, delivery.MakeRespErrNotFoundByID(err))
			return
		}
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("Order deleted")
	delivery.RenderResponse(w, req, http.StatusOK, "Order deleted")
}

// listClientOrders gets list of client orders
func (s Server) listClientOrders(w http.ResponseWriter, req *http.Request) {
	id, err := delivery.GetClientIDFromURL(req)
	if err != nil {
		log.Printf("[listClientOrders] GetClientIDFromURL: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	list, err := s.OrderService.ListClientOrders(req.Context(), id)
	if err != nil {
		log.Printf("[listClientOrders] s.OrderService.ListClientOrders: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("Got list of clients orders! Length = %d.\n", len(list))
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespOrderList(list))
}

// giveOutOrders gives out a list of orders
func (s Server) giveOutOrders(w http.ResponseWriter, req *http.Request) {
	clientID, ids, err := delivery.GetDataForGiveOut(req)
	if err != nil {
		log.Printf("[giveOutOrders] GetDataForGiveOut: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.OrderService.GiveOutOrders(req.Context(), clientID, ids)
	if err != nil {
		log.Printf("[giveOutOrders] s.OrderService.GiveOutOrders: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("Orders are given out")
	delivery.RenderResponse(w, req, http.StatusOK, "Orders are given out")
}

// returnOrder returns order
func (s Server) returnOrder(w http.ResponseWriter, req *http.Request) {
	clientID, id, err := delivery.GetDataForReturnOrder(req)
	if err != nil {
		log.Printf("[returnOrder] GetDataForGiveOut: %v", err)
		delivery.RenderResponse(w, req, http.StatusBadRequest, delivery.MakeRespErrInvalidData(err))
		return
	}

	err = s.OrderService.ReturnOrder(req.Context(), clientID, id)
	if err != nil {
		log.Printf("[returnOrder] s.OrderService.ReturnOrder: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Println("Order is returned")
	delivery.RenderResponse(w, req, http.StatusOK, "Order is returned")
}

// returnOrder returns order
func (s Server) listReturnedOrders(w http.ResponseWriter, req *http.Request) {
	list, err := s.OrderService.ListReturnedOrders(req.Context())
	if err != nil {
		log.Printf("[listReturnedOrders] s.OrderService.ListReturnedOrders: %v\n", err)
		delivery.RenderResponse(w, req, http.StatusInternalServerError, delivery.MakeRespErrInternalServer(err))
		return
	}

	log.Printf("Got list of returned orders! Length = %d.\n", len(list))
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespOrderList(list))
}
