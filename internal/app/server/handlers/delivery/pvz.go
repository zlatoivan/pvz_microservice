package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

func PrepToPrintPVZ(pvz model.PVZ) string {
	if pvz.ID == uuid.Nil {
		return fmt.Sprintf("   Name: %s\n   Address: %s\n   Contacts: %s\n", pvz.Name, pvz.Address, pvz.Contacts)
	}
	return fmt.Sprintf("   Id: %s\n   Name: %s\n   Address: %s\n   Contacts: %s\n", pvz.ID, pvz.Name, pvz.Address, pvz.Contacts)
}

func GetPVZFromReq(req *http.Request) (model.PVZ, error) {
	var pvz RequestPVZ
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("io.ReadAll: %w", err)
	}
	err = req.Body.Close()
	if err != nil {
		log.Printf("[error] req.Body.Close: %v", err)
	}
	err = json.Unmarshal(data, &pvz)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))
	newPVZ := model.PVZ{
		ID:       pvz.ID,
		Name:     pvz.Name,
		Address:  pvz.Address,
		Contacts: pvz.Contacts,
	}

	return newPVZ, nil
}
