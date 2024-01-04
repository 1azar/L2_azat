package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type config struct {
	fieldsStr *string
	delimiter *string
	separated *bool
}

func main() {
	// Определение структуры с настройками
	conf := config{}
	// Определение флагов
	conf.fieldsStr = flag.String("f", "", "выбрать поля (колонки)")
	conf.delimiter = flag.String("d", "\t", "использовать другой разделитель")
	conf.separated = flag.Bool("s", false, "только строки с разделителем")
	//fieldsStr := flag.String("f", "", "выбрать поля (колонки)")
	//delimiter := flag.String("d", "\t", "использовать другой разделитель")
	//separated := flag.Bool("s", false, "только строки с разделителем")

	// Парсинг флагов
	flag.Parse()

	//
	app(&conf, os.Stdin, os.Stdout)
}

func app(conf *config, reader io.Reader, writer io.Writer) {
	// обработка значений флага -f
	var fields []int
	if *conf.fieldsStr != "" {
		var err error
		fields, err = parseFields(*conf.fieldsStr)
		if err != nil {
			log.Fatal("ошибка: невозможно обработать значение флага f: ", err)
		}
	}

	// Чтение ввода по строкам
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		// Проверка, содержит ли строка разделитель
		if *conf.separated && !strings.Contains(line, *conf.delimiter) {
			continue
		}

		// Разбивка строки на колонки
		columns := strings.Split(line, *conf.delimiter)

		// Выбор запрошенных полей
		if len(fields) > 0 {
			for _, field := range fields {
				if field >= 1 && field <= len(columns) {
					fmt.Fprint(writer, columns[field-1])
					if field < len(columns) {
						fmt.Fprint(writer, *conf.delimiter)
					}
				}
			}
			fmt.Fprintln(writer)
		} else {
			// Вывод всей строки, если не указаны поля
			fmt.Fprintln(writer, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Ошибка ввода:", err)
	}
}

// parseFields преобразует строку с номерами полей в срез целых чисел
func parseFields(fieldsStr string) ([]int, error) {
	// Разбиваем строку на подстроки, разделенные запятыми
	fieldsList := strings.Split(fieldsStr, ",")

	// Слайс для хранения результата
	fieldsSlice := make([]int, 0, len(fieldsList))

	// Преобразование каждой подстроки в число
	for _, valueStr := range fieldsList {
		value, err := strconv.Atoi(strings.TrimSpace(valueStr))
		if err != nil {
			return nil, err
		}
		fieldsSlice = append(fieldsSlice, value)
	}
	return fieldsSlice, nil
}
