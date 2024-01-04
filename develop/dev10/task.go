package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {

	// Парсим аргументы командной строки
	if len(os.Args) != 3 {
		fmt.Println("Использование: go run .\\task.go 1.1.1.1. 123")
		return
	}

	timeout := flag.Duration("timeout", 10*time.Second, "Таймаут подключения")
	flag.Parse()

	// Формируем адрес
	address := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])

	// Подключаемся к серверу
	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v\n", err)
	}
	defer conn.Close()

	// Создаем канал для отслеживания событий завершения
	done := make(chan struct{})

	// Запускаем горутину для чтения данных из сокета и вывода их в STDOUT
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Printf("Ошибка чтения данных: %v\n", err)
				close(done)
				return
			}
			fmt.Print(string(buffer[:n]))
		}
	}()

	// Запускаем горутину для обработки сигнала Ctrl+C
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		<-sigCh
		close(done)
	}()

	// Ожидаем ввод данных из STDIN и отправляем их в сокет
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buffer)
			if err != nil {
				fmt.Printf("Ошибка чтения из STDIN: %v\n", err)
				break
			}

			_, err = conn.Write(buffer[:n])
			if err != nil {
				fmt.Printf("Ошибка отправки данных: %v\n", err)
				break
			}
		}
	}()

	// Ожидаем завершения горутин и закрываем сокет
	<-done
	fmt.Println("Программа завершена.")
}
