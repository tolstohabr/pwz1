package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("Введите команду или 'help' для списка доступных команд.")

	storage := NewFileStorage("orders.json")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" {
			fmt.Println("Завершение работы.")
			break
		}

		args := strings.Split(input, " ")
		switch args[0] {
		case "help":
			printHelp()
		case "accept-order":
			handleAcceptOrder(storage, args[1:])
		case "return-order":
			handleReturnOrder(storage, args[1:])
		case "process-order":
			handleProcessOrders(storage, args[1:])
		default:
			fmt.Println("Неизвестная команда. Введите 'help' для списка.")
		}
	}
}

func printHelp() {
	fmt.Println("Список команд:")
	fmt.Println("  help             Показать эту справку")
	fmt.Println("  accept-order     Принять заказ от курьера")
	fmt.Println("  return-order     Вернуть заказ") //выбросить на помойку
	fmt.Println("  process-order    Выдать или принять возврат")
	fmt.Println("  exit             Выйти из приложения")
}

func handleAcceptOrder(storage Storage, args []string) {
	var orderID, userID, expiresStr string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--order-id":
			if i+1 < len(args) {
				orderID = args[i+1]
				i++
			}
		case "--user-id":
			if i+1 < len(args) {
				userID = args[i+1]
				i++
			}
		case "--expires":
			if i+1 < len(args) {
				expiresStr = args[i+1]
				i++
			}
		}
	}

	if orderID == "" || userID == "" || expiresStr == "" {
		fmt.Println("ERROR: VALIDATION_FAILED: Все параметры обязательны.")
		return
	}

	expiresAt, err := time.Parse("2006-01-02", expiresStr)
	if err != nil {
		fmt.Println("ERROR: VALIDATION_FAILED: Неверный формат даты (нужен yyyy-mm-dd)")
		return
	}

	err = AcceptOrder(storage, orderID, userID, expiresAt)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		fmt.Println("ORDER_ACCEPTED:", orderID)
	}
}

func handleReturnOrder(storage Storage, args []string) {
	var orderID string

	for i := 0; i < len(args); i++ {
		if args[i] == "--order-id" && i+1 < len(args) {
			orderID = args[i+1]
			i++
		}
	}

	if orderID == "" {
		fmt.Println("ERROR: VALIDATION_FAILED: Параметр --order-id обязателен.")
		return
	}

	err := ReturnOrder(storage, orderID)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		fmt.Println("ORDER_RETURNED:", orderID)
	}
}

func handleProcessOrders(storage Storage, args []string) {
	var userID, action, orderIDsStr string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--user-id":
			if i+1 < len(args) {
				userID = args[i+1]
				i++
			}
		case "--action":
			if i+1 < len(args) {
				action = args[i+1]
				i++
			}
		case "--order-ids":
			if i+1 < len(args) {
				orderIDsStr = args[i+1]
				i++
			}
		}
	}

	if userID == "" || action == "" || orderIDsStr == "" {
		fmt.Println("ERROR: VALIDATION_FAILED: Требуются --user-id, --action и --order-ids")
		return
	}

	orderIDs := strings.Split(orderIDsStr, ",")
	results := ProcessOrders(storage, userID, action, orderIDs)

	for _, res := range results {
		fmt.Println(res)
	}
}
