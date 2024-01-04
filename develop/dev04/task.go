package main

import (
	"fmt"
	"slices"
	"strings"
)

func main() {
	dict1 := []string{"пятак", "пятка", "тяпка"}
	fmt.Println(*AnagramLookup(dict1))
}

func AnagramLookup(input []string) *map[string][]string {
	// реализация множества: ключи map
	type mySet map[string]struct{}

	// map для хранения анаграмм
	anagrams := make(map[string]mySet, len(input))

	// ключами для anagrams будут отсортированные по буквам слова и для восстановления исходного слово используется aliases
	aliases := make(map[string]string)

	for _, word := range input {
		// приводим к нижнему регистру
		word = strings.ToLower(word)

		// сортируем слово по буквам
		sortedString := sortString(word)

		// если множество для отсортированного слова уже существует, то добавляем текущее слово к нему.
		if set, ok := anagrams[sortedString]; ok {
			set[word] = struct{}{}
		} else {
			// иначе создаем новое множество анаграмм
			anagrams[sortedString] = make(mySet)
			// запоминаем исходное слово для отсортированной версии, которая является ключом множества
			aliases[sortedString] = word
		}
	}

	// результирующая map
	result := make(map[string][]string)

	// итерация по множествам анаграмм, удаление неподходящих множеств, запись в результирующую map и сортировка слов
	for key, set := range anagrams {
		if len(set) < 1 {
			// удаление множеств с 1 элементом
			delete(anagrams, key)
			continue
		}
		// записываем слова из множества в результирующий map, предварительно восстановив исходное слова для ключа
		key = aliases[key]
		result[key] = make([]string, 0, len(set))
		for wrd, _ := range set {
			result[key] = append(result[key], wrd)
		}
		// сортируем выходной слайс слов
		slices.Sort(result[key])
	}
	return &result
}

func sortString(s string) string {
	runes := []rune(s)
	slices.Sort(runes)
	return string(runes)
}
