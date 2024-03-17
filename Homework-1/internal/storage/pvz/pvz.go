package pvz

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/model/pvz"
)

type PVZStorage struct {
	Mutex sync.RWMutex       `json:"mutex"`
	Pvzs  map[string]pvz.PVZ `json:"pvzs,omitempty"`
}

const PVZStoragePath = "db/pvz_db.json"

func New() (*PVZStorage, error) {
	file, err := os.OpenFile(PVZStoragePath, os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	store := &PVZStorage{
		Pvzs: make(map[string]pvz.PVZ),
	}
	bytes, err := os.ReadFile(PVZStoragePath)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return store, nil
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&store.Pvzs)
	return store, err
}

func (s *PVZStorage) CreatePVZ(pvz pvz.PVZ) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	for k := range s.Pvzs {
		if k == pvz.Name {
			return fmt.Errorf("PVZ with this name is already in the storage")
		}
	}
	s.Pvzs[pvz.Name] = pvz

	file, err := os.Open(PVZStoragePath)
	if err != nil {
		return fmt.Errorf("s.store.CreatePVZ: %w", err)
	}
	defer file.Close()
	bytes, err := json.MarshalIndent(s.Pvzs, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(PVZStoragePath, bytes, 0644)
	return err
}

func (s *PVZStorage) GetPVZ(name string) (pvz.PVZ, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	ok := false
	for k := range s.Pvzs {
		if k == name {
			ok = true
		}
	}
	if !ok {
		return pvz.PVZ{}, fmt.Errorf("no PVZ with this name was found")
	}

	return s.Pvzs[name], nil
}
