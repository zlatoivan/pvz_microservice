package service

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

const dateLayoutConst = "02.01.2006"
const helpPath = "help.txt"

type storage interface {
	Create(order model.Order) error
	Delete(id int) error
	GiveOut(ids []int) error
	List(id int, lastN int, inPVZ bool) ([]int, error)
	Return(id int, clientID int) error
	ListOfReturned(pagenum int, itemsonpage int) ([]int, error)
	CreatePVZ(pvz model.PVZ) error
	GetPVZ(title string) (model.PVZ, error)
}

type Service struct {
	store  storage
	wg     sync.WaitGroup
	infoWg sync.WaitGroup
}

func New(s storage) (Service, error) {
	return Service{store: s}, nil
}

func (s *Service) Help() error {
	data, err := os.ReadFile(helpPath)
	if err != nil {
		return err
	}
	fmt.Println(string(data) + "\n")
	return nil
}

func validateCreate(id int, clientID int, storesTillStr string) (time.Time, error) {
	if id == -1 || clientID == -1 || storesTillStr == "-" {
		return time.Time{}, fmt.Errorf("incorrect input data format")
	}

	// Привести срок хранения к типу даты
	storesTill, err := time.Parse(dateLayoutConst, storesTillStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("incorrect date format")
	}

	return storesTill, nil
}

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

	fmt.Println("Orders have been given out to the client")

	return nil
}

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

func (s *Service) Writer(writeCh <-chan model.PVZ, infoCh chan<- string) {
	for {
		select {
		case newPVZ := <-writeCh:
			s.infoWg.Add(1)
			infoCh <- "\t[INFO]: goroutine Writer is activated and is writing now"
			s.infoWg.Wait()
			err := s.store.CreatePVZ(newPVZ)
			if err != nil {
				log.Printf("s.Writer: %v", err)
			} else {
				fmt.Println("PVZ added")
			}
			s.infoWg.Add(1)
			infoCh <- "\t[INFO]: goroutine Writer wrote everything and is waiting again"
			s.wg.Done()
		}
	}
}

func (s *Service) Reader(readCh <-chan string, infoCh chan<- string) {
	for {
		select {
		case title := <-readCh:
			s.infoWg.Add(1)
			infoCh <- "\t[INFO]: goroutine Reader is activated and is reading now"
			s.infoWg.Wait()
			pvz, err := s.store.GetPVZ(title)
			if err != nil {
				log.Printf("s.Reader: %v", err)
			} else {
				fmt.Println("Information about the PVZ:")
				fmt.Println("Title:", pvz.Title)
				fmt.Println("Address:", pvz.Address)
				fmt.Println("Contacts:", pvz.Contacts)
			}
			s.infoWg.Add(1)
			infoCh <- "\t[INFO]: goroutine Reader read everything and is waiting again"
			s.wg.Done()
		}
	}
}

func (s *Service) CreatePVZ(writeCh chan<- model.PVZ) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the title of PVZ: ")
	title, err := reader.ReadString('\n')
	title = title[:len(title)-1]
	if err != nil {
		return err
	}

	fmt.Print("Enter the address of PVZ: ")
	address, err := reader.ReadString('\n')
	address = address[:len(address)-1]
	if err != nil {
		return err
	}

	fmt.Print("Enter the contacts of PVZ: ")
	contacts, err := reader.ReadString('\n')
	contacts = contacts[:len(contacts)-1]
	if err != nil {
		return err
	}
	fmt.Println()

	newPVZ := model.PVZ{Title: title, Address: address, Contacts: contacts}
	writeCh <- newPVZ

	return nil
}

func (s *Service) GetPVZ(readCh chan<- string) error {
	fmt.Print("Enter the title of PVZ: ")
	reader := bufio.NewReader(os.Stdin)
	title, err := reader.ReadString('\n')
	title = title[:len(title)-1]
	if err != nil {
		return err
	}
	fmt.Println()

	readCh <- title

	return nil
}

func (s *Service) Signal(signalCh <-chan os.Signal) {
	sig := <-signalCh
	fmt.Println("\nThe program termination signal has been received:", sig)
	fmt.Print("Shutting down the tool...\n\n")
	os.Exit(0)
}

func (s *Service) Info(infoCh <-chan string) {
	for {
		info := <-infoCh
		fmt.Println(info)
		s.infoWg.Done()
	}
}

func (s *Service) InteractiveMode() error {
	infoCh := make(chan string)
	writeCh := make(chan model.PVZ)
	readCh := make(chan string)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go s.Info(infoCh)
	go s.Writer(writeCh, infoCh)
	s.infoWg.Add(1)
	infoCh <- "\t[INFO]: goroutine Writer is launched (waiting)"
	go s.Reader(readCh, infoCh)
	s.infoWg.Add(1)
	infoCh <- "\t[INFO]: goroutine Reader is launched (waiting)"
	go s.Signal(signalCh)
	s.infoWg.Wait()

	for {
		fmt.Print("\n> ")
		reader := bufio.NewReader(os.Stdin)

		command, err := reader.ReadString('\n')
		command = command[:len(command)-1]
		if err != nil {
			return err
		}

		s.wg.Add(1)
		switch command {
		case "add":
			err = s.CreatePVZ(writeCh)
			if err != nil {
				return err
			}

		case "get":
			err = s.GetPVZ(readCh)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown command")
		}
		s.wg.Wait()
		s.infoWg.Wait()
	}
}

func (s *Service) Work() error {
	id, clientID, storesTillStr, idsStr, lastN, inPVZ, pagenum, itemsonpage := parseFlags()
	command := flag.Args()[len(flag.Args())-1]

	switch command {
	case "help":
		return s.Help()
	case "create":
		return s.Create(id, clientID, storesTillStr)
	case "delete":
		return s.Delete(id)
	case "giveout":
		return s.GiveOut(idsStr)
	case "list":
		return s.List(clientID, lastN, inPVZ)
	case "return":
		return s.Return(id, clientID)
	case "listofreturned":
		return s.ListOfReturned(pagenum, itemsonpage)
	case "interactive_mode":
		return s.InteractiveMode()
	default:
		return fmt.Errorf("unknown command")
	}
}
