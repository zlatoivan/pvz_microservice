package delivery

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func MakeRespErrInvalidData(err error) ResponseError {
	return ResponseError{
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid data",
		ErrorText:      err.Error(),
	}
}

func MakeRespErrAlreadyExists(err error) ResponseError {
	return ResponseError{
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Already exists",
		ErrorText:      err.Error(),
	}
}

func MakeRespErrNotFoundByID(err error) ResponseError {
	return ResponseError{
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Not found by ID",
		ErrorText:      err.Error(),
	}
}

func MakeRespErrInternalServer(err error) ResponseError {
	return ResponseError{
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server error",
		ErrorText:      err.Error(),
	}
}

func MakeRespId(id uuid.UUID) ResponseID {
	return ResponseID{ID: id}
}

func MakeRespComment(comment string) ResponseComment {
	return ResponseComment{Comment: comment}
}

func MakeRespPVZList(list []model.PVZ) []ResponsePVZ {
	respList := make([]ResponsePVZ, 0)
	for _, pvz := range list {
		respPVZ := ResponsePVZ{
			ID:       pvz.ID,
			Name:     pvz.Name,
			Address:  pvz.Address,
			Contacts: pvz.Contacts,
		}
		respList = append(respList, respPVZ)
	}
	return respList
}

func MakeRespPVZ(pvz model.PVZ) ResponsePVZ {
	respOrder := ResponsePVZ{
		ID:       pvz.ID,
		Name:     pvz.Name,
		Address:  pvz.Address,
		Contacts: pvz.Contacts,
	}
	return respOrder
}

func MakeRespOrderList(list []model.Order) []ResponseOrder {
	respList := make([]ResponseOrder, 0)
	for _, order := range list {
		respOrder := ResponseOrder{
			ID:            order.ID,
			ClientID:      order.ClientID,
			Weight:        order.Weight,
			Cost:          order.Cost,
			StoresTill:    order.StoresTill,
			GiveOutTime:   order.GiveOutTime,
			IsReturned:    order.IsReturned,
			IsDeleted:     order.IsDeleted,
			PackagingType: order.PackagingType,
		}
		respList = append(respList, respOrder)
	}
	return respList
}

func MakeRespOrder(order model.Order) ResponseOrder {
	respOrder := ResponseOrder{
		ID:            order.ID,
		ClientID:      order.ClientID,
		Weight:        order.Weight,
		Cost:          order.Cost,
		StoresTill:    order.StoresTill,
		GiveOutTime:   order.GiveOutTime,
		IsReturned:    order.IsReturned,
		PackagingType: order.PackagingType,
	}
	return respOrder
}

func RenderResponse(w http.ResponseWriter, req *http.Request, status int, data interface{}) {
	render.Status(req, status)
	render.JSON(w, req, data)
}
