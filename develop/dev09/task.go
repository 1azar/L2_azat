package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Определение флага "url"
	urlFlag := flag.String("url", "", "URL веб-сайта для скачивания")
	flag.Parse()

	// Проверка наличия значения флага "url"
	if *urlFlag == "" {
		fmt.Println("Используйте флаг -url для указания URL веб-сайта.")
		return
	}

	// Выполнение HTTP-запроса
	resp, err := http.Get(*urlFlag)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при чтении тела ответа: %v\n", err)
		return
	}

	// Извлечение имени файла из URL
	fileName := "output.html"
	urlParts := strings.Split(*urlFlag, "/")
	if len(urlParts) > 2 {
		fileName = urlParts[len(urlParts)-2] + ".html"
	}

	// Сохранение HTML-кода в файл
	err = os.WriteFile(fileName, body, 0644)
	if err != nil {
		fmt.Printf("Ошибка при сохранении в файл: %v\n", err)
		return
	}

	fmt.Printf("HTML-код веб-сайта сохранен в файл: %s\n", fileName)
}
