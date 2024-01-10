package validators

import (
	"regexp"
	"time"
)

func IsValidDate(date string) bool {
	// Проверка формата даты с использованием регулярного выражения
	datePattern := `^\d{4}-\d{2}-\d{2}$`
	if match, _ := regexp.MatchString(datePattern, date); !match {
		return false
	}

	// Парсинг даты
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}

	// Проверка, что дата после парсинга соответствует исходной строке
	return parsedDate.Format("2006-01-02") == date
}
