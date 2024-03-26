package order

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const helpPath = "help.txt"
const dateLayout = "02.01.2006"

type Storage interface {
	Create(order model.Order) error
	Delete(id int) error
	GiveOut(ids []int) error
	List(id int, lastN int, inPVZ bool) ([]int, error)
	Return(id int, clientID int) error
	ListOfReturned(pagenum int, itemsonpage int) ([]int, error)
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
func validateCreate(id int, clientID int, storesTillStr string) (time.Time, error) {
	if id == -1 || clientID == -1 || storesTillStr == "-" {
		return time.Time{}, fmt.Errorf("incorrect input data format")
	}

	// Привести срок хранения к типу даты
	storesTill, err := time.Parse(dateLayout, storesTillStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("incorrect date format")
	}

	return storesTill, nil
}

// Create creates a new order in the Storage
func (s *Service) Create(id int, clientID int, storesTillStr string) error {
	storesTill, err := validateCreate(id, clientID, storesTillStr)
	if err != nil {
		return err
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
func (s *Service) Delete(id int) error {
	if id == -1 {
		return fmt.Errorf("incorrect input data format")
	}

	err := s.store.Delete(id)
	if err != nil {
		return fmt.Errorf("s.store.Delete: %w", err)
	}

	fmt.Println("The order has been deleted")

	return nil
}

// GiveOut gives out an orders to the client by orders ids
func (s *Service) GiveOut(idsStr string) error {
	idsToSplit := idsStr[1 : len(idsStr)-1]
	idsToInt := strings.Split(idsToSplit, ",")
	ids := make([]int, len(idsToInt))
	for i := range idsToInt {
		idInt, err := strconv.Atoi(idsToInt[i])
		if err != nil {
			return fmt.Errorf("invalid order IDs " + idsToInt[i])
		}
		ids[i] = idInt
	}

	err := s.store.GiveOut(ids)
	if err != nil {
		return fmt.Errorf("s.store.GiveOut: %w", err)
	}

	fmt.Println("orders have been given out to the client")

	return nil
}

// List gets a list of all orders of a specific client
func (s *Service) List(clientID int, lastN int, inPVZ bool) error {
	if clientID == -1 {
		return fmt.Errorf("incorrect input data format")
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
func (s *Service) Return(id int, clientID int) error {
	if id == -1 || clientID == -1 {
		return fmt.Errorf("incorrect input data format")
	}

	err := s.store.Return(id, clientID)
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
