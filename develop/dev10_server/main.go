package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Новый клиент подключен: %s\n", clientAddr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("Получено от %s: %s\n", clientAddr, message)

		resp := fmt.Sprintf("[%s]\tСообшение принято.\n", time.Now().Format(time.TimeOnly))
		conn.Write([]byte(resp))
		fmt.Println("Ответ отправлен клиенту")
	}

	fmt.Printf("Клиент %s отключился.\n", clientAddr)
}

func main() {
	// Запускаем сервер на порту 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Сервер запущен. Ожидание подключений на ", listener.Addr())

	// Обработка сигнала Ctrl+C для корректного завершения сервера
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		<-sigCh
		fmt.Println("\nПринят сигнал завершения. Завершение сервера.")
		os.Exit(0)
	}()

	// Бесконечный цикл ожидания и обработки подключений
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Ошибка при принятии подключения: %v\n", err)
			continue
		}

		// Запускаем горутину для обработки клиента
		go handleClient(conn)
	}
}
