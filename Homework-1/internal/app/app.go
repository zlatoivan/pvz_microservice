package app

import (
	"flag"
	"fmt"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/service/order"
	"gitlab.ozon.dev/zlatoivan4/homework/internal/service/pvz"
)

type App struct {
	orderService order.OrderService
	pvzService   pvz.PVZService
}

func New(os order.OrderStorage, ps pvz.PvzStorage) (*App, error) {
	orderServ, err := order.New(os)
	if err != nil {
		return nil, err
	}
	pvzServ, err := pvz.New(ps)
	if err != nil {
		return nil, err
	}
	//_ = pvzServ

	serv := &App{
		orderService: *orderServ,
		pvzService:   *pvzServ,
		//pvzService:   pvz.PVZService{store: ps},
	}
	return serv, nil
}

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

func (s *App) Work() error {
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
		return s.pvzService.InteractiveMode()
	default:
		return fmt.Errorf("unknown command")
	}
}
