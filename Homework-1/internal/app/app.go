package app

import (
	"context"
	"flag"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/service/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/service/pvz"
)

type App struct {
	orderService order.Service
	pvzService   pvz.Service
}

// New creates a new app
func New(orderService *order.Service, pvzService *pvz.Service) (*App, error) {
	service := &App{
		orderService: *orderService,
		pvzService:   *pvzService,
	}
	return service, nil
}

// parseFlags parses flags from console
func parseFlags() (int, int, string, string, int, bool, int, int) {
	id := flag.Int("id", -1, "id of order")
	clientID := flag.Int("clientid", -1, "id of client")
	storesTillStr := flag.String("storestill", "-", "shelf life of order")
	idsStr := flag.String("ids", "-", "ids of orders to give out")
	lastN := flag.Int("lastn", -1, "last n orders of client")
	inPVZ := flag.Bool("inpvz", false, "client's orders that are in pvz")
	pagenum := flag.Int("pagenum", -1, "number of pages")
	itemsonpage := flag.Int("itemsonpage", -1, "number of items on page")
	flag.Parse()

	return *id, *clientID, *storesTillStr, *idsStr, *lastN, *inPVZ, *pagenum, *itemsonpage
}

// Work starts the application work
func (s *App) Work(ctx context.Context, cancel context.CancelFunc) error {
	id, clientID, storesTillStr, idsStr, lastN, inPVZ, pagenum, itemsonpage := parseFlags()
	command := flag.Args()[len(flag.Args())-1]

	switch command {
	case "help":
		return s.orderService.Help()
	case "create":
		return s.orderService.Create(id, clientID, storesTillStr)
	case "delete":
		return s.orderService.Delete(id)
	case "giveout":
		return s.orderService.GiveOut(idsStr)
	case "list":
		return s.orderService.List(clientID, lastN, inPVZ)
	case "return":
		return s.orderService.Return(id, clientID)
	case "listofreturned":
		return s.orderService.ListOfReturned(pagenum, itemsonpage)
	case "interactive_mode":
		return s.pvzService.InteractiveMode(ctx, cancel)
	default:
		return fmt.Errorf("unknown command")
	}
}
