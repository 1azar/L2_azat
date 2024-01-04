package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
)

/*
Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).


Реализовать поддержку утилитой следующих ключей:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", напечатать номер строки

*/

type Configuration struct {
	afterFlag      int
	beforeFlag     int
	contextFlag    int
	countFlag      bool
	ignoreCaseFlag bool
	invertFlag     bool
	fixedFlag      bool
	lineNumFlag    bool
	regexPattern   string
	files          []string
}

// main содержит только код для обработки флагов и записи этих значений в Configuration.
// Остальной код перемещен в функцию app, чтобы облегчить написание тестов.
// Будем считать, что код в main в тестировании не нуждается
func main() {
	// Структура, хранящая значения всех флагов
	conf := Configuration{}

	// Объявляем флаги.
	conf.afterFlag = *flag.Int("A", 0, "\"after\" печатать +N строк после совпадения")
	conf.beforeFlag = *flag.Int("B", 0, "\"before\" печатать +N строк до совпадения")
	conf.contextFlag = *flag.Int("C", 0, "\"context\" (A+B) печатать ±N строк вокруг совпадения")
	conf.countFlag = *flag.Bool("c", false, "\"count\" (количество строк)")
	conf.ignoreCaseFlag = *flag.Bool("i", false, "\"ignore-case\" (игнорировать регистр)")
	conf.invertFlag = *flag.Bool("v", false, "\"invert\" (вместо совпадения, исключать)\n")
	conf.fixedFlag = *flag.Bool("F", false, "\"fixed\", точное совпадение со строкой, не паттерн")
	conf.lineNumFlag = *flag.Bool("n", false, "\"line num\", напечатать номер строки")

	// Обрабатываем флаги.
	flag.Parse()

	// Проверяем, указан ли хотя бы один файл.
	if flag.NArg() == 0 {
		log.Fatal("Использование: grep [OPTIONS] \"PATTERN\" [FILE...]\n" +
			"Пример: grep -i -n \"apple|banana\" fruitList.txt")
	}

	// Получаем паттерны для поиска.
	conf.regexPattern = flag.Args()[0]

	// Имена файлов для поиска.
	conf.files = flag.Args()[1:]

	// Выполнение основной логики
	app(&conf, os.Stdout)
}

// App логика вынесена в отдельную функцию с целью облегчения написания тестов
func app(conf *Configuration, writer io.Writer) {
	printedAnyLines = false
	// Создаем объект для вывода текста
	var logger *log.Logger
	// если задан флаг -c, то найденные строки не будут выводится
	if conf.countFlag {
		logger = log.New(io.Discard, "", 0)
	} else {
		logger = log.New(writer, "", 0)
	}

	// Экранирование всех спец символов при использовании флага -F
	if conf.fixedFlag {
		conf.regexPattern = regexp.QuoteMeta(conf.regexPattern)
	}

	// игнорировать регистр, если указан флаг -i
	if conf.ignoreCaseFlag {
		conf.regexPattern = "(?i)" + conf.regexPattern
	}

	// если задан флаг -C 3 -> -B 3 -A 3
	if conf.contextFlag > 0 {
		conf.beforeFlag = conf.contextFlag
		conf.afterFlag = conf.contextFlag
	}

	re, err := regexp.Compile(conf.regexPattern)
	if err != nil {
		logger.Fatal("ошибка обработки шаблона")
	}

	// Счетчик подходящих строк
	InclusionsCount := 0

	// Чтение файлов.
	for _, filename := range conf.files {
		// Открытие файла для чтения. Close() вызывается в конце тела цикла
		file, err := os.Open(filename)
		if err != nil {
			logger.Println("ошибка при открытии файла: ", filename, err)
			continue
		}

		// Сканирование файла по строкам
		scanner := bufio.NewScanner(file)

		// Валидация значения флага -B
		if conf.beforeFlag < 0 {
			conf.beforeFlag = 0
		}

		// Буфер для предыдущих строк (-B). +тк текущая строка тоже будет пушиться
		prevLines := NewCircularBuffer(conf.beforeFlag + 1)

		// Счетчик, который сообщает сколько строк выводить для случая когда задан флаг (-A)
		lines2print := 0

		// Флаг, сообщающий, что данная строка уже напечатана
		alreadyPrinted := false

		// Счетчик всех строк в файле
		lineNum := 0

		// Перебор строк файла
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// сохранение строки для флага -B
			if conf.beforeFlag > 0 {
				prevLines.Push(printLine(lineNum, line, conf.lineNumFlag))
			}

			// выводим строку, если это нужно согласно -A
			if lines2print > 0 {
				// если строка уже выведена, сохранять ее в буфере для последующего вывода не нужно
				logger.Println(printLine(lineNum, line, conf.lineNumFlag))
				alreadyPrinted = true
				lines2print--

				_ = prevLines.Pop()
			}

			// Применяем фильтры.
			matches := re.MatchString(line)

			// Если -v флаг задан
			if conf.invertFlag {
				matches = !matches
			}

			// если это неподходящая строка -> следующая
			if !matches {
				continue
			}

			// Найдена нужная строка:
			InclusionsCount++

			// Определяем сколько следующих строк печатать
			if conf.afterFlag > 0 {
				lines2print = conf.afterFlag
			}

			// печатаем предыдущие строки согласно флагу -B
			if conf.beforeFlag > 0 {
				pl := prevLines.Retrieve()
				if len(pl) > 0 {
					for _, i := range pl[:len(pl)-1] {
						logger.Println(i)
					}
				}
				// обновляем буфер, чтобы не повторялись строки при наложении областей
				prevLines = NewCircularBuffer(conf.beforeFlag + 1)
			}

			// Выводим строку если ее не вывели в начале
			if !alreadyPrinted {
				logger.Println(printLine(lineNum, line, conf.lineNumFlag))
			}

			alreadyPrinted = false
		}

		_ = file.Close()
	}

	// если был задан флаг -c, "чиним" вывод в консоль
	logger.SetOutput(writer)

	// Выводим количество совпадений (если был указан флаг -c).
	if conf.countFlag {
		logger.Println(InclusionsCount)
	}

	// Если не было ни одной совпадающей строки, выводим сообщение.
	if !printedAnyLines {
		logger.Println("Совпадений не найдено.")
	}
}

// Выведено ли что-нибудь?
var printedAnyLines bool

// printLine печатает строку с номером, если указан флаг -n.
func printLine(lineNum int, line string, printLineNum bool) string {
	printedAnyLines = true
	if printLineNum {
		return fmt.Sprintf("%d:%s", lineNum, line)
	} else {
		return fmt.Sprintf(line)
	}
}

// CircularBuffer простоя реализация циклического буфера
type CircularBuffer struct {
	buf  []any
	size int
	tail int
	len  int
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		buf:  make([]any, size),
		size: size,
		tail: 0,
		len:  0,
	}
}

func (c *CircularBuffer) Push(data any) {
	if c.size < 1 {
		return
	}

	if c.len < c.size {
		c.len++
	}

	c.buf[c.tail] = data
	c.tail = (c.tail + 1) % c.size
}

// Retrieve возвращает
func (c *CircularBuffer) Retrieve() []any {
	result := make([]any, 0, c.len)
	for {
		if v := c.Pop(); v != nil {
			result = append(result, v)
		} else {
			break
		}
	}
	slices.Reverse(result)
	return result
}

func (c *CircularBuffer) Pop() any {
	if c.len == 0 {
		return nil
	}
	c.tail--
	if c.tail < 0 {
		c.tail = c.size - 1
	}
	item := c.buf[c.tail]
	c.len--
	return item
}

func (c *CircularBuffer) cleanBuf() {
	c.buf = make([]any, c.size)
}
