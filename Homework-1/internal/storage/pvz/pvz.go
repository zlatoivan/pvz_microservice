package pvz

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model"
)

type Storage struct {
	mu   sync.RWMutex
	pvzs map[int]model.PVZ
}

const storagePath = "db_files/pvz_db.json"

// New Creates a new pvz storage
func New() (*Storage, error) {
	file, err := os.OpenFile(storagePath, os.O_CREATE, 0777)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Println("[pvz][Storage] file.Close:", err)
		}
	}()

	store := &Storage{
		pvzs: make(map[int]model.PVZ),
	}
	bytes, err := os.ReadFile(storagePath)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return store, nil
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&store.pvzs)
	if err != nil {
		return nil, fmt.Errorf("decoder.Decode: %w", err)
	}

	return store, nil
}

// Close closes the storage connection, updating and saving data in a file
func (s *Storage) Close() error {
	file, err := os.Open(storagePath)
	if err != nil {
		return fmt.Errorf("pvz.Storage os.OpenFile: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Println("pvz.Storage file.Close:", err)
		}
	}()

	s.mu.Lock()
	defer s.mu.Unlock()
	bytes, err := json.MarshalIndent(s.pvzs, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(storagePath, bytes, 0644)
	if err != nil {
		return fmt.Errorf("pvz.Storage os.WriteFile: %w", err)
	}

	return nil
}

// CreatePVZ creates a new PVZ
func (s *Storage) CreatePVZ(pvz model.PVZ) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pvzs[pvz.ID] = pvz
	return nil
}

// GetPVZ gets information about the PVZ
func (s *Storage) GetPVZ(name string) ([]model.PVZ, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	pvzs := make([]model.PVZ, 0)
	for _, p := range s.pvzs {
		if p.Name == name {
			pvzs = append(pvzs, p)
		}
	}
	if len(pvzs) == 0 {
		return nil, fmt.Errorf("no PVZ with this name was found")
	}

	return pvzs, nil
}

// GetListOfPVZs gets information about all PVZs
func (s *Storage) GetListOfPVZs() ([]model.PVZ, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pvzsList := make([]model.PVZ, 0)
	for _, v := range s.pvzs {
		pvzsList = append(pvzsList, v)
	}

	return pvzsList, nil
}
