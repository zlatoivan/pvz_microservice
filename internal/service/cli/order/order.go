package order

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const helpPath = "help.txt"
const dateLayout = "02.01.2006"

type Storage interface {
	Create(order model.Order) error
	Delete(id uuid.UUID) error
	GiveOut(ids []uuid.UUID) error
	List(id uuid.UUID, lastN int, inPVZ bool) ([]uuid.UUID, error)
	Return(id uuid.UUID, clientID uuid.UUID) error
	ListOfReturned(pagenum int, itemsonpage int) ([]uuid.UUID, error)
}

type Service struct {
	store Storage
}

// New creates a new order service
func New(orderStore Storage) (*Service, error) {
	return &Service{store: orderStore}, nil
}

// Help displays auxiliary information about the functionality of the program
func (s *Service) Help() error {
	data, err := os.ReadFile(helpPath)
	if err != nil {
		return err
	}
	fmt.Println(string(data) + "\n")
	return nil
}

// validateCreate checks the validation of the data for the create method
func validateCreate(idStr string, clientIDStr string, storesTillStr string) (uuid.UUID, uuid.UUID, time.Time, error) {
	if idStr == "-" || clientIDStr == "-" || storesTillStr == "-" {
		return uuid.UUID{}, uuid.UUID{}, time.Time{}, fmt.Errorf("incorrect input data format")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, time.Time{}, fmt.Errorf("uuid.Parse: %w", err)
	}
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, time.Time{}, fmt.Errorf("uuid.Parse: %w", err)
	}
	storesTill, err := time.Parse(dateLayout, storesTillStr)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, time.Time{}, fmt.Errorf("time.Parse: %w", err)
	}

	return id, clientID, storesTill, nil
}

// Create creates a new order in the Storage
func (s *Service) Create(idStr string, clientIDStr string, storesTillStr string) error {
	id, clientID, storesTill, err := validateCreate(idStr, clientIDStr, storesTillStr)
	if err != nil {
		return fmt.Errorf("validateCreate: %w", err)
	}

	newOrder := model.Order{
		ID:          id,
		ClientID:    clientID,
		StoresTill:  storesTill,
		IsDeleted:   false,
		GiveOutTime: time.Time{}, // zero value
		IsReturned:  false,
	}
	err = s.store.Create(newOrder)
	if err != nil {
		return fmt.Errorf("s.store.Create: %w", err)
	}

	fmt.Println("The order is accepted")

	return nil
}

// Delete deletes an order from the Storage
func (s *Service) Delete(idStr string) error {
	if idStr == "-" {
		return fmt.Errorf("incorrect input data format")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", err)
	}

	err = s.store.Delete(id)
	if err != nil {
		return fmt.Errorf("s.store.Delete: %w", err)
	}

	fmt.Println("The order has been deleted")

	return nil
}

// GiveOut gives out an orders to the client by orders ids
func (s *Service) GiveOut(idsStr string) error {
	idsToSplit := idsStr[1 : len(idsStr)-1]
	idsToUUID := strings.Split(idsToSplit, ",")
	ids := make([]uuid.UUID, len(idsToUUID))
	for i := range idsToUUID {
		id, err := uuid.Parse(idsToUUID[i])
		if err != nil {
			return fmt.Errorf("invalid order IDs " + idsToUUID[i])
		}
		ids[i] = id
	}

	err := s.store.GiveOut(ids)
	if err != nil {
		return fmt.Errorf("s.store.GiveOut: %w", err)
	}

	fmt.Println("orders have been given out to the client")

	return nil
}

// List gets a list of all orders of a specific client
func (s *Service) List(clientIDStr string, lastN int, inPVZ bool) error {
	if clientIDStr == "-" {
		return fmt.Errorf("incorrect input data format")
	}
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", err)
	}

	list, err := s.store.List(clientID, lastN, inPVZ)
	if err != nil {
		return fmt.Errorf("s.store.List: %w", err)
	}

	if len(list) == 0 {
		return fmt.Errorf("this client does not have orders with such parameters")
	}

	fmt.Println("Customer orders:", list)

	return nil
}

// Return returns the order by order id and client id
func (s *Service) Return(idStr string, clientIDStr string) error {
	if idStr == "-" || clientIDStr == "-" {
		return fmt.Errorf("incorrect input data format")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", err)
	}
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", err)
	}

	err = s.store.Return(id, clientID)
	if err != nil {
		return fmt.Errorf("s.store.Return: %w", err)
	}

	fmt.Println("Order return accepted")

	return nil
}

// ListOfReturned gets a list of returned orders. With pagination.
func (s *Service) ListOfReturned(pagenum int, itemsonpage int) error {
	if pagenum == -1 || itemsonpage == -1 {
		return fmt.Errorf("incorrect input data format")
	}

	list, err := s.store.ListOfReturned(pagenum, itemsonpage)
	if err != nil {
		return fmt.Errorf("s.store.ListOfReturned: %w", err)
	}

	if len(list) == 0 {
		return fmt.Errorf("there are no returned orders")
	}

	fmt.Printf("Returned orders (page %d; %d orders on page):\n%v\n", pagenum, itemsonpage, list)

	return nil
}
