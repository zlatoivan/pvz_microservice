package pvz

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model/pvz"
)

type PvzStorage interface {
	CreatePVZ(pvz pvz.PVZ) error
	GetPVZ(name string) ([]pvz.PVZ, error)
}

type PVZService struct {
	store  PvzStorage
	wg     sync.WaitGroup
	infoWg sync.WaitGroup
}

func New(ps PvzStorage) (*PVZService, error) {
	return &PVZService{store: ps}, nil
}

func (s *PVZService) Writer(writeCh <-chan pvz.PVZ, infoCh chan<- string) {
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

func (s *PVZService) Reader(readCh <-chan string, infoCh chan<- string) {
	for {
		select {
		case name := <-readCh:
			s.infoWg.Add(1)
			infoCh <- "\t[INFO]: goroutine Reader is activated and is reading now"
			s.infoWg.Wait()
			pvzs, err := s.store.GetPVZ(name)
			if err != nil {
				log.Printf("s.Reader: %v", err)
			} else {
				fmt.Println("Information about the PVZs:")
				for _, p := range pvzs {
					fmt.Printf("Name: %s\nAddress: %s\nContacts: %s\n", p.Name, p.Address, p.Contacts)
				}
			}
			s.infoWg.Add(1)
			infoCh <- "\t[INFO]: goroutine Reader read everything and is waiting again"
			s.wg.Done()
		}
	}
}

func (s *PVZService) CreatePVZ(writeCh chan<- pvz.PVZ) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the name of PVZ: ")
	name, err := reader.ReadString('\n')
	name = name[:len(name)-1]
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

	newPVZ := pvz.PVZ{Name: name, Address: address, Contacts: contacts}
	writeCh <- newPVZ

	return nil
}

func (s *PVZService) GetPVZ(readCh chan<- string) error {
	fmt.Print("Enter the name of PVZ: ")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	name = name[:len(name)-1]
	if err != nil {
		return err
	}
	fmt.Println()

	readCh <- name

	return nil
}

func (s *PVZService) Signal(signalCh <-chan os.Signal) {
	sig := <-signalCh
	fmt.Println("\nThe program termination signal has been received:", sig)
	fmt.Print("Shutting down the tool...\n\n")
	os.Exit(0)
}

func (s *PVZService) Info(infoCh <-chan string) {
	for {
		info := <-infoCh
		fmt.Println(info)
		s.infoWg.Done()
	}
}

func (s *PVZService) InteractiveMode() error {
	infoCh := make(chan string)
	writeCh := make(chan pvz.PVZ)
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
