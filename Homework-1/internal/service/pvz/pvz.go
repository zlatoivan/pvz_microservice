package pvz

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type Storage interface {
	CreatePVZ(pvz model.PVZ) error
	GetPVZ(name string) ([]model.PVZ, error)
}

type Service struct {
	store Storage
}

// New creates a new PVZ service
func New(ps Storage) (*Service, error) {
	return &Service{store: ps}, nil
}

func (s *Service) Writer(writeCh <-chan model.PVZ, printCh chan<- string, errCh chan<- error) {
	for newPVZ := range writeCh {
		printCh <- "\t[INFO]: goroutine Writer is activated and is writing now\n"
		err := s.store.CreatePVZ(newPVZ)
		if err != nil {
			//_, _ = fmt.Fprintf(os.Stderr, "s.Writer: %v", err)
			errCh <- fmt.Errorf("s.Writer: %w", err)
		} else {
			printCh <- "PVZ added\n"
		}
		printCh <- "\t[INFO]: goroutine Writer wrote everything and is waiting again\n"
	}
}

func (s *Service) Reader(readCh <-chan string, printCh chan<- string, errCh chan<- error) {
	for name := range readCh {
		printCh <- "\t[INFO]: goroutine Reader is activated and is reading now\n"
		pvzs, err := s.store.GetPVZ(name)
		if err != nil {
			//log.Printf("s.Reader: %v", err)
			errCh <- fmt.Errorf("s.Writer: %w", err)
		} else {
			printCh <- "Information about the PVZs:"
			for _, p := range pvzs {
				pr := fmt.Sprintf("Name: %s\nAddress: %s\nContacts: %s\n", p.Name, p.Address, p.Contacts)
				printCh <- pr
			}
		}
		printCh <- "\t[INFO]: goroutine Reader read everything and is waiting again\n"
	}
}

func (s *Service) CreatePVZ(writeCh chan<- model.PVZ, printCh chan<- string) error {
	scanner := bufio.NewScanner(os.Stdin)

	printCh <- "Enter the name of PVZ: "
	scanner.Scan()
	name := scanner.Text()

	printCh <- "Enter the address of PVZ: "
	scanner.Scan()
	address := scanner.Text()

	printCh <- "Enter the contacts of PVZ: "
	scanner.Scan()
	contacts := scanner.Text()
	printCh <- "\n"

	newPVZ := model.PVZ{ID: uuid.NewString(), Name: name, Address: address, Contacts: contacts}
	fmt.Println(newPVZ)
	writeCh <- newPVZ

	return nil
}

func (s *Service) GetPVZ(readCh chan<- string, printCh chan<- string) error {
	printCh <- "Enter the name of PVZ: "
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()
	printCh <- "\n"

	readCh <- name

	return nil
}

func (s *Service) Print(printCh <-chan string) {
	for pr := range printCh {
		fmt.Print(pr)
	}
}

func (s *Service) Signal(ctx context.Context, cancel context.CancelFunc, signalCh <-chan os.Signal, printCh chan<- string) {
	for sig := range signalCh {
		printCh <- fmt.Sprintf("\nThe program termination signal has been received: %v\n", sig)
		printCh <- "Shutting down the tool...\n\n"
		cancel()
	}
}

func (s *Service) Errors(ctx context.Context, cancel context.CancelFunc, errCh <-chan error) {
	for err := range errCh {
		log.Printf("Error: %v", err)
		cancel()
	}
}

func (s *Service) Work(ctx context.Context, cancel context.CancelFunc, writeCh chan model.PVZ, readCh chan string, printCh chan<- string, errCh chan<- error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	//wgForWork := sync.WaitGroup{}
	for {
		//wgForWork.Add(1)
		select {
		case <-ctx.Done():
			close(writeCh)
			close(readCh)
			//if errors.Is(ctx.Err(), context.Canceled) {}  // Тип ctx.Done может быть разный
			fmt.Println("STOP")
			//errCh <- fmt.Errorf("%w", ctx.Err())
		case <-ticker.C:
			printCh <- "\n> "
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			command := scanner.Text()

			switch command {
			case "add":
				err := s.CreatePVZ(writeCh, printCh)
				if err != nil {
					errCh <- fmt.Errorf("pvz.Service.CreatePVZ: %w", err)
				}

			case "get":
				err := s.GetPVZ(readCh, printCh)
				if err != nil {
					errCh <- fmt.Errorf("pvz.Service.GetPVZ: %w", err)
				}

			default:
				errCh <- fmt.Errorf("unknown command")
			}
			//wgForWork.Done()
		}
		//wgForWork.Wait()
	}
}

// InteractiveMode starts an interactive mod of the program
func (s *Service) InteractiveMode(ctx context.Context, cancel context.CancelFunc) error {
	printCh := make(chan string)
	writeCh := make(chan model.PVZ)
	readCh := make(chan string)
	signalCh := make(chan os.Signal)
	errCh := make(chan error)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		s.Print(printCh)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.Writer(writeCh, printCh, errCh)
		wg.Done()
	}()
	printCh <- "\t[INFO]: goroutine Reader is launched (waiting)\n"

	wg.Add(1)
	go func() {
		s.Reader(readCh, printCh, errCh)
		wg.Done()
	}()
	printCh <- "\t[INFO]: goroutine Writer is launched (waiting)\n"

	wg.Add(1)
	go func() {
		s.Signal(ctx, cancel, signalCh, printCh)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.Errors(ctx, cancel, errCh)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.Work(ctx, cancel, writeCh, readCh, printCh, errCh)
		wg.Done()
	}()

	wg.Wait()
	close(printCh)

	return nil
}
