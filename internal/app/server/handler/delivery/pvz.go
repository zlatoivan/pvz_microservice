package delivery

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func GetRawDataFromReq(req *http.Request) ([]byte, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	err = req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("req.Body.Close: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))
	return data, nil
}

func GetPVZFromReq(req *http.Request) (model.PVZ, error) {
	var reqPVZ RequestPVZ
	err := render.DecodeJSON(req.Body, &reqPVZ)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("render.DecodeJSON: %w", err)
	}
	pvz := model.PVZ{
		ID:       reqPVZ.ID,
		Name:     reqPVZ.Name,
		Address:  reqPVZ.Address,
		Contacts: reqPVZ.Contacts,
	}

	return pvz, nil
}
