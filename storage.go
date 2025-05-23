package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Storage interface {
	SaveOrder(order Order) error
	GetOrder(id string) (Order, error)
	DeleteOrder(id string) error
	ListOrders() ([]Order, error)
}

type FileStorage struct {
	filePath string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{filePath: path}
}

func (fs *FileStorage) load() ([]Order, error) {
	file, err := os.Open(fs.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Order{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var orders []Order
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (fs *FileStorage) save(orders []Order) error {
	file, err := os.Create(fs.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(orders)
}

func (fs *FileStorage) SaveOrder(order Order) error {
	orders, err := fs.load()
	if err != nil {
		return err
	}

	for _, o := range orders {
		if o.ID == order.ID {
			return errors.New("ORDER_ALREADY_EXISTS")
		}
	}

	orders = append(orders, order)
	return fs.save(orders)
}

func (fs *FileStorage) GetOrder(id string) (Order, error) {
	orders, err := fs.load()
	if err != nil {
		return Order{}, err
	}

	for _, o := range orders {
		if o.ID == id {
			return o, nil
		}
	}

	return Order{}, errors.New("ORDER_NOT_FOUND")
}

func (fs *FileStorage) DeleteOrder(id string) error {
	orders, err := fs.load()
	if err != nil {
		return err
	}

	updated := []Order{}
	found := false

	for _, o := range orders {
		if o.ID == id {
			found = true
			continue
		}
		updated = append(updated, o)
	}

	if !found {
		return errors.New("ORDER_NOT_FOUND")
	}

	return fs.save(updated)
}

func (fs *FileStorage) ListOrders() ([]Order, error) {
	return fs.load()
}
