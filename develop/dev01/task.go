package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"os"
	time2 "time"
)

/*
Создать программу, печатающую точное время с использованием NTP -библиотеки. Инициализировать как go module.
Использовать библиотеку github.com/beevik/ntp.
Требования:
Программа должна быть оформлена как go module
Программа должна корректно обрабатывать ошибки библиотеки: выводить их в STDERR и возвращать ненулевой код выхода в OS
*/

// Создание и настройка логгера
var errLog = log.New(os.Stderr, "", log.Ldate|log.Ltime)

// Адрес сервера ntp
var ntpServerAddr = "pool.ntp.org"

func main() {
	PrintTime()
}

func PrintTime() {
	time, err := ntp.Time(ntpServerAddr)
	if err != nil {
		// вывод логов ошибки и os.Exit(1)
		errLog.Fatal(err)
	}
	fmt.Println(time.Format(time2.DateTime))
}
