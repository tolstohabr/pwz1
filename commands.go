package main

import (
	"errors"
	"fmt"
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

	//если заказ в ПВЗ после возврата
	if order.Status == StatusReturned {
		return storage.DeleteOrder(orderID)
	}

	//если время хранения не истекло
	if time.Now().Before(order.ExpiresAt) {
		return errors.New("STORAGE_NOT_EXPIRED: время хранения не истекло")
	}

	return storage.DeleteOrder(orderID)
}

// обработать выдачу или возврат заказа
func ProcessOrders(storage Storage, userID string, action string, orderIDs []string) []string {
	var results []string
	for _, id := range orderIDs {
		order, err := storage.GetOrder(id)
		if err != nil {
			results = append(results, fmt.Sprintf("ERROR %s: ORDER_NOT_FOUND", id))
			continue
		}

		if order.UserID != userID {
			results = append(results, fmt.Sprintf("ERROR %s: USER_MISMATCH", id))
			continue
		}

		if action == "issue" {
			if time.Now().After(order.ExpiresAt) {
				results = append(results, fmt.Sprintf("ERROR %s: STORAGE_EXPIRED", id))
				continue
			}
			now := time.Now()
			order.Status = StatusIssued
			order.IssuedAt = &now
		} else if action == "return" {
			if order.IssuedAt == nil || time.Since(*order.IssuedAt) > 48*time.Hour {
				results = append(results, fmt.Sprintf("ERROR %s: RETURN_WINDOW_EXPIRED", id))
				continue
			}
			order.Status = StatusReturned
		} else {
			results = append(results, fmt.Sprintf("ERROR %s: INVALID_ACTION", id))
			continue
		}

		// Перезапись заказа
		_ = storage.DeleteOrder(order.ID)
		_ = storage.SaveOrder(order)

		results = append(results, fmt.Sprintf("PROCESSED: %s", id))
	}
	return results
}
