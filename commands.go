package main

import (
	"errors"
	"strings"
	"time"
)

// добавить заказ в ПВЗ
func AcceptOrder(storage Storage, orderID, userID string, expiresAt time.Time) error {
	if strings.TrimSpace(orderID) == "" || strings.TrimSpace(userID) == "" {
		return errors.New("VALIDATION_FAILED: ID заказа и пользователя обязательны")
	}

	//если срок хранения в прошлом
	if expiresAt.Before(time.Now()) {
		return errors.New("VALIDATION_FAILED: срок хранения в прошлом")
	}

	//если такой заказ уже есть
	_, err := storage.GetOrder(orderID)
	if err == nil {
		return errors.New("ORDER_ALREADY_EXISTS: заказ уже есть")
	}

	newOrder := Order{
		ID:        orderID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Status:    StatusAccepted,
	}

	return storage.SaveOrder(newOrder)
}

// удалить заказ
func ReturnOrder(storage Storage, orderID string) error {
	order, err := storage.GetOrder(orderID)
	if err != nil {
		return errors.New("ORDER_NOT_FOUND")
	}

	//если заказ у клиента
	if order.Status == StatusIssued {
		return errors.New("ORDER_ALREADY_ISSUED")
	}

	//если время хранения не истекло
	if time.Now().Before(order.ExpiresAt) {
		return errors.New("STORAGE_NOT_EXPIRED: время хранения не истекло")
	}

	return storage.DeleteOrder(orderID)
}
