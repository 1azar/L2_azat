package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

/*
Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны, например:
"a4bc2d5e" => "aaaabccddddde"
"abcd" => "abcd"
"45" => "" (некорректная строка)
"" => ""

Дополнительно
Реализовать поддержку escape-последовательностей.
Например:
qwe\4\5 => qwe45 (*)
qwe\45 => qwe44444 (*)
qwe\\5 => qwe\\\\\ (*)
В случае если была передана некорректная строка, функция должна возвращать ошибку. Написать unit-тесты.
*/

var InvalidString = errors.New("некорректная строка")

var escapeChar rune = '\u005C'

func main() {
	//str := "a4bc2d5e"
	//str := "abcd"
	//str := "45"
	//str := ""
	//str := "a30"
	str := "qwe\\4\\5"
	res, err := Unpack(str)
	fmt.Println(res, err)
}

// Unpack осуществляет примитивную распаковку строки, содержащую повторяющиеся символы/руны, например:
// "a4bc2d5e" => "aaaabccddddde"
// "abcd" => "abcd"
// "45" => "" (некорректная строка)
// "" => ""
func Unpack(str string) (string, error) {
	// переводим строку в массив рун, чтобы была возможность итерации по символам в unicode кодировке
	chars := []rune(str)
	// результат, емкость заведомо больше, чтобы исключить переаллокацию в некоторых случаях
	result := make([]rune, 0, len(chars)*2)

	runesLen := len(chars)
	for i := 0; i < runesLen; i++ {
		escapeFlag := false
		// если это экранирующий символ, устанавливаем флаг и увеличиваем итератор на 1.
		if chars[i] == escapeChar {
			escapeFlag = true
			i++
		}

		// если символ не цифра - append в результирующий слайс и переход к следующей итерации
		// если установлен флаг escapeFlag, то добавляем символ
		if !(unicode.IsDigit(chars[i])) || escapeFlag {
			result = append(result, chars[i])
			continue
		}
		// если число первое в строке (перед ним нет других символов) -> ошибка: некорректная строка
		if i == 0 {
			return "", InvalidString
		}

		// тк символ - цифра, то извлекаем число (число может быть больше 9: "abc56sad" -> 56)
		number, digitsCount, err := extractNumber(chars[i:])
		if err != nil {
			return "", err
		}
		// если число = 0, то символ перед числом повторяется 0 раз - удалям последний элемент
		if number == 0 {
			result = result[:len(result)-1]
		}
		// добавить в result букву, находящуюся перед число number, number раз
		for j := 0; j < number-1; j++ {
			result = append(result, chars[i-1])
		}
		// индекс символа, который следует за числом
		nextCharIdx := i + digitsCount
		// смещение итератора. -1 тк в произойдет автоматическая
		i = nextCharIdx - 1
	}
	return string(result), nil
}

// extractNumber принимает []rune и определяет число, состоящее из цифр с начала слайса.
// number - число, digitsCount - количество цифр в числе
// Пример:
// {"5", "9", "a", "s", ...} -> number: 59, digitsCount: 2, err: nil;
// {"8", "a", "s", ...} -> number: 8, digitsCount: 1, err: nil;
func extractNumber(r []rune) (number int, digitsCount int, err error) {
	// Слайс для сбора всех цифр, идущих подряд в слайсе с начала. Емкость = 10, чтобы исключить переаллокацию в некоторых случаях
	runeRes := make([]rune, 0, 10)
	// Цикл по слайсу пока идут цифры - ["1", "1","3","f","d"] -> 1 1 3 break
	for _, d := range r {
		if unicode.IsDigit(d) {
			runeRes = append(runeRes, d)
			continue
		}
		break
	}

	// перевод слайса рун (набор цифр) в строку (число).
	stringRes := string(runeRes)
	digitsCount = len(stringRes)
	// преобразование строки в число
	number, err = strconv.Atoi(stringRes)
	if err != nil {
		return 0, digitsCount, err
	}

	return number, digitsCount, nil
}
