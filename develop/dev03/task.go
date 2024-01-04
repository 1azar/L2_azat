package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры):
на входе подается файл c несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительно
Реализовать поддержку утилитой следующих ключей:
-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учетом суффиксов
*/

func main() {
	errLog := log.New(os.Stderr, "ERRO", log.Ltime|log.Ldate)

	// Определение флагов
	column := flag.Int("k", 0, "Указание колонки для сортировки (по умолчанию 0 - вся строка)")
	numeric := flag.Bool("n", false, "Сортировать по числовому значению")
	reverse := flag.Bool("r", false, "Сортировать в обратном порядке")
	unique := flag.Bool("u", false, "Не выводить повторяющиеся строки")
	monthSort := flag.Bool("M", false, "Сортировать по названию месяца")
	ignoreTrailingSpace := flag.Bool("b", false, "Игнорировать хвостовые пробелы")
	checkSorted := flag.Bool("c", false, "Проверять отсортированы ли данные")
	numericWithSuffix := flag.Bool("h", false, "Сортировать по числовому значению с учетом суффиксов")

	flag.Parse()

	var inputFilename, outputFilename string
	// имя исходного файла, если не дано - ошибка
	if inputFilename = flag.Arg(0); inputFilename == "" {
		log.Fatal("ошибка: не указано имя файла")
	}
	// имя результирующего отсортированного файла, если не задано - имя исходного файла с суффиксом _SORTED
	if outputFilename = flag.Arg(1); outputFilename == "" {
		extension := filepath.Ext(inputFilename)
		basename := filepath.Base(inputFilename)
		outputFilename = basename[:len(basename)-len(extension)] + "_SORTED" + filepath.Ext(inputFilename)
	}

	// Открытие файла для чтения
	file, err := os.Open(inputFilename)
	if err != nil {
		errLog.Fatal("ошибка открытия файла:", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			errLog.Fatal("ошибка закрытия файла", err)
		}
	}()

	// Считывание строк из файла в слайс строк.
	// todo: При недостатке памяти хорошо бы применить https://en.wikipedia.org/wiki/External_sorting
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Обработка ошибок сканера
	if err := scanner.Err(); err != nil {
		errLog.Fatal("ошибка чтения файла", err)
	}

	// Функция сравнения для сортировки строк
	comparator := func(i, j int) bool {
		return compare(lines[i], lines[j], *column, *numeric, *monthSort, *ignoreTrailingSpace, *numericWithSuffix)
	}

	// Проверка сортировки, если указан флаг -c
	if *checkSorted && !sort.SliceIsSorted(lines, comparator) {
		errLog.Fatal("ошибка: данные не отсортированы")
	}

	// Сортировка строк
	sort.Slice(lines, comparator)

	// Уникальные строки, если указан флаг -u
	if *unique {
		lines = uniqueLines(lines)
	}

	// Обратный порядок, если указан флаг -r
	if *reverse {
		reverseLines(lines)
	}

	// сохранение отсортированных строк
	outfile, err := os.Create(outputFilename)
	if err != nil {
		errLog.Fatal("ошибка создания файла", err)
	}
	defer func() {
		if err = outfile.Close(); err != nil {
			errLog.Fatal("ошибка закрытия файла", err)
		}
	}()
	writer := bufio.NewWriter(outfile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			errLog.Fatal("ошибка записи данных в файл")
		}
	}
	// Вызываем Flush, чтобы убедиться, что все данные записаны в файл
	err = writer.Flush()
	if err != nil {
		errLog.Fatal(err)
	}

	// Вывод отсортированных строк
	for _, line := range lines {
		fmt.Println(line)
	}
}

// Функция сравнения строк для сортировки
func compare(a, b string, column int, numeric, monthSort, ignoreTrailingSpace, numericWithSuffix bool) bool {
	// Разбиение строк на колонки
	fieldsA := strings.Fields(a)
	fieldsB := strings.Fields(b)

	// Выбор колонки для сравнения
	var colA, colB string
	if column < len(fieldsA) {
		colA = fieldsA[column]
	}
	if column < len(fieldsB) {
		colB = fieldsB[column]
	}

	// Игнорирование хвостовых пробелов, если указан флаг -b
	if ignoreTrailingSpace {
		colA = strings.TrimSpace(colA)
		colB = strings.TrimSpace(colB)
	}

	// Сортировка по числовому значению с учетом суффиксов, если указан флаг -h
	if numericWithSuffix {
		numA, errA := parseNumericWithSuffix(colA)
		numB, errB := parseNumericWithSuffix(colB)

		// В случае ошибки преобразования, считаем строки равными
		if errA != nil || errB != nil {
			return a < b
		}

		// Сравнение числовых значений
		return numA < numB
	}

	// Сортировка по числовому значению, если указан флаг -n
	if numeric {
		numA, errA := strconv.Atoi(colA)
		numB, errB := strconv.Atoi(colB)

		// Если числовое преобразование не удалось, сравниваем строки как строки
		if errA != nil || errB != nil {
			return a < b
		}

		// Сравнение числовых значений
		return numA < numB
	}

	// Сортировка по названию месяца, если указан флаг -M
	if monthSort {
		monthA, errA := parseMonth(colA)
		monthB, errB := parseMonth(colB)

		// Если преобразование названия месяца не удалось, сравниваем строки как строки
		if errA != nil || errB != nil {
			return a < b
		}

		// Сравнение месяцев
		return monthA < monthB
	}

	// Сравнение строк
	return colA < colB
}

// Парсинг числовых значений с учетом суффиксов (например, 10M)
func parseNumericWithSuffix(s string) (int, error) {

	// Извлечение числового значения и суффикса
	var numStr, suffix string
	for i, ch := range s {
		if ch < '0' || ch > '9' { // если не число в таблице ASCII
			numStr = s[:i]
			suffix = s[i:]
			break
		}
	}

	// Преобразование числового значения
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	// Умножение на множитель суффикса
	if multiplier, exists := suffixes[suffix]; exists {
		num *= multiplier
	}

	return num, nil
}

// Парсинг названия месяца
func parseMonth(s string) (time.Month, error) {
	month, ok := monthMap[s]
	if !ok {
		return 0, fmt.Errorf("неверный месяц")
	}

	return month, nil
}

// Уникальные строки
func uniqueLines(lines []string) []string {
	var uniqueLines []string
	seen := make(map[string]bool)

	for _, line := range lines {
		if _, exists := seen[line]; !exists {
			seen[line] = true
			uniqueLines = append(uniqueLines, line)
		}
	}

	return uniqueLines
}

// Обратный порядок строк
func reverseLines(lines []string) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
}

var suffixes = map[string]int{
	"k":   1e3,
	"K":   1e3,
	"m":   1e6,
	"M":   1e6,
	"g":   1e9,
	"G":   1e9,
	"t":   1e12,
	"T":   1e12,
	"p":   1e15,
	"P":   1e15,
	"e1":  1e1,
	"e2":  1e2,
	"e3":  1e3,
	"e4":  1e4,
	"e5":  1e5,
	"e6":  1e6,
	"e7":  1e7,
	"e8":  1e8,
	"e9":  1e9,
	"e10": 1e10,
	"e11": 1e11,
	"e12": 1e12,
	"e13": 1e13,
	"e14": 1e14,
	"e15": 1e15,
}
var monthMap = map[string]time.Month{
	"Jan": time.January,
	"Feb": time.February,
	"Mar": time.March,
	"Apr": time.April,
	"May": time.May,
	"Jun": time.June,
	"Jul": time.July,
	"Aug": time.August,
	"Sep": time.September,
	"Oct": time.October,
	"Nov": time.November,
	"Dec": time.December,
}
